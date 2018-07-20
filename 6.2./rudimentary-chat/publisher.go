package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"syscall"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/chat-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	memberlist chatgroup.MemberList
)

// Start publisher service to provide member registration and message publishing
func startPublisher(ip string, port string) error {

	// Create listener
	l, err := net.Listen("tcp", ":"+port)

	// Exit on unexpected error
	if err != nil && !strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {
		log.Fatalf("could not listen to %s:%s: %v\n", ip, port, err)
	}

	// Do not start publisher on "address already in use"
	if err != nil {
		return err
	}

	// Create gRPC server
	srv := grpc.NewServer()

	// Register publisher
	var publisher publishServer
	chatgroup.RegisterPublisherServer(srv, publisher)

	// Start gRPC server
	go func() {
		srv.Serve(l)
	}()
	return nil
}

// Receiver to implement the publisher service interface PublisherServer
type publishServer struct{}

// PublisherServer's Subscribe implementation
func (ps publishServer) Subscribe(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {
	log.Printf("SUBSCRIBE: %v\n", subscr)

	// Check subscriber for uniqueness
	for _, recipient := range memberlist.Member {
		if recipient.Name == subscr.Name {
			return nil, fmt.Errorf("name %q already used", subscr.Name)
		}
		if recipient.Ip == subscr.Ip && recipient.Port == subscr.Port {
			return nil, fmt.Errorf("address %s:%s already used by %s", recipient.Ip, recipient.Port, recipient.Name)
		}
	}

	// Add subscriber
	memberlist.Member = append(memberlist.Member, subscr)

	// Inform other subscribers via gRPC Displayer service
	for _, recipient := range memberlist.Member {

		if recipient.Name != subscr.Name {

			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			client := chatgroup.NewDisplayerClient(conn)

			_, err = client.DisplaySubscription(ctx, &chatgroup.Member{Name: subscr.Name, Ip: subscr.Ip, Port: subscr.Port})
			if err != nil {
				return nil, fmt.Errorf("could not display subscription: %v", err)
			}
		}
	}
	return subscr, nil
}

// PublisherServer's Unsubscribe implementation
func (ps publishServer) Unsubscribe(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {
	log.Printf("UNSUBSCRIBE: %v\n", subscr)

	// Remove subscriber
	for i, s := range memberlist.Member {
		if s.Name == subscr.Name {
			memberlist.Member = append(memberlist.Member[:i], memberlist.Member[i+1:]...)
			break
		}
	}

	// Inform other subscribers via gRPC Displayer service
	for _, recipient := range memberlist.Member {

		conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		client := chatgroup.NewDisplayerClient(conn)

		_, err = client.DisplayUnsubscription(ctx, &chatgroup.Member{Name: subscr.Name})
		if err != nil {
			return nil, fmt.Errorf("could not display message: %v", err)
		}
	}
	return subscr, nil
}

// PublisherServer's Send implementation
func (ps publishServer) Publish(ctx context.Context, message *chatgroup.Message) (*chatgroup.MemberList, error) {
	log.Printf("PUBLISH: %v\n", message)

	sender := message.Sender

	// Send message to other subscribers via gRPC Displayer service
	for _, recipient := range memberlist.Member {

		if recipient.Name != sender.Name {
			log.Printf("From %s to %s (%s:%s): %q\n", sender.Name, recipient.Name, recipient.Ip, recipient.Port, message.Text)

			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			client := chatgroup.NewDisplayerClient(conn)

			sender := chatgroup.Member{Name: sender.Name, Ip: sender.Ip, Port: sender.Port}
			message := chatgroup.Message{Sender: &sender, Text: message.Text}

			_, err = client.DisplayText(ctx, &message)
			if err != nil {
				return nil, fmt.Errorf("could not display message: %v", err)
			}
		}
	}
	return &memberlist, nil
}
