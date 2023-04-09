package client

import (
	"net/url"
	"testing"
)

func TestSetURL(t *testing.T) {
	c := Client{URL: &url.URL{}}
	c.SetURL("https://example.com/", "test")

	result := c.URL.String()
	want := "https://example.com/test"

	if result != want {
		t.Errorf("c.URL.String() = %v, want %v", result, want)
	}
}

func TestJoinPath(t *testing.T) {
	c := NewClient()
	c.SetURL("https://example.com/test")
	c.JoinPath("join", "path")

	want := "https://example.com/test/join/path"
	result := c.URL.String()

	if result != want {
		t.Errorf("c.URL.String() = %v, want %v", result, want)
	}
}

func TestAddQuery(t *testing.T) {
	c := NewClient()
	c.SetURL("https://example.com/test")
	c.AddQuery("vals", "name", "age", "name")

	result := c.URL.String()
	want := "https://example.com/test?vals=name&vals=age"

	if result != want {
		t.Errorf("c.URL.String() = %v, want %v", result, want)
	}
}

func TestSetQuery(t *testing.T) {
	c := NewClient()
	c.SetURL("https://example.com/test?vals=name&vals=age")
	c.SetQuery("vals", "dob")

	result := c.URL.String()
	want := "https://example.com/test?vals=dob"

	if result != want {
		t.Errorf("c.URL.String() = %v, want %v", result, want)
	}
}

func TestClearQuery(t *testing.T) {
	c := NewClient()
	c.SetURL("https://example.com/test?vals=name&vals=age&occupation=student")
	c.ClearQuery()

	result := c.URL.String()
	want := "https://example.com/test"

	if result != want {
		t.Errorf("c.URL.String() = %v, want %v", result, want)
	}
}
