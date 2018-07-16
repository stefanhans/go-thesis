package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-subscriber/display"
	"flag"
)

func main() {
	// Check command args
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing port")
		os.Exit(1)
	}

	// Create client with insecure connection
	conn, err := grpc.Dial(":" + flag.Arg(0), grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to backend: %v", err)
	}
	fmt.Printf("Dial to localhost:%s\n", flag.Arg(0))
	client := display.NewDisplayTweetsClient(conn)

	sender := display.Sender{Name: "stefan", Ip: "localhost", Port: "22365"}

	tweet := display.Tweet{Sender: &sender, Text: "hello from stefan"}

	err = show(context.Background(), client, &tweet)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Tweet wrapper function
func show(ctx context.Context, client display.DisplayTweetsClient, tweet *display.Tweet) error {

	// Write to gRPC client
	_, err := client.Display(ctx, tweet)
	if err != nil {
		return fmt.Errorf("could not display tweet: %v", err)
	}
	return nil
}