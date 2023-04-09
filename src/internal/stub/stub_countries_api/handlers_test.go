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

}
