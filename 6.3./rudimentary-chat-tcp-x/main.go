package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
	"github.com/jroimartin/gocui"
)

func main() {

	// Check command args and set own chatgroup.Member
	flag.Parse()
	if flag.NArg() < 3 {
		fmt.Fprintln(os.Stderr, "missing parameter: <name> <ip> <port>")
		os.Exit(1)
	}
	selfMember = &chatgroup.Member{Name: flag.Arg(0), Ip: flag.Arg(1), Port: flag.Arg(2), Leader: false}

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename = fmt.Sprintf("rudimentary-chat-tcp-%s-%v%02d%02d%02d%02d%02d.log", selfMember.Name,
		year, int(month), int(day), int(hour), int(minute), int(second))

	var err error
	logfile, err = os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile %v: %v", logfilename, err)
	}
	defer logfile.Close()

	// Config logging
	if debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("DEBUG: ")

		//debug = log.New(f, "DEBUG: ", log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetPrefix("LOG: ")
	}

	// Switch logging to logfile
	log.SetOutput(logfile)

	// Initialize chat command usage
	commandUsageInit()

	// Create the TUI
	clientGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		fmt.Errorf("could not create tui: %v\n", err)
	}
	defer clientGui.Close()

	// Try to start publishing service and subscribe accordingly
	go func() {

		err := startPublisher()

		// Check if Publisher is "already in use"
		if err != nil && strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {

			// Subscribe to the already running publishing service
			err = Subscribe()
			if err != nil {
				log.Fatalf("Failed to subscribe to running publishing service: %v", err)
			}
			log.Printf("Subscribed to the already running publishing service\n")
		}
	}()

	// Start displaying service
	go func() {

		displayingService = selfMember.Ip + ":" + selfMember.Port

		err = startDisplayer()
		if err != nil {
			log.Fatalf("Failed to start displaying service on %q: %v", displayingService, err)
		}
	}()

	// Start text-based UI
	err = runTUI()
	if err != nil {
		log.Fatalf("runTUI: %v", err)
	}
}
