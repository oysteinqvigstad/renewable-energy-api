package db

type YearRecord struct {
	Name       string  `json:"name"`
	ISO        string  `json:"isoCode"`
	Year       string  `json:"year,omitempty"`
	Percentage float64 `json:"percentage"`
}

type RenewableDB map[string][]YearRecord
