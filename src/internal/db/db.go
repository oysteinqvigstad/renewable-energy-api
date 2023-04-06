package db

// getNeighbours returns a slice of 3-letter country codes for that specific country
func GetNeighbours(name string) []string {
	// TODO: implementation
	return []string{}
}

// getLatestEnergyData gets the newest data on record for a specific country
func GetLatestEnergyData(countryName string) RenewablesAPIData {
	// TODO: implementation
	return RenewablesAPIData{}
}

// getCurrentEnegeryData retrieves the latest energy data for a single country, and returns
// it as a list of structs. If the `includeNeighbours` flag has been set, then the energy data
// for the country's neighbour will be appended to the list
func GetCurrentEnergyData(countryName string, includeNeighbours bool) []RenewablesAPIData {
	data := []RenewablesAPIData{GetLatestEnergyData(countryName)}
	if includeNeighbours {
		for _, neighbour := range GetNeighbours(countryName) {
			data = append(data, GetLatestEnergyData(neighbour))
		}
	}
	return data
}
