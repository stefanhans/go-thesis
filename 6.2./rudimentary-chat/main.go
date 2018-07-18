package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

		err = startServerServer(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal("startServerServer: %v\n", err)
		}

	case "client":
		// Check command args
		flag.Parse()
		if flag.NArg() < 4 {
			fmt.Fprintln(os.Stderr, "missing parameter: client <name> <ip> <port>")
			os.Exit(1)
		}

		err = subscribeClient(flag.Arg(1), flag.Arg(2), flag.Arg(3))
		if err != nil {
			log.Fatal("subscribeClient: %v\n", err)
		}
		fmt.Printf("Client %q (%s:%s) has subscribed\n", flag.Arg(1), flag.Arg(2), flag.Arg(3))

		err = startClientServer(flag.Arg(1), flag.Arg(2), flag.Arg(3))
		if err != nil {
			log.Fatal("startClientServer: %v\n", err)
		}

	case "list":
		// Create client with insecure connection
		conn, err := grpc.Dial(":8888", grpc.WithInsecure())
		if err != nil {
			log.Fatal("could not connect to backend: %v", err)
		}
		client := subscribergroup.NewSubscribersClient(conn)

		err = list(context.Background(), client)

	case "subscribe":

		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: subscribe <name> <ip> <port>")
			os.Exit(2)
		}

		err = subscribeClient(flag.Arg(1), flag.Arg(2), flag.Arg(3))
		if err != nil {
			log.Fatal("subscribeClient: %v\n", err)
		}

	case "unsubscribe":

		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "missing parameter: unsubscribe <name>")
			os.Exit(2)
		}

		err = unsubscribeClient(flag.Arg(1))
		if err != nil {
			log.Fatal("unsubscribeClient: %v\n", err)
		}

	case "send":

		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "missing parameter: send <name> <text>...")
			os.Exit(1)
		}

		err = sendMessage(flag.Arg(1), strings.Join(flag.Args()[2:], " "))
		if err != nil {
			log.Fatal("sendMessage: %v\n", err)
		}

	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
		if err != nil {
			log.Fatal("%v\n", err)
		}
	}
}
