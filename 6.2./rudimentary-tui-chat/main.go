package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	memberName string
	memberIp   string
	memberPort string
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

		err = startPublisher(flag.Arg(1), flag.Arg(2))
		if err != nil {
			log.Fatal("startPublisher: %v\n", err)
		}

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

		err = subscribeClient(memberName, memberIp, memberPort)
		if err != nil {
			log.Fatal("subscribeClient: %v\n", err)
		}
		fmt.Printf("Client %q (%s:%s) has subscribed\n", memberName, memberIp, memberPort)

		err = startDisplayer(memberName, memberIp, memberPort)
		if err != nil {
			log.Fatal("startDisplayer: %v\n", err)
		}

		err = startTui()
		if err != nil {
			log.Fatal("startTui: %v\n", err)
		}
		displayText(fmt.Sprintf("<%s (%s:%s) has joined>", memberName, memberIp, memberPort))


	case "list":
		err = listMembers()
		if err != nil {
			log.Fatal("listMembers: %v\n", err)
		}

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
