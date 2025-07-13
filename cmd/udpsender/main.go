package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const address = "localhost:42069"

func main() {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("error resolving UDP: %s", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("error dialing UDP: %s", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")

		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("Error reading the input: %s", err)
		}

		// FIXME: Assuming isPrefix is always false
		if isPrefix == true {
			log.Fatalf("Line is too long: %s\n", line)
			continue
		}

		_, err = conn.Write(line)
		if err != nil {
			log.Fatalf("Error writing to connection")
		}

	}
}
