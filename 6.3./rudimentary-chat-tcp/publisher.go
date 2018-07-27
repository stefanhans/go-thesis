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

			// Append text messages in "messages" view of subscriber
			displayText(fmt.Sprintf("<%s (%s:%s) has joined>", selfMember.Name, selfMember.Ip, selfMember.Port))

			return nil
		}

		// Exit on unexpected error
		log.Fatalf("could not listen to %q: %v\n", publishingService, err)
	}
	defer listener.Close()

	log.Printf("Started publishing service listening on %q\n", publishingService)

	// Append text messages in "messages" view of publisher
	displayText(fmt.Sprintf("<publishing service running: %s (%s:%s)>", selfMember.Name, serverIp, serverPort))

	// Subscribe directly at started publishing service
	selfMember.Leader = true
	cgMember = append(cgMember, selfMember)
	log.Printf("Subscribed directly at started publishing service: %v\n", cgMember[0])

	// Append text messages in "messages" view of publisher
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", selfMember.Name, selfMember.Ip, selfMember.Port))

	// Endless loop in foreground of goroutine
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

	case chatgroup.Message_CMD_LIST:

		log.Printf("CMD_LIST: %v\n", msg)

		// Handle the protobuf message: Member
		err := handleCmdList(&msg, addr)
		if err != nil {
			fmt.Printf("could not handleCmdList from %v: %v", addr, err)
		}

		//_, err = conn.Write([]byte(""))
		//if err != nil {
		//	return
		//}

	default:

		log.Printf("publisher: unknown message type %v\n", msg.MsgType)
	}
}

func handleSubscribe(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	// Check subscriber for uniqueness
	for _, recipient := range cgMember {
		if recipient.Name == message.Sender.Name {
			return fmt.Errorf("name %q already used", message.Sender.Name)
		}
		if recipient.Ip == message.Sender.Ip && recipient.Port == message.Sender.Port {
			return fmt.Errorf("address %s:%s already used by %s", recipient.Ip, recipient.Port, recipient.Name)
		}
	}

	// Add subscriber
	log.Printf("Add subscriber: %v\n", message.Sender)
	cgMember = append(cgMember, message.Sender)
	log.Printf("Current members registered: %v\n", cgMember)

	err := publishDisplayerRequest(message, chatgroup.Message_DISPLAY_SUBSCRIPTION)
	if err != nil {
		fmt.Errorf("Failed to publish Message_DISPLAY_SUBSCRIPTION", err)
	}

	return nil
}

func handleUnsubscribe(message *chatgroup.Message) error {

	log.Printf("Unregister: %v\n", message.Sender)

	// Remove subscriber
	for i, s := range cgMember {
		if s.Name == message.Sender.Name {
			cgMember = append(cgMember[:i], cgMember[i+1:]...)
			break
		}
	}
	log.Printf("Current members registered: %v\n", cgMember)

	err := publishDisplayerRequest(message, chatgroup.Message_DISPLAY_UNSUBSCRIPTION)
	if err != nil {
		fmt.Errorf("Failed to publish Message_DISPLAY_UNSUBSCRIPTION", err)
	}

	return nil
}

func handlePublish(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	log.Printf("Publish from %v: %q\n", message.Sender.Name, message.Text)

	err := publishDisplayerRequest(message, chatgroup.Message_DISPLAY_TEXT)
	if err != nil {
		fmt.Errorf("Failed to publish Message_DISPLAY_TEXT", err)
	}

	return nil
}

func handleCmdList(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	log.Printf("List request from %v: %q\n", message.Sender.Name, message.Text)

	err := executeCmdList(message)
	if err != nil {
		fmt.Errorf("Failed to execute list request", err)
	}

	err = replyCmdRequest(message)
	if err != nil {
		fmt.Errorf("Failed to reply to list request", err)
	}

	return nil
}

func updateRemoteIP(msg *chatgroup.Message, addr net.Addr) {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		log.Printf("Remote Ip address update from %v to %v\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]
	}
}

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

func publishDisplayerRequest(message *chatgroup.Message, msgType chatgroup.Message_MessageType) error {

	message.MsgType = msgType

	// Forward message to other chat group members
	for _, recipient := range cgMember {

		// Exclude sender and publisher from message forwarding
		if recipient.Name != message.Sender.Name {
			log.Printf("From %s to %s (%s:%s): %q\n",
				message.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, message.Sender)

			err := sendDisplayerRequest(message, recipient.Ip+":"+recipient.Port)
			if err != nil {
				fmt.Errorf("Failed send displayer request", err)
			}

			//conn, err := net.Dial("tcp", recipient.Ip+":"+recipient.Port)
			//if err != nil {
			//	return fmt.Errorf("could not connect to displaying service: %v", err)
			//}
			//
			//// Marshal into binary format
			//byteArray, err := proto.Marshal(message)
			//if err != nil {
			//	return fmt.Errorf("could not encode message: %v", err)
			//}
			//
			//n, err := conn.Write(byteArray)
			//log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, message)
			//
			//err = conn.Close()
			//if err != nil {
			//	return fmt.Errorf("could not close connection: %v", err)
			//}
		}
	}
	return nil
}

func replyCmdRequest(message *chatgroup.Message) error {

	message.MsgType = chatgroup.Message_CMD_REPLY

	err := sendDisplayerRequest(message, message.Sender.Ip+":"+message.Sender.Port)
	if err != nil {
		return fmt.Errorf("could not send command reply: %v", err)
	}

	return nil
}

func executeCmdList(message *chatgroup.Message) error {

	text := ""
	for i, member := range cgMember {
		text += fmt.Sprintf("<LIST>: %v: %s\n", i, member)
	}
	message.Text = strings.Trim(text, "\n")

	message.Sender.Name = selfMember.Name

	return nil
}
