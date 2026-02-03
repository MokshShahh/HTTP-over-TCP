package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MokshShahh/HTTP-over-TCP/internal/request"
	"github.com/MokshShahh/HTTP-over-TCP/internal/response"
	"github.com/MokshShahh/HTTP-over-TCP/internal/server"
)

const port = 42069
const html400 = `<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>`
const html500 = `<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>`
const html200 = `<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>`

func mainHandler(w *response.Writer, req *request.Request) *server.HandlerError {
	var body string
	var status response.StatusCode

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		body = html400
		status = response.StatusBadReq
	case "/myproblem":
		body = html500
		status = response.StatusServerErr
	default:
		body = html200
		status = response.StatusOK
	}

	bodyBytes := []byte(body)
	h := response.GetDefaultHeaders(len(bodyBytes))
	h["content-type"] = "text/html"

	w.WriteStatusLine(status)
	w.WriteHeaders(h)
	w.WriteBody(bodyBytes)

	return nil
}

func main() {
	server, err := server.Serve(port, mainHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
