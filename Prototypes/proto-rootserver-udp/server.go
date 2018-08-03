package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-udp/member-group"
	"os"
	"time"
)

var mbStorage []*membergroup.Member

var (
	err  error
	l    net.PacketConn
	n    int
	addr net.Addr
)

func main() {

	// Prepare logfile for logging
	year, month, day := time.Now().Date()
	hour, minute, second := time.Now().Clock()
	logfilename := fmt.Sprintf("rudimentary-chat-udp-%s-%v%02d%02d%02d%02d%02d.log", "test",
		year, int(month), int(day), int(hour), int(minute), int(second))

	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening logfile %v: %v", logfilename, err)
	}
	defer f.Close()

	// Config logging to logfile
	log.SetPrefix("DEBUG: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(f)

	// Create listener
	l, err = net.ListenPacket("udp", ":"+Port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	fmt.Printf("Rootserver does listen on %s:%s via udp\n", IpAddr, Port)

	buffer := make([]byte, 1024)

	go func() {
		for {
			n, addr, err = l.ReadFrom(buffer)
			if err != nil {
				log.Printf("cannot read from buffer:%v", err)
			} else {
				go func(buffer []byte, addr net.Addr) {
					reply := handleMembergroupConnection(buffer, addr)
					if reply != nil {
						l.WriteTo(reply, addr)
					}
				}(buffer[:n], addr)
			}
		}
	}()

	for {
	}
}

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func handleMembergroupConnection(data []byte, addr net.Addr) []byte {

	var (
		reply []byte
		err   error
	)

	// Switch according to the message type
	switch data[0] {

	case MemberMessage:

		// Handle the protobuf message: Member
		reply, err = handleMemberRequest(data[1:], addr)
		if err != nil {
			fmt.Printf("could not handle MemberMessage from %v: %v", addr, err)
		}

	case MemberListMessage:

		// Handle the protobuf message: MemberList
		reply, err = handleMemberListRequest(data[1:], addr)
		if err != nil {
			fmt.Printf("could not handle MemberListMessage from %v: %v", addr, err)
		}

	default:

		fmt.Printf("unknown MemberMessage")
	}
	return reply
}
