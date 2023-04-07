package db

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var GlobalRenewableDB RenewableDB

type RenewableDB struct {
	data map[string][]RenewablesAPIData
}

func (db *RenewableDB) ParseCSV(filepath string) {
	// open file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file"+filepath, err)
	}
	reader := csv.NewReader(file)

	// initiate data structure
	db.initiate()

	// discards header line of the file
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	// read each record until EOF and appends them to the global structure
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		db.insert(record)
	}

	// sort each struct for every country by year in case the CSV file is in incorrect order
	db.sortYearAsc()

	err = file.Close()
	if err != nil {
		log.Fatal("could not close file" + filepath + "maybe it has been closed already")
	}
}

func (db *RenewableDB) initiate() {
	if db.data == nil {
		db.data = make(map[string][]RenewablesAPIData)
	} else {
		log.Fatal("globalRenewableDB should not be initialized twice")
	}
}

// GetLatest retrieves the latest energy data for a single country, and returns
// it as a list of structs. If the `includeNeighbours` flag has been set, then the energy data
// for the country's neighbour will be appended to the list
func (db *RenewableDB) GetLatest(countryName string, includeNeighbours bool) []RenewablesAPIData {
	data := db.retrieveLatest(countryName)
	if len(countryName) > 0 && includeNeighbours {
		// TODO: Add support for neighbours
		//for _, neighbour := range GetNeighbours(countryName) {
		//	data = append(data, db.retrieveLatest(neighbour)...)
		//}
	}
	return data
}

//func (db *RenewableDB) GetHistoric(countryCode string, start, end int, sort bool) []RenewablesAPIData {
//	TODO: implementation
//countryCode = strings.ToUpper(countryCode)
//if len(countryCode) == 0 {
//	return db.GetHistoricAvg(start, end, sort)
//} else {
//	return db.GetHistoric(countryCode, start, end)
//}
//
// TODO: if {sortByValue} is set -> Sort all the
// TODO: if {country} IS set -> return list of structs for that country
// TODO: if {country} IS NOT set -> return all data? Will be very large return

//}

func (db *RenewableDB) GetHistoricAvg(start, end int, sort bool) []RenewablesAPIData {
	var data []RenewablesAPIData
	for _, recordList := range db.data {
		for _, record := range recordList {
			if yearInRange(record, start, end) {
				data = append(data, record)
			}
		}
	}
	if sort == true {
		// TODO: Sorting
	}
	return data

}

//sum := 0.0
//for _, record := range recordList {
//if yearInRange(record, start, end) {
//sum = sum + record.Percentage
//}
//}
//data = append(data, RenewablesAPIData{
//Name:       recordList[0].Name,
//ISO:        recordList[0].ISO,
//Percentage: sum / float64(len(recordList)),
//})

func (db *RenewableDB) GetHistoric(countryCode string, start, end int) []RenewablesAPIData {
	var data []RenewablesAPIData
	countryCode = strings.ToUpper(countryCode)
	recordList, ok := db.data[countryCode]
	if ok {
		for _, record := range recordList {
			if yearInRange(record, start, end) {
				data = append(data, record)
			}
		}
	}
	return data
}

// insert will append single struct into the map
func (db *RenewableDB) insert(record []string) {
	isoCode := strings.ToUpper(record[1])
	if len(isoCode) == 3 {
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

		// allocating room for 60 historical data for each country. Will speed up append slightly
		if _, ok := db.data[isoCode]; !ok {
			db.data[isoCode] = make([]RenewablesAPIData, 0, 60)
		}
		db.data[isoCode] = append(db.data[isoCode], entry)
	}
}

// sortYearAsc goes through every country and sorts the struct in ascending order by year
func (db *RenewableDB) sortYearAsc() {
	for _, val := range db.data {
		sort.Slice(val, func(i, j int) bool {
			first, err := strconv.Atoi(val[i].Year)
			second, err2 := strconv.Atoi(val[j].Year)
			if err != nil || err2 != nil {
				log.Fatal("Inconsistent data in globalRenewAbleDB. Could not parse string to int")
			}
			return first < second
		})
	}
}

// GetLatestEnergyData gets the newest data on record for a specific country
// if an empty string is given then all countries should be returned
func (db *RenewableDB) retrieveLatest(countryCode string) []RenewablesAPIData {
	var data []RenewablesAPIData

	// check if all countries should be retrieved
	if len(countryCode) == 0 {
		for _, country := range db.data {
			if len(country) > 0 {
				data = append(data, country[len(country)-1])
			}
		}
		// otherwise just fetch the single country
	} else {
		countryCode = strings.ToUpper(countryCode)
		country, ok := db.data[countryCode]
		if ok {
			if len(country) > 0 {
				data = append(data, country[len(country)-1])
			}
		}
	}
	return data
}

func yearInRange(data RenewablesAPIData, start, end int) bool {
	year, err := strconv.Atoi(data.Year)
	if err != nil {
		log.Fatal("Inconsistent data in RenewableDB, could not convert to int")
	}
	if start == 0 {
		// no start specified
		return end == 0 || year <= end
	} else if end == 0 {
		// no end specified
		return year >= start
	} else {
		// otherwise check range at both ends
		return start <= year && year <= end
	}
}
