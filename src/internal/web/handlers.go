package web

import (
	"assignment2/internal/utils"
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
			http.Error(w, info, http.StatusBadRequest)
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
		if cache, err := s.firestoreMode.GetCacheFromFirebase(r.URL); err == nil {
			httpRespondJSON(w, cache, s)
			return
		}

		segments := utils.GetSegments(r.URL, RenewablesCurrentPath)
		neighbours, _ := utils.GetQueryStr(r.URL, "neighbours")

		switch len(segments) {
		case 0:
			// Return the latest data for all countries
			httpCacheAndRespondJSON(w, r.URL, s.getCurrentRenewable("", false), s)
		case 1:
			returnData := s.getCurrentRenewable(segments[0], neighbours == "true")
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
		if cache, err := s.firestoreMode.GetCacheFromFirebase(r.URL); err == nil {
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
			httpCacheAndRespondJSON(w, r.URL, s.db.GetHistoricAvg(begin, end, sort == "true"), s)
		case 1:
			// Return the historical data for a specific country
			returnData := s.db.GetHistoric(segments[0], begin, end, sort == "true")
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
			ListWebhooksByID(w, segments[0], s)
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
		http.Error(w, "Only GET, POST and DELETE Method is supported", http.StatusBadRequest)
	}
}

// StatusHandler serves the status endpoint, providing availability info for dependent services,
// the number of registered webhooks, API version, and uptime.
func (s *State) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Get URL segments after the StatusPath
	segments := utils.GetSegments(r.URL, StatusPath)

	switch r.Method {
	case http.MethodGet:
		switch len(segments) {
		case 0:
			// Create a struct to hold the API status information
			httpRespondJSON(w, APIStatus{
				Countriesapi:    s.countriesAPIMode.getRestCountriesStatus(), // HTTP status code for *REST Countries API*
				Notification_db: s.firestoreMode.getNotificationDBStatus(),   // HTTP status code for *Notification DB* in Firebase
				Webhooks:        s.getNumberOfRegistrations(),                // Number of registered webhooks
				Version:         Version,                                     // API version
				Uptime:          utils.GetUptime(),                           // Uptime in seconds since the last service restart
			}, s)
		default:
			// Handle any other cases with URL segments
			http.Error(w, "Usage: energy/v1/status/", http.StatusBadRequest)
		}
	default:
		// Handle any unsupported HTTP methods
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}
