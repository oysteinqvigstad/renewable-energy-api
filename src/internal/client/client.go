package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RestClient struct {
	request *http.Request
}

func NewRestClient(url string) RestClient {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Error in request:", err.Error())
	}

	return RestClient{r}
}

// Sends request and returns status code
// Request should be easy for external to handle
func (client *RestClient) Prod() (string, error) {
	status := ""

	res, err := client.Get()
	if err == nil {
		status = res.Status
		res.Body.Close()
	}

	return status, err
}

func (client *RestClient) SetPath(value string) {
	client.request.URL = client.request.URL.JoinPath(value)
}

func (client *RestClient) AddQuery(key string, value string) {
	query := client.request.URL.Query()
	query.Add(key, value)
	client.request.URL.RawQuery = query.Encode()
}

func (client *RestClient) SetQuery(key string, value string) {
	query := client.request.URL.Query()
	query.Set(key, value)
	client.request.URL.RawQuery = query.Encode()
}

func (client *RestClient) ClearQuery() {
	query := client.request.URL.Query()
	for k := range query {
		delete(query, k)
	}
}

func (client *RestClient) Get() (*http.Response, error) {
	// instantiate client
	c := &http.Client{}
	defer c.CloseIdleConnections()

	fmt.Println("Performing query: ", client.request.URL.String())

	// issue request
	res, err := c.Do(client.request)
	return res, err
}

func (client *RestClient) GetContent(output any) error {
	// issue request
	res, err := client.Get()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// fail if not ok
	if res.StatusCode != http.StatusOK {
		output = nil
		return fmt.Errorf("restclient: expected status code 200 OK but got %s instead. output set to nil", res.Status)
	}

	// decode json data
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(output)

	return err
}
