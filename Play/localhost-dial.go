package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:22365")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not dial: %v", err)
		os.Exit(1)
	}
	fmt.Fprintf(conn, "Hi it's me :)\n")
	status, err := bufio.NewReader(conn).ReadString('\n')

	fmt.Print(status)
}
