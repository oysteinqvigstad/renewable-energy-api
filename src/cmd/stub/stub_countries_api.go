package main

import (
	"assignment2/internal/stub/stub_countries_api"
	"log"
	"net/http"
	"path"
)

func main() {
	port := stub_countries_api.StubServicePort

	domainNamePort := "http://localhost:" + port

	CountriesData := stub_countries_api.ParseJSON(path.Join("res", stub_countries_api.JSONFileName))
	handler := http.HandlerFunc(stub_countries_api.StubHandler(&CountriesData))

	log.Println("Started services on:")
	log.Println(domainNamePort + stub_countries_api.StubServicePath)
	log.Println("Supported queries:")
	log.Println(domainNamePort + stub_countries_api.StubServicePath + "all/")
	log.Println(domainNamePort + stub_countries_api.StubServicePath + "alpha/{cca3}?fields=field1,field2")
	log.Println(domainNamePort + stub_countries_api.StubServicePath + "name/{partial_name}")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
