package web

import (
	"assignment2/internal/db"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"testing"
)

func TestEnergyCurrentHandler(t *testing.T) {
	energyData := db.ParseCSV(path.Join("..", "..", "res", db.CSVFilePath))
	server := httptest.NewServer(http.HandlerFunc(EnergyCurrentHandler(energyData)))
	defer server.Close()

	client := http.Client{}
	res, err := client.Get(server.URL + RenewablesCurrentPath)
	if err != nil {
		t.Fatal("Get request to URL failed:", err.Error())
	}

	data := db.YearRecordList{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Fatal("Error during decoding", err.Error())
	}

	if len(data) < 50 || len(data) > 100 {
		t.Fatal("Unexpected number of records:", strconv.Itoa(len(data)))
	}

}
