package web

import (
	"assignment2/api"
	"assignment2/internal/firebase_client"
	"assignment2/internal/types"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// State represents the application state and holds the necessary data and channels.
type State struct {
	db               types.RenewableDB
	invocationCounts map[string]int64
	registrations    map[string]types.InvocationRegistration
	firestoreMode    firestoreMode
	countriesAPIMode restCountriesMode
	lock             sync.RWMutex
	chInvocation     chan string
	chRegistration   chan types.RegistrationAction
	chCache          chan map[string]types.YearRecordList
}

// NewService initializes a new State with the provided CSV filepath and mode.
func NewService(filepath string, countriesMode restCountriesMode, firebaseMode firestoreMode) *State {
	s := State{
		db:               types.ParseCSV(filepath),
		invocationCounts: firebaseMode.GetAllInvocationCounts(),
		registrations:    firebaseMode.GetAllInvocationRegistrations(),
		firestoreMode:    firebaseMode,
		countriesAPIMode: countriesMode,
	}

	// Initialize channels and start the worker for updating Firebase in WithFirestore firebaseMode
	switch firebaseMode.(type) {
	case WithFirestore:
		s.chInvocation = make(chan string, 1000)
		s.chRegistration = make(chan types.RegistrationAction, 10)
		s.chCache = make(chan map[string]types.YearRecordList, 100)
		go firebaseUpdateWorker(&s)
	}

	return &s
}

// WithoutFirestore represents a mode where the Firestore interaction is disabled.
type WithoutFirestore struct{}

// WithFirestore represents a mode where the Firestore interaction is enabled.
type WithFirestore struct{}

// StubRestCountries represents a mode where the Countries API is stubbed.
type StubRestCountries struct{}

// UseRestCountries represents a mode where the 3rd party API is used.
type UseRestCountries struct{}

// firestoreMode defines an interface for managing cache and invocation counts. Updates are handled
// by channels and go routine instead.
type firestoreMode interface {
	GetCacheFromFirebase(url *url.URL) (types.YearRecordList, error)
	GetAllInvocationCounts() map[string]int64
	GetAllInvocationRegistrations() map[string]types.InvocationRegistration
	getNotificationDBStatus() int
}

// deleteRegistration removes a registration from the state's registrations map by its webhookID.
// If the registration is found, it is deleted and the deletion is updated in Firestore. If the
// registration is not found, it returns an error.
func (s *State) deleteRegistration(webhookID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if registration, ok := s.registrations[webhookID]; ok {
		delete(s.registrations, webhookID)
		updateFirestore(s.chRegistration, types.RegistrationAction{Add: false, Registration: registration})
		return nil
	} else {
		return errors.New("Could not find the webhookID: " + webhookID)
	}
}

// getNumberOfRegistrations returns the number of registrations
// currently stored in the state's registrations map.
func (s *State) getNumberOfRegistrations() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.registrations)
}

// getRegistration retrieves a registration from the state's registrations map by its webhookID.
// If the registration is found, it returns the registration and true, otherwise, it returns an
// empty registration object and false.
func (s *State) getRegistration(webhookID string) (types.InvocationRegistration, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	value, ok := s.registrations[webhookID]
	return value, ok
}

// getAllRegistrations returns a slice of all registrations stored in the state's registrations map.
func (s *State) getAllRegistrations() []types.InvocationRegistration {
	s.lock.RLock()
	defer s.lock.RUnlock()
	list := make([]types.InvocationRegistration, 0, len(s.registrations))
	for _, value := range s.registrations {
		list = append(list, value)
	}
	return list
}

// newRegistration adds a new registration to the state's registrations map and updates Firestore
// with the new entry.
func (s *State) newRegistration(registration types.InvocationRegistration) {
	newEntry := types.RegistrationAction{Add: true, Registration: registration}
	updateFirestore(s.chRegistration, newEntry)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.registrations[registration.WebhookID] = registration
}

// incrementInvocationCount increments the invocation count for a given countryCode and returns
// the updated count. The invocation count is stored in the state's invocationCounts map.
func (s *State) incrementInvocationCount(countryCode string) int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.invocationCounts[countryCode]++
	return s.invocationCounts[countryCode]
}

// getInvocationCount retrieves the invocation count for a given countryCode from the state's
// invocationCounts map and returns the count.
func (s *State) getInvocationCount(countryCode string) int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.invocationCounts[countryCode]
}

func (s *State) getCurrentRenewable(countryCode string, includeNeighbours bool) types.YearRecordList {
	data := s.db.RetrieveLatest(countryCode)
	if len(countryCode) > 0 && includeNeighbours {
		neighbours, err := s.countriesAPIMode.getNeighboursCca(countryCode)
		if err == nil {
			for _, neighbour := range neighbours {
				data = append(data, s.db.RetrieveLatest(neighbour)...)
			}
		}
	}
	return data
}

// restCountriesMode defines an interface for either using the stubbed RESTCountries service or the
// real 3rd party service
type restCountriesMode interface {
	getNeighboursCca(cca string) ([]string, error)
	getRestCountriesStatus() int
}

// getNeighboursCca returns neighbouring list from stubbed country api
func (t StubRestCountries) getNeighboursCca(cca string) ([]string, error) {
	val, err := api.GetNeighboursCca(cca, api.STUB_BASE)
	return val, err
}

// getNeighboursCca returns neighbouring list from 3rd party country api
func (p UseRestCountries) getNeighboursCca(cca string) ([]string, error) {
	val, err := api.GetNeighboursCca(cca, api.API_BASE)
	return val, err
}

func (t StubRestCountries) getRestCountriesStatus() int {
	return getStatusCode(api.STUB_BASE + api.API_VERSION + "/alpha/nor")
}

func (p UseRestCountries) getRestCountriesStatus() int {
	return getStatusCode(api.API_BASE + api.API_VERSION + "/alpha/nor")
}

// GetCacheFromFirebase returns an error as caching is disabled in WithoutFirestore mode.
func (t WithoutFirestore) GetCacheFromFirebase(_ *url.URL) (types.YearRecordList, error) {
	return types.YearRecordList{}, errors.New("firebase disabled")
}

// GetCacheFromFirebase retrieves the cached data from Firestore in WithFirestore mode.
func (p WithFirestore) GetCacheFromFirebase(url *url.URL) (types.YearRecordList, error) {
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Println("Could not start firebase client")
		return types.YearRecordList{}, errors.New("could not start firebase client")
	}
	defer client.Close()
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
		if err != nil {
			log.Println("Could not start firebase client")
			return data
		}
		defer client.Close()
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
		if err != nil {
			log.Println("Could not start firebase client")
			return result
		}
		defer client.Close()
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

// GetNotificationDBStatus for WithFirestore type returns Notification DB status code.
// It handles client creation errors and returns appropriate status codes.
func (p WithFirestore) getNotificationDBStatus() int {
	client, err := firebase_client.NewFirebaseClient()
	defer client.Close()
	if err != nil {
		// Handle any error from the Firebase client
		return http.StatusInternalServerError
	}
	// Set the Notification DB status to OK if no error
	return http.StatusOK
}

// GetNotificationDBStatus for WithoutFirestore type returns a
// service unavailable status code, as Notification DB isn't used in this mode.
func (t WithoutFirestore) getNotificationDBStatus() int {
	return http.StatusServiceUnavailable
}
