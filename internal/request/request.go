package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	requestLine, err := parseRequestLine(string(buf))

	if err != nil {
		return nil, errors.New("Error parsing the request line")
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(buf string) (*RequestLine, error) {
	lines := strings.Split(buf, "\r\n")
	requestLine := lines[0]

	attributes := strings.Split(requestLine, " ")
	if len(attributes) < 3 {
		return nil, errors.New("Not enough attributes in request line")
	}

	// method
	if attributes[0] != strings.ToUpper(attributes[0]) {
		return nil, errors.New("Method is not all uppercase.")
	}

	// target
	if (len(attributes[1]) > 0 && attributes[1][0] != '/') || (len(attributes[1]) > 1 && attributes[1][1:] != strings.ToLower(string(attributes[1][1:]))) {
		return nil, errors.New("Request target isn't all lowercase")
	}

	// version
	if attributes[2] != "HTTP/1.1" {
		return nil, errors.New("HTTP version needs to be HTTP/1.1")
	}

	return &RequestLine{
		Method:        attributes[0],
		RequestTarget: attributes[1],
		HttpVersion:   "1.1",
	}, nil

}
