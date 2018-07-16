package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-isomorphic/info"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"time"
)

var (
	myself     info.Member
	memberlist []*info.Member
)

func main() {
	// Check command args
	flag.Parse()
	if flag.NArg() < 3 {
		fmt.Fprintln(os.Stderr, "missing subcommand: read or write")
		os.Exit(1)
	}

	myself.Name = flag.Arg(1)
	myself.Port = flag.Arg(2)

	memberlist = append(memberlist, &myself)

	// Create and register server
	var infos infoServer
	srv := grpc.NewServer()
	info.RegisterInfosServer(srv, infos)

	// Create listener
	l, err := net.Listen("tcp", ":"+myself.Port)
	if err != nil {
		log.Fatal("could not listen to :%v: \v", myself.Port, err)
	}
	// Serve messages via listener
	go func() { log.Fatal(srv.Serve(l)) }()

	// Create client with insecure connection
	conn, err := grpc.Dial(":"+myself.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to backend: %v", err)
	}
	client := info.NewInfosClient(conn)

	// Switch subcommands and call wrapper function
	switch cmd := flag.Arg(0); cmd {
	case "read":
		err = read(context.Background(), client)
	case "write":
		if flag.NArg() < 3 {
			fmt.Fprintln(os.Stderr, "missing parameter: write <from> <text>...")
			os.Exit(1)
		}
		err = write(context.Background(), client, flag.Arg(1), strings.Join(flag.Args()[2:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Subscribe gRPC client
	err = subscribe(context.Background(), client, myself.Name, myself.Port)
	if err != nil {
		fmt.Errorf("could not Subscribe: %v", err)
	}
	//time.Sleep(time.Second)

	err = Publish(context.Background(), client, myself.Name, myself.Port)
	if err != nil {
		fmt.Errorf("could not Publish: %v", err)
	}
	//
	//for _, member := range memberlist {
	//	conn, err := grpc.Dial(":" + member.Port, grpc.WithInsecure())
	//	if err != nil {
	//		log.Fatal("could not connect to backend: %v", err)
	//	}
	//	client = info.NewInfosClient(conn)
	//	err = write(context.Background(), client, flag.Arg(1), strings.Join(flag.Args()[2:], " "))
	//	if err != nil {
	//		fmt.Errorf("could not Publish: %v", err)
	//	}
	//}
	time.Sleep(time.Second)
}
