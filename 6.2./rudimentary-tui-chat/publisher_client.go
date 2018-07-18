package main

import (
	"fmt"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-tui-chat/chat-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func listMembers() error {

	client, err := dialPublisher()
	if err != nil {
		return err
	}

	// List from gRPC client
	l, err := client.ListSubscriber(context.Background(), &chatgroup.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch membergroup: %v", err)
	}

	// Print members
	for _, t := range l.Member {
		fmt.Printf("%s %s %s\n", t.Name, t.Ip, t.Port)
	}

	return nil
}

func subscribeClient(name string, ip string, port string) error {

	client, err := dialPublisher()
	if err != nil {
		return err
	}

	// Write to gRPC client
	_, err = client.Subscribe(context.Background(), &chatgroup.Member{Name: name, Ip: ip, Port: port})
	if err != nil {
		return fmt.Errorf("could not add member in the membergroup: %v", err)
	}
	return nil
}

func unsubscribeClient(name string) error {

	client, err := dialPublisher()
	if err != nil {
		return err
	}

	// Write to gRPC client
	_, err = client.Unsubscribe(context.Background(), &chatgroup.Member{Name: name})
	if err != nil {
		return fmt.Errorf("could not add member in the membergroup: %v", err)
	}
	return nil
}

func sendMessage(name string, text ...string) error {

	client, err := dialPublisher()
	if err != nil {
		return err
	}

	msg := chatgroup.Message{Sender: &chatgroup.Member{Name: name}, Text: strings.Join(text[:], " ")}

	// List from gRPC client
	_, err = client.Publish(context.Background(), &msg)
	if err != nil {
		return fmt.Errorf("could not send to subscribergroup: %v", err)
	}
	//fmt.Printf("Sent to subscribergroup from %v: %q\n", msg.Sender, msg.Text)

	// Print members
	//for _, t := range l.Member {
	//	fmt.Printf("Sent receipt: %s %s %s\n", t.Name, t.Ip, t.Port)
	//}
	return nil
}

func dialPublisher() (chatgroup.PublisherClient, error) {

	// Create client with insecure connection
	conn, err := grpc.Dial(":"+serverPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to backend: %v", err)
	}
	return chatgroup.NewPublisherClient(conn), nil
}
