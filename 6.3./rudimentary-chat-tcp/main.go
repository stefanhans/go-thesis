package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"syscall"
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

	// Start publishing service, if not running already
	go func() {
		err := startPublisher()

		// Check if Publisher "already in use"
		if err != nil && strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) {

			// Subscribe application via running publisher
			log.Printf("Publisher 'already in use'\n")
			err = Subscribe()
			if err != nil {
				log.Fatalf("Subscribe: %v", err)
			}
		}
	}()
	// Start displaying service
	go func() {
		err = startDisplayer()
		if err != nil {
			log.Fatalf("startDisplayer on %q: %v", displayingService, err)
		}
	}()

	// Start text-based UI
	err = runTUI()
	if err != nil {
		log.Fatalf("runTUI: %v", err)
	}
}