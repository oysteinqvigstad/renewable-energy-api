package app

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// GetSegments returns a slice of strings from a url
func GetSegments(url *url.URL, prefix string) []string {
	var segments []string
	split := strings.Split(strings.TrimPrefix(url.Path, prefix), "/")
	for _, str := range split {
		if len(str) > 0 {
			segments = append(segments, str)
		}
	}
	return segments
}

// GetQueryStr returns the value of a query key as string
func GetQueryStr(url *url.URL, name string) (string, error) {
	for key, value := range url.Query() {
		if name == key {
			return strings.Join(value, ","), nil
		}
	}
	return "", errors.New("could not find key")
}

// GetQueryInt returns the value of a query key as an integer
func GetQueryInt(url *url.URL, name string) (int, error) {
	value, err := GetQueryStr(url, name)
	if err != nil {
		return 0, err
	} else {
		number, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		} else {
			return number, nil
		}
	}
}
