package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
	"github.com/golang/protobuf/proto"
)

// Start displayer service to provide displaying messages in the text-based UI
func startDisplayer() error {

	// Create listener
	listener, err := net.Listen("tcp", displayingService)

	if err != nil {
		log.Fatalf("could not listen to %q: %v\n", displayingService, err)
	}
	defer listener.Close()

	log.Printf("Started displaying service listening on %q\n", displayingService)

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			continue //log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleDisplayerRequest(conn)
	}

	return nil
}

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func handleDisplayerRequest(conn net.Conn) {

	defer conn.Close()

	var buf [512]byte
	var data []byte
	addr := conn.RemoteAddr()

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
		data = append(data, buf[0:n]...)
	}

	log.Printf("Displayer received (%v bytes): %q\n", len(data), data)

	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Errorf("could not unmarshall msg: %v", err)
	}

	log.Printf("msg from %v: %v\n", addr, msg)

	// Switch according to the message type
	switch msg.MsgType {

	case chatgroup.Message_DISPLAY_SUBSCRIPTION:

		// Handle the protobuf message: Member
		err := handleDisplaySubscription(&msg)
		if err != nil {
			log.Printf("could not handleDisplaySubscription from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	case chatgroup.Message_DISPLAY_UNSUBSCRIPTION:

		// Handle the protobuf message: Member
		err := handleDisplayUnsubscription(&msg)
		if err != nil {
			log.Printf("could not handleDisplayUnsubscription from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	case chatgroup.Message_DISPLAY_TEXT:

		// Handle the protobuf message: Member
		err := handleDisplayText(&msg)
		if err != nil {
			log.Printf("could not handleDisplayText from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	case chatgroup.Message_CMD_REPLY:

		// Handle the protobuf message: Member
		err := handleCmdReply(&msg)
		if err != nil {
			log.Printf("could not handleDisplayText from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}


	case chatgroup.Message_MEMBERLIST_REPLY:

		// Handle the protobuf message: Member
		err := handleMemberListReply(&msg)
		if err != nil {
			log.Printf("could not handleDisplayText from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}


	default:

		log.Printf("Displayer: unknown message type %v\n", msg.MsgType)
	}
}

func handleDisplaySubscription(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", msg.Sender.Name, msg.Sender.Ip, msg.Sender.Port))

	return nil
}

func handleDisplayUnsubscription(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s has left>", msg.Sender.Name))

	return nil
}

func handleDisplayText(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s: %s", msg.Sender.Name, msg.Text))

	return nil
}

func handleCmdReply(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s", msg.Text))

	return nil
}

func handleMemberListReply(msg *chatgroup.Message) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%v", msg.MemberList))

	return nil
}


