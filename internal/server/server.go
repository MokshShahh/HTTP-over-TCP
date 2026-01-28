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
	handler Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError{

}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
		handler: handler
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

	_, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}
	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		return
	}
	h := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, h)
	if err != nil {
		return
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}
