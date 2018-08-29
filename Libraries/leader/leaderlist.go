package leader

import (
	"fmt"
	"net"
	"bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// Leaderlist has only one leader with status WORKING, which is the actual leader
type Leaderlist struct {
	name string

	serviceIp   string
	servicePort int

	member *leaderlist.Leader
	List    []*leaderlist.Leader

	Message          *leaderlist.Message
	actionMap        map[leaderlist.Message_MessageType]func(*leaderlist.Message, net.Addr) error
}


func (leaderlist *Leaderlist) TcpListen(addr string) {

	// Create listener
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		glog.Fatalf("could not listen to %q: %v\n", addr, err)
	}
	defer listener.Close()

	glog.Infof("Started List service listening on %q\n", addr)

	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go leaderlist.handleRequest(conn)
	}
}

func (leaderlist *Leaderlist) TcpSend(message *leaderlist.Message, recipient string) error {

	glog.V(2).Infof("TcpSend: %v", message)

	// Connect to the recipient
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		return fmt.Errorf("could not connect to recipient %q: %v", recipient, err)
	}
	glog.V(2).Infof("TcpSend: %v", message)

	// Marshal into binary format
	byteArray, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("could not encode message: %v", err)
	}

	// Write the bytes to the connection
	n, err := conn.Write(byteArray)
	if err != nil {
		return fmt.Errorf("could not write message to the connection: %v", err)
	}
	glog.Infof("%v, %q: %d bytes sent\n", message, recipient, n)

	// Close connection
	return conn.Close()
}
