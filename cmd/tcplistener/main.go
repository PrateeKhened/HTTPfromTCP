package main

import (
	"fmt"
	"log"
	"net"

	"github.com/PrateeKhened/HTTPfromTCP/internal/request"
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

		rq, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", rq.RequestLine.Method, rq.RequestLine.RequestTarget, rq.RequestLine.HttpVersion)
	}
}
