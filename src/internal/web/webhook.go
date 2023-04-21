package web

import (
	"assignment2/internal/firebase_client"
	"assignment2/internal/types"
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

// ProcessWebhookByCountry is a function that processes a list of country codes, increments their invocation count,
// and triggers webhooks accordingly.
func ProcessWebhookByCountry(ccna3 []string, s *State) {
	for _, code := range ccna3 {
		s.InvocationCounts[code]++
		updateFirestore(s.ChInvocation, code)
		triggerWebhooksForCountry(code, s.InvocationCounts[code], s.DB.GetName(code), s)
	}
}

// generateWebhookID is a function that generates a unique 13-character
// alphanumeric string to be used as a webhook ID.
func generateWebhookID(s *State) string {
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
		if _, ok := s.Registrations[string(b)]; !ok {
			break
		}
	}
	return string(b)
}

// firebaseUpdateWorker is a function that runs in the background and periodically
// sends updates to Firebase. It listens for invocation count and registration
// updates and sends them to Firebase in bulk every FirebaseUpdateFreq seconds.
func firebaseUpdateWorker(s *State) {
	// Create a new ticker to trigger Firebase updates at a fixed interval
	ticker := time.NewTicker(FirebaseUpdateFreq * time.Second)

	// Create a new Firebase client
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Fatal("Could not start firebase client")
	}
	defer client.Close()

	// Initialize an empty bundled update to store updates before sending them to Firebase
	updates := firebase_client.NewBundledUpdate()

	// Start an infinite loop to listen for updates and send them to Firebase
	for {
		select {
		// If there's a new invocation count update, add it to the updates
		// and set the Ready flag to indicate that updates are available
		case countryCode := <-s.ChInvocation:
			updates.InvocationCount[countryCode] = s.InvocationCounts[countryCode]
			updates.Ready = true

		// If there's a new registration update, add it to the updates
		// and set the Ready flag to indicate that updates are available
		case action := <-s.ChRegistration:
			updates.Registrations[action.Registration.WebhookID] = action
			updates.Ready = true

		case cache := <-s.ChCache:
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
func registerWebhook(w http.ResponseWriter, r *http.Request, s *State) {
	// create a new InvocationRegistration struct to store the decoded data
	var data types.InvocationRegistration

	// decode the request body into the InvocationRegistration struct
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error during decoding: " + err.Error())
		return
	}

	// validating the JSON input
	if err := validateRegistrationData(data, s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// adding registration to data structure and notifying firestore that the registration
	// can be backup up
	data.WebhookID = generateWebhookID(s)
	newEntry := types.RegistrationAction{Add: true, Registration: data}
	s.Registrations[data.WebhookID] = data
	updateFirestore(s.ChRegistration, newEntry)
	w.WriteHeader(http.StatusCreated)
	httpRespondJSON(w, map[string]interface{}{"webhook_id": data.WebhookID}, nil)
}

// validateRegistrationData is a function that takes an InvocationRegistration
// struct as input and returns an error if the registration data is invalid.
func validateRegistrationData(registration types.InvocationRegistration, s *State) error {
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
	if _, ok := s.DB[registration.Country]; !ok {
		return errors.New("country not recognized")
	}

	// If all checks pass, return nil, indicating no error occurred.
	return nil
}

// listAllWebhooks is a function that retrieves all registered webhooks
// and sends them as a JSON response to the client.
func listAllWebhooks(w http.ResponseWriter, s *State) {
	registrationList := make([]types.InvocationRegistration, 0, len(s.Registrations))
	for _, registration := range s.Registrations {
		registrationList = append(registrationList, registration)
	}
	httpRespondJSON(w, registrationList, nil)
}

// listAllWebhooksByID is a function that retrieves a registered webhook by its ID
// and sends it as a JSON response to the client if found, otherwise it sends an error.
func listAllWebhooksByID(w http.ResponseWriter, webhookID string, s *State) {
	if reg, ok := s.Registrations[webhookID]; ok {
		httpRespondJSON(w, reg, nil)
		return
	}
	http.Error(w, "Could not find the webhook ID: "+webhookID, http.StatusBadRequest)
}

// RemoveWebhookByID is a function that removes a webhook registration by its ID.
// If the webhook is found and deleted successfully, it sends a "No Content" status,
// otherwise, it sends an error.
func RemoveWebhookByID(w http.ResponseWriter, webhookID string, s *State) {
	if record, ok := s.Registrations[webhookID]; ok {
		delete(s.Registrations, webhookID)
		update := types.RegistrationAction{Add: false, Registration: record}
		updateFirestore(s.ChRegistration, update)
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
func triggerWebhooksForCountry(countrycode string, count int64, name string, s *State) {
	for _, reg := range s.Registrations {
		if reg.Country == countrycode {
			if count > 0 && count%reg.Calls == 0 {
				postToWebhook(reg.URL, WebhookResponse{
					WebhookID: reg.WebhookID,
					Country:   name,
					Calls:     count,
				})
			}
		}
	}
}

// postToWebhook is a function that sends a POST request to the specified webhook URL
// with the provided registration data as the request body.
func postToWebhook(url string, registration WebhookResponse) {
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
