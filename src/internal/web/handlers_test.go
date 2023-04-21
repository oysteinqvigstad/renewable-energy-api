package web

import (
	"assignment2/internal/types"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestEnergyDefaultHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(DefaultHandler))
	defer server.Close()
	statusCode := HttpGetStatusCode(t, server.URL)
	if statusCode != http.StatusInternalServerError {
		t.Fatalf("Wrong status code, expected: %d, got: %d", http.StatusInternalServerError, statusCode)
	}
}

// TestEnergyCurrentHandler tests the EnergyCurrentHandler function.
func TestEnergyCurrentHandler(t *testing.T) {
	wd, _ := os.Getwd()
	if err := os.Chdir(filepath.Join(wd, "..", "..")); err != nil {
		t.Fatal("Could not change working directory. CSV file will probably not be found")
	}
	s := NewService(path.Join("res", types.CSVFilePath), WithoutFirestore{})
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
		t.Fatal("Expected more than 1 country, got : ", len(dataList))
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

}

// Tests the invalid Method for EnergyCurrentHandler
func TestEnergyCurrentHandler_InvalidMethod(t *testing.T) {
	s := NewService(path.Join("res", types.CSVFilePath), WithoutFirestore{})
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
	s := NewService(path.Join("res", types.CSVFilePath), WithoutFirestore{})
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

// Tests the invalid Method for EnergyHistoryHandler
func TestEnergyHistoryHandler_InvalidMethod(t *testing.T) {
	s := NewService(path.Join("res", types.CSVFilePath), WithoutFirestore{})
	server := httptest.NewServer(http.HandlerFunc(s.EnergyHistoryHandler))
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
