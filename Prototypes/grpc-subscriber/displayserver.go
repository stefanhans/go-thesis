package main

import (
	"fmt"

	"google.golang.org/grpc"
	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-subscriber/display"
	"golang.org/x/net/context"
	"net"
	"log"
)

func main() {

	// Create and register server
	srv := grpc.NewServer()

	// Register server for tweets
	var tweets displayServer
	display.RegisterDisplayTweetsServer(srv, tweets)

	// Create listener
	l, err := net.Listen("tcp", ":22365")
	fmt.Printf("Displayserver does listen on localhost:22365\n", )
	if err != nil {
		log.Fatal("could not listen to :22365: \v", err)
	}
	// Serve via listener
	log.Fatal(srv.Serve(l))
}

type displayServer struct {}

func (ds displayServer) Display(ctx context.Context, tweet *display.Tweet) (*display.Tweet, error) {

	fmt.Printf("Display: %q from %s (%s:%s)\n", tweet.Text, tweet.Sender.Name, tweet.Sender.Ip, tweet.Sender.Port)

	return tweet, nil
}