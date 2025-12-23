package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MokshShahh/HTTP-over-TCP/internal/request"
)

func main() {
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal("error: ", err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("err", err)
		}
		fmt.Printf("Request Line\n- Method: %s\n- Target: %s\n- Version: %s", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)

	}
}
