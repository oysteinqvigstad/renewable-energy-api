package api

import (
	"assignment2/internal/web_client"
)

const (
	API_BASE      = "http://129.241.150.113:8080/"
	API_VERSION   = "v3.1"
	ENDPOINT_NAME = "name"
	ENDPOINT_CCA  = "alpha"
)

type country struct {
	Borders []string          `json:"borders"`
	Name    map[string]string `json:"name"`
}

// GetNeighbours takes a name string
// and returns an array of CCA3 code strings
func GetNeighbours(name string) ([]string, error) {
	// Instantiate client
	cl := web_client.NewClient()
	err := cl.SetURL(API_BASE, API_VERSION, ENDPOINT_NAME, name)
	if err != nil {
		return nil, err
	}

	// Add queries
	cl.AddQuery("fullText", "true")
	cl.AddQuery("fields", "borders")

	// Perform get request
	resp := []country{}
	err = cl.GetAndDecode(&resp)
	if err != nil {
		return nil, err
	}

	return resp[0].Borders, nil
}

// GetNeighboursCca takes a cca3 code and returns
// an array of cca3 codes for bordering countries.
func GetNeighboursCca(cca string) ([]string, error) {
	// Instantiate client
	cl := web_client.NewClient()
	err := cl.SetURL(API_BASE, API_VERSION, ENDPOINT_CCA, cca)
	if err != nil {
		return nil, err
	}

	// Add query
	cl.AddQuery("fields", "borders")

	// Perform GET request
	resp := country{}
	err = cl.GetAndDecode(&resp)

	if err != nil {
		return nil, err
	}

	return resp.Borders, nil
}
