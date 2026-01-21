package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/MokshShahh/HTTP-over-TCP/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parserState int

const (
	StateInit           parserState = 0
	StateParsingHeaders parserState = 1
	StateDone           parserState = 2
)

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	for {
		switch r.state {
		case StateDone:
			return read, nil
		case StateInit:
			rl, n, err := parseRequestLine(string(data[read:]))
			if n == 0 {
				return read, nil
			}
			if err != nil {
				return 0, err
			}
			r.RequestLine = *rl
			read += n
			r.state = StateParsingHeaders
		case StateParsingHeaders:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			read += n
			if done {
				r.state = StateDone
			}
			return read, nil

		}
	}

}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: make(map[string]string),
	}
}

func parseRequestLine(data string) (*RequestLine, int, error) {
	idx := strings.Index(data, "\r\n")
	if idx == -1 {
		return nil, 0, nil
	}
	line := data[:idx]
	requestLine := strings.Fields(line)
	method := requestLine[0]
	for _, r := range method {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return nil, 0, fmt.Errorf("method is not in all uppercase")
		}
		if !unicode.IsLetter(r) {
			return nil, 0, fmt.Errorf("method is not alphabetical")
		}
	}
	httpVer := strings.TrimSpace(requestLine[2])
	if httpVer != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("incorrect http ver: expected HTTP/1.1, got %s", httpVer)
	}
	rl := &RequestLine{
		Method:        requestLine[0],
		RequestTarget: requestLine[1],
		HttpVersion:   "1.1",
	}
	return rl, idx + 2, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0
	for {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
		if request.state == StateDone {
			break
		}

	}
	return request, nil

}
