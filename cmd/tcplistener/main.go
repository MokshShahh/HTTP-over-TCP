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
			log.Print("error: ", err)
			continue
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("err", err)
		}
		fmt.Printf("Request Line\n- Method: %s\n- Target: %s\n- Version: %s\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range r.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Printf("Body:%s", string(r.Body))
	}
}
