package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-subscriber/display"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Check command args
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "missing parameter: <ip> <port>")
		os.Exit(1)
	}

	// Create and register server
	srv := grpc.NewServer()

	// Register server for tweets
	var tweets displayServer
	display.RegisterDisplayTweetsServer(srv, tweets)

	// Create listener
	l, err := net.Listen("tcp", ":"+flag.Arg(1))
	fmt.Printf("Subscriber does listen on %s:%s\n", flag.Arg(0), flag.Arg(1))
	if err != nil {
		log.Fatal("could not listen to %s:%s: \v", flag.Arg(0), flag.Arg(1), err)
	}
	// Serve via listener
	log.Fatal(srv.Serve(l))
}

type displayServer struct{}

func (ds displayServer) Display(ctx context.Context, tweet *display.Tweet) (*display.Tweet, error) {

	fmt.Printf("%s (%s:%s): %q\n", tweet.Sender.Name, tweet.Sender.Ip, tweet.Sender.Port, tweet.Text)

	return tweet, nil
}
