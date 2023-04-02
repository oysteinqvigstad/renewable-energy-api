package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	URL *url.URL // url to perform requests on.
}

// Construct a new instance of Client.
func NewClient(rawURL string) (*Client, error) {
	URL, e := url.Parse(rawURL)
	if e != nil {
		return nil, e
	}

	return &Client{URL}, nil
}

// Sends request and returns status code.
// User should make sure to send a request that is
// easy for the server to handle.
func (client *Client) Prod() (string, error) {
	status := ""

	res, err := client.Get()
	if err == nil {
		status = res.Status
		res.Body.Close()
	}

	return status, err
}

// Set path component of URL
func (client *Client) SetPath(value ...string) {
	client.URL = client.URL.JoinPath(value...)
}

// Add a key/value pair to the query.
// Duplicate values are discarded.
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

// Set key/value pair for query, replacing any existing values.
func (client *Client) SetQuery(key string, value string) {
	query := client.URL.Query()
	query.Set(key, value)
	client.URL.RawQuery = query.Encode()
}

// Clear all queries from URL.
func (client *Client) ClearQuery() {
	query := client.URL.Query()
	for k := range query {
		delete(query, k)
	}
}

// Instantiates a request, http client, and performs request.
func (client *Client) Do(method string, reader io.Reader) (*http.Response, error) {
	// make request
	r, e := http.NewRequest(method, client.URL.String(), reader)
	if e != nil {
		return nil, e
	}

	// instantiate client
	c := &http.Client{}
	defer c.CloseIdleConnections()

	fmt.Printf("Performing %v request on URL: %v\n", method, r.URL.String())

	// issue request
	res, err := c.Do(r)
	return res, err
}

// Wrapper method for Do.
// Performs a GET request with no body.
func (client *Client) Get() (*http.Response, error) {
	return client.Do(http.MethodGet, nil)
}

// Wrapper method for Do.
// Performs a POST request with provided body.
func (client *Client) Post(body io.Reader) (*http.Response, error) {
	if body == nil {
		return nil, fmt.Errorf("client: body cannot be nil when performing POST request")
	}

	return client.Do(http.MethodPost, body)
}

// Perform a get request and decode response body to the provided struct.
func (client *Client) GetAndDecode(output any) error {
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
