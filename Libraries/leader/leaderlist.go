package leader

import (
	"fmt"
	"net"
	"syscall"

	"bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// Leaderlist has only one leader with status WORKING, which is the actual leader
type Leaderlist struct {
	name string
	leaderVersion int

	serviceIp   string
	servicePort int

	member *leaderlist.Leader
	List   []*leaderlist.Leader

	Message   *leaderlist.Message
	actionMap map[leaderlist.Message_MessageType]func(*leaderlist.Message, net.Addr) error
}

func (leaderlist *Leaderlist) TcpListen(addr string, acceptAddrInUse bool) (bool, error) {

	glog.V(3).Infof("TcpListen(%v)", addr)

	// Try to create listener
	listener, err := net.Listen("tcp", addr)

	if err != nil {

		// "address already in use" error
		if acceptAddrInUse {
			if err == syscall.EADDRINUSE {
				return false, nil
			}
		}

		// Unexpected error
		return false, err
	}

	// Goroutine for handling requests
	go func(listener net.Listener) {
		defer listener.Close()
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
	}(listener)
	return true, nil
}

func (leaderlist *Leaderlist) TcpSend(message *leaderlist.Message, recipient string) error {

	glog.V(2).Infof("TcpSend(%v, %v)", message, recipient)
	glog.V(2).Infof("message.MsgType: %v", message.MsgType)

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
