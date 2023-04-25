package web

import (
	"assignment2/internal/types"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// httpRespondJSON takes any type of data and attempts to encode it as JSON to the response writer
func httpRespondJSON(w http.ResponseWriter, data any, s *State) {
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
	}
	go invocate(data, s)
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

// invocate processes webhooks for the given data and the provided application state.
func invocate(data any, s *State) {
	var invocationList []string
	switch data.(type) {
	case types.YearRecordList:
		invocationList = data.(types.YearRecordList).MakeUniqueCCNACodes()
	case types.YearRecord:
		invocationList = types.YearRecordList{data.(types.YearRecord)}.MakeUniqueCCNACodes()
	}

	// Process webhooks if there are any country codes in the invocation list
	if len(invocationList) > 0 {
		println("Invocated: " + strings.Join(invocationList, ","))
		ProcessWebhookByCountry(invocationList, s)
	}
}

// updateFirestore sends the given data to the specified channel.
func updateFirestore[T any](channel chan T, data T) {
	// checking if channel is open
	if channel != nil {
		select {
		case channel <- data:
			// successful
		default:
			println("channel is full! dropping data")
		}
	}
}

// httpCacheAndRespondJSON updates the cache, and sends a JSON response with the provided data and application state.
func httpCacheAndRespondJSON(w http.ResponseWriter, url *url.URL, data types.YearRecordList, s *State) {
	value := make(map[string]types.YearRecordList)
	value[url.String()] = data
	updateFirestore(s.chCache, value)
	httpRespondJSON(w, data, s)
}

// HttpGetStatusCode returns the statuscode of a POST request
func HttpPostStatusCode(t *testing.T, url string, payload string) int {
	client := http.Client{}
	defer client.CloseIdleConnections()
	res, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatal("Failed to send POST", err)
	}
	return res.StatusCode
}

// HttpPostAndDecode returns the response of a POST request
func HttpPostAndDecode(t *testing.T, url string, payload string, response any) int {
	client := http.Client{}
	defer client.CloseIdleConnections()
	res, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatal("Failed to send POST", err)
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	return res.StatusCode
}

// HttpDeleteStatusCode returns the statuscode of a DELETE request
func HttpDeleteStatusCode(t *testing.T, url string, payload string) int {
	client := http.Client{}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(payload))
	if err != nil {
		log.Fatalf("Error creating DELETE request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending DELETE request: %v", err)
	}
	return res.StatusCode
}
