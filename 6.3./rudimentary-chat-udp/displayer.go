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
	displayingListener, err := net.ListenPacket("udp", displayingService)

	if err != nil {
		log.Fatalf("could not listen to %q: %v\n", displayingService, err)
	}
	defer displayingListener.Close()

	log.Printf("Started displaying service listening on %q\n", displayingService)

	buffer := make([]byte, bufferSize)

	for {
		n, addr, err := displayingListener.ReadFrom(buffer)
		if err != nil {
			log.Printf("cannot read from buffer: %v", err)
		} else {
			//log.Printf("Read %v bytes from %v: %v\n", n, addr, buffer)
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

	// Fetch the handler from a map by the message type and call it accordingly
	if replyAction, ok := replyActionMap[msg.MsgType]; ok {
		log.Printf("%v\n", msg)
		err := replyAction(&msg)
		if err != nil {
			fmt.Printf("could not handle %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		log.Printf("displayer: unknown message type %v\n", msg.MsgType)
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

