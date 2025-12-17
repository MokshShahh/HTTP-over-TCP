package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
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
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\r\n")
	requestLine := strings.Fields(lines[0])
	method := requestLine[0]
	for _, r := range method {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return nil, fmt.Errorf("method is not in all uppercase")
		}
		if !unicode.IsLetter(r) {
			return nil, fmt.Errorf("method is not alphabetical")
		}
	}
	httpVer := strings.TrimSpace(requestLine[2])
	if httpVer != "HTTP/1.1" {
		return nil, fmt.Errorf("incorrect http ver: expected HTTP/1.1, got %s", httpVer)
	}
	return &Request{
		RequestLine: RequestLine{
			Method:        requestLine[0],
			RequestTarget: requestLine[1],
			HttpVersion:   "1.1",
		},
	}, nil
}
