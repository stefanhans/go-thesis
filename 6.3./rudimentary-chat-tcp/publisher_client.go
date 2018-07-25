package main

import (
	"fmt"
	"log"
	"net"

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

	// Append subscription message in "messages" view
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", memberName, memberIp, memberPort))

	return sendPublisherRequest(newMember)
}

func Unsubscribe(memberName string) error {

	leavingMember := &chatgroup.Message{
		MsgType: chatgroup.Message_UNSUBSCRIBE,
		Sender: &chatgroup.Member{
			Name: memberName}}

	return sendPublisherRequest(leavingMember)
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

	return sendPublisherRequest(message)
}

// Dial publisher and return connection
func sendPublisherRequest(message *chatgroup.Message) error {

	// Connect to publishing service
	conn, err := net.Dial("tcp", publishingService)
	if err != nil {
		return fmt.Errorf("could not connect to publishing service: %v", err)
	}

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	// Write message into connection
	n, err := conn.Write(byteArray)
	if err != nil {
		return fmt.Errorf("could not write message: %v", err)
	}
	log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)

	// Receive reply
	//conn.Read(byteArray)
	//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

	return conn.Close()
}
