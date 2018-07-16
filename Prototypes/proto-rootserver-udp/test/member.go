package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-udp/member-group"
	"github.com/golang/protobuf/proto"
	"os"
)

const (
	MemberMessage = iota
	MemberListMessage
)

func main() {

	member := &membergroup.Member{
		Name:   "Karl",
		Ip:     "localhost",
		Port:   "12345",
		Leader: false,
	}







	// Create listener
	conn, err := net.Dial("udp", "localhost:22365")
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
	msgType := []byte{MemberMessage}
	//fmt.Printf("Message Type: %T %#v\n", msgType, msgType)
	byteArray = append(msgType, byteArray...)

	conn.Write(byteArray)
	fmt.Printf("Member sent (%v byte): %v\n", len(byteArray), member)

	n, err := conn.Read(byteArray)
	fmt.Printf("Reply: %q\n", byteArray[:n])

	conn.Close()
}

