package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/MokshShahh/HTTP-over-TCP/internal/headers"
)

type Response struct {
}

type StatusCode int

const (
	StatusOK        StatusCode = 200
	StatusBadReq    StatusCode = 400
	StatusServerErr StatusCode = 500
)

type writerState int

const (
	stateStatus writerState = iota
	stateHeaders
	stateBody
	stateDone
)

type Writer struct {
	w     io.Writer
	state writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:     w,
		state: stateStatus,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != stateStatus {
		return fmt.Errorf("not in status line state")
	}
	statusLine := []byte{}
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadReq:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusServerErr:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("unrecognised code")
	}
	_, err := w.w.Write(statusLine)
	w.state = stateHeaders
	return err

}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != stateHeaders {
		return fmt.Errorf("state is not header state")
	}
	for key, value := range headers {
		_, err := w.w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.w.Write([]byte("\r\n"))
	w.state = stateBody
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != stateBody {
		return 0, fmt.Errorf("state is not state body")
	}
	n, err := w.w.Write(p)
	if err != nil {
		return 0, err
	}
	w.state = stateDone
	return n, nil

}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = strconv.Itoa(contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h

}
