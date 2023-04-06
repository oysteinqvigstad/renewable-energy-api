package utils

import (
	"net/url"
	"testing"
)

func TestGetSegments(t *testing.T) {
	path := url.URL{Path: "/example/path/100"}
	segments := GetSegments(&path, "/example")
	if len(segments) != 2 {
		t.Fatal("Wrong number of segments")
	}
	if segments[0] != "path" || segments[1] != "100" {
		t.Fatal("One of the segments are returned incorrectly")
	}
	path.Path = "/example/path//"
	segments = GetSegments(&path, "/example/")
	if len(segments) != 1 {
		t.Fatal("GetSegments returns empty segment, that should not happen")
	}
	path.Path = "/example/path//"
	segments = GetSegments(&path, "/example/path/100")
	if len(segments) != 2 {
		t.Fatal("GetSegments should not apply prefix if it does not match")
	}
}

func TestGetQuery(t *testing.T) {
	path := url.URL{RawQuery: url.Values{"foo": {"bar"}, "num": {"1"}}.Encode()}

	val, err := GetQueryStr(&path, "foo")
	if val != "bar" {
		t.Fatal("GetQueryStr returned wrong value")
	}
	if err != nil {
		t.Fatal("GetQueryStr returned err, should have been nil")
	}
	val, err = GetQueryStr(&path, "num")
	if val != "1" {
		t.Fatal("GetQueryStr returned wrong value")
	}
	if err != nil {
		t.Fatal("GetQueryStr returned err, should have been nil")
	}
	val, err = GetQueryStr(&path, "numb")
	if err == nil {
		t.Fatal("GetQueryStr did not return err as expected")
	}
	var number int
	number, err = GetQueryInt(&path, "num")
	if number != 1 {
		t.Fatal("GetQueryInt did not return integer")
	}
	if err != nil {
		t.Fatal("GetQueryInt returned error unexpectedly")
	}
	number, err = GetQueryInt(&path, "numb")
	if err == nil {
		t.Fatal("GetQueryInt did not return error as expected")
	}
}
