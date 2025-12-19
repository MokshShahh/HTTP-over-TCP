# Custom Go HTTP Server

This project is a from-scratch implementation of an HTTP/1.1 server and request parser in Go. It operates directly on raw TCP streams, avoiding the `net/http` standard library to demonstrate the underlying mechanics of web communication.

## Core Features

* **Direct TCP Handling**: Uses the `net` package to manage raw socket connections and data transmission.
* **Manual Request Parsing**: Implements a custom parser to extract methods, paths, and headers from byte streams.
* **RFC 9112 Alignment**: Follows standard HTTP/1.1 message formatting and structure.
* **Concurrency**: Leverages Go routines to handle multiple simultaneous client connections.

## Technical Architecture

The server manages the full lifecycle of an HTTP transaction:

1. **Connection Acceptance**: Listens on a specified port and accepts incoming TCP connections.
2. **Byte Stream Processing**: Reads raw data into buffers, handling the transition from transport layer segments to application layer messages.
3. **State Machine Parsing**: Iterates through bytes to identify the Request Line, Header blocks (separated by `\r\n`), and the message body.
4. **Response Serialization**: Manually constructs status lines and headers into a valid HTTP response format before writing to the socket.
