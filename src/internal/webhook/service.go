package webhook

import (
	"assignment2/internal/firebase_client"
	"log"
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
	go UpdateInvocationCountWorker()

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
}

func AddWebhookRegistration(newReg firebase_client.InvocationRegistration) bool {
	if _, ok := registrations[newReg.WebhookID]; ok {
		data := firebase_client.RegistrationAction{
			Add:          true,
			Registration: newReg,
		}
		registrations[newReg.WebhookID] = newReg
		registrationChannel <- data
	}
	return false
}

func DelWebhookRegistration(webhook_id string) bool {
	if _, ok := registrations[webhook_id]; ok {
		delete(registrations, webhook_id)
		data := firebase_client.RegistrationAction{
			Add: false,
			Registration: firebase_client.InvocationRegistration{
				WebhookID: webhook_id,
			},
		}
		registrationChannel <- data
		return true
	}
	return false
}
