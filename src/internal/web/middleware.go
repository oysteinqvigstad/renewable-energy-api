package web

import (
	"encoding/json"
	"net/http"
	"testing"
)

// httpRespondJSON takes any type of data and attempts to encode it as JSON to the response writer
func httpRespondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
	}
}

// HttpGetAndDecode is a helper function that retrieves and returns the JSON data
// from a specific url
func HttpGetAndDecode(t *testing.T, url string, data any) {
	client := http.Client{}
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:", err.Error())
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Fatal("Error during decoding", err.Error())
	}
}
