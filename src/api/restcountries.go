package api

import (
	"assignment2/internal/client"
	"log"
	"net/url"
)

const (
	API_BASE      = "http://129.241.150.113:8080/"
	API_VERSION   = "v3.1"
	ENDPOINT_NAME = "name"
	ENDPOINT_CCA  = "alpha"
)

// Search for countries with name component.
// TODO: Populate and return response struct.
func SearchByName(name string) {
	// make client object
	cl := client.Client{URL: &url.URL{}}
	cl.SetURL(API_BASE, API_VERSION, ENDPOINT_NAME, name)

	// perform get request
	resp, e := cl.Get()
	if e != nil {
		log.Fatal(e.Error())
	}
	defer resp.Body.Close()

	// TEMP: print status until response struct is implemented
	log.Printf("Response status: %v\n", resp.Status)
}
