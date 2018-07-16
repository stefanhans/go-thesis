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

	// Declare array with protobuffer messages
	members := &membergroup.MemberList{Member: []*membergroup.Member{&membergroup.Member{
		Name:   "Stefan",
		Ip:     "localhost",
		Port:   "12345",
		Leader: false,
	}, &membergroup.Member{
		Name:   "I am a painter",
		Ip:     "Marc Chagall",
		Port:   "12345",
		Leader: false,
	}}}

	// Create listener
	conn, err := net.Dial("udp", "localhost:22365")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Dial to Rootserver on localhost:22365\n")

	// Marshal into binary format
	byteArray, err := proto.Marshal(members)
	if err != nil {
		fmt.Errorf("could not encode info: %v", err)
		os.Exit(1)
	}

	// Prepend message type
	msgType := []byte{MemberListMessage}
	//fmt.Printf("Message Type: %T %#v\n", msgType, msgType)
	byteArray = append(msgType, byteArray...)

	conn.Write(byteArray)
	fmt.Printf("Member sent (%v byte): %v\n", len(byteArray), members)

	n, err := conn.Read(byteArray)
	fmt.Printf("Reply: %q\n", byteArray[:n])

	conn.Close()
}

