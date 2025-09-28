package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {

	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer f.Close()
		currentLine := ""
		for {
			buffer := make([]byte, 8)
			bytesRead, err := f.Read(buffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}
				if err == io.EOF {
					break
				}
				return
			}
			parts := strings.Split(string(buffer[:bytesRead]), "\n")
			for _, part := range parts[:len(parts)-1] {
				ch <- fmt.Sprintf("%s%s", currentLine, part)
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return ch
}
