package web

import (
	"assignment2/internal/db"
	"assignment2/internal/utils"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			http.Error(w, "Not implemented", http.StatusInternalServerError)
			// TODO: implement simple informational handler
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

func EnergyCurrentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
		//return // uncomment this for testing

		segments := utils.GetSegments(r.URL, RenewablesCurrentPath)
		neighbours, _ := utils.GetQueryStr(r.URL, "neighbours")

		switch len(segments) {
		case 0:
			httpRespondJSON(w, db.GetCurrentEnergyData("", false))
		case 1:
			httpRespondJSON(w, db.GetCurrentEnergyData(segments[0], neighbours == "true"))
		default:
			http.Error(w, "Usage: {country?}{?neighbours=bool?}", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

func EnergyHistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
		// TODO: get URL segment {country} (3 letter code)
		// TODO: get query {begin}, {end} and {sortByValue}
		// TODO: write fetch function for dataset
		// TODO: if {sortByValue} is set -> Sort all the
		// TODO: if {begin} IS set -> omit year attribute in country struct (returns single average)
		// TODO: if {country} IS set -> return list of structs for that country
		// TODO: if {country} IS NOT set -> return all data? Will be very large return
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	case http.MethodPost:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	case http.MethodDelete:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, "Unimplemented", http.StatusServiceUnavailable)
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}
