package api

import (
	"assignment2/internal/client"
)

const (
	API_BASE      = "http://129.241.150.113:8080/"
	API_VERSION   = "v3.1"
	ENDPOINT_NAME = "name"
	ENDPOINT_CCA  = "alpha"
)

type bordersResp struct {
	Borders []string `json:"borders"`
}

func GetNeighbours(name string) ([]string, error) {
	// Instantiate client
	cl := client.NewClient()
	err := cl.SetURL(API_BASE, API_VERSION, ENDPOINT_NAME, name)
	if err != nil {
		return nil, err
	}

	// Add queries
	cl.AddQuery("fullText", "true")
	cl.AddQuery("fields", "borders")

	// Perform get request
	resp := []bordersResp{}
	err = cl.GetAndDecode(&resp)
	if err != nil {
		return nil, err
	}

	return resp[0].Borders, nil
}

// GetBorders takes a cca3 code and returns
// an array of cca3 codes for bordering countries.
func GetNeighboursCca(cca string) ([]string, error) {
	// Instantiate client
	cl := client.NewClient()
	cl.SetURL(API_BASE, API_VERSION, ENDPOINT_CCA, cca)
	cl.AddQuery("fields", "borders")

	// Perform GET request
	resp := bordersResp{}
	err := cl.GetAndDecode(&resp)

	return resp.Borders, err
}
