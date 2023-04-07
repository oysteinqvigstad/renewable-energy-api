package db

type YearRecord struct {
	Name       string  `json:"name"`
	ISO        string  `json:"isoCode"`
	Year       string  `json:"year,omitempty"`
	Percentage float64 `json:"percentage"`
}

// RenewableDB is a collection of YearRecords
type RenewableDB struct {
	data map[string][]YearRecord
}
