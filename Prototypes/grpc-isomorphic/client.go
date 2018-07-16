package main

import (
	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-isomorphic/info"
	"fmt"
	"golang.org/x/net/context"
)

//func main() {
//	// Check command args
//	flag.Parse()
//	if flag.NArg() < 1 {
//		fmt.Fprintln(os.Stderr, "missing subcommand: read or write")
//		os.Exit(1)
//	}
//
//	// Create client with insecure connection
//	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
//	if err != nil {
//		log.Fatal("could not connect to backend: %v", err)
//	}
//	client := info.NewInfosClient(conn)
//
//	// Switch subcommands and call wrapper function
//	switch cmd := flag.Arg(0); cmd {
//	case "read":
//		err = read(context.Background(), client)
//	case "write":
//		if flag.NArg() < 3 {
//			fmt.Fprintln(os.Stderr, "missing parameter: write <from> <text>...")
//			os.Exit(1)
//		}
//		err = write(context.Background(), client, flag.Arg(1), strings.Join(flag.Args()[2:], " "))
//	default:
//		err = fmt.Errorf("unknown subcommand %s", cmd)
//	}
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//}

// Write wrapper function
func write(ctx context.Context, client info.InfosClient, from string, text string) error {

	// Write to gRPC client
	_, err := client.Write(ctx, &info.Info{From: from, Text: text})
	if err != nil {
		return fmt.Errorf("could not add info in the backend: %v", err)
	}
	return nil
}

// Read wrapper function
func read(ctx context.Context, client info.InfosClient) error {

	// Read from gRPC client
	l, err := client.Read(ctx, &info.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch info: %v", err)
	}

	// Print messages
	for _, t := range l.Infos {
		fmt.Printf("%s: %s\n", t.From, t.Text)
	}
	return nil
}

// Subscribe wrapper function
func subscribe(ctx context.Context, client info.InfosClient, name string, port string) error {
	// Subscribe to gRPC client
	_, err := client.Subscribe(ctx, &info.Member{Name: name, Port: port})
	if err != nil {
		return fmt.Errorf("could not subscribe: %v", err)
	}
	return nil
}

// Publish wrapper function
func Publish(ctx context.Context, client info.InfosClient, from string, text string) error {
	// Publish via gRPC server
	_, err := client.Publish(ctx, &info.Info{From: from, Text: text})
	if err != nil {
		return fmt.Errorf("could not publish: %v", err)
	}
	return nil
}
