package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.4./rudimentary-chat-tcp/chat-group"
	"github.com/golang/protobuf/proto"
)

// Start displayer service to provide displaying messages in the text-based UI
func startServer() error {

	// Create displayingListener
	displayingListener, err := net.Listen("tcp", displayingService)

	if err != nil {
		log.Fatalf("could not listen to %q: %v\n", displayingService, err)
	}
	defer displayingListener.Close()

	log.Printf("Started displaying service listening on %q\n", displayingService)

	for {
		// Wait for a connection.
		conn, err := displayingListener.Accept()
		if err != nil {
			continue //log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleRequest(conn)
	}

	return nil
}

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func handleRequest(conn net.Conn) {

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

	log.Printf("Displayer received %v bytes\n", len(data))

	// Unmarshall message
	var msg chatgroup.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Errorf("could not unmarshall message: %v", err)
	}

	log.Printf("msg from %v: %v\n", addr, msg)

	// Fetch the handler from a map by the message type and call it accordingly
	if replyAction, ok := actionMap[msg.MsgType]; ok {
		log.Printf("%v\n", msg)
		err := replyAction(&msg, nil)
		if err != nil {
			fmt.Printf("could not handle %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		log.Printf("server: unknown message type %v\n", msg.MsgType)
	}
}