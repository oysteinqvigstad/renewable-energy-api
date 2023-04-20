package web

import (
	"assignment2/internal/datastore"
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
			http.Error(w, info, http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
	}
}

func EnergyCurrentHandler(energyData datastore.RenewableDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			segments := utils.GetSegments(r.URL, RenewablesCurrentPath)
			neighbours, _ := utils.GetQueryStr(r.URL, "neighbours")

			switch len(segments) {
			case 0:
				httpRespondJSON(w, energyData.GetLatest("", false))
			case 1:
				returnData := energyData.GetLatest(segments[0], neighbours == "true")
				switch len(returnData) {
				case 0:
					http.Error(w, "Could not find specified country code", http.StatusBadRequest)
				case 1:
					httpRespondJSON(w, returnData[0])
				default:
					httpRespondJSON(w, returnData)
				}
			default:
				http.Error(w, "Usage: {country?}{?neighbours=bool?}", http.StatusBadRequest)
			}
		default:
			http.Error(w, "Only GET Method is supported", http.StatusBadRequest)
		}
	}
}

func EnergyHistoryHandler(energyData datastore.RenewableDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			segments := utils.GetSegments(r.URL, RenewablesHistoryPath)
			begin, _ := utils.GetQueryInt(r.URL, "begin")
			end, _ := utils.GetQueryInt(r.URL, "end")
			sort, _ := utils.GetQueryStr(r.URL, "sortByValue")

			switch len(segments) {
			case 0:
				httpRespondJSON(w, energyData.GetHistoricAvg(begin, end, sort == "true"))
			case 1:
				returnData := energyData.GetHistoric(segments[0], begin, end, sort == "true")
				if len(returnData) > 0 {
					httpRespondJSON(w, returnData)
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
}

func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	segments := utils.GetSegments(r.URL, NotificationsPath)

	switch r.Method {
	case http.MethodGet:
		switch len(segments) {
		case 0:
			viewAllWebhooks(w)
		case 1:
			viewWebhookByID(w, segments[0])
		default:
			http.Error(w, "Usage: "+NotificationsPath+"{?webhook_id}", http.StatusBadRequest)
		}

	case http.MethodPost:
		if len(segments) == 0 {
			registerWebhook(w, r)
		} else {
			http.Error(w, "Expected POST in JSON on "+NotificationsPath, http.StatusBadRequest)
		}
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
