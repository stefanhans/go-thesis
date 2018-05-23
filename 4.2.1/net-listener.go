package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {

	// listen for TCP connections on localhost port 22365
	listener, err := net.Listen("tcp", "localhost:22365")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// wait for connections
	for {
		// accept connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// create a goroutine for connection
		go func(conn net.Conn) {

			// read and print the message
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Message received: %s", msg)

			// send reply
			conn.Write([]byte(fmt.Sprintf("Message accepted: %s", msg)))

			// close connection
			conn.Close()
		}(conn)
	}
}
