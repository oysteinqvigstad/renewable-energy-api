package web

import (
	"assignment2/internal/firebase_client"
	"log"
	"math/rand"
	"time"
)

var (
	invocationCount     map[string]int
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
	invocationCount = make(map[string]int)
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

}

func Invocate(ccna3 []string) {
	for _, code := range ccna3 {
		invocationCount[code]++
		invocateChannel <- code
	}
	// TODO: write go routine for triggering webhook
}

func AddWebhookRegistration(newReg firebase_client.InvocationRegistration) (string, bool) {
	newReg.WebhookID = generateWebhookID()
	if _, ok := registrations[newReg.WebhookID]; !ok {
		data := firebase_client.RegistrationAction{
			Add:          true,
			Registration: newReg,
		}
		registrations[newReg.WebhookID] = newReg
		registrationChannel <- data
		return newReg.WebhookID, true
	}
	return "", false
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
	ticker := time.NewTicker(5 * time.Second)
	client, err := firebase_client.NewFirebaseClient()
	if err != nil {
		log.Fatal("Could not start firebase client")
	}

	updates := firebase_client.NewBundledUpdate()

	for {
		select {
		// increment the invocation count and add it to the list of
		// countries that should be updated
		case code := <-invocateChannel:
			updates.InvocationCount[code] = invocationCount[code]
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
