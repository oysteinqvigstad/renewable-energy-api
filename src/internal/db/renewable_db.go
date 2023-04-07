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

// ParseCSV will load a CSV file into GlobalRenewableDB
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

// initiate will make the GlobalRenewableDB structure
func (db *RenewableDB) initiate() {
	if db.data == nil {
		db.data = make(map[string][]YearRecord)
	} else {
		log.Fatal("globalRenewableDB should not be initialized twice")
	}
}

// GetLatest retrieves the latest energy data for a single country, and returns
// it as a list of structs. If the `includeNeighbours` flag has been set, then the energy data
// for the country's neighbour will be appended to the list
func (db *RenewableDB) GetLatest(countryName string, includeNeighbours bool) []YearRecord {
	data := db.retrieveLatest(countryName)
	if len(countryName) > 0 && includeNeighbours {
		// TODO: Add support for neighbours
		//for _, neighbour := range GetNeighbours(countryName) {
		//	data = append(data, db.retrieveLatest(neighbour)...)
		//}
	}
	return data
}

// GetHistoricAvg will calculate the average renewable energy percentage for all countries.
// If a year range is specified then it will only calculate the average for those years,
// and if `shouldSort` is enabled then it will sort by percentage in descending order
func (db *RenewableDB) GetHistoricAvg(start, end int, shouldSort bool) []YearRecord {
	var data []YearRecord
	for _, recordList := range db.data {
		sum := 0.0
		numOfYears := 0
		for _, record := range recordList {
			if yearInRange(record, start, end) {
				sum = sum + record.Percentage
				numOfYears = numOfYears + 1
			}
		}
		if numOfYears > 0 {
			data = append(data, YearRecord{
				Name:       recordList[0].Name,
				ISO:        recordList[0].ISO,
				Percentage: sum / float64(numOfYears),
			})
		}

	}
	if shouldSort {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Percentage > data[j].Percentage
		})
	}
	return data
}

// GetHistoric will get the historic records of renewable energy share for a specific country,
// and filter by years if `start` and/or `end` has been provided. If `shouldSort` is true,
// then the results will be sorted in descending order
func (db *RenewableDB) GetHistoric(countryCode string, start, end int, shouldSort bool) []YearRecord {
	var data []YearRecord
	countryCode = strings.ToUpper(countryCode)
	recordList, ok := db.data[countryCode]
	if ok {
		for _, record := range recordList {
			if yearInRange(record, start, end) {
				data = append(data, record)
			}
		}
	}
	if shouldSort {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Percentage > data[j].Percentage
		})
	}

	return data
}

// insert will append single record into GlobalRenewableDB
func (db *RenewableDB) insert(record []string) {
	isoCode := strings.ToUpper(record[1])
	if len(isoCode) == 3 {
		// converts percentage to float64
		percentage, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}
		entry := YearRecord{
			Name:       record[0],
			ISO:        record[1],
			Year:       record[2],
			Percentage: percentage,
		}

		if _, ok := db.data[isoCode]; !ok {
			// allocating room for 60 historical data for each country. Will speed up append slightly
			db.data[isoCode] = make([]YearRecord, 0, 60)
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
func (db *RenewableDB) retrieveLatest(countryCode string) []YearRecord {
	var data []YearRecord

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

// yearInRange will return true/false depending on a record is within the range between `start` and `end`.
// Should any of them be 0, then no limit is set
func yearInRange(data YearRecord, start, end int) bool {
	year, err := strconv.Atoi(data.Year)
	if err != nil {
		log.Fatal("Inconsistent data in RenewableDB, could not convert to int")
	}
	noStartSpecified := start == 0
	noEndSpecified := end == 0

	return (noStartSpecified && (noEndSpecified || year <= end)) || // no start
		(noEndSpecified && year >= start) || // no end
		(start <= year && year <= end) // valid range

}
