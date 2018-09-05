package leader

import (
	"fmt"
	"net"
	"strings"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// Leaderlist should have only one leader with status WORKING, which is the actual leader
type Leaderlist struct {
	name          string
	leaderVersion int

	serviceIp   string
	servicePort string

	member *leadlist.Leader
	list   []*leadlist.Leader

	message   *leadlist.Message
	actionMap map[leadlist.Message_MessageType]func(*leadlist.Message, net.Addr) error
}

// tcpListen is a method to start a listener and handle incoming requests in separate goroutines
func (leaderlist *Leaderlist) tcpListen(addr string, acceptAddrInUse bool) (bool, error) {

	glog.V(3).Infof("tcpListen(%v)", addr)

	// Try to create listener
	listener, err := net.Listen("tcp", addr)

	if err != nil {

		// "address already in use" error
		if acceptAddrInUse {
			if strings.Contains(err.Error(), "address already in use") {
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

// tcpSend
func (leaderlist *Leaderlist) tcpSend(message *leadlist.Message, recipient string) error {

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
	glog.V(2).Infof("%v, %q: %d bytes sent\n", message, recipient, n)

	// Close connection
	return conn.Close()
}

// handleRequest reads all incoming data, take the leading byte as message type,
// and use the appropriate handler for the rest
func (leaderlist *Leaderlist) handleRequest(conn net.Conn) {

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

	glog.V(1).Infof("%q->%q: %d bytes received\n", conn.RemoteAddr().String(), conn.LocalAddr().String(), len(data))

	// Unmarshall message
	var msg leadlist.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		glog.Errorf("could not unmarshall leadergroup.Message: %v\n", err)
		return
	}

	glog.V(3).Info(msg)

	// Fetch the handler from a map by the message type and call it accordingly
	if replyAction, ok := leaderlist.actionMap[msg.MsgType]; ok {
		err := replyAction(&msg, conn.RemoteAddr())
		if err != nil {
			glog.Errorf("could not handle request %v from %v: %v", msg.MsgType, addr, err)
		}
	} else {
		glog.V(1).Infof("server: unknown message type %v\n", msg.MsgType)
	}
}

// handleLeaderSyncRequest
func (leaderlist *Leaderlist) handleLeaderSyncRequest(message *leadlist.Message, addr net.Addr) error {
	glog.V(3).Infof("handleLeaderSyncRequest(%v, %v)", message, addr)

	// Update remote IP address, if changed
	updateVersion := false
	if leaderlist.updateRemoteIP(message, addr) {
		updateVersion = true
	}

	// Is sender already in service list?
	exists := false
	for _, l := range leaderlist.list {
		if l.Name == message.Sender.Name {
			exists = true
			break
		}
	}

	// Add sender as candidate, if needed
	if !exists {
		message.Sender.Status = leadlist.Leader_CANDIDATE
		leaderlist.list = append(leaderlist.list, message.Sender)

		updateVersion = true

		glog.V(3).Infof("Added sender: %v", leaderlist.list)
	}

	// No leader found
	if leaderlist.leaderCount() == 0 {
		isLeaderSet := false

		// Set leader from message, if already existing
		for _, msglist := range message.LeaderList.Leader {

			// Is member of message list a leader?
			if msglist.Status == leadlist.Leader_WORKING {

				// Set leader in leaderlist
				for _, l := range leaderlist.list {
					if l.Name == msglist.Name {
						l.Status = leadlist.Leader_WORKING
						isLeaderSet = true
						break
					}
				}
			}
		}

		// Set sender as leader, if needed
		if !isLeaderSet {
			for _, l := range leaderlist.list {
				if l.Name == message.Sender.Name {
					l.Status = leadlist.Leader_WORKING
					updateVersion = true
					break
				}
			}
		}
	}

	// Increase version of leaderlist
	if updateVersion {
		leaderlist.leaderVersion++
	}

	leaderlist.message.MsgType = leadlist.Message_LEADER_SYNC_REPLY
	leaderlist.message.Sender = message.Sender

	// Add only service and leaders
	leaderlist.message.LeaderList.Leader = leaderlist.collectWorking()
	glog.V(3).Infof("leaderlist.Message.LeaderList.Leader: %v", leaderlist.message.LeaderList.Leader)

	err := leaderlist.sendSyncReply(leaderlist.message)
	if err != nil {
		glog.Errorf("could not send synchronization reply: %v: %v", leaderlist.message, err)
	}

	return nil
}

// handleLeaderSyncReply
func (leaderlist *Leaderlist) handleLeaderSyncReply(message *leadlist.Message, addr net.Addr) error {
	glog.V(3).Info(message)

	leaderlist.list = message.LeaderList.Leader
	leaderlist.leaderVersion++

	glog.V(2).Infof("Replace List with (new) List: %q", leaderlist.list)
	glog.V(3).Infof("chatleaders: %v", leaderlist)

	return nil
}

// handlePingRequest
func (leaderlist *Leaderlist) handlePingRequest(message *leadlist.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	leaderlist.updateRemoteIP(message, addr)

	leaderlist.message.MsgType = leadlist.Message_PING_REPLY
	leaderlist.message.Sender = leaderlist.member

	err := leaderlist.tcpSend(leaderlist.message, message.Sender.Ip+":"+message.Sender.Port)
	if err != nil {
		return fmt.Errorf("tcpSend: %v: %v", message, err)
	}

	return nil
}

// handlePingReply is not really implemented yet
func (leaderlist *Leaderlist) handlePingReply(message *leadlist.Message, addr net.Addr) error {
	glog.V(3).Info(message)

	// todo introduce resilience

	return nil
}

// updateRemoteIP
func (leaderlist *Leaderlist) updateRemoteIP(msg *leadlist.Message, addr net.Addr) bool {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		glog.V(2).Infof("Remote Ip address update from %q to %q\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])

		// Update message's sender
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]

		// Update leader List
		for i, l := range leaderlist.list {
			if l.Name == msg.Sender.Name {
				if leaderlist.list[i].Ip != strings.Split(addr.String(), ":")[0] {
					leaderlist.list[i].Ip = strings.Split(addr.String(), ":")[0]
					leaderlist.leaderVersion++

					return true
				}
				break
			}
		}
	}
	return false
}

// sendSyncReply sends the message as LEADER_SYNC_REPLY to its original sender
func (leaderlist *Leaderlist) sendSyncReply(message *leadlist.Message) error {

	leaderlist.message.MsgType = leadlist.Message_LEADER_SYNC_REPLY

	return leaderlist.tcpSend(leaderlist.message, message.Sender.Ip+":"+message.Sender.Port)
}

// leaderCount
func (leaderlist *Leaderlist) leaderCount() int {
	cnt := 0
	for _, l := range leaderlist.list {
		if l.Status == leadlist.Leader_WORKING {
			cnt++
		}
	}
	return cnt
}

// collectWorking collect list items with status SERVICE or WORKING
func (leaderlist *Leaderlist) collectWorking() []*leadlist.Leader {
	var list []*leadlist.Leader

	for _, member := range leaderlist.list {

		if member.Status == leadlist.Leader_WORKING || member.Status == leadlist.Leader_SERVICE {
			list = append(list, member)
		}
	}

	return list
}
