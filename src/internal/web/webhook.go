package web

import (
	"assignment2/internal/datastore"
	"assignment2/internal/firebase_client"
	"assignment2/internal/web_client"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	invocationCount     map[string]int64
	registrations       map[string]firebase_client.InvocationRegistration
	invocateChannel     chan string
	registrationChannel chan firebase_client.RegistrationAction
	cacheChannel        chan map[string]datastore.YearRecordList
	serviceStarted      bool
)

// InitializeWebhookService is a function that initializes the webhook service.
// It prevents the service from being started multiple times, initiates data structures,
// retrieves registrations from Firestore, and starts the firebaseWorker goroutine.
func InitializeWebhookService() {
	if serviceStarted == true {
		log.Fatal("Webhook service cannot be started twice")
	}
	initializeDataStructures()
	loadRegistrationsFromFirestore()
	go firebaseWorker()

}

// initializeDataStructures is a function that initializes the necessary data structures
// for the webhook service, such as invocation count, registrations, and channels.
func initializeDataStructures() {
	invocationCount = make(map[string]int64)
	registrations = make(map[string]firebase_client.InvocationRegistration)
	invocateChannel = make(chan string)
	registrationChannel = make(chan firebase_client.RegistrationAction)
	cacheChannel = make(chan map[string]datastore.YearRecordList)
}

// loadRegistrationsFromFirestore is a function that retrieves the invocation counts
// and registrations from Firestore using the Firebase client and populates the respective data structures.
func loadRegistrationsFromFirestore() {
	// initiate data structures
	println("Downloading counters and registrations from firestore")
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Println("Could not initiate firebase client")
	}
	invocationCount = client.GetAllInvocationCounts()
	registrations = client.GetAllInvocationRegistrations()
}

// ProcessWebhookByCountry is a function that processes a list of country codes, increments their invocation count,
// and triggers webhooks accordingly.
func ProcessWebhookByCountry(ccna3 []string, db *datastore.RenewableDB) {
	for _, code := range ccna3 {
		// TODO: Add mutex lock or other mechanism
		invocationCount[code]++
		invocateChannel <- code
		triggerWebhooksForCountry(code, invocationCount[code], db.GetName(code))
	}
}

// generateWebhookID is a function that generates a unique 13-character
// alphanumeric string to be used as a webhook ID.
func generateWebhookID() string {
	// Define the character set for the webhook ID
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Set the desired length of the webhook ID
	length := 13
	b := make([]byte, length)

	// Continue as long as the newly generated Webhook ID has not been proven to be unique
	for {
		var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}

		// Check if the generated webhook ID already exists in the registrations map.
		// If it doesn't exist (i.e., it's unique), break the loop and return the ID.
		if _, ok := registrations[string(b)]; !ok {
			break
		}
	}
	return string(b)
}

// firebaseWorker is a function that runs in the background and periodically
// sends updates to Firebase. It listens for invocation count and registration
// updates and sends them to Firebase in bulk every FirebaseUpdateFreq seconds.
func firebaseWorker() {
	// Create a new ticker to trigger Firebase updates at a fixed interval
	ticker := time.NewTicker(FirebaseUpdateFreq * time.Second)

	// Create a new Firebase client
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Fatal("Could not start firebase client")
	}

	// Initialize an empty bundled update to store updates before sending them to Firebase
	updates := firebase_client.NewBundledUpdate()

	// Start an infinite loop to listen for updates and send them to Firebase
	for {
		select {
		// If there's a new invocation count update, add it to the updates
		// and set the Ready flag to indicate that updates are available
		case countryCode := <-invocateChannel:
			updates.InvocationCount[countryCode] = invocationCount[countryCode]
			updates.Ready = true

		// If there's a new registration update, add it to the updates
		// and set the Ready flag to indicate that updates are available
		case action := <-registrationChannel:
			updates.Registrations[action.Registration.WebhookID] = action
			updates.Ready = true

		case cache := <-cacheChannel:
			for key, value := range cache {
				updates.Cache[key] = value
			}
			updates.Ready = true

		// When the ticker triggers, check if there are updates to send and
		// send them to Firebase in bulk if there are. Reset the updates
		// and Ready flag afterward.
		case <-ticker.C:
			if updates.Ready {
				client.BulkWrite(updates)
				updates = firebase_client.NewBundledUpdate()
			}
		}
	}
}

