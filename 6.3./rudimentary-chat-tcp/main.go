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
)

func main() {

	// Check command args and set identity
	flag.Parse()
	if flag.NArg() < 3 {
		fmt.Fprintln(os.Stderr, "missing parameter: <name> <ip> <port>")
		os.Exit(1)
	}
	memberName = flag.Arg(0)
	memberIp = flag.Arg(1)
	memberPort = flag.Arg(2)

	displayingService = memberIp + ":" + memberPort

	selfMember = &chatgroup.Member{Name:memberName, Ip:memberIp, Port:memberPort, Leader:false}

	isPublisher := false

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename := fmt.Sprintf("rudimentary-chat-tcp-%s-%v%02d%02d%02d%02d%02d.log", memberName,
		year, int(month), int(day), int(hour), int(minute), int(second))

	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile %v: %v", logfilename, err)
	}
	defer f.Close()

	if debug {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.SetPrefix("DEBUG: ")

		//debug = log.New(f, "DEBUG: ", log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetPrefix("LOG: ")
	}

	// Switch logging to logfile
	log.SetOutput(f)

	// Try to start publishing service and subscribe accordingly
	go func() {

		err := startPublisher()

		// Check if Publisher is "already in use"
		if err != nil && strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {

			isPublisher = true

			// Subscribe to the already running publishing service
			err = Subscribe()
			if err != nil {
				log.Fatalf("Failed to subscribe to running publishing service: %v", err)
			}
			log.Printf("Subscribed to the already running publishing service\n")
		} else {
			isPublisher = true
		}
	}()

	// Start displaying service
	go func() {
		err = startDisplayer()
		if err != nil {
			log.Fatalf("Failed to start displaying service on %q: %v", displayingService, err)
		}
	}()

	// Start text-based UI
	go func() {
		err = runTUI()
		if err != nil {
			log.Fatalf("runTUI: %v", err)
		}
	}()

	// todo: waitgroup
	time.Sleep(time.Second)

	if isPublisher {
		// Append text messages in "messages" view of publisher
		displayText(fmt.Sprintf("<publishing service running: %s (%s:%s)", memberName, memberIp, memberPort))
		displayText(fmt.Sprintf("<%s (%s:%s) has joined>", memberName, memberIp, memberPort))
	}

	for {
	}
}
