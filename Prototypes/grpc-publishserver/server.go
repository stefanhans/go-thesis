package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-publishserver/member-group"
	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-publishserver/subscriber-group"
	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-subscriber/display"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
)

var (
	memberlist     membergroup.MemberList
	subscriberlist subscribergroup.SubscriberList
)

func main() {
	// Create and register server
	srv := grpc.NewServer()

	// Register server for membergroup
	var members memberServer
	membergroup.RegisterMembersServer(srv, members)

	// Register server for subscribergroup
	var subscribers subscriberServer
	subscribergroup.RegisterSubscribersServer(srv, subscribers)

	// Create listener
	l, err := net.Listen("tcp", ":8888")
	fmt.Printf("Publishserver does listen on localhost:8888\n")
	if err != nil {
		log.Fatal("could not listen to :8888: \v", err)
	}
	// Serve via listener
	log.Fatal(srv.Serve(l))
}

// ********* MEMBERS *********

// Receiver for implementing the server service interface Members
type memberServer struct{}

// Server's Subscribe implementation
func (s memberServer) Register(ctx context.Context, info *membergroup.Member) (*membergroup.Member, error) {
	memberlist.Member = append(memberlist.Member, info)
	return info, nil
}

// Server's List implementation
func (s memberServer) List(ctx context.Context, void *membergroup.Void) (*membergroup.MemberList, error) {
	return &memberlist, nil
}

// ********* SUBSCRIBERS *********

// Receiver for implementing the server service interface Subscribers
type subscriberServer struct{}

// Server's Subscribe implementation
func (s subscriberServer) Subscribe(ctx context.Context, subscriber *subscribergroup.Subscriber) (*subscribergroup.Subscriber, error) {
	subscriberlist.Subscriber = append(subscriberlist.Subscriber, subscriber)
	return subscriber, nil
}

// Server's List implementation
func (s subscriberServer) List(ctx context.Context, void *subscribergroup.Void) (*subscribergroup.SubscriberList, error) {
	return &subscriberlist, nil
}

// Server's Send implementation
func (s subscriberServer) Send(ctx context.Context, tweet *subscribergroup.Tweet) (*subscribergroup.SubscriberList, error) {
	fmt.Printf("Send request received from %v\n", tweet.Sender)
	sender := tweet.Sender

	for _, recipient := range subscriberlist.Subscriber {
		fmt.Printf("Check recipient: %v\n", recipient)
		if recipient.Name != sender.Name {
			fmt.Printf("Send: %q from sender %s to recipient %s (%s:%s)\n", tweet.Text, sender.Name, recipient.Name, recipient.Ip, recipient.Port)

			// Create client with insecure connection
			conn, err := grpc.Dial(":"+recipient.Port, grpc.WithInsecure())
			if err != nil {
				log.Fatal("could not connect to backend: %v", err)
			}
			fmt.Printf("Dial to localhost:%s\n", recipient.Port)
			client := display.NewDisplayTweetsClient(conn)

			source := display.Sender{Name: sender.Name, Ip: sender.Ip, Port: sender.Port}

			tweet := display.Tweet{Sender: &source, Text: tweet.Text}

			err = show(context.Background(), client, &tweet)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			//// create TCP connection to recipient via 'Ip:Port'
			//conn, err := net.Dial("tcp", recipient.Ip+":"+recipient.Port)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Printf("Send: connection to %s:%s dialed\n", recipient.Ip, recipient.Port)
			//
			//// send message
			//fmt.Fprintf(conn, tweet.Text)
			//fmt.Printf("Send: message sent\n")
			//// receive and print reply
			//reply, err := bufio.NewReader(conn).ReadString('\n')
			//fmt.Printf("Reply: %q", reply)
			//// close connection
			//conn.Close()
		}
	}
	return &subscriberlist, nil
}

// Tweet wrapper function
func show(ctx context.Context, client display.DisplayTweetsClient, tweet *display.Tweet) error {

	// Write to gRPC client
	_, err := client.Display(ctx, tweet)
	if err != nil {
		return fmt.Errorf("could not display tweet: %v", err)
	}
	return nil
}
