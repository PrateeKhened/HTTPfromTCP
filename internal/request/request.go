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
	state       parseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parseState int

const (
	stateInitialised parseState = iota
	stateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{state: stateInitialised}

	buf := make([]byte, bufferSize)
	readToIndex := 0

	for r.state != stateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		nRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF && nRead == 0 {
			// no more data; let parser decide if done
		} else if err != nil && err != io.EOF {
			return nil, err
		}
		readToIndex += nRead

		nParsed, perr := r.parse(buf[:readToIndex])
		if perr != nil {
			return nil, perr
		}

		if nParsed > 0 {
			copy(buf, buf[nParsed:readToIndex])
			readToIndex -= nParsed
		}

		if err == io.EOF && nRead == 0 {
			break
		}
	}

	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialised:
		n, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.state = stateDone
		return n, nil
	case stateDone:
		return 0, fmt.Errorf("parser in done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return 0, nil, err
	}
	return idx + len(crlf), requestLine, nil
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
