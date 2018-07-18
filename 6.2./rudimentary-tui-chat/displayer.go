package main

import (
	"fmt"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-tui-chat/chat-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func startDisplayer(name string, ip string, port string) error {

	// Create and register server
	srv := grpc.NewServer()

	// Register server for messages
	var messageServer displayServer
	chatgroup.RegisterDisplayerServer(srv, messageServer)

	// Create listener
	l, err := net.Listen("tcp", ":"+port)
	fmt.Printf("Subscriber %q does listen on %s:%s\n", name, ip, port)
	if err != nil {
		return fmt.Errorf("%q could not listen to %s:%s: %v\n", name, ip, port, err)
	}

	// todo: handle possible error
	// Serve via listener
	go func() {
		srv.Serve(l)
	}()

	return nil
}

type displayServer struct{}

func (ds displayServer) DisplayText(ctx context.Context, message *chatgroup.Message) (*chatgroup.Message, error) {

	displayText(fmt.Sprintf("%s: %s", message.Sender.Name, message.Text))

	return message, nil
}

func (ds displayServer) DisplaySubscription(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {

	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", subscr.Name, subscr.Ip, subscr.Port))

	return subscr, nil
}

func (ds displayServer) DisplayUnsubscription(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {

	displayText(fmt.Sprintf("<%s has left>", subscr.Name))

	return subscr, nil
}
