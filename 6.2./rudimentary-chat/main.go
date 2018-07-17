package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/subscriber"
	"bitbucket.org/stefanhans/go-thesis/6.2./rudimentary-chat/subscriber-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	// Check command args
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: server or client")
		os.Exit(1)
	}

	// Switch subcommands and call wrapper function
	var err error
	switch cmd := flag.Arg(0); cmd {
	case "server":

		// Check command args
		flag.Parse()
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: <ip> <port>")
			os.Exit(1)
		}

		// Create and register server
		srv := grpc.NewServer()

		// Register server for subscribergroup
		var subscribers subscriberServer
		subscribergroup.RegisterSubscribersServer(srv, subscribers)

		// Create listener
		l, err := net.Listen("tcp", ":"+flag.Arg(2))
		fmt.Printf("subscriber-group server does listen on %s:%s\n", flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal("could not listen to %s:%s: \v", flag.Arg(1), flag.Arg(2), err)
		}
		// Serve via listener
		log.Fatal(srv.Serve(l))

	case "client":
		// Check command args
		flag.Parse()
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: client <ip> <port>")
			os.Exit(1)
		}

		// Create and register server
		srv := grpc.NewServer()

		// Register server for tweets
		var tweets displayServer
		subscriber.RegisterDisplayMessagesServer(srv, tweets)

		// Create listener
		l, err := net.Listen("tcp", ":"+flag.Arg(2))
		fmt.Printf("Subscriber does listen on %s:%s\n", flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal("could not listen to %s:%s: \v", flag.Arg(1), flag.Arg(2), err)
		}
		// Serve via listener
		log.Fatal(srv.Serve(l))

	case "list":
		// Create client with insecure connection
		conn, err := grpc.Dial(":8888", grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		client := subscribergroup.NewSubscribersClient(conn)

		err = list(context.Background(), client)

	case "subscribe":
		// Create client with insecure connection
		conn, err := grpc.Dial(":8888", grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		client := subscribergroup.NewSubscribersClient(conn)

		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: subscribe <name> <ip> <port>")
			os.Exit(2)
		}
		err = subscribe(context.Background(), client, flag.Arg(1), flag.Arg(2), flag.Arg(3))


	case "unsubscribe":
		// Create client with insecure connection
		conn, err := grpc.Dial(":8888", grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		client := subscribergroup.NewSubscribersClient(conn)

		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "missing parameter: unsubscribe <name>")
			os.Exit(2)
		}
		err = unsubscribe(context.Background(), client, flag.Arg(1))

	case "send":
		// Create client with insecure connection
		conn, err := grpc.Dial(":8888", grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		client := subscribergroup.NewSubscribersClient(conn)

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
