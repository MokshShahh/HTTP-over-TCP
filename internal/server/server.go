package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/MokshShahh/HTTP-over-TCP/internal/request"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: l,
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

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!"

	conn.Write([]byte(response))
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}
