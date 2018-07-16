package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-publishserver/subscriber-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Check command args
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: subscribe, send, or list")
		os.Exit(1)
	}

	// Create client with insecure connection
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to backend: %v", err)
	}
	client := subscribergroup.NewSubscribersClient(conn)

	// Switch subcommands and call wrapper function
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "subscribe":
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: subscribe <name> <ip> <port>")
			os.Exit(1)
		}
		err = subscribe(context.Background(), client, flag.Arg(1), flag.Arg(2), flag.Arg(3))
	case "send":
		if flag.NArg() < 4 {
			fmt.Fprintln(os.Stderr, "missing parameter: send <name> <ip> <port> <text>...")
			os.Exit(1)
		}
		sender := &subscribergroup.Subscriber{Name: flag.Arg(1), Ip: flag.Arg(2), Port: flag.Arg(3)}
		err = send(context.Background(), client, sender, strings.Join(flag.Args()[4:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Subscribe wrapper function
func subscribe(ctx context.Context, client subscribergroup.SubscribersClient, name string, ip string, port string) error {

	// Write to gRPC client
	_, err := client.Subscribe(ctx, &subscribergroup.Subscriber{Name: name, Ip: ip, Port: port, Leader: false})
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
		fmt.Printf("%s %s %s %v\n", t.Name, t.Ip, t.Port, t.Leader)
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
		fmt.Printf("Sent receipt: %s %s %s %v\n", t.Name, t.Ip, t.Port, t.Leader)
	}
	return nil
}
