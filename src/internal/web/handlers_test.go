package web

import (
	"assignment2/internal/types"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"
)

const (
	testWithFirestore = false
)

func TestSetup(t *testing.T) {
	wd, _ := os.Getwd()
	if err := os.Chdir(filepath.Join(wd, "..", "..")); err != nil {
		t.Fatal("Could not change working directory. CSV file will probably not be found")
	}
}

func TestEnergyDefaultHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(DefaultHandler))
	defer server.Close()
	// Test 1: Send a GET request to the DefaultHandler
	statusCode := HttpGetStatusCode(t, server.URL)
	if statusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode)
	}
	//Test 2: Testing for use an invalid method: POST
	statusCode2 := HttpPostStatusCode(t, server.URL, "")
	if statusCode2 != http.StatusBadRequest {
		t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode2)
	}
}

// TestEnergyCurrentHandler tests the EnergyCurrentHandler function.

func TestEnergyCurrentHandler(t *testing.T) {
	runTests := func(t *testing.T, s *State) {
		server := httptest.NewServer(http.HandlerFunc(s.EnergyCurrentHandler))
		defer server.Close()

		//dataSingle := datastore.YearRecord{}
		dataList := types.YearRecordList{}

		// Test 1: Get all current country data
		HttpGetAndDecode(t, server.URL+RenewablesCurrentPath, &dataList)
		if len(dataList) != 79 {
			t.Fatal("expected 79 countries to be returned, got: ", len(dataList))
		}

		// Test 2: Get only norway
		HttpGetAndDecode(t, server.URL+RenewablesCurrentPath+"nor", &dataList)
		if dataList[0].Name != "Norway" {
			t.Fatal("Expected Norway to be returned")
		}
		if dataList[0].Year != "2021" {
			t.Fatal("Expected 2021 got: ", dataList[0].Year)
		}
		if dataList[0].Percentage != 71.558365 {
			t.Fatal("Expected 2021 got: ", dataList[0].Percentage)
		}

		// Test 3: Verify that the API returns more than one neighboring country for Norway when the 'neighbours' query parameter is set to 'true'.
		HttpGetAndDecode(t, server.URL+RenewablesCurrentPath+"nor?neighbours=true", &dataList)
		if len(dataList) <= 1 {
			t.Fatal("Expected more than 1 country, got :", len(dataList),
				"Have you remembered to start the stub service at /cmd/stub/stub_countries_api.go ?")
		}

		// Status codes tests:

		// Test 1: Testing whether a Bad Request error is returned when an invalid country code is provided in the URL.
		statusCode1 := HttpGetStatusCode(t, server.URL+RenewablesCurrentPath+"norr")
		if statusCode1 != http.StatusBadRequest {
			t.Fatal("Wrong status code, expected: 400 got: ", statusCode1)
		}

		//Test 2: Testing whether a Bad Request error is returned when 2 country codes is provided in the URL.
		statusCode2 := HttpGetStatusCode(t, server.URL+RenewablesHistoryPath+"nor/swe")
		if statusCode2 != http.StatusBadRequest {
			t.Fatal("Wrong status code, expected: 400, got: ", statusCode2)
		}
		//
		//Test 3: Testing for use an invalid method: POST
		statusCode3 := HttpPostStatusCode(t, server.URL, "")
		if statusCode3 != http.StatusBadRequest {
			t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode3)
		}
	}

	runTests(t, NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithoutFirestore{}))
	if testWithFirestore {
		runTests(t, NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithFirestore{}))
	}

}

// Tests the invalid Method for EnergyCurrentHandler

