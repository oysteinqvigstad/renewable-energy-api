package web

import (
	"assignment2/internal/firebase_client"
	"encoding/json"
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
	serviceStarted      bool
)

func StartWebhookService() {
	if serviceStarted == true {
		log.Fatal("Webhook service cannot be started twice")
	}
	initiateDatastructures()
	retrieveRegistrationsFromFirestore()
	go firebaseWorker()

}

func initiateDatastructures() {
	invocationCount = make(map[string]int64)
	registrations = make(map[string]firebase_client.InvocationRegistration)
	invocateChannel = make(chan string)
	registrationChannel = make(chan firebase_client.RegistrationAction)
}

func retrieveRegistrationsFromFirestore() {
	// initiate data structures
	println("Downloading counters and registrations from firestore")
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Println("Could not initiate firebase client")
	}
	invocationCount = client.GetAllInvocationCounts()
	registrations = client.GetAllInvocationRegistrations()
}

func Invocate(ccna3 []string) {
	for _, code := range ccna3 {
		invocationCount[code]++
		invocateChannel <- code
	}
	// TODO: write go routine for triggering webhook
}

func DelWebhookRegistration(webhookId string) bool {
	if _, ok := registrations[webhookId]; ok {
		delete(registrations, webhookId)
		data := firebase_client.RegistrationAction{
			Add: false,
			Registration: firebase_client.InvocationRegistration{
				WebhookID: webhookId,
			},
		}
		registrationChannel <- data
		return true
	}
	return false
}

func generateWebhookID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 13
	b := make([]byte, length)
	for {
		var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
		// ensuring WebhookID is unique
		if _, ok := registrations[string(b)]; !ok {
			break
		}
	}
	return string(b)
}

func firebaseWorker() {
	ticker := time.NewTicker(FirebaseUpdateFreq * time.Second)
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Fatal("Could not start firebase client")
	}

	updates := firebase_client.NewBundledUpdate()

	for {
		select {
		// increment the invocation count and add it to the list of
		// countries that should be updated
		case countryCode := <-invocateChannel:
			updates.InvocationCount[countryCode] = invocationCount[countryCode]
			updates.Ready = true

		case action := <-registrationChannel:
			updates.Registrations[action.Registration.WebhookID] = action
			updates.Ready = true

			// every 5 seconds, send updates in bulk
		case <-ticker.C:
			if updates.Ready {
				client.BulkWrite(updates)
				updates = firebase_client.NewBundledUpdate()

			}

		}
	}
}

func registerWebhook(w http.ResponseWriter, r *http.Request) {
	var data firebase_client.InvocationRegistration
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Error during decoding: " + err.Error())
		return
	}

	if data.Calls < 1 {
		http.Error(w, "Number of calls must be 1 or higher", http.StatusBadRequest)
	}

	if !strings.HasPrefix(data.URL, "http://") && !strings.HasPrefix(data.URL, "https://") {
		http.Error(w, "URL must be prefixed by http:// or https://", http.StatusBadRequest)
	}

	if _, ok := invocationCount[data.Country]; !ok {
		http.Error(w, "Country not recognized", http.StatusBadRequest)
	}

	data.WebhookID = generateWebhookID()
	newEntry := firebase_client.RegistrationAction{
		Add:          true,
		Registration: data,
	}
	registrations[data.WebhookID] = data
	registrationChannel <- newEntry

	for key, val := range registrations {
		fmt.Println(key, val)
	}

	httpRespondJSON(w, map[string]interface{}{"webhook_id": data.WebhookID})
}
