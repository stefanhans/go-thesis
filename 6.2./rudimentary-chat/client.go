package main

import (
	"fmt"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/subscriber"
	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/subscriber-group"
	"golang.org/x/net/context"
)

type displayServer struct{}

func (ds displayServer) DisplayText(ctx context.Context, message *subscriber.Message) (*subscriber.Message, error) {

	fmt.Printf("%s (%s:%s): %q\n", message.Sender.Name, message.Sender.Ip, message.Sender.Port, message.Text)

	return message, nil
}

func (ds displayServer) DisplaySubscribe(ctx context.Context, subscr *subscriber.Sender) (*subscriber.Sender, error) {

	fmt.Printf("<%s (%s:%s) has joined>\n", subscr.Name, subscr.Ip, subscr.Port)

	return subscr, nil
}

func (ds displayServer) DisplayUnsubscribe(ctx context.Context, subscr *subscriber.Sender) (*subscriber.Sender, error) {

	fmt.Printf("<%s has left>\n", subscr.Name)

	return subscr, nil
}


// Subscribe wrapper function
func subscribe(ctx context.Context, client subscribergroup.SubscribersClient, name string, ip string, port string) error {

	// Write to gRPC client
	_, err := client.Subscribe(ctx, &subscribergroup.Subscriber{Name: name, Ip: ip, Port: port})
	if err != nil {
		return fmt.Errorf("could not add member in the membergroup: %v", err)
	}
	return nil
}

// Unsubscribe wrapper function
func unsubscribe(ctx context.Context, client subscribergroup.SubscribersClient, name string) error {

	// Write to gRPC client
	_, err := client.Unsubscribe(ctx, &subscribergroup.Subscriber{Name: name})
	if err != nil {
		return fmt.Errorf("could not add member in the membergroup: %v", err)
	}
	return nil
}

// List wrapper function
func list(ctx context.Context, client subscribergroup.SubscribersClient) error {

	// List from gRPC client
	l, err := client.List(ctx, &subscribergroup.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch membergroup: %v", err)
	}

	// Print members
	for _, t := range l.Subscriber {
		fmt.Printf("%s %s %s\n", t.Name, t.Ip, t.Port)
	}
	return nil
}

// Send wrapper function
func send(ctx context.Context, client subscribergroup.SubscribersClient, sender *subscribergroup.Subscriber, text string) error {

	msg := subscribergroup.Tweet{Sender: sender, Text: text}

	// List from gRPC client
	l, err := client.Send(ctx, &msg)
	if err != nil {
		return fmt.Errorf("could not send to subscribergroup: %v", err)
	}
	fmt.Printf("Sent to subscribergroup from %v: %q\n", msg.Sender, msg.Text)

	// Print members
	for _, t := range l.Subscriber {
		fmt.Printf("Sent receipt: %s %s %s\n", t.Name, t.Ip, t.Port)
	}
	return nil
}