func TestEnergyCurrentHandler_InvalidMethod(t *testing.T) {
	s := NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithoutFirestore{})
	server := httptest.NewServer(http.HandlerFunc(s.EnergyHistoryHandler))
	defer server.Close()
	// Test: Send a POST request to the EnergyHistoryHandler
	req, err := http.NewRequest(http.MethodPost, server.URL+RenewablesCurrentPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	// Check if the returned status code is http.StatusBadRequest
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

// TestEnergyHistoryHandler tests the EnergyHistoryHandler function.

func TestEnergyHistoryHandler(t *testing.T) {
	runTests := func(t *testing.T, s *State) {
		server := httptest.NewServer(http.HandlerFunc(s.EnergyHistoryHandler))
		defer server.Close()

		//dataSingle := datastore.YearRecord{}
		dataList := types.YearRecordList{}
		//Test 1: Get the current countries renewable energy for a specific year: 2005
		HttpGetAndDecode(t, server.URL+RenewablesHistoryPath+"nor?begin=2005&end=2005", &dataList)
		if dataList[0].Percentage != 69.73603 {
			t.Fatal("Expected: 69.73603, got:", dataList[0].Percentage)
		}

		// Test 2: get countries renewable energy between a specific range of years: 2001 - 2008
		HttpGetAndDecode(t, server.URL+RenewablesHistoryPath+"nor?begin=2001&end=2008", &dataList)
		if len(dataList) != 8 {
			t.Fatal("Expected , got: ", len(dataList))
		}

		// Test 3: find the average renewable energy between a specific range of years: 2001 - 2008
		if calculateAverage(dataList) != 67.56928574999999 {
			t.Fatal("Expected: 67.56928574999999, got:", calculateAverage(dataList))
		}

		// Test 4: if you don't specify an alpha code
		HttpGetAndDecode(t, server.URL+RenewablesHistoryPath, &dataList)
		if len(dataList) != 79 {
			t.Fatal("Expected : 79, got: ", len(dataList))
		}

		// Test 5: Testing if the sort by value is working
		HttpGetAndDecode(t, server.URL+RenewablesHistoryPath+"nor?begin=2004&end=2009&sortByValue=true", &dataList)
		if dataList[0].Percentage < dataList[5].Percentage {
			t.Fatal("Wrong sort order: expected the first value to be greater than the last value")
		}

		// Status codes tests:
		// Test 1: Testing whether a Bad Request error is returned when an invalid country code is provided in the URL.
		statusCode1 := HttpGetStatusCode(t, server.URL+RenewablesHistoryPath+"norr")
		if statusCode1 != http.StatusBadRequest {
			t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode1)
		}

		// Test 2: status code for invalid year range
		statusCode2 := HttpGetStatusCode(t, server.URL+RenewablesHistoryPath+"nor?begin=2005&end=2001")
		if statusCode2 != http.StatusBadRequest {
			t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode2)
		}

		//Test 3: Testing whether a Bad Request error is returned when 2 country codes is provided in the URL.
		statusCode3 := HttpGetStatusCode(t, server.URL+RenewablesHistoryPath+"nor/swe")
		if statusCode3 != http.StatusBadRequest {
			t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode3)
		}
	}

	runTests(t, NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithoutFirestore{}))
	if testWithFirestore {
		runTests(t, NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithFirestore{}))
	}
}

// Tests the invalid Method for EnergyHistoryHandler

