package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"syscall"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
	"github.com/golang/protobuf/proto"
)

// Start publisher service to provide member registration and message publishing
func startPublisher() error {

	// Create listener
	listener, err := net.Listen("tcp", publishingService)

	if err != nil {

		// Check if publisher error is "address already in use"
		if strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {

			// Subscribe at already running Publisher
			err = Subscribe()
			if err != nil {
				log.Fatalf("failed to subscribe at already running Publisher: %v", err)
			}
			return nil
		}

		// Exit on unexpected error
		log.Fatalf("could not listen to %q: %v\n", publishingService, err)
	}
	defer listener.Close()

	log.Printf("Started publishing service listening on %q\n", publishingService)

	// Subscribe directly at started publishing service
	cgMember = append(cgMember, &chatgroup.Member{Name: memberName, Ip: memberIp, Port: memberPort, Leader: true})
	log.Printf("Subscribed directly at started publishing service: %v\n", cgMember[0])

	selfMember.Leader = true

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection from publishing service listener: %s\n", err)
			continue
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handlePublisherRequest(conn)
	}

	return nil
}

// Read all incoming data, take the message type,
// and use the appropriate handler for the rest
func handlePublisherRequest(conn net.Conn) {

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

	log.Printf("Publisher received %v bytes\n", len(data))

	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Errorf("could not unmarshall message: %v", err)
	}

	// Switch according to the message type
	switch msg.MsgType {

	case chatgroup.Message_SUBSCRIBE:

		log.Printf("SUBSCRIBE: %v\n", msg)

		err := handleSubscribe(&msg, addr)
		if err != nil {
			fmt.Printf("could not handleSubscribe from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	case chatgroup.Message_UNSUBSCRIBE:

		log.Printf("UNSUBSCRIBE: %v\n", msg)

		// Handle the protobuf message: Member
		err := handleUnsubscribe(&msg)
		if err != nil {
			fmt.Printf("could not handleUnsubscribe from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	case chatgroup.Message_PUBLISH:

		log.Printf("PUBLISH: %v\n", msg)

		// Handle the protobuf message: Member
		err := handlePublish(&msg, addr)
		if err != nil {
			fmt.Printf("could not handlePublish from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	default:

		log.Printf("publisher: unknown message type %v\n", msg.MsgType)
	}
}

func handleSubscribe(msg *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(msg, addr)

	// Check subscriber for uniqueness
	for _, recipient := range cgMember {
		if recipient.Name == msg.Sender.Name {
			return fmt.Errorf("name %q already used", msg.Sender.Name)
		}
		if recipient.Ip == msg.Sender.Ip && recipient.Port == msg.Sender.Port {
			return fmt.Errorf("address %s:%s already used by %s", recipient.Ip, recipient.Port, recipient.Name)
		}
	}

	// Add subscriber
	log.Printf("Add subscriber: %v\n", msg.Sender)
	cgMember = append(cgMember, msg.Sender)
	log.Printf("Current members registered: %v\n", cgMember)

	// Forward message to other chat group members
	for _, recipient := range cgMember {

		msg.MsgType = chatgroup.Message_DISPLAY_SUBSCRIPTION

		// Exclude sender and publisher from message forwarding
		if recipient.Name != msg.Sender.Name && recipient.Name != memberName {
			log.Printf("From %s to %s (%s:%s): %q\n", msg.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, msg.Sender)

			err := sendDisplayerRequest(msg, recipient.Ip+":"+recipient.Port)
			if err != nil {
				fmt.Errorf("Failed send displayer request", err)
			}
		}
	}

	// Append text message in "messages" view of publisher
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", msg.Sender.Name, msg.Sender.Ip, msg.Sender.Port))

	return nil
}

func handleUnsubscribe(msg *chatgroup.Message) error {

	log.Printf("Unregister: %v\n", msg.Sender)

	// Remove subscriber
	for i, s := range cgMember {
		if s.Name == msg.Sender.Name {
			cgMember = append(cgMember[:i], cgMember[i+1:]...)
			break
		}
	}
	log.Printf("Current members registered: %v\n", cgMember)

	// Send message to other subscribers via gRPC Displayer service
	for _, recipient := range cgMember {

		msg.MsgType = chatgroup.Message_DISPLAY_UNSUBSCRIPTION

		// Exclude sender and publisher from message forwarding
		if recipient.Name != msg.Sender.Name && recipient.Name != memberName {
			log.Printf("From %s to %s (%s:%s): %q\n", msg.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, msg.Sender)

			err := sendDisplayerRequest(msg, recipient.Ip+":"+recipient.Port)
			if err != nil {
				fmt.Errorf("Failed send displayer request", err)
			}
		}
	}

	// Append text message in "messages" view of publisher
	displayText(fmt.Sprintf("<%s has left>", msg.Sender.Name))

	return nil
}

func handlePublish(msg *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(msg, addr)

	log.Printf("Publish from %v: %q\n", msg.Sender.Name, msg.Text)

	msg.MsgType = chatgroup.Message_DISPLAY_TEXT

	// Send message to other subscribers via gRPC Displayer service
	for _, recipient := range cgMember {

		// Exclude sender and publisher from message forwarding
		if recipient.Name != msg.Sender.Name {
			log.Printf("From %s to %s (%s:%s): %q\n", msg.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, msg.Sender)

			err := sendDisplayerRequest(msg, recipient.Ip+":"+recipient.Port)
			if err != nil {
				fmt.Errorf("Failed send displayer request", err)
			}
		}
	}

	// Append text message in "messages" view of publisher
	//displayText(fmt.Sprintf("%s: %s", msg.Sender.Name, msg.Text))

	return nil
}

func updateRemoteIP(msg *chatgroup.Message, addr net.Addr) {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		log.Printf("Remote Ip address update from %v to %v\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]
	}
}

//
func sendDisplayerRequest(message *chatgroup.Message, service string) error {

	conn, err := net.Dial("tcp", service)
	if err != nil {
		return fmt.Errorf("could not connect to displaying service: %v", err)
	}

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	n, err := conn.Write(byteArray)
	log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)

	//conn.Read(byteArray)
	//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

	// Receive reply
	return conn.Close()
}
