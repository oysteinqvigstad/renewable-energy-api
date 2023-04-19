package web

import (
	"assignment2/internal/datastore"
	"encoding/json"
	"net/http"
	"strings"
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
	go invocate(data)
}

// HttpGetAndDecode is a helper function that retrieves and returns the JSON data
// from a specific url
func HttpGetAndDecode(t *testing.T, url string, data any) {
	client := http.Client{}
	defer client.CloseIdleConnections()
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:", err.Error())
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Fatal("Error during decoding", err.Error())
	}
}

// HttpGetStatusCode returns the statuscode of a GET request
func HttpGetStatusCode(t *testing.T, url string) int {
	client := http.Client{}
	defer client.CloseIdleConnections()
	res, err := client.Get(url)
	if err != nil {
		t.Fatal("Get request to URL failed:", err.Error())
	}
	return res.StatusCode
}

func invocate(data any) {
	var invocationList []string
	switch data.(type) {
	case datastore.YearRecordList:
		invocationList = data.(datastore.YearRecordList).Invocate()
	case datastore.YearRecord:
		invocationList = datastore.YearRecordList{data.(datastore.YearRecord)}.Invocate()
	}

	if len(invocationList) > 0 {
		println("Invocated: " + strings.Join(invocationList, ","))
		Invocate(invocationList)
	}
}
