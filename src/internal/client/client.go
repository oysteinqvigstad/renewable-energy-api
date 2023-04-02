package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	URL *url.URL
}

func NewClient(rawURL string) (*Client, error) {
	URL, e := url.Parse(rawURL)
	if e != nil {
		return nil, e
	}

	return &Client{URL}, nil
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

func (client *Client) SetPath(value ...string) {
	client.URL = client.URL.JoinPath(value...)
}

func (client *Client) AddQuery(key string, value string) {
	query := client.URL.Query()
	vals := query[key]

	//check if duplicate key/value pair
	duplicate := false
	for i := range vals {
		duplicate = duplicate || vals[i] == value
	}

	if !duplicate {
		query.Add(key, value)
		client.URL.RawQuery = query.Encode()
	}
}

func (client *Client) SetQuery(key string, value string) {
	query := client.URL.Query()
	query.Set(key, value)
	client.URL.RawQuery = query.Encode()
}

func (client *Client) ClearQuery() {
	query := client.URL.Query()
	for k := range query {
		delete(query, k)
	}
}

func (client *Client) Get() (*http.Response, error) {
	// make request
	r, e := http.NewRequest(http.MethodGet, client.URL.String(), nil)
	if e != nil {
		return nil, e
	}

	// instantiate client
	c := &http.Client{}
	defer c.CloseIdleConnections()

	fmt.Println("Performing query: ", r.URL.String())

	// issue request
	res, err := c.Do(r)
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
