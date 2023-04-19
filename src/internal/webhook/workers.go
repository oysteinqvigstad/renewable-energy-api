package webhook

import (
	"assignment2/internal/firebase_client"
	"log"
	"time"
)

func UpdateInvocationCountWorker() {
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
