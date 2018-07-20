package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

// Publishing service on a commonly known address
const (
	serverIp   string = "localhost"
	serverPort string = "22365"
)

// Identity set by command args
var (
	memberName string
	memberIp   string
	memberPort string
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

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename := fmt.Sprintf("rudimentary-chat-%s-%v%02d%02d%02d%02d%02d.log", memberName,
		year, int(month), int(day), int(hour), int(minute), int(second))

	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile %v: %v", logfilename, err)
	}
	defer f.Close()

	// Switch logging to logfile
	log.SetOutput(f)

	// Start publishing service, if not running already
	startPublisher(serverIp, serverPort)

	// Subscribe client to publisher
	err = Subscribe(memberName, memberIp, memberPort)
	if err != nil {
		log.Fatalf("Subscribe: %v", err)
	}

	// Start displaying service for text-based UI
	err = startDisplayer(memberName, memberIp, memberPort)
	if err != nil {
		log.Fatalf("startDisplayer: %v", err)
	}

	// Start text-based UI
	err = runTUI()
	if err != nil {
		log.Fatalf("runTUI: %v", err)
	}
}
