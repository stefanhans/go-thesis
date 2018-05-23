package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {

	// create TCP connection to localhost port 22365
	conn, err := net.Dial("tcp", "localhost:22365")
	if err != nil {
		log.Fatal(err)
	}

	// send message
	fmt.Fprintf(conn, "Hi, it's me :)\n")

	// receive and print reply
	reply, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(reply)

	// close connection
	conn.Close()
}
