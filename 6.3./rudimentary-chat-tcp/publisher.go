package main

import (
	"fmt"
	"log"
	"net"
	"os"
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
				log.Fatalf("Subscribe: %v", err)
			}
			return nil
		}

		// Exit on unexpected error
		log.Fatalf("could not listen to %q: %v\n", publishingService, err)
	}
	defer listener.Close()

	// Subscribe directly
	cgMember = append(cgMember, &chatgroup.Member{Name: memberName, Ip: memberIp, Port: memberPort, Leader: true})
	log.Printf("Subscribed directly: %v\n", cgMember[0])

	log.Printf("Start publishing service listening on %s\n", publishingService)

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			continue //log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleChatgroup(conn)
	}

	return nil
}

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func handleChatgroup(conn net.Conn) {
	log.Printf("handleChatgroup(conn net.Conn)\n")

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

	log.Printf("Publisher received (%v bytes): %q\n", len(data), data)

	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Errorf("could not unmarshall msg: %v", err)
	}

	log.Printf("msg from %v: %v\n", addr, msg)

	// Switch according to the message type
	switch msg.MsgType {

	case chatgroup.Message_SUBSCRIBE:

		// Handle the protobuf message: Member
		err := handleSubscribe(&msg, addr)
		if err != nil {
			fmt.Printf("could not handleSubscribe from %v: %v", addr, err)
		}
		_, err = conn.Write([]byte("12345678901234567890123456789012345678901234567890"))
		if err != nil {
			return
		}

	case chatgroup.Message_UNSUBSCRIBE:

		// Handle the protobuf message: Member
		err := handleUnsubscribe(&msg)
		if err != nil {
			fmt.Printf("could not handleSubscribe from %v: %v", addr, err)
		}
		_, err = conn.Write([]byte("12345678901234567890123456789012345678901234567890"))
		if err != nil {
			return
		}

	case chatgroup.Message_PUBLISH:

		// Handle the protobuf message: Member
		err := handlePublish(&msg, addr)
		if err != nil {
			fmt.Printf("could not handleSubscribe from %v: %v", addr, err)
		}
		_, err = conn.Write([]byte("12345678901234567890123456789012345678901234567890"))
		if err != nil {
			return
		}

	default:

		fmt.Printf("unknown MemberMessage")
	}
}

func handleSubscribe(msg *chatgroup.Message, addr net.Addr) error {

	log.Printf("handleSubscribe: %v\n", msg.Sender)
	msg.Sender.Ip = strings.Split(addr.String(), ":")[0]
	log.Printf("handleSubscribe: %v\n", msg.Sender)

	// Check subscriber for uniqueness
	log.Printf("Check subscriber for uniqueness: %v\n", msg.Sender)
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
	log.Printf("memberlist: %v", cgMember)

	log.Printf("Current members registered: %v\n", cgMember)

	// Send message to other subscribers via gRPC Displayer service
	for _, recipient := range cgMember {

		msg.MsgType = chatgroup.Message_DISPLAY_SUBSCRIPTION

		if recipient.Name != msg.Sender.Name && recipient.Name != memberName {
			log.Printf("From %s to %s (%s:%s): %q\n", msg.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, msg.Sender)

			conn, err := net.Dial("tcp", recipient.Ip+":"+recipient.Port)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Dial to displaying server on %q\n", recipient.Ip+":"+recipient.Port)

			// Marshal into binary format
			byteArray, err := proto.Marshal(msg)
			if err != nil {
				fmt.Errorf("could not encode new member: %v", err)
				os.Exit(1)
			}

			n, err := conn.Write(byteArray)
			log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, msg)

			//conn.Read(byteArray)
			//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

			// Receive reply
			conn.Close()
		}
	}

	// Append text message in "messages" view
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

		if recipient.Name != msg.Sender.Name && recipient.Name != memberName  {
			log.Printf("From %s to %s (%s:%s): %q\n", msg.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, msg.Sender)

			conn, err := net.Dial("tcp", recipient.Ip+":"+recipient.Port)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Dial to displaying server on %q\n", recipient.Ip+":"+recipient.Port)

			// Marshal into binary format
			byteArray, err := proto.Marshal(msg)
			if err != nil {
				fmt.Errorf("could not encode new member: %v", err)
				os.Exit(1)
			}

			n, err := conn.Write(byteArray)
			log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, msg)

			//conn.Read(byteArray)
			//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

			// Receive reply
			conn.Close()
		}
	}



	displayText(fmt.Sprintf("<%s has left>", msg.Sender.Name))

	return nil
}
func handlePublish(msg *chatgroup.Message, addr net.Addr) error {

	log.Printf("handleSubscribe: %v\n", msg.Sender)
	msg.Sender.Ip = strings.Split(addr.String(), ":")[0]
	log.Printf("handleSubscribe: %v\n", msg.Sender)

	log.Printf("Publish from %v: %q\n", msg.Sender.Name, msg.Text)

	// Send message to other subscribers via gRPC Displayer service
	for _, recipient := range cgMember {

		msg.MsgType = chatgroup.Message_DISPLAY_TEXT

		if recipient.Name != msg.Sender.Name {
			log.Printf("From %s to %s (%s:%s): %q\n", msg.Sender.Name, recipient.Name, recipient.Ip, recipient.Port, msg.Text)

			conn, err := net.Dial("tcp", recipient.Ip+":"+recipient.Port)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Dial to Rootserver on 127.0.0.1:22365\n")

			// Marshal into binary format
			byteArray, err := proto.Marshal(msg)
			if err != nil {
				fmt.Errorf("could not encode new member: %v", err)
				os.Exit(1)
			}

			n, err := conn.Write(byteArray)
			log.Printf("Message (%v byte) sent (%v byte): %v\n", len(byteArray), n, msg)

			//conn.Read(byteArray)
			//fmt.Printf("New member (%v byte) red: %v\n", len(byteArray), byteArray)

			// Receive reply
			conn.Close()
		}
	}

	return nil
}
