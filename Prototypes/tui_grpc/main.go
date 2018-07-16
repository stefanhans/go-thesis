package main

import (
	"bitbucket.org/stefanhans/go-thesis/Prototypes/tui_grpc/chat-message-pb"
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	msg             string
	nickname        string
	initialMessages string
	client          chatmessage.ChatMessagesClient
	portnumber      string
	portDial        string
)

func init() {
	nickname = "me"
	initialMessages = "init()"
}

func main() {
	// Check command args
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing portnumber")
		os.Exit(1)
	}
	portnumber = flag.Arg(0)
	portDial = flag.Arg(1)

	// Create and register server
	var messages chatMessageServer
	srv := grpc.NewServer()
	chatmessage.RegisterChatMessagesServer(srv, messages)

	// Create listener
	l, err := net.Listen("tcp", ":"+portnumber)
	if err != nil {
		log.Fatal("could not listen to :%v: %v", portnumber, err)
	}
	// Serve messages via listener
	go func() {
		log.Fatal(srv.Serve(l))
	}()

	// Create client with insecure connection
	conn, err := grpc.Dial(":"+portDial, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to backend: %v", err)
	}
	client = chatmessage.NewChatMessagesClient(conn)

	err = write(context.Background(), client, nickname, "test it")

	// Switch subcommands and call wrapper function
	//switch cmd := flag.Arg(0); cmd {
	//case "read":
	//	err = read(context.Background(), client)
	//case "write":
	//	if flag.NArg() < 4 {
	//		fmt.Fprintln(os.Stderr, "missing parameter: write <from> <text>...")
	//		os.Exit(1)
	//	}
	//	err = write(context.Background(), client, flag.Arg(1), strings.Join(flag.Args()[2:], " "))
	//default:
	//	err = fmt.Errorf("unknown subcommand %s", cmd)
	//}
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	os.Exit(1)
	//}

	// Open file for reading info
	byteArray, _ := OpenStorageForRead("storage")

	fmt.Printf("1: %q\n", initialMessages)

	initialMessages = readMessage(byteArray)
	fmt.Printf("2: %q\n", initialMessages)

	// Create the terminal GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	// Set function to manage all views and keybindings
	g.SetManagerFunc(layout)

	// Bind keys with functions
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send)

	// Start main event loop of the GUI
	g.MainLoop()
}
