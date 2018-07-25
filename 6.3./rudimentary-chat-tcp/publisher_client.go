package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
	"github.com/golang/protobuf/proto"
)

func Subscribe() error {

	newMember := &chatgroup.Message{
		MsgType: chatgroup.Message_SUBSCRIBE,
		Sender: &chatgroup.Member{
			Name:   memberName,
			Ip:     memberIp,
			Port:   memberPort,
			Leader: false}}

	conn, err := net.Dial("tcp", "127.0.0.1:22365")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Dial to Rootserver on 127.0.0.1:22365\n")

	// Marshal into binary format
	byteArray, err := proto.Marshal(newMember)
	if err != nil {
		fmt.Errorf("could not encode new member: %v", err)
		os.Exit(1)
	}

	n, err := conn.Write(byteArray)
	log.Printf("New member (%v byte) sent (%v byte): %v\n", len(byteArray), n, newMember)

	//conn.Read(byteArray)
	//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

	// Receive reply
	conn.Close()

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", memberName, memberIp, memberPort))

	return nil
}

func Unsubscribe(memberName string) error {

	leavingMember := &chatgroup.Message{
		MsgType: chatgroup.Message_UNSUBSCRIBE,
		Sender: &chatgroup.Member{
			Name: memberName}}

	conn, err := net.Dial("tcp", publishingService)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Dial to Rootserver on 127.0.0.1:22365\n")

	// Marshal into binary format
	byteArray, err := proto.Marshal(leavingMember)
	if err != nil {
		fmt.Errorf("could not encode leaving member: %v", err)
		os.Exit(1)
	}

	n, err := conn.Write(byteArray)
	fmt.Printf("New member (%v byte) sent (%v byte): %v\n", len(byteArray), n, leavingMember)

	//conn.Read(byteArray)
	//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

	// Receive reply
	conn.Close()

	return nil
}

func Publish(text string) error {
	message := &chatgroup.Message{
		MsgType: chatgroup.Message_PUBLISH,
		Sender: &chatgroup.Member{
			Name:   memberName,
			Ip:     memberIp,
			Port:   memberPort,
			Leader: false},
		Text: text}

	conn, err := net.Dial("tcp", "127.0.0.1:22365")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Dial to Rootserver on 127.0.0.1:22365\n")

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		fmt.Errorf("could not encode new member: %v", err)
		os.Exit(1)
	}

	n, err := conn.Write(byteArray)
	log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)

	//conn.Read(byteArray)
	//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

	// Receive reply
	conn.Close()

	return nil
}

// Dial publisher and return connection
func dialPublisher() (net.Conn, error) {

	conn, err := net.Dial("tcp", ":"+serverPort)
	if err != nil {
		return nil, fmt.Errorf("could not connect to publisher: %v", err)
	}
	return conn, nil
}
