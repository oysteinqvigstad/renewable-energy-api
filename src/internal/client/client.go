package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct {
	request *http.Request
}

func NewClient(url string) Client {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Error in request:", err.Error())
	}

	return Client{r}
}

// Sends request and returns status code
// Request should be easy for external to handle
func (client *Client) Prod() (string, error) {
	status := ""

	res, err := client.Get()
	if err == nil {
		status = res.Status
		res.Body.Close()
	}

	return status, err
}

func (client *Client) SetPath(value string) {
	client.request.URL = client.request.URL.JoinPath(value)
}

func (client *Client) AddQuery(key string, value string) {
	query := client.request.URL.Query()
	vals := query[key]

	//check if duplicate key/value pair
	duplicate := false
	for i := range vals {
		duplicate = duplicate || vals[i] == value
	}

	if !duplicate {
		query.Add(key, value)
		client.request.URL.RawQuery = query.Encode()
	}
}

func (client *Client) SetQuery(key string, value string) {
	query := client.request.URL.Query()
	query.Set(key, value)
	client.request.URL.RawQuery = query.Encode()
}

func (client *Client) ClearQuery() {
	query := client.request.URL.Query()
	for k := range query {
		delete(query, k)
	}
}

func (client *Client) Get() (*http.Response, error) {
	// instantiate client
	c := &http.Client{}
	defer c.CloseIdleConnections()

	fmt.Println("Performing query: ", client.request.URL.String())

	// issue request
	res, err := c.Do(client.request)
	return res, err
}

func (client *Client) GetContent(output any) error {
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
