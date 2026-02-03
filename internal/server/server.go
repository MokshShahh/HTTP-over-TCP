package server

import (
	"fmt"
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

type Handler func(w *response.Writer, req *request.Request) *HandlerError

func writeError(w *response.Writer, err *HandlerError) {
	w.WriteStatusLine(err.StatusCode)
	h := response.GetDefaultHeaders(len(err.Message))
	w.WriteHeaders(h)
	w.WriteBody([]byte(err.Message))
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

	resWriter := response.NewWriter(conn)
	handlerErr := s.handler(resWriter, req)

	if handlerErr != nil {
		writeError(resWriter, handlerErr)
		return
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}
