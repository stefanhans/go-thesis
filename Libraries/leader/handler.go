package leader

import (
	"fmt"
	"net"
	"strings"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// Read all incoming data, take the leading byte as message type,
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

	glog.Infof("%q->%q: %d bytes received\n", conn.RemoteAddr().String(), conn.LocalAddr().String(), len(data))

	// Unmarshall message
	var msg leadlist.Message
	err := proto.Unmarshal(data, &msg)
	if err != nil {
		fmt.Printf("could not unmarshall leadergroup.Message: %v", err)
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

func (leaderlist *Leaderlist) handleLeaderSyncRequest(message *leadlist.Message, addr net.Addr) error {
	glog.V(3).Info(message)

	// Add sender, if not present
	exists := false
	for _, l := range leaderlist.list {
		if l.Name == message.Sender.Name {
			exists = true
			break
		}
	}
	glog.V(3).Infof("Added sender: %v", leaderlist.list)

	if !exists {
		message.Sender.Status = leadlist.Leader_UNKNOWN
		leaderlist.list = append(leaderlist.list, message.Sender)
	}

	// No leader found
	if leaderlist.leaderCount() == 0 {

		// Set leader from message
		leaderSet := false
		for _, msgl := range message.LeaderList.Leader {
			if msgl.Status == leadlist.Leader_WORKING {
				for _, l := range leaderlist.list {
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
			for _, l := range leaderlist.list {
				if l.Name == message.Sender.Name {
					l.Status = leadlist.Leader_WORKING
					break
				}
			}
		}
	}

	// Update remote IP address, if changed
	leaderlist.updateRemoteIP(message, addr)

	leaderlist.Message.MsgType = leadlist.Message_LEADER_SYNC_REPLY
	leaderlist.Message.Sender = message.Sender
	err := leaderlist.sendSyncReply(leaderlist.Message)
	if err != nil {
		fmt.Printf("sendSyncReply: %v: %v", message, err)
	}

	return nil
}

func (leaderlist *Leaderlist) handleLeaderSyncReply(message *leadlist.Message, addr net.Addr) error {
	glog.Info(message)

	leaderlist.list = message.LeaderList.Leader

	glog.V(2).Infof("Replace list with (new) list: %q", leaderlist.list)
	glog.V(3).Infof("chatleaders: %v", leaderlist)

	return nil
}

func (leaderlist *Leaderlist) updateRemoteIP(msg *leadlist.Message, addr net.Addr) {

	// Check remote Ip address change of message
	if msg.Sender.Ip != strings.Split(addr.String(), ":")[0] {
		glog.Infof("Remote Ip address update from %q to %q\n", msg.Sender.Ip, strings.Split(addr.String(), ":")[0])

		// Update message's sender
		msg.Sender.Ip = strings.Split(addr.String(), ":")[0]

		// Update leader list
		for i, l := range leaderlist.list {
			if l.Name == msg.Sender.Name {
				leaderlist.list[i].Ip = strings.Split(addr.String(), ":")[0]
				break
			}
		}

	}
}

func (leaderlist *Leaderlist) sendSyncReply(message *leadlist.Message) error {

	leaderlist.Message.MsgType = leadlist.Message_LEADER_SYNC_REPLY
	leaderlist.Message.LeaderList.Leader = leaderlist.list

	return leaderlist.TcpSend(leaderlist.Message, message.Sender.Ip+":"+message.Sender.Port)
}

func (leaderlist *Leaderlist) leaderCount() int {
	cnt := 0
	for _, l := range leaderlist.list {
		if l.Status == leadlist.Leader_WORKING {
			cnt++
		}
	}
	return cnt
}
