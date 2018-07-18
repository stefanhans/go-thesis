package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/subscriber"
	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/subscriber-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	subscriberlist subscribergroup.SubscriberList
)

func startServerServer(ip string, port string) error {

	// Create and register server
	srv := grpc.NewServer()

	// Register server for subscribergroup
	var subscribers subscriberServer
	subscribergroup.RegisterSubscribersServer(srv, subscribers)

	// Create listener
	l, err := net.Listen("tcp", ":"+port)
	fmt.Printf("subscriber-group server does listen on %s:%s\n", ip, port)
	if err != nil {
		return fmt.Errorf("could not listen to %s:%s: \v", ip, port, err)
	}
	// Serve via listener
	return srv.Serve(l)
}

func startClientServer(name string, ip string, port string) error {

	// Create and register server
	srv := grpc.NewServer()

	// Register server for messages
	var messageServer displayServer
	subscriber.RegisterDisplayMessagesServer(srv, messageServer)

	// Create listener
	l, err := net.Listen("tcp", ":"+port)
	fmt.Printf("Subscriber %q does listen on %s:%s\n", name, ip, port)
	if err != nil {
		return fmt.Errorf("%q could not listen to %s:%s: %v\n", name, ip, port, err)
	}
	// Serve via listener
	return srv.Serve(l)
}

func subscribeClient(name string, ip string, port string) error {

	client, err := dialServer()
	if err != nil {
		return err
	}

	return subscribe(context.Background(), client, name, ip, port)
}

func unsubscribeClient(name string) error {

	client, err := dialServer()
	if err != nil {
		return err
	}

	return unsubscribe(context.Background(), client, name)
}

func sendMessage(name string, text ...string) error {

	client, err := dialServer()
	if err != nil {
		return err
	}
	sender := &subscribergroup.Subscriber{Name: name}

	return send(context.Background(), client, sender, strings.Join(text[:], " "))

}

func dialServer() (subscribergroup.SubscribersClient, error) {

	// Create client with insecure connection
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to backend: %v", err)
	}
	return subscribergroup.NewSubscribersClient(conn), nil
}

// Receiver for implementing the server service interface Subscribers
type subscriberServer struct{}

// Server's Subscribe implementation
func (s subscriberServer) Subscribe(ctx context.Context, subscr *subscribergroup.Subscriber) (*subscribergroup.Subscriber, error) {
	fmt.Printf("SUBSCRIBE: %v\n", subscr)
	subscriberlist.Subscriber = append(subscriberlist.Subscriber, subscr)

	for _, recipient := range subscriberlist.Subscriber {
		//fmt.Printf("Check recipient: %v\n", recipient)
		if recipient.Name != subscr.Name {
			//fmt.Printf("From %s to %s (%s:%s)\n", subscr.Name, recipient.Name, recipient.Ip, recipient.Port)

			// Create client with insecure connection
			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			//fmt.Printf("Dial to localhost:%s\n", recipient.Port)
			client := subscriber.NewDisplayMessagesClient(conn)

			// Write to gRPC client
			_, err = client.DisplaySubscribe(ctx, &subscriber.Sender{Name: subscr.Name, Ip: subscr.Ip, Port: subscr.Port})
			if err != nil {
				return nil, fmt.Errorf("could not display subscription: %v", err)
			}
		}
	}

	return subscr, nil
}

// Server's Subscribe implementation
func (s subscriberServer) Unsubscribe(ctx context.Context, subscr *subscribergroup.Subscriber) (*subscribergroup.Subscriber, error) {
	fmt.Printf("UNSUBSCRIBE: %v\n", subscr)

	for i, s := range subscriberlist.Subscriber {
		if s.Name == subscr.Name {
			subscriberlist.Subscriber = append(subscriberlist.Subscriber[:i], subscriberlist.Subscriber[i+1:]...)
			break
		}
	}

	for _, recipient := range subscriberlist.Subscriber {
		// Create client with insecure connection
		conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		//fmt.Printf("Dial to localhost:%s\n", recipient.Port)
		client := subscriber.NewDisplayMessagesClient(conn)

		_, err = client.DisplayUnsubscribe(ctx, &subscriber.Sender{Name: subscr.Name})
		if err != nil {
			return nil, fmt.Errorf("could not display tweet: %v", err)
		}
	}
	return subscr, nil
}

// Server's List implementation
func (s subscriberServer) List(ctx context.Context, void *subscribergroup.Void) (*subscribergroup.SubscriberList, error) {
	fmt.Printf("LIST: %v\n", subscriberlist)

	return &subscriberlist, nil
}

// Server's Send implementation
func (s subscriberServer) Send(ctx context.Context, tweet *subscribergroup.Tweet) (*subscribergroup.SubscriberList, error) {
	fmt.Printf("SEND: %v\n", tweet)
	sender := tweet.Sender

	for _, recipient := range subscriberlist.Subscriber {
		//fmt.Printf("Check recipient: %v\n", recipient)
		if recipient.Name != sender.Name {
			fmt.Printf("From %s to %s (%s:%s): %q\n", sender.Name, recipient.Name, recipient.Ip, recipient.Port, tweet.Text)

			// Create client with insecure connection
			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			//fmt.Printf("Dial to localhost:%s\n", recipient.Port)
			client := subscriber.NewDisplayMessagesClient(conn)

			source := subscriber.Sender{Name: sender.Name, Ip: sender.Ip, Port: sender.Port}
			tweet := subscriber.Message{Sender: &source, Text: tweet.Text}

			_, err = client.DisplayText(ctx, &tweet)
			if err != nil {
				return nil, fmt.Errorf("could not display tweet: %v", err)
			}
		}
	}
	return &subscriberlist, nil
}
