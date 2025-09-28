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
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read the io.Reader input:%s", err.Error())
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse request line:%s", err.Error())
	}
	return &Request{*requestLine}, nil

}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, fmt.Errorf("could not find CRLF in request-line")
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}
	return requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	var r RequestLine

	requestLine := str
	if requestLine == "" {
		return nil, errors.New("the request line is empty")
	}
	requestLineParts := strings.Fields(requestLine)

	if len(requestLineParts) != 3 {
		return nil, errors.New("check the request line semantic")
	}

	requestMethod := requestLineParts[0]
	for _, r := range requestMethod {
		if r < 'A' || r > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", requestMethod)
		}
	}
	r.Method = requestMethod

	requestTarget := requestLineParts[1]
	r.RequestTarget = requestTarget

	parts := strings.Split(requestLineParts[2], "/")
	if len(parts) != 2 {
		return nil, errors.New("HTTP version must be in format HTTP/1.1")
	}
	if parts[0] != "HTTP" {
		return nil, errors.New("the protocol must be HTTP")
	}
	if parts[1] != "1.1" {
		return nil, fmt.Errorf("HTTP version unsupported: %s", parts[1])
	}
	r.HttpVersion = parts[1]

	return &r, nil
}
