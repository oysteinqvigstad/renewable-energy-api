package web

import (
	"assignment2/api"
	"assignment2/internal/firebase_client"
	"assignment2/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			info := "Usage:\n" +
				"/energy/v1/renewables/current/{country?}{?neighbours=bool?}\n" +
				"/energy/v1/renewables/history/{country?}{?begin=year&end=year?}{?sortByValue=bool?}\n" +
				"/energy/v1/notifications\n" +
				"/energy/v1/status\n"
			http.Error(w, info, http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

// EnergyCurrentHandler handles the request for the current energy data.
// It retrieves the data from cache if available, otherwise it fetches the data from the database.
func (s *State) EnergyCurrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		// Checking cache first
		if cache, err := s.Mode.GetCacheFromFirebase(r.URL); err == nil {
			httpRespondJSON(w, cache, s)
			return
		}

		segments := utils.GetSegments(r.URL, RenewablesCurrentPath)
		neighbours, _ := utils.GetQueryStr(r.URL, "neighbours")

		switch len(segments) {
		case 0:
			// Return the latest data for all countries
			httpCacheAndRespondJSON(w, r.URL, s.DB.GetLatest("", false), s)
		case 1:
			returnData := s.DB.GetLatest(segments[0], neighbours == "true")
			switch len(returnData) {
			case 0:
				// Return the latest data for a specific country
				http.Error(w, "Could not find specified country code", http.StatusBadRequest)
			default:
				httpCacheAndRespondJSON(w, r.URL, returnData, s)
			}
		default:
			http.Error(w, "Usage: {country?}{?neighbours=bool?}", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

// EnergyHistoryHandler handles the request for the historical energy data.
// It retrieves the data from cache if available, otherwise it fetches the data from the database.
func (s *State) EnergyHistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		// Check cache first
		if cache, err := s.Mode.GetCacheFromFirebase(r.URL); err == nil {
			println("got from cache!!!!")
			httpRespondJSON(w, cache, s)
			return
		}

		segments := utils.GetSegments(r.URL, RenewablesHistoryPath)
		begin, _ := utils.GetQueryInt(r.URL, "begin")
		end, _ := utils.GetQueryInt(r.URL, "end")
		sort, _ := utils.GetQueryStr(r.URL, "sortByValue")

		switch len(segments) {
		case 0:
			// Return the historical average data for all countries
			httpCacheAndRespondJSON(w, r.URL, s.DB.GetHistoricAvg(begin, end, sort == "true"), s)
		case 1:
			// Return the historical data for a specific country
			returnData := s.DB.GetHistoric(segments[0], begin, end, sort == "true")
			if len(returnData) > 0 {
				httpCacheAndRespondJSON(w, r.URL, returnData, s)
			} else {
				http.Error(w, "Could not find specified country code", http.StatusBadRequest)
			}
		default:
			http.Error(w, "Usage: {country?}{?neighbours=bool?}", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

// NotificationHandler handles the request for managing webhook notifications.
// It supports GET, POST, and DELETE methods for listing, registering, and removing webhooks, respectively.
func (s *State) NotificationHandler(w http.ResponseWriter, r *http.Request) {
	segments := utils.GetSegments(r.URL, NotificationsPath)

	switch r.Method {
	case http.MethodGet:
		switch len(segments) {
		case 0:
			// List all registered webhooks
			listAllWebhooks(w, s)
		case 1:
			// List a specific webhook by its ID
			listAllWebhooksByID(w, segments[0], s)
		default:
			http.Error(w, "Usage: "+NotificationsPath+"{?webhook_id}", http.StatusBadRequest)
		}
	case http.MethodPost:
		switch len(segments) {
		case 0:
			// Register a new webhook
			registerWebhook(w, r, s)
		default:
			http.Error(w, "Expected POST in JSON on "+NotificationsPath, http.StatusBadRequest)
		}
	case http.MethodDelete:
		switch len(segments) {
		case 1:
			// Remove a webhook by its ID
			RemoveWebhookByID(w, segments[0], s)
		default:
			http.Error(w, "Usage: "+NotificationsPath+"{id}", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

// StatusHandler serves the status endpoint, providing availability info for dependent services,
// the number of registered webhooks, API version, and uptime.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Get URL segments after the StatusPath
	segments := utils.GetSegments(r.URL, StatusPath)
	// Calculate the uptime in seconds since the last service restart
	uptime := utils.GetUptime()

	switch r.Method {
	case http.MethodGet:
		switch len(segments) {
		case 0:
			// Define the countries API URL
			countriesAPI := api.API_BASE + api.API_VERSION + "/" + "all"
			// Send a request to the countries API
			resp, err := http.Get(countriesAPI)
			if err != nil {
				// Handle any error from the countries API request
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Safely close the response body and log any errors that occur
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}

			// Create a Firebase client to get registered webhooks
			client, err := firebase_client.NewFirebaseClient()
			var notificationDBStatus int
			if err != nil {
				// Handle any error from the Firebase client
				http.Error(w, err.Error(), http.StatusInternalServerError)
				notificationDBStatus = http.StatusInternalServerError
				return
			} else {
				// Set the Notification DB status to OK if no error
				notificationDBStatus = http.StatusOK
			}
			// Get the invocation count for all registered webhooks
			invocationCount := client.GetAllInvocationCounts()
			// Create a struct to hold the API status information
			allAPI := APIStatus{
				Countriesapi:    resp.StatusCode,      // HTTP status code for *REST Countries API*
				Notification_db: notificationDBStatus, // HTTP status code for *Notification DB* in Firebase
				Webhooks:        len(invocationCount), // Number of registered webhooks
				Version:         api.API_VERSION,      // API version
				Uptime:          uptime,               // Uptime in seconds since the last service restart
			}
			// Set the response content type to JSON
			w.Header().Set("Content-Type", "application/json")
			// Encode and return the API status as JSON, handling any errors that occur during encoding
			if err := json.NewEncoder(w).Encode(allAPI); err != nil {
				http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			}

		default:
			// Handle any other cases with URL segments
			w.Header().Set("content-type", "text/html")
			output := "Usage: energy/v1/status/"
			_, err := fmt.Fprintf(w, "%v", output)
			if err != nil {
				// Handle any error when returning the output
				http.Error(w, "Error when returning output", http.StatusInternalServerError)
			}
			return
		}
	default:
		// Handle any unsupported HTTP methods
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}
