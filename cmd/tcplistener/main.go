package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("could not open %s: %s\n", port, err)
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Connection error: %s", err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()
		currentRead := ""

		for {
			buffer := make([]byte, 8, 8)

			n, err := f.Read(buffer)
			if err != nil {
				if currentRead != "" {
					lines <- currentRead
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			currentRead += string(buffer[:n])

			for len(currentRead) > 0 {
				index := strings.Index(currentRead, "\n")
				if index == -1 {
					break
				}

				line := currentRead[:index]
				lines <- line

				currentRead = currentRead[index+1:]
			}

		}

	}()

	return lines
}
