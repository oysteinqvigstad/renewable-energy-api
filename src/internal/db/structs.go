package db

type RenewablesAPIData struct {
	Name       string `json:"name"`
	ISO        string `json:"isoCode"`
	Year       string `json:"year,omitempty"`
	Percentage int    `json:"percentage"`
}
