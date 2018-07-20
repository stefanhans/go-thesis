package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-tcp/member-group"
	"github.com/golang/protobuf/proto"
	"os"
)

const (
	// ReadBytes delimiter
	EOF byte = '\x08'

	// API
	Join = iota
	Members
	Update
	Leave
)

func main() {

	member := &membergroup.Member{
		Name:   "Karl",
		Ip:     "localhost",
		Port:   "12345",
		Leader: false,
	}

	fmt.Printf("%b\n", EOF)

	// Create listener
	conn, err := net.Dial("tcp", ":22365")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Dial to Rootserver on localhost:22365\n")

	// Marshal into binary format
	byteArray, err := proto.Marshal(member)
	if err != nil {
		fmt.Errorf("could not encode info: %v", err)
		os.Exit(1)
	}

	// Prepend message type
	msgType := []byte{Join}
	//fmt.Printf("Message Type: %T %#v\n", msgType, msgType)
	byteArray = append(msgType, byteArray...)
	byteArray = append(byteArray, EOF)

	conn.Write(byteArray)
	fmt.Printf("Member sent (%v byte): %v\n", len(byteArray), member)

	// Receive reply
	conn.Close()
}
