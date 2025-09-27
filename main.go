package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
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
				fmt.Printf("error: %s\n", err.Error())
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