func TestEnergyHistoryHandler_InvalidMethod(t *testing.T) {
	s := NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithoutFirestore{})
	//server := httptest.NewServer(http.HandlerFunc(s.EnergyHistoryHandler))
	server := httptest.NewServer(SetupRoutes("8081", s))
	defer server.Close()
	// Test: Send a POST request to the EnergyHistoryHandler
	req, err := http.NewRequest(http.MethodPost, server.URL+RenewablesHistoryPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	// Check if the returned status code is http.StatusBadRequest
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func calculateAverage(dataList types.YearRecordList) float64 {
	sum := 0.0
	for _, data := range dataList {
		sum += data.Percentage
	}
	avg := sum / float64(len(dataList))
	return avg
}
func TestNotificationHandler(t *testing.T) {
	runTests := func(t *testing.T, s *State) {
		server := httptest.NewServer(http.HandlerFunc(s.NotificationHandler))
		defer server.Close()

		// synchronize webhook response
		wg := sync.WaitGroup{}
		wg.Add(1)

		singleResponse := WebhookResponse{}
		multiResponse := make([]WebhookResponse, 0, 5)

		// setting up receiving webhook server
		receiverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// test 10: receiving webhook method is Post
			if r.Method != http.MethodPost {
				t.Fatal("expected method to be post")
			}

			// test 11: country name is correct
			data := WebhookResponse{}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				t.Fatal("Could not decode json")
			}
			if data.Country != "Germany" {
				t.Fatal("Expected country name to be Germany")
			}

			// POST received and tested, allowing test to continue
			wg.Done()
		}))
		defer receiverServer.Close()

		// test 1: Initiate two valid webhook registrations
		jsonBodyCorrect := "{ \"url\": \"" + receiverServer.URL + "\", \"country\": \"DEU\", \"calls\": 5 }"
		if HttpPostAndDecode(t, server.URL+NotificationsPath, jsonBodyCorrect, &singleResponse) != http.StatusCreated {
			t.Fatal("Expected 201 created")
		}
		webhookID1 := singleResponse.WebhookID
		if HttpPostAndDecode(t, server.URL+NotificationsPath, jsonBodyCorrect, &singleResponse) != http.StatusCreated {
			t.Fatal("Expected 201 created")
		}
		webhookID2 := singleResponse.WebhookID

		// test 2: Invalid webhook registration (no http:// or https:// prefix)
		jsonBodyInvalid := "{ \"url\": \"webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c\", \"country\": \"DEU\", \"calls\": 5 }"
		if HttpPostStatusCode(t, server.URL+NotificationsPath, jsonBodyInvalid) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 3: Invalid webhook registration (invalid country code)
		jsonBodyInvalid = "{ \"url\": \"http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c\", \"country\": \"SVE\", \"calls\": 5 }"
		if HttpPostStatusCode(t, server.URL+NotificationsPath, jsonBodyInvalid) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 4: Invalid webhook registration (invalid calls digit)
		jsonBodyInvalid = "{ \"url\": \"http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c\", \"country\": \"DEU\", \"calls\": -1 }"
		if HttpPostStatusCode(t, server.URL+NotificationsPath, jsonBodyInvalid) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 2: Invalid webhook registration (no http:// or https:// prefix)
		jsonBodyInvalid = "{ \"url\": \"webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c\", \"country\": \"DEU\", \"calls\": 5 }"
		if HttpPostStatusCode(t, server.URL+NotificationsPath, jsonBodyInvalid) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 3: Invalid webhook registration (invalid country code)
		jsonBodyInvalid = "{ \"url\": \"http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c\", \"country\": \"SVE\", \"calls\": 5 }"
		if HttpPostStatusCode(t, server.URL+NotificationsPath, jsonBodyInvalid) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 4: Invalid webhook registration (invalid calls digit)
		jsonBodyInvalid = "{ \"url\": \"http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c\", \"country\": \"DEU\", \"calls\": -1 }"
		if HttpPostStatusCode(t, server.URL+NotificationsPath, jsonBodyInvalid) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 5: Show all webhook registrations
		switch s.firestoreMode.(type) {
		case WithFirestore:
			HttpGetAndDecode(t, server.URL+NotificationsPath, &multiResponse)
			if len(multiResponse) < 2 {
				t.Fatal("expected at least two registrations")
			}
		case WithoutFirestore:
			HttpGetAndDecode(t, server.URL+NotificationsPath, &multiResponse)
			if len(multiResponse) != 2 {
				t.Fatal("expected exactly two registrations")
			}

		}

		// test 6: Show a specific webhook registration based on ID
		HttpGetAndDecode(t, server.URL+NotificationsPath+webhookID1, &singleResponse)
		if singleResponse.WebhookID != webhookID1 {
			t.Fatal("expected service to return webhook data")
		}

		// test 7: Delete a webhook
		if HttpDeleteStatusCode(t, server.URL+NotificationsPath+webhookID1, "") != http.StatusAccepted {
			t.Fatal("Expected webhook to be deleted")
		}

		// test 8: Attempting to delete an invalid webhook ID
		if HttpDeleteStatusCode(t, server.URL+NotificationsPath+"my_webhook_id", "") != http.StatusBadRequest {
			t.Fatal("Expected server to respond with Bad Request")
		}

		// test 9: Trigger a webhook
		renewablesServer := httptest.NewServer(http.HandlerFunc(s.EnergyCurrentHandler))
		defer renewablesServer.Close()
		for i := 0; i < 5; i++ {
			if HttpGetStatusCode(t, renewablesServer.URL+RenewablesCurrentPath+"deu") != http.StatusOK {
				t.Fatal("Expected 200 OK")
			}
		}

		// Waiting for response to be processed by the receiving mock server
		wg.Wait()

		// test 12: Invalid request with second segment after webhook_ID
		if HttpGetStatusCode(t, server.URL+NotificationsPath+"UHVvwFuc4xLTt/extra_segment") != http.StatusBadRequest {
			t.Fatal("Expected status Bad Request")
		}

		// test 13: Invalid request with second segment after webhook_ID
		if HttpGetStatusCode(t, server.URL+NotificationsPath+"UHVvwFuc4xLTt/extra_segment") != http.StatusBadRequest {
			t.Fatal("Expected status Bad Request")
		}

		// test 14: Invalid webhook registration with extra URL segment
		if HttpPostStatusCode(t, server.URL+NotificationsPath+"extra_segment", jsonBodyCorrect) != http.StatusBadRequest {
			t.Fatal("Expected 400 Bad Request")
		}

		// test 15: Attempting to delete webhook without webhook ID
		if HttpDeleteStatusCode(t, server.URL+NotificationsPath, "") != http.StatusBadRequest {
			t.Fatal("Expected server to respond with Bad Request")
		}

		//  clean up the second webhook
		if HttpDeleteStatusCode(t, server.URL+NotificationsPath+webhookID2, "") != http.StatusAccepted {
			t.Fatal("Expected webhook to be deleted")
		}

	}

	runTests(t, NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithoutFirestore{}))
	if testWithFirestore {
		runTests(t, NewService(path.Join("res", types.CSVFilePath), StubRestCountries{}, WithFirestore{}))
	}
}

