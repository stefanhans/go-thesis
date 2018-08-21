package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/6.4./rudimentary-chat-tcp/chat-group"
	"github.com/golang/protobuf/proto"
)

// Publish a message to all members except the sender
func publishMessage(message *chatgroup.Message, msgType chatgroup.Message_MessageType) error {

	// Set the reply message type
	message.MsgType = msgType

	// Forward message to other chat group members
	for _, recipient := range selfMemberList {

		// Exclude sender
		if recipient.Name != message.Sender.Name {

			// Send message to recipient
			log.Printf("From %s to %s (%s:%s): %q\n",
				message.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, message.Sender)
			err := sendMessage(message, recipient.Ip+":"+recipient.Port)
			if err != nil {
				return fmt.Errorf("failed send reply: %v", err)
			}
		}
	}
	return nil
}

// Send reply to the sender of the message
func sendMessage(message *chatgroup.Message, recipient string) error {

	// Connect to the recipient
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		return fmt.Errorf("could not connect to recipient %q: %v", recipient, err)
	}

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	// Write the bytes to the connection
	n, err := conn.Write(byteArray)
	if err != nil {
		return fmt.Errorf("could not write message to the connection: %v", err)
	}
	log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)

	// Close connection
	return conn.Close()
}

func updateRemoteIP(msg *chatgroup.Message, addr net.Addr) {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		log.Printf("Remote Ip address update from %v to %v\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]
	}
}
