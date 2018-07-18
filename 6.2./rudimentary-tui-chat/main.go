package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	serverIp	string = "localhost"
	serverPort	string = "22365"
)

var (
	memberName string
	memberIp   string
	memberPort string
	isServer	bool
)

func main() {

	// Check command args
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: server, client, list, subscribe, unsubscribe, or send")
		os.Exit(1)
	}

	// Switch subcommands and call wrapper function
	var err error
	switch cmd := flag.Arg(0); cmd {
	case "server":

		// Check command args
		flag.Parse()
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: server <ip> <port>")
			os.Exit(1)
		}
		isServer = true

		startPublisher(flag.Arg(1), flag.Arg(2), isServer)

	case "client":
		// Check command args
		flag.Parse()
		if flag.NArg() < 4 {
			fmt.Fprintln(os.Stderr, "missing parameter: client <name> <ip> <port>")
			os.Exit(1)
		}
		memberName = flag.Arg(1)
		memberIp = flag.Arg(2)
		memberPort = flag.Arg(3)
		isServer = false

		startPublisher(serverIp, serverPort, isServer)


		err = subscribeClient(memberName, memberIp, memberPort)
		if err != nil {
			log.Fatalf("subscribeClient: %v", err)
		}
		fmt.Printf("Client %q (%s:%s) has subscribed\n", memberName, memberIp, memberPort)

		err = startDisplayer(memberName, memberIp, memberPort)
		if err != nil {
			log.Fatalf("startDisplayer: %v", err)
		}

		err = startTui()
		if err != nil {
			log.Fatalf("startTui: %v", err)
		}
		displayText(fmt.Sprintf("<%s (%s:%s) has joined>", memberName, memberIp, memberPort))


	case "list":
		err = listMembers()
		if err != nil {
			log.Fatalf("listMembers: %v", err)
		}

	case "subscribe":

		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: subscribe <name> <ip> <port>")
			os.Exit(2)
		}

		err = subscribeClient(flag.Arg(1), flag.Arg(2), flag.Arg(3))
		if err != nil {
			log.Fatalf("subscribeClient: %v", err)
		}

	case "unsubscribe":

		if flag.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "missing parameter: unsubscribe <name>")
			os.Exit(2)
		}

		err = unsubscribeClient(flag.Arg(1))
		if err != nil {
			log.Fatalf("unsubscribeClient: %v", err)
		}

	case "send":

		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "missing parameter: send <name> <text>...")
			os.Exit(1)
		}

		err = sendMessage(flag.Arg(1), strings.Join(flag.Args()[2:], " "))
		if err != nil {
			log.Fatalf("sendMessage: %v", err)
		}

	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
}
