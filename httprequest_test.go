package main

import (
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestNewHttpRequest(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  HttpRequest
	}{
		{
			name: "GET",
			input: `GET https://example.com/
`,
			want: HttpRequest{
				Method:      "GET",
				URL:         "https://example.com/",
				HttpVersion: DEFAULT_HTTP_VERSION,
				Header:      []HeaderField{},
				Body:        "",
			},
		},
		{
			name: "POST",
			input: `POST https://example.com/ HTTP/1.1
User-Agent: github.com/shirokurostone/curl-template
Content-Type: application/json

{"key":"value"}`,
			want: HttpRequest{
				Method:      "POST",
				URL:         "https://example.com/",
				HttpVersion: HTTP1_1,
				Header: []HeaderField{
					HeaderField{"User-Agent", "github.com/shirokurostone/curl-template"},
					HeaderField{"Content-Type", "application/json"},
				},
				Body: "{\"key\":\"value\"}\n",
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(*testing.T) {
			pattern := regexp.MustCompile(`^\s+`)
			s := pattern.ReplaceAllString(testcase.input, "")

			req, err := NewHttpRequest(strings.NewReader(s))
			if err != nil {
				t.Fatalf("err = %#v, want %#v", err, nil)
			}

			want := testcase.want
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
		})
	}

}
