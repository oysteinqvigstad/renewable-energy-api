package stub_countries_api

import (
	"assignment2/internal/web"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"testing"
)

func TestEnergyCurrentHandler(t *testing.T) {
	jsonData := ParseJSON(path.Join("..", "..", "..", "res", JSONFileName))
	server := httptest.NewServer(http.HandlerFunc(StubHandler(&jsonData)))
	defer server.Close()

	var data []interface{}
	var statusCode int

	// Test 1: checking that "all/" returns 250 countries
	web.HttpGetAndDecode(t, server.URL+StubServicePath+"all/", &data)
	if len(data) != 250 {
		t.Fatal("Unexpected number of records:", strconv.Itoa(len(data)))
	}

	// Test 2: checking that "alpha/nor" returns only one country
	web.HttpGetAndDecode(t, server.URL+StubServicePath+"alpha/nor", &data)
	if len(data) != 1 {
		t.Fatal("Unexpected number of records:", strconv.Itoa(len(data)))
	}

	// Test 2: checking that "alpha/nor" returns only one country
	web.HttpGetAndDecode(t, server.URL+StubServicePath+"alpha/nor", &data)
	if len(data) != 1 {
		t.Fatal("Unexpected number of records:", strconv.Itoa(len(data)))
	}

	// Test 3: checking that "alpha/nor" returns only one country
	web.HttpGetAndDecode(t, server.URL+StubServicePath+"name/norway", &data)
	if len(data) != 1 {
		t.Fatal("Unexpected number of records:", strconv.Itoa(len(data)))
	}

	// Test 4: checking that not giving a cca3 code returns status bad request
	statusCode = web.HttpGetStatusCode(t, server.URL+StubServicePath+"alpha/")
	if statusCode != http.StatusBadRequest {
		t.Fatal("Expected HTTP status bad request")
	}

	// Test 5: checking that giving no segments returns bad request
	statusCode = web.HttpGetStatusCode(t, server.URL+StubServicePath)
	if statusCode != http.StatusBadRequest {
		t.Fatal("Expected HTTP status bad request")
	}

	// Test 6: checking for a given segment that is not implemented
	statusCode = web.HttpGetStatusCode(t, server.URL+StubServicePath+"region/europe")
	if statusCode != http.StatusNotImplemented {
		t.Fatal("Expected HTTP status not implemented")
	}
}
