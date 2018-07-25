package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-tcp/member-group"
	"bufio"
	_ "github.com/golang/protobuf/proto"
	"io"
)

var mbStorage []*membergroup.Member

func main() {

	fmt.Printf("%b\n", EOT)

	// Create listener
	l, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	fmt.Printf("Rootserver does listen on %s:%s\n", IpAddr, Port)

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go handleMembergroupConnection(conn)
	}
}

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func handleMembergroupConnection(conn net.Conn) {
	fmt.Printf("Handling new connection from %v\n", conn.RemoteAddr())

	// Close connection when this function ends
	defer func() {
		fmt.Println("Closing connection...\n")
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	for {
		buf, err := reader.ReadBytes(EOT)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received %d bytes\n", len(buf))

		// Switch according to the message type
		switch buf[0] {

		case Join:

			// Handle the protobuf message: Member
			err := handleMemberRequest(buf[1:len(buf)-1], conn.RemoteAddr())
			if err != nil {
				fmt.Printf("could not handle MemberMessage from %v: %v", conn, err)
				return
			}

		case Members:

			// Handle the protobuf message: MemberList
			err := handleMemberListRequest(buf[1 : len(buf)-1])
			if err != nil {
				fmt.Printf("could not handle MemberListMessage from %v: %v", conn, err)
				return
			}

		default:

			fmt.Printf("unknown MemberMessage\n")
		}
	}

	//data, err := ioutil.ReadAll(conn)
	//if err != nil {
	//	fmt.Printf("could not read data from %v: %v", conn, err)
	//	return
	//}

}
