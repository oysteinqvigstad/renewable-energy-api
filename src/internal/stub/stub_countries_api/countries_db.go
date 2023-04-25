package stub_countries_api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type JSONdata []interface{}

// ParseJSON will read a JSON file and return a data structure for it
func ParseJSON(filepath string) JSONdata {
	// open file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("could not load json file", err)
	}

	// read and decode
	var data []interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Fatal("error during decoding of json", err)
	}

	// closing file
	err = file.Close()
	if err != nil {
		log.Fatal("could not close file" + filepath + "maybe it has been closed already")
	}

	return data
}

// filterByCCA3Code will return the country matching the 3-letter country code
func (data *JSONdata) filterByCCA3Code(countryCode string, fields []string) interface{} {
	//filteredJSON := JSONdata{}
	countryCode = strings.ToUpper(countryCode)
	for _, each := range *data {
		if country, ok := each.(map[string]interface{}); ok {
			if cca3, ok := country["cca3"].(string); ok {
				if cca3 == countryCode {
					if len(fields) > 0 {
						filteredJSON := make(map[string]interface{})
						fmt.Println(fields)
						for _, field := range fields {
							if record, ok := country[field].(interface{}); ok {
								filteredJSON[field] = record
							}
						}
						return filteredJSON
					} else {
						return append([]interface{}{}, country)
					}
				}
			}
		}
	}
	var filteredJSON []interface{}
	return filteredJSON
}

// filterByName will return all countries where the name contains `partialName`
func (data *JSONdata) filterByName(partialName string) JSONdata {
	filteredJSON := JSONdata{}
	for _, each := range *data {
		if country, ok := each.(map[string]interface{}); ok {
			if nameList, ok := country["name"].(map[string]interface{}); ok {
				if name, ok := nameList["common"].(string); ok {
					if strings.Contains(strings.ToUpper(name), strings.ToUpper(partialName)) {
						filteredJSON = append(filteredJSON, country)
					}
				}
			}
		}
	}
	return filteredJSON
}
