package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/MokshShahh/HTTP-over-TCP/internal/request"
	"github.com/MokshShahh/HTTP-over-TCP/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
	handler  Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func writeError(w io.Writer, err *HandlerError) {
	response.WriteStatusLine(w, err.StatusCode)
	h := response.GetDefaultHeaders(len(err.Message))
	response.WriteHeaders(w, h)
	w.Write([]byte(err.Message))
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		handler:  handler,
	}

	go s.listen()
	return s, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}
	buf := new(bytes.Buffer)
	handlerErr := s.handler(buf, req)
	if handlerErr != nil {
		writeError(conn, handlerErr)
		return
	}
	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		return
	}
	h := response.GetDefaultHeaders(len(buf.Bytes()))
	err = response.WriteHeaders(conn, h)
	if err != nil {
		return
	}
	conn.Write(buf.Bytes())
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}