// registerWebhook is an HTTP handler function that registers a new webhook
// by decoding an incoming JSON request and validating the data.
func registerWebhook(w http.ResponseWriter, r *http.Request) {
	// create a new InvocationRegistration struct to store the decoded data
	var data firebase_client.InvocationRegistration

	// decode the request body into the InvocationRegistration struct
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error during decoding: " + err.Error())
		return
	}

	// validating the JSON input
	if err := validateRegistrationData(data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// adding registration to data structure and notifying firestore that the registration
	// can be backup up
	data.WebhookID = generateWebhookID()
	newEntry := firebase_client.RegistrationAction{Add: true, Registration: data}
	registrations[data.WebhookID] = data
	registrationChannel <- newEntry
	w.WriteHeader(http.StatusCreated)
	httpRespondJSON(w, map[string]interface{}{"webhook_id": data.WebhookID}, nil)
}

// validateRegistrationData is a function that takes an InvocationRegistration
// struct as input and returns an error if the registration data is invalid.
func validateRegistrationData(registration firebase_client.InvocationRegistration) error {
	// Check if the number of calls is less than 1. If so, return an error.
	if registration.Calls < 1 {
		return errors.New("number of calls must be 1 or higher")
	}

	// Check if the URL is properly formatted with http:// or https:// prefix.
	// If not, return an error.
	if !strings.HasPrefix(registration.URL, "http://") && !strings.HasPrefix(registration.URL, "https://") {
		return errors.New("URL must be prefixed by http:// or https://")
	}

	// Check if the country is recognized (i.e., if it exists in the invocationCount map).
	// If not, return an error.
	if _, ok := invocationCount[registration.Country]; !ok {
		return errors.New("country not recognized")
	}

	// If all checks pass, return nil, indicating no error occurred.
	return nil
}

// listAllWebhooks is a function that retrieves all registered webhooks
// and sends them as a JSON response to the client.
func listAllWebhooks(w http.ResponseWriter) {
	registrationList := make([]firebase_client.InvocationRegistration, 0, len(registrations))
	for _, registration := range registrations {
		registrationList = append(registrationList, registration)
	}
	httpRespondJSON(w, registrationList, nil)
}

// listAllWebhooksByID is a function that retrieves a registered webhook by its ID
// and sends it as a JSON response to the client if found, otherwise it sends an error.
func listAllWebhooksByID(w http.ResponseWriter, webhookID string) {
	if reg, ok := registrations[webhookID]; ok {
		httpRespondJSON(w, reg, nil)
		return
	}
	http.Error(w, "Could not find the webhook ID: "+webhookID, http.StatusBadRequest)
}

// RemoveWebhookByID is a function that removes a webhook registration by its ID.
// If the webhook is found and deleted successfully, it sends a "No Content" status,
// otherwise, it sends an error.
func RemoveWebhookByID(w http.ResponseWriter, webhookID string) {
	if record, ok := registrations[webhookID]; ok {
		delete(registrations, webhookID)
		update := firebase_client.RegistrationAction{Add: false, Registration: record}
		registrationChannel <- update
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "Could not find the webhookID: "+webhookID, http.StatusBadRequest)
}

// triggerWebhooksForCountry is a function that iterates through all registered webhooks
// and triggers them if the specified country code matches and the call count
// reaches the specified threshold.
//
// TODO: Refactor to use a worker thread!
// This is not the best way to handle webhooks, as it needs to check
// every registration every time a country code is invoked.
func triggerWebhooksForCountry(countrycode string, count int64, name string) {
	for _, reg := range registrations {
		if reg.Country == countrycode {
			if count > 0 && count%reg.Calls == 0 {
				go postToWebhook(reg.URL, map[string]interface{}{
					"webhook_id": reg.WebhookID,
					"country":    name,
					"calls":      count,
				})
			}
		}
	}
}

// postToWebhook is a function that sends a POST request to the specified webhook URL
// with the provided registration data as the request body.
func postToWebhook(url string, registration map[string]interface{}) {
	client := web_client.NewClient()
	err := client.SetURL(url)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(registration)
	if err != nil {
		log.Println("Could not encode to json, ", err.Error())
	}
	_, _ = client.Post(&buf)
}
