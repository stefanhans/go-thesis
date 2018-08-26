package leader

import (
	"fmt"
	"net"
	"os"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

// Leaderlist has only one leader with status WORKING, which is the actual leader
type Leaderlist struct {
	name string

	service *leaderlist.Leader
	list    []*leaderlist.Leader

	Message          *leaderlist.Message
	defaultTransport string
	actionMap        map[leaderlist.Message_MessageType]func(*leaderlist.Message, net.Addr) error
}

func NewLeaderlist(config *Config) (*Leaderlist, error) {

	// Resolve IP string and update accordingly
	addr, err := net.ResolveIPAddr("ip", config.Ip)
	if err != nil {
		fmt.Printf("no valid ip address of client %q for publishing service: %v\n", config.Ip, err.Error())
		os.Exit(1)
	}
	config.Ip = addr.String()

	//fmt.Println("CHANGED")
	self := &leaderlist.Leader{
		Name:   config.Member,
		Ip:     config.Ip,
		Port:   fmt.Sprint(config.Port),
		Status: config.Status,
	}
	var list []*leaderlist.Leader
	list = append(list, self)

	actionMap := make(map[leaderlist.Message_MessageType]func(*leaderlist.Message, net.Addr) error)

	leader := &Leaderlist{
		name:    config.Leaderlist,
		service: self,
		list:    list,
		Message: &leaderlist.Message{
			MsgType: leaderlist.Message_LEADER_SYNC_REQUEST,
			Sender:  self,
			LeaderList: &leaderlist.LeaderList{
				Leader: list,
			},
		},
		defaultTransport: config.Transport,
		actionMap:        actionMap,
	}

	actionMap[leaderlist.Message_LEADER_SYNC_REQUEST] = leader.handleLeaderSyncRequest
	actionMap[leaderlist.Message_LEADER_SYNC_REPLY] = leader.handleLeaderSyncReply

	go leader.TcpListen()

	return leader, nil
}

func (leaderlist *Leaderlist) String() string {
	out := "leadergroup.Leaderlist:\n"
	out += fmt.Sprintf("\tName: %q\n", leaderlist.name)
	out += fmt.Sprintf("\tSelf: name:%q ip:%q port:%q status:%v\n",
		leaderlist.service.Name, leaderlist.service.Ip, leaderlist.service.Port, leaderlist.service.Status)
	for i, l := range leaderlist.list {
		out += fmt.Sprintf("\tList[%d]: name:%q ip:%q port:%q status:%v\n",
			i, l.Name, l.Ip, l.Port, l.Status)
	}
	out += fmt.Sprintf("\tMessage: %v\n", leaderlist.Message)
	out += fmt.Sprintf("\tDefaultTransport: %v\n", leaderlist.defaultTransport)
	out += fmt.Sprintf("\tActionMap: %v\n", leaderlist.actionMap)

	return out
}

func (leaderlist *Leaderlist) TcpListen() {

	leaderService := leaderlist.service.Ip + ":" + leaderlist.service.Port

	// Create listener
	listener, err := net.Listen(strings.ToLower(leaderlist.defaultTransport), leaderService)

	if err != nil {
		glog.Fatalf("could not listen to %q: %v\n", leaderService, err)
	}
	defer listener.Close()

	glog.Infof("Started list service listening on %q\n", leaderService)

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
