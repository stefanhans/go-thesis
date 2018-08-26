package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.4./rudimentary-chat-tcp/chat-group"
	"github.com/golang/protobuf/proto"
	"syscall"
	"strings"
)


// Start publisher service to provide member registration and message publishing
func startLeaderService() error {

	// Create publishingListener
	publishingListener, err := net.Listen("tcp", leaderService)

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
		log.Fatalf("could not listen to %q: %v\n", leaderService, err)
	}
	defer publishingListener.Close()

	log.Printf("Started publishing service listening on %q\n", leaderService)

	// Append text messages in "messages" view of publisher
	displayText(fmt.Sprintf("<publishing service running: %s (%s:%s)>", selfMember.Name, serverIp, serverPort))

	// Subscribe directly at started publishing service
	selfMember.Leader = true
	selfMemberList = append(selfMemberList, selfMember)
	log.Printf("Subscribed directly at started publishing service: %v\n", selfMemberList[0])

	// Append text messages in "messages" view of publisher
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", selfMember.Name, selfMember.Ip, selfMember.Port))

	// Endless loop in foreground of goroutine
	for {
		// Wait for a connection.
		conn, err := publishingListener.Accept()
		if err != nil {
			log.Printf("failed to accept connection from publishing service publishingListener: %s\n", err)
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

	// Read all data from the connection
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

	// Unmarshall message
	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Errorf("could not unmarshall message: %v", err)
	}

	// Fetch the handler from a map by the message type and call it accordingly
	if requestAction, ok := actionMap[msg.MsgType]; ok {
		log.Printf("%v\n", msg)
		err := requestAction(&msg, addr)
		if err != nil {
			fmt.Printf("could not handle %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		log.Printf("publisher: unknown message type %v\n", msg.MsgType)
	}
}

func RequestLeaderlist() error {

	newMember := &chatgroup.Message{
		MsgType: chatgroup.Message_LEADERLIST_REQUEST,
		Sender:  selfMember}

	return sendMessage(newMember, leaderService)
}

func Subscribe() error {

	newMember := &chatgroup.Message{
		MsgType: chatgroup.Message_SUBSCRIBE_REQUEST,
		Sender:  selfMember}

	return sendPublisherRequest(newMember)
}

func Unsubscribe(memberName string) error {

	leavingMember := &chatgroup.Message{
		MsgType: chatgroup.Message_UNSUBSCRIBE_REQUEST,
		Sender: &chatgroup.Member{
			Name: memberName}}

	return sendPublisherRequest(leavingMember)
}

func Publish(text string) error {

	message := &chatgroup.Message{
		MsgType: chatgroup.Message_PUBLISH_REQUEST,
		Sender:  selfMember,
		Text:    text}

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s: %s", selfMember.Name, message.Text))

	return sendPublisherRequest(message)
}

// Dial publisher and return connection
func sendPublisherRequest(message *chatgroup.Message) error {

	// Connect to publishing service
	conn, err := net.Dial("tcp", leaderService)
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

	// Close connection
	return conn.Close()
}
