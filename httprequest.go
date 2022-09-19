package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/textproto"
	"os"
	"regexp"
	"strings"
)

type HttpVersion int

const (
	DEFAULT_HTTP_VERSION HttpVersion = iota
	HTTP1_0
	HTTP1_1
	HTTP2
	HTTP3
)

func (h HttpVersion) CurlOption() string {
	switch h {
	case HTTP1_0:
		return "--http1.0"
	case HTTP1_1:
		return "--http1.1"
	case HTTP2:
		return "--http2"
	case HTTP3:
		return "--http3"
	}
	return ""
}

type HeaderField struct {
	Name  string
	Value string
}

type HttpRequest struct {
	Method      string
	URL         string
	HttpVersion HttpVersion
	Header      []HeaderField
	Body        string
}

var requestLine = regexp.MustCompile(`^([0-9A-Za-z]+)\s+(.*?)\s*(HTTP/\d(\.\d)?)?$`)
var fieldLine = regexp.MustCompile("^([!#$%&*+-.^_`|~0-9A-Za-z]+):\\s*(.*)\\s*$")
var commentLine = regexp.MustCompile(`^\s*#.*$`)

func NewHttpRequest(r io.Reader) (*HttpRequest, error) {

	req := new(HttpRequest)

	var line string
	var err error

	reader := textproto.NewReader(bufio.NewReader(r))
	for {
		line, err = reader.ReadLine()
		if err != nil {
			return nil, err
		}

		if commentLine.MatchString(line) {
			continue
		}

		match := requestLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		req.Method = match[1]
		req.URL = match[2]
		switch match[3] {
		case "HTTP/1.0":
			req.HttpVersion = HTTP1_0
		case "HTTP/1.1":
			req.HttpVersion = HTTP1_1
		case "HTTP/2":
			req.HttpVersion = HTTP2
		case "HTTP/3":
			req.HttpVersion = HTTP3
		case "":
			req.HttpVersion = DEFAULT_HTTP_VERSION
		default:
		}
		break
	}

	for {
		line, err = reader.ReadContinuedLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if line == "" {
			break
		}
		if commentLine.MatchString(line) {
			continue
		}

		match := fieldLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		req.Header = append(req.Header, HeaderField{Name: match[1], Value: match[2]})
	}

	var sb strings.Builder
	for {
		line, err = reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if commentLine.MatchString(line) {
			continue
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}
	req.Body = sb.String()

	return req, nil
}

type JsonConfig struct {
	Method string            `json:"method"`
	URL    string            `json:"url"`
	Header map[string]string `json:"header"`
	Body   any               `json:"body"`
}

func NewHttpRequestJson(r io.Reader) (*HttpRequest, error) {

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var d JsonConfig
	if err = json.Unmarshal(b, &d); err != nil {
		return nil, err
	}

	req := new(HttpRequest)

	req.Method = d.Method
	req.URL = d.URL
	for n, v := range d.Header {
		req.Header = append(req.Header, HeaderField{Name: n, Value: v})
	}

	if d.Body == nil {
		req.Body = ""
	} else {
		switch v := d.Body.(type) {
		case string:
			req.Body = v
		default:
			body, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			req.Body = string(body)
		}
	}
	return req, nil
}

func (req *HttpRequest) ExpandEnv() {
	req.Method = os.ExpandEnv(req.Method)
	req.URL = os.ExpandEnv(req.URL)
	for i := 0; i < len(req.Header); i++ {
		req.Header[i].Name = os.ExpandEnv(req.Header[i].Name)
		req.Header[i].Value = os.ExpandEnv(req.Header[i].Value)
	}
	req.Body = os.ExpandEnv(req.Body)
}

func shellstring(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, `'`, `'\''`))
}

func (req *HttpRequest) CurlCommand(prettyprint bool, flags []string) string {
	command := []string{"curl"}

	if flags != nil {
		command = append(command, flags...)
	}

	command = append(command, shellstring(req.URL))

	if prettyprint {
		command = append(command, "\\\n")
	}

	if req.Method != "GET" || req.Body != "" {
		command = append(command, "-X", shellstring(req.Method))
	}
	if req.HttpVersion != DEFAULT_HTTP_VERSION {
		command = append(command, req.HttpVersion.CurlOption())
	}

	for _, field := range req.Header {
		if prettyprint {
			command = append(command, "\\\n")
		}
		command = append(command, "-H", shellstring(fmt.Sprintf("%s: %s", field.Name, field.Value)))
	}

	if req.Body != "" {
		if prettyprint {
			command = append(command, "\\\n")
		}
		command = append(command, "-d", shellstring(req.Body))
	}

	return strings.Join(command, " ")
}
