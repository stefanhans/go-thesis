package leader

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// leaderlist has only one leader with status WORKING, which is the actual leader
type leaderlist struct {
	name          string
	leaderVersion int

	serviceIp   string
	servicePort int

	member *leadlist.Leader
	List   []*leadlist.Leader

	Message   *leadlist.Message
	actionMap map[leadlist.Message_MessageType]func(*leadlist.Message, net.Addr) error
}

func (leaderlist *leaderlist) tcpListen(addr string, acceptAddrInUse bool) (bool, error) {

	glog.V(3).Infof("tcpListen(%v)", addr)

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

func (leaderlist *leaderlist) tcpSend(message *leadlist.Message, recipient string) error {

	glog.V(2).Infof("tcpSend(%v, %v)", message, recipient)
	glog.V(2).Infof("message.MsgType: %v", message.MsgType)

	// Connect to the recipient
	conn, err := net.Dial("tcp", recipient)
	if err != nil {
		return fmt.Errorf("could not connect to recipient %q: %v", recipient, err)
	}
	glog.V(2).Infof("tcpSend: %v", message)

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

// Read all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func (leaderlist *leaderlist) handleRequest(conn net.Conn) {

	defer conn.Close()

	// Read all data from the connection
	var buf [512]byte
	var data []byte
	addr := conn.RemoteAddr()

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
		data = append(data, buf[0:n]...)
	}

	glog.Infof("%q->%q: %d bytes received\n", conn.RemoteAddr().String(), conn.LocalAddr().String(), len(data))

	// Unmarshall message
	var msg leadlist.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("could not unmarshall leadergroup.Message: %v\n", err)
		return
	}

	glog.Info(msg)

	// Fetch the handler from a map by the message type and call it accordingly
	if replyAction, ok := leaderlist.actionMap[msg.MsgType]; ok {
		err := replyAction(&msg, conn.RemoteAddr())
		if err != nil {
			fmt.Printf("could not handle %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		fmt.Printf("server: unknown message type %v\n", msg.MsgType)
	}
}

func (leaderlist *leaderlist) handleLeaderSyncRequest(message *leadlist.Message, addr net.Addr) error {
	glog.V(3).Infof("handleLeaderSyncRequest(%v, %v)", message, addr)

	updateVersion := false

	// Add sender, if not present
	exists := false
	for _, l := range leaderlist.List {
		if l.Name == message.Sender.Name {
			exists = true
			break
		}
	}

	if !exists {
		message.Sender.Status = leadlist.Leader_CANDIDATE
		leaderlist.List = append(leaderlist.List, message.Sender)

		updateVersion = true

		glog.V(3).Infof("Added sender: %v", leaderlist.List)
	}

	// No leader found
	if leaderlist.leaderCount() == 0 {

		// Set leader from message
		leaderSet := false
		for _, msgl := range message.LeaderList.Leader {
			if msgl.Status == leadlist.Leader_WORKING {
				for _, l := range leaderlist.List {
					if l.Name == msgl.Name {
						l.Status = leadlist.Leader_WORKING
						leaderSet = true
						break
					}
				}
				break
			}
		}

		// Set sender as leader, if needed
		if !leaderSet {
			for _, l := range leaderlist.List {
				if l.Name == message.Sender.Name {
					l.Status = leadlist.Leader_WORKING
					updateVersion = true
					break
				}
			}
		}
	}

	// Update remote IP address, if changed
	if leaderlist.updateRemoteIP(message, addr) {
		updateVersion = true
	}

	if updateVersion {
		leaderlist.leaderVersion++
	}

	leaderlist.Message.MsgType = leadlist.Message_LEADER_SYNC_REPLY
	leaderlist.Message.Sender = message.Sender
	leaderlist.Message.LeaderList.Leader = leaderlist.List

	err := leaderlist.tcpSend(leaderlist.Message, message.Sender.Ip+":"+message.Sender.Port)
	if err != nil {
		fmt.Printf("tcpSend: %v: %v", message, err)
	}

	return nil
}

func (leaderlist *leaderlist) handleLeaderSyncReply(message *leadlist.Message, addr net.Addr) error {
	glog.Info(message)

	leaderlist.List = message.LeaderList.Leader
	leaderlist.leaderVersion++

	glog.V(2).Infof("Replace List with (new) List: %q", leaderlist.List)
	glog.V(3).Infof("chatleaders: %v", leaderlist)

	return nil
}

func (leaderlist *leaderlist) handlePingRequest(message *leadlist.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	leaderlist.updateRemoteIP(message, addr)

	leaderlist.Message.MsgType = leadlist.Message_PING_REPLY
	leaderlist.Message.Sender = leaderlist.member

	err := leaderlist.tcpSend(leaderlist.Message, message.Sender.Ip+":"+message.Sender.Port)
	if err != nil {
		fmt.Printf("tcpSend: %v: %v", message, err)
	}

	return nil
}

func (leaderlist *leaderlist) handlePingReply(message *leadlist.Message, addr net.Addr) error {
	glog.Info(message)

	return nil
}

func (leaderlist *leaderlist) updateRemoteIP(msg *leadlist.Message, addr net.Addr) bool {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		glog.Infof("Remote Ip address update from %q to %q\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])

		// Update message's sender
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]

		// Update leader List
		for i, l := range leaderlist.List {
			if l.Name == msg.Sender.Name {
				if leaderlist.List[i].Ip != strings.Split(addr.String(), ":")[0] {

					leaderlist.List[i].Ip = strings.Split(addr.String(), ":")[0]
					leaderlist.leaderVersion++

					return true
				}
				break
			}
		}
	}
	return false
}

func (leaderlist *leaderlist) sendSyncReply(message *leadlist.Message) error {

	leaderlist.Message.MsgType = leadlist.Message_LEADER_SYNC_REPLY
	leaderlist.Message.LeaderList.Leader = leaderlist.List

	return leaderlist.tcpSend(leaderlist.Message, message.Sender.Ip+":"+message.Sender.Port)
}

func (leaderlist *leaderlist) leaderCount() int {
	cnt := 0
	for _, l := range leaderlist.List {
		if l.Status == leadlist.Leader_WORKING {
			cnt++
		}
	}
	return cnt
}
