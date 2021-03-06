package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-tui-chat/chat-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strings"
	"syscall"
)

var (
	memberlist chatgroup.MemberList
)

// Start publisher service to provide member registration and message publishing
func startPublisher(ip string, port string, foreground bool) error {

	// Create listener
	l, err := net.Listen("tcp", ":"+port)

	// Exit application on unexpected error
	if err != nil && !strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {
		log.Fatalf("could not listen to %s:%s: %v\n", ip, port, err)
	}

	// Do not start publisher on "address already in use"
	if err != nil {
		return err
	}

	// Create gRPC server
	srv := grpc.NewServer()

	// Register Publisher
	var publisher publishServer
	chatgroup.RegisterPublisherServer(srv, publisher)

	// Start gRPC server in fore- or background respectively
	if foreground {
		return srv.Serve(l)
	} else {
		go func() {
			srv.Serve(l)
		}()
		return nil
	}
}

// Receiver for implementing the server service interface Subscribers
type publishServer struct{}

// Server's Subscribe implementation
func (s publishServer) Subscribe(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {
	if isServer {
		fmt.Printf("SUBSCRIBE: %v\n", subscr)
	}
	//todo: check for uniquness of member

	for _, recipient := range memberlist.Member {
		if recipient.Name == subscr.Name {
			return nil, fmt.Errorf("name %q already subscribed", subscr.Name)
		}
		if recipient.Ip == subscr.Ip && recipient.Port == subscr.Port {
			return nil, fmt.Errorf("adress %s:%s already used by %s", recipient.Ip, recipient.Port, recipient.Name)
		}
	}

	memberlist.Member = append(memberlist.Member, subscr)

	for _, recipient := range memberlist.Member {
		//fmt.Printf("Check recipient: %v\n", recipient)
		if recipient.Name != subscr.Name {
			//fmt.Printf("From %s to %s (%s:%s)\n", subscr.Name, recipient.Name, recipient.Ip, recipient.Port)

			// Create client with insecure connection
			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			//fmt.Printf("Dial to localhost:%s\n", recipient.Port)
			client := chatgroup.NewDisplayerClient(conn)

			// Write to gRPC client
			_, err = client.DisplaySubscription(ctx, &chatgroup.Member{Name: subscr.Name, Ip: subscr.Ip, Port: subscr.Port})
			if err != nil {
				return nil, fmt.Errorf("could not display subscription: %v", err)
			}
		}
	}

	return subscr, nil
}

// Server's Subscribe implementation
func (s publishServer) Unsubscribe(ctx context.Context, subscr *chatgroup.Member) (*chatgroup.Member, error) {
	if isServer {
		fmt.Printf("UNSUBSCRIBE: %v\n", subscr)
	}

	for i, s := range memberlist.Member {
		if s.Name == subscr.Name {
			memberlist.Member = append(memberlist.Member[:i], memberlist.Member[i+1:]...)
			break
		}
	}

	for _, recipient := range memberlist.Member {
		// Create client with insecure connection
		conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		//fmt.Printf("Dial to localhost:%s\n", recipient.Port)
		client := chatgroup.NewDisplayerClient(conn)

		_, err = client.DisplayUnsubscription(ctx, &chatgroup.Member{Name: subscr.Name})
		if err != nil {
			return nil, fmt.Errorf("could not display message: %v", err)
		}
	}
	return subscr, nil
}

// Server's List implementation
func (s publishServer) ListSubscriber(ctx context.Context, void *chatgroup.Void) (*chatgroup.MemberList, error) {
	if isServer {
		fmt.Printf("LIST: %v\n", memberlist)
	}

	return &memberlist, nil
}

// Server's Send implementation
func (s publishServer) Publish(ctx context.Context, message *chatgroup.Message) (*chatgroup.MemberList, error) {
	if isServer {
		fmt.Printf("PUBLISH: %v\n", message)
	}
	sender := message.Sender

	for _, recipient := range memberlist.Member {
		//fmt.Printf("Check recipient: %v\n", recipient)
		if recipient.Name != sender.Name {
			if isServer {
				fmt.Printf("From %s to %s (%s:%s): %q\n", sender.Name, recipient.Name, recipient.Ip, recipient.Port, message.Text)
			}

			// Create client with insecure connection
			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			//fmt.Printf("Dial to localhost:%s\n", recipient.Port)
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
