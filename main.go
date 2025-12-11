package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)
		line := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				line += string(data[:i])
				data = data[i+1:]
				out <- line
				line = ""
			}
			line += string(data)
			if err != nil {
				break
			}
		}

	}()

	return out

}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", err)
	}
	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("Read: %s\n", line)

	}
}