// TestStatusHandler verifies the behavior of the StatusHandler function by testing
// various scenarios, such as sending requests with different HTTP methods and
// checking the expected values in the APIStatus struct.
func TestStatusHandler(t *testing.T) {
	runTests := func(t *testing.T, s *State) {
		server := httptest.NewServer(http.HandlerFunc(s.StatusHandler))
		defer server.Close()

		var apiStatus APIStatus

		// Test 1: Check if the APIStatus struct fields have expected values
		HttpGetAndDecode(t, server.URL+StatusPath, &apiStatus)
		if apiStatus.Countriesapi != http.StatusOK {
			t.Errorf("Have you started the stub service for REST countries?"+
				"Unexpected countries API status: got %v want %v", apiStatus.Countriesapi, http.StatusOK)
		}
		// Test 2: Testing whether a Bad Request error is returned when an unsupported HTTP method is used
		statusCode1 := HttpPostStatusCode(t, server.URL+StatusPath, "")
		if statusCode1 != http.StatusBadRequest {
			t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, statusCode1)
		}

		// Test 3: Send a GET request to the StatusHandler with additional URL segments
		resp, err := http.Get(server.URL + StatusPath + "/extra/segments")
		if err != nil {
			t.Fatalf("Error making GET request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusBadRequest, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}
		expected := "Usage: energy/v1/status/\n"
		if string(body) != expected {
			t.Fatalf("Unexpected response body, expected: [%s], got: [%s]", expected, string(body))
		}

	}

	filePath := path.Join("res", types.CSVFilePath)
	state1 := NewService(filePath, StubRestCountries{}, WithoutFirestore{})
	runTests(t, state1)

	if testWithFirestore {
		state2 := NewService(filePath, StubRestCountries{}, WithFirestore{})
		runTests(t, state2)
	}
}
