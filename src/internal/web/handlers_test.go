package web

import (
	"assignment2/internal/db"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func TestEnergyCurrentHandler(t *testing.T) {
	energyData := db.ParseCSV(path.Join("..", "..", "res", db.CSVFilePath))
	server := httptest.NewServer(http.HandlerFunc(EnergyCurrentHandler(energyData)))
	defer server.Close()

	dataSingle := db.YearRecord{}
	dataList := db.YearRecordList{}

	// Test 1: Get all current country data
	HttpGetAndDecode(t, server.URL+RenewablesCurrentPath, &dataList)
	if len(dataList) != 79 {
		t.Fatal("expected 79 countries to be returned")
	}

	// Test 2: Get only norway
	HttpGetAndDecode(t, server.URL+RenewablesCurrentPath+"nor", &dataSingle)
	if dataSingle.Name != "Norway" {
		t.Fatal("Expected Norway to be returned")
	}
}
