package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewHttpRequest(t *testing.T) {
	test := `
GET https://example.com/ HTTP/1.1
User-Agent: github.com/shirokurostone/curl-template
	`

	want := HttpRequest{
		Method:      "GET",
		URL:         "https://example.com/",
		HttpVersion: HTTP1_1,
		Header: []HeaderField{
			HeaderField{"User-Agent", "github.com/shirokurostone/curl-template"},
		},
		Body: "",
	}

	req, err := NewHttpRequest(strings.NewReader(test))
	if err != nil {
		t.Fatalf("err = %#v, want %#v", err, nil)
	}

	if req.Method != want.Method {
		t.Errorf("req.Method = %#v, want %#v", req.Method, want.Method)
	}

	if req.URL != want.URL {
		t.Errorf("req.URL = %#v, want %#v", req.URL, want.URL)
	}

	if req.HttpVersion != want.HttpVersion {
		t.Errorf("req.HttpVersion = %#v, want %#v", req.HttpVersion, want.HttpVersion)
	}

	if !reflect.DeepEqual(req.Header, want.Header) {
		t.Errorf("req.Header = %#v, want %#v", req.Header, want.Header)
	}

	if req.Body != want.Body {
		t.Errorf("req.Body = %#v, want %#v", req.Body, want.Body)
	}

}
