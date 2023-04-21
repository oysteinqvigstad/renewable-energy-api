package types

import (
	"assignment2/api"
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	CSVFilePath = "renewable-share-energy.csv"
)

// YearRecord represents a record of renewable energy data for a specific country and year.
type YearRecord struct {
	Name       string  `json:"name"`
	ISO        string  `json:"isoCode"`
	Year       string  `json:"year,omitempty"`
	Percentage float64 `json:"percentage"`
}

// YearRecordList is a list of YearRecord instances.
type YearRecordList []YearRecord

// RenewableDB is a map representing renewable energy data organized by country code.
type RenewableDB map[string]YearRecordList

// ParseCSV will load a CSV file into RenewableDB
func ParseCSV(filepath string) RenewableDB {
	db := make(RenewableDB)
	// open file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open file"+filepath, err)
	}
	reader := csv.NewReader(file)

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
	for _, countryList := range db {
		countryList.sortByYear(true)
	}

	err = file.Close()
	if err != nil {
		log.Fatal("could not close file" + filepath + "maybe it has been closed already")
	}
	return db
}

// GetLatest retrieves the latest energy data for a single country, and returns
// it as a list of structs. If the `includeNeighbours` flag has been set, then the energy data
// for the country's neighbour will be appended to the list
func (db *RenewableDB) GetLatest(countryName string, includeNeighbours bool) YearRecordList {
	data := db.retrieveLatest(countryName)
	// adding neighbours if applicable
	if len(countryName) > 0 && includeNeighbours {
		neighbours, err := api.GetNeighboursCca(countryName)
		if err == nil {
			for _, neighbour := range neighbours {
				data = append(data, db.retrieveLatest(neighbour)...)
			}
		}
	}
	return data
}

// GetHistoricAvg will calculate the average renewable energy percentage for all stub_countries_api.
// If a year range is specified then it will only calculate the average for those years,
// and if `shouldSort` is enabled then it will sort by percentage in descending order
func (db *RenewableDB) GetHistoricAvg(start, end int, sortByPercentage bool) YearRecordList {
	var data YearRecordList
	for _, recordList := range *db {
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
	if sortByPercentage {
		data.sortByPercentage(true)
	} else {
		data.sortByName(true)
	}
	return data
}

// GetHistoric will get the historic records of renewable energy share for a specific country,
// and filter by years if `start` and/or `end` has been provided. If `shouldSort` is true,
// then the results will be sorted in descending order
func (db *RenewableDB) GetHistoric(countryCode string, start, end int, sortByPercentage bool) YearRecordList {
	var data YearRecordList
	countryCode = strings.ToUpper(countryCode)
	recordList, ok := (*db)[countryCode]
	if ok {
		for _, record := range recordList {
			if yearInRange(record, start, end) {
				data = append(data, record)
			}
		}
	}
	if sortByPercentage {
		data.sortByPercentage(true)
	}

	return data
}

// insert will append single record into the renewableDB
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

		if _, ok := (*db)[isoCode]; !ok {
			// allocating room for 60 historical data for each country. Will speed up append slightly
			(*db)[isoCode] = make(YearRecordList, 0, 60)
		}
		(*db)[isoCode] = append((*db)[isoCode], entry)
	}
}

// retrieveLatest gets the newest data on record for a specific country
// if an empty string is given then all stub_countries_api should be returned
func (db *RenewableDB) retrieveLatest(countryCode string) YearRecordList {
	var data YearRecordList

	// check if all stub_countries_api should be retrieved
	if len(countryCode) == 0 {
		for _, country := range *db {
			if len(country) > 0 {
				data = append(data, country[len(country)-1])
			}
		}
		// otherwise just fetch the single country
	} else {
		countryCode = strings.ToUpper(countryCode)
		country, ok := (*db)[countryCode]
		if ok {
			if len(country) > 0 {
				data = append(data, country[len(country)-1])
			}
		}
	}
	data.sortByName(true)
	return data
}

// GetName returns the country name for the given countryCode in the RenewableDB.
func (db *RenewableDB) GetName(countryCode string) string {
	return (*db)[countryCode][0].Name
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

// sortByYear takes a YearRecordList and sorts it by percentage
func (list YearRecordList) sortByPercentage(descending bool) {
	sort.Slice(list, func(i, j int) bool {
		if descending {
			return list[i].Percentage > list[j].Percentage
		} else {
			return list[i].Percentage < list[j].Percentage

		}
	})
}

// sortByYear takes a YearRecordList and sorts it by name
func (list YearRecordList) sortByName(ascending bool) {
	sort.Slice(list, func(i, j int) bool {
		if ascending {
			return list[i].Name < list[j].Name
		} else {
			return list[i].Name > list[j].Name
		}
	})
}

// sortByYear takes a YearRecordList and sorts it by year
func (list YearRecordList) sortByYear(ascending bool) {
	sort.Slice(list, func(i, j int) bool {
		first, err1 := strconv.Atoi(list[i].Year)
		second, err2 := strconv.Atoi(list[j].Year)
		if err1 != nil || err2 != nil {
			log.Fatal("Inconsistent data in YearRecordList. Could not parse year string to int")
		}
		if ascending {
			return first < second
		} else {
			return first > second
		}
	})
}

// MakeUniqueCCNACodes extracts unique country codes from a list of YearRecord instances.
func (list YearRecordList) MakeUniqueCCNACodes() []string {
	seen := map[string]bool{}
	var result []string
	for _, record := range list {
		if !seen[record.ISO] {
			seen[record.ISO] = true
			result = append(result, record.ISO)
		}
	}
	return result
}
