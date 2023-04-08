package stub_countries_api

import (
	"assignment2/internal/utils"
	"encoding/json"
	"net/http"
)

// StubHandler is a `catch all` handler for all stub queries. Should be divided if
// more queries are to be supported
func StubHandler(data *JSONdata) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			segments := utils.GetSegments(r.URL, StubServicePath)

			switch len(segments) {
			case 1:
				switch segments[0] {
				case "all":
					httpRespondJSON(w, data)
				default:
					http.Error(w, "Unsupported URL segment, Usage: alpha/{ccn3}", http.StatusNotImplemented)
				}
			case 2:
				switch segments[0] {
				case "alpha":
					httpRespondJSON(w, data.filterByCCA3Code(segments[1]))
				case "name":
					httpRespondJSON(w, data.filterByName(segments[1]))
				default:
					http.Error(w, "segment: "+segments[0]+" not supported", http.StatusNotImplemented)
				}
			default:
				http.Error(w, "Unsupported URL segment, Usage: alpha/{cca3}", http.StatusNotImplemented)
			}
		default:
			http.Error(w, "Method not supported", http.StatusNotImplemented)
		}
	}
}

// httpRespondJSON takes any type of data and attempts to encode it as JSON to the response writer
func httpRespondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
	}
}
