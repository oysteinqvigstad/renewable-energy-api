package db

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
