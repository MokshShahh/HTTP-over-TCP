package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", err)
	}
	line := ""
	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		data = data[:n]
		fmt.Print(data)
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			line += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", line)
			line = ""
		}
		line += string(data)
		if err != nil {
			break
		}
	}
}
