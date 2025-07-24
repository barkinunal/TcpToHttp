package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	parserState parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState int

const (
	initialised parserState = iota
	done
)

const CRLF = "\r\n"
const BUFFER_SIZE = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		parserState: initialised,
	}
	buffer := make([]byte, BUFFER_SIZE, BUFFER_SIZE)
	readToIndex := 0

	for request.parserState != done {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.parserState = done
				break
			}
			return nil, err
		}

		readToIndex += numBytesRead

		numBytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func (request *Request) parse(data []byte) (int, error) {
	numBytesConsumed, requestLine, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}

	if numBytesConsumed != 0 {
		request.RequestLine = *requestLine
		request.parserState = done
	}

	return numBytesConsumed, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return 0, nil, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return 0, nil, err
	}
	return idx, requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}
