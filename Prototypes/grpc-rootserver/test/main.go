package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-rootserver/member-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Check command args
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: register or list")
		os.Exit(1)
	}

	// Create client with insecure connection
	conn, err := grpc.Dial(":22365", grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to rootserver on localhost:22365: %v", err)
	}
	fmt.Printf("Dialed to rootserver on localhost:22365\n")
	client := membergroup.NewMembersClient(conn)

	// Switch subcommands and call wrapper function
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "register":
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: register <name> <ip> <port> [true]")
			os.Exit(1)
		}
		var leader bool
		if flag.NArg() > 3 && flag.Arg(4) == "true" {
			leader = true
		}
		err = register(context.Background(), client, flag.Arg(1), flag.Arg(2), flag.Arg(3), leader)
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Subscribe wrapper function
func register(ctx context.Context, client membergroup.MembersClient, name string, ip string, port string, leader bool) error {

	// Write to gRPC client
	_, err := client.Register(ctx, &membergroup.Member{Name: name, Ip: ip, Port: port, Leader: leader})
	if err != nil {
		return fmt.Errorf("could not add member in the membergroup: %v", err)
	}
	return nil
}

// List wrapper function
func list(ctx context.Context, client membergroup.MembersClient) error {

	// List from gRPC client
	l, err := client.List(ctx, &membergroup.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch membergroup: %v", err)
	}

	// Print members
	for _, t := range l.Member {
		fmt.Printf("%s %s %s %v\n", t.Name, t.Ip, t.Port, t.Leader)
	}
	return nil
}
