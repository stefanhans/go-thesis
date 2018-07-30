package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-udp/chat-group"
	"github.com/golang/protobuf/proto"
)

// Start displayer service to provide displaying messages in the text-based UI
func startDisplayer() error {

	// Create listener
	listener, err := net.ListenPacket("udp", displayingService)

	if err != nil {
		log.Fatalf("could not listen to %q: %v\n", displayingService, err)
	}
	defer listener.Close()

	log.Printf("Started displaying service listening on %q\n", displayingService)

	buffer := make([]byte, bufferSize)

	for {
		n, addr, err = listener.ReadFrom(buffer)
		if err != nil {
			log.Printf("cannot read from buffer:%v", err)
		} else {
			go func(buffer []byte, addr net.Addr) {
				handleDisplayerRequest(buffer, addr)

			}(buffer[:n], addr)
		}
	}

	return nil
}

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func handleDisplayerRequest(data []byte, addr net.Addr) {

	// Unmarshall message
	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Errorf("could not unmarshall message: %v", err)
	}

	log.Printf("msg from %v: %v\n", addr, msg)

	// Switch according to the message type and call appropriate handler
	switch msg.MsgType {

	case chatgroup.Message_SUBSCRIBE_REPLY:

		err := handleSubscribeReply(&msg)
		if err != nil {
			log.Printf("could not handleSubscribeReply from %v: %v", addr, err)
		}

	case chatgroup.Message_UNSUBSCRIBE_REPLY:

		err := handleUnsubscribeReply(&msg)
		if err != nil {
			log.Printf("could not handleUnsubscribeReply from %v: %v", addr, err)
		}

	case chatgroup.Message_PUBLISH_REPLY:

		// Handle the protobuf message: Member
		err := handlePublishReply(&msg)
		if err != nil {
			log.Printf("could not handlePublishReply from %v: %v", addr, err)
		}

	default:

		log.Printf("Reply: unknown message type %v\n", msg.MsgType)
	}
}

// Display new member
func handleSubscribeReply(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", msg.Sender.Name, msg.Sender.Ip, msg.Sender.Port))

	return nil
}

func handleUnsubscribeReply(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s has left>", msg.Sender.Name))

	return nil
}

func handlePublishReply(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s: %s", msg.Sender.Name, msg.Text))

	return nil
}

