package web

import (
	"assignment2/internal/firebase_client"
	"assignment2/internal/types"
	"errors"
	"log"
	"net/url"
	"time"
)

// State represents the application state and holds the necessary data and channels.
type State struct {
	DB               types.RenewableDB
	InvocationCounts map[string]int64
	Registrations    map[string]types.InvocationRegistration
	Mode             Mode
	ChInvocation     chan string
	ChRegistration   chan types.RegistrationAction
	ChCache          chan map[string]types.YearRecordList
}

// NewService initializes a new State with the provided CSV filepath and mode.
func NewService(filepath string, mode Mode) *State {
	s := State{
		DB:               types.ParseCSV(filepath),
		InvocationCounts: mode.GetAllInvocationCounts(),
		Registrations:    mode.GetAllInvocationRegistrations(),
		Mode:             mode,
	}

	// Initialize channels and start the worker for updating Firebase in WithFirestore mode
	switch mode.(type) {
	case WithFirestore:
		s.ChInvocation = make(chan string, 1000)
		s.ChRegistration = make(chan types.RegistrationAction, 10)
		s.ChCache = make(chan map[string]types.YearRecordList, 100)
		go firebaseUpdateWorker(&s)
	}

	return &s
}

// WithoutFirestore represents a mode where the Firestore interaction is disabled.
type WithoutFirestore struct{}

// WithFirestore represents a mode where the Firestore interaction is enabled.
type WithFirestore struct{}

// Mode defines an interface for managing cache and invocation counts. Updates are handled
// by channels and go routine instead.
type Mode interface {
	GetCacheFromFirebase(url *url.URL) (types.YearRecordList, error)
	GetAllInvocationCounts() map[string]int64
	GetAllInvocationRegistrations() map[string]types.InvocationRegistration
}

// GetCacheFromFirebase returns an error as caching is disabled in WithoutFirestore mode.
func (t WithoutFirestore) GetCacheFromFirebase(_ *url.URL) (types.YearRecordList, error) {
	return types.YearRecordList{}, errors.New("firebase disabled")
}

// GetCacheFromFirebase retrieves the cached data from Firestore in WithFirestore mode.
func (p WithFirestore) GetCacheFromFirebase(url *url.URL) (types.YearRecordList, error) {
	println("attempting to get from cache ", url.String())
	client, err := firebase_client.NewFirebaseClient()
	// TODO: Handle timestamp?
	data, timestamp, err := client.GetRenewablesCache(url.String())
	twoMinutesAgo := time.Now().Add(-24 * 7 * time.Hour)
	if timestamp.Before(twoMinutesAgo) {
		client.DeleteRenewablesCache(url.String())
		return data, errors.New("old cache, deleted object")

	}
	if err != nil {
		return data, errors.New("cache not found")
	}
	return data, nil
}

// GetAllInvocationCounts retrieves all invocation counts from Firestore in WithFirestore mode.
func (p WithFirestore) GetAllInvocationCounts() map[string]int64 {
	data := map[string]int64{}
	if client, err := firebase_client.NewFirebaseClient(); err == nil {
		docs, err := client.GetAllDocuments(firebase_client.CollectionInvocationCounts)
		if err != nil {
			log.Printf("Could not fetch invocation counts from firestore")
			return data
		}
		for _, docField := range docs {
			if count, err := docField.DataAt("count"); err == nil {
				data[docField.Ref.ID] = count.(int64)
			}
		}
	}
	return data
}

// GetAllInvocationCounts returns an empty map as the feature is disabled in WithoutFirestore mode.
func (t WithoutFirestore) GetAllInvocationCounts() map[string]int64 {
	return map[string]int64{}
}

// GetAllInvocationRegistrations retrieves all invocation registrations from Firestore in WithFirestore mode.
func (p WithFirestore) GetAllInvocationRegistrations() map[string]types.InvocationRegistration {
	result := map[string]types.InvocationRegistration{}
	if client, err := firebase_client.NewFirebaseClient(); err == nil {
		docs, err := client.GetAllDocuments(firebase_client.CollectionInvocationRegistrations)
		if err != nil {
			log.Printf("Could not fetch data from firestore")
			return result
		}
		for _, doc := range docs {
			var registration types.InvocationRegistration
			err = doc.DataTo(&registration)
			if err != nil {
				return result
			}
			result[doc.Ref.ID] = registration
		}
	}
	return result
}

// GetAllInvocationRegistrations returns an empty map as invocation registrations are not supported in WithoutFirestore mode.
func (t WithoutFirestore) GetAllInvocationRegistrations() map[string]types.InvocationRegistration {
	return map[string]types.InvocationRegistration{}

}
