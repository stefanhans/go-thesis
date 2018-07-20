package main

import (
	"fmt"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/chat-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Start displayer service to provide displaying messages in the text-based UI
func startDisplayer(name string, ip string, port string) error {

	// Create listener
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("%q could not listen to %s:%s: %v\n", name, ip, port, err)
	}

	// Create gRPC server
	srv := grpc.NewServer()

	// Register displayer
	var displayer displayServer
	chatgroup.RegisterDisplayerServer(srv, displayer)

	// Start gRPC server
	go func() {
		srv.Serve(l)
	}()

	return nil
}

// Receiver to implement the displayer service interface DisplayerServer
type displayServer struct{}

// DisplayerServer's DisplayText implementation
func (ds displayServer) DisplayText(ctx context.Context, message *chatgroup.Message) (*chatgroup.Message, error) {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s: %s", message.Sender.Name, message.Text))

	return message, nil
}

// DisplayerServer's DisplaySubscription implementation
func (ds displayServer) DisplaySubscription(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {

	// Append subscription message in "messages" view
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", subscr.Name, subscr.Ip, subscr.Port))

	return subscr, nil
}

// DisplayerServer's DisplayUnsubscription implementation
func (ds displayServer) DisplayUnsubscription(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {

	// Append unsubscription message in "messages" view
	displayText(fmt.Sprintf("<%s has left>", subscr.Name))

	return subscr, nil
}
