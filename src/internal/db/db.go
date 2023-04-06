package db

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

var globalRenewableDB []RenewablesAPIData

// GetNeighbours returns a slice of 3-letter country codes for that specific country
func GetNeighbours(name string) []string {
	// TODO: implementation
	return []string{}
}

// GetLatestEnergyData gets the newest data on record for a specific country
// if an empty string is given then all countries should be returned
func GetLatestEnergyData(countryName string) RenewablesAPIData {
	// TODO: implementation
	return RenewablesAPIData{}
}

// GetCurrentEnergyData retrieves the latest energy data for a single country, and returns
// it as a list of structs. If the `includeNeighbours` flag has been set, then the energy data
// for the country's neighbour will be appended to the list
func GetCurrentEnergyData(countryName string, includeNeighbours bool) []RenewablesAPIData {
	data := []RenewablesAPIData{GetLatestEnergyData(countryName)}
	if len(countryName) > 0 && includeNeighbours {
		for _, neighbour := range GetNeighbours(countryName) {
			data = append(data, GetLatestEnergyData(neighbour))
		}
	}
	return data
}

func GetHistoricEnergyData(countryCode string, start, end int, sort bool) []RenewablesAPIData {
	// TODO: implementation
	// TODO: if {sortByValue} is set -> Sort all the
	// TODO: if {begin} IS set -> omit year attribute in country struct (returns single average)
	// TODO: if {country} IS set -> return list of structs for that country
	// TODO: if {country} IS NOT set -> return all data? Will be very large return
	return []RenewablesAPIData{}

}

// ParseRenewableCSV reads from the Renewable Share Enegery CSV from the
// specified filepath into the global data structure "globalRenewableDB"
func ParseRenewableCSV(filepath string) {
	// opens file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file"+filepath, err)
	}
	reader := csv.NewReader(file)
	defer file.Close()

	if len(globalRenewableDB) != 0 {
		log.Fatal("The CSV has already been parsed. Should not be parsed again")
	}

	// the CSV file is ~5600 records long, allocating room for 6000 records
	globalRenewableDB = make([]RenewablesAPIData, 0, 6000)

	// discards header line
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	// reads each record until EOF and appends them to the global structure
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// converts percentage to float64
		percentage, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		entry := RenewablesAPIData{
			Name:       record[0],
			ISO:        record[1],
			Year:       record[2],
			Percentage: percentage,
		}

		globalRenewableDB = append(globalRenewableDB, entry)
	}
}
