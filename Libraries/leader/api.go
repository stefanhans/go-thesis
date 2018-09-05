package leader

import (
	"fmt"
	"net"
	"strconv"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

// MemberNotFound is an error containing the name of the missing member
type MemberNotFound struct {
	Name string
}

// Error returns the error message with the missing member
func (e *MemberNotFound) Error() string {
	return fmt.Sprintf("member %q not found", e.Name)
}

// NewLeaderlist creates a struct which mainly contains
//
// - an identifier
//
// - an IP address of its service
//
// - an IP address of itself as a member
//
// - a list of members
//
// - a message to be customized and sent
//
// - a map of functions for the appropriate message type to be handled
//
func NewLeaderlist(
	name string,
	serviceIp string,
	servicePort string,
	memberName string,
	memberIp string,
	memberPort string,
	memberStatus leadlist.Leader_LeaderStatus) (*Leaderlist, error) {

	// Resolve IP string of service and update accordingly
	addr, err := net.ResolveIPAddr("ip", serviceIp)
	if err != nil {
		return nil, fmt.Errorf("no valid ip address %q for service: %v\n", serviceIp, err.Error())
	}
	serviceIp = addr.String()

	// Resolve IP string of member and update accordingly
	addr, err = net.ResolveIPAddr("ip", memberIp)
	if err != nil {
		return nil, fmt.Errorf("no valid ip address of member %q for publishing service: %v\n", memberIp, err.Error())
	}
	memberIp = addr.String()

	member := &leadlist.Leader{
		Name:   memberName,
		Ip:     memberIp,
		Port:   memberPort,
		Status: memberStatus,
	}
	var list []*leadlist.Leader
	list = append(list, member)

	actionMap := make(map[leadlist.Message_MessageType]func(*leadlist.Message, net.Addr) error)

	leader := &Leaderlist{
		name:          name,
		leaderVersion: 0,
		serviceIp:     serviceIp,
		servicePort:   servicePort,
		member:        member,
		list:          list,
		message: &leadlist.Message{
			MsgType: leadlist.Message_LEADER_SYNC_REQUEST,
			Sender:  member,
			LeaderList: &leadlist.LeaderList{
				Leader: list,
			},
		},
		actionMap: actionMap,
	}

	actionMap[leadlist.Message_LEADER_SYNC_REQUEST] = leader.handleLeaderSyncRequest
	actionMap[leadlist.Message_LEADER_SYNC_REPLY] = leader.handleLeaderSyncReply
	actionMap[leadlist.Message_PING_REQUEST] = leader.handlePingRequest
	actionMap[leadlist.Message_PING_REPLY] = leader.handlePingReply

	return leader, nil
}

// String() shows a textual representation of a leaderlist
func (leaderlist *Leaderlist) String() string {
	out := "leadergroup.leaderlist:\n"
	out += fmt.Sprintf("\tName: %q\n", leaderlist.name)
	out += fmt.Sprintf("\tService: ip:%q port:%q\n",
		leaderlist.serviceIp, leaderlist.servicePort)
	for i, l := range leaderlist.list {
		out += fmt.Sprintf("\tList[%d]: name:%q ip:%q port:%q status:%v\n",
			i, l.Name, l.Ip, l.Port, l.Status)
	}
	out += fmt.Sprintf("\tMessage: %v\n", leaderlist.message)
	out += fmt.Sprintf("\tActionMap: %v\n", leaderlist.actionMap)

	return out
}

// SetServiceIp sets the IP address of its service
func (leaderlist *Leaderlist) SetServiceIp(ip string) {
	leaderlist.serviceIp = ip
}

// ServiceIp returns the IP address of its service
func (leaderlist *Leaderlist) ServiceIp() string {

	return leaderlist.serviceIp
}

// SetServicePort sets the port number of its service
func (leaderlist *Leaderlist) SetServicePort(port string) error {

	// Port number is an integer
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	// Within free port number range without root access, i.e. [1024, 65535]
	if p < 1024 || p > 65535 {
		return fmt.Errorf("portnumber %d not between 1024 and 65535", p)
	}

	leaderlist.servicePort = port
	return nil
}

// ServicePort returns the port number of its service
func (leaderlist *Leaderlist) ServicePort() string {

	return leaderlist.servicePort
}

// SetMemberStatus sets the status of itself as a member
func (leaderlist *Leaderlist) SetMemberStatus(memberStatus leadlist.Leader_LeaderStatus) {
	leaderlist.member.Status = memberStatus
}

// MemberStatus returns the status of itself as a member
func (leaderlist *Leaderlist) MemberStatus() leadlist.Leader_LeaderStatus {

	return leaderlist.member.Status
}

// LeaderAddress returns the IP address of its leader
func (leaderlist *Leaderlist) LeaderAddress() string {

	for _, m := range leaderlist.list {
		if m.Status == leadlist.Leader_WORKING {
			return net.JoinHostPort(m.Ip, m.Port)
		}
	}
	return ""
}

// LeaderVersion returns the version of the leaderlist
func (leaderlist *Leaderlist) LeaderVersion() int {

	return leaderlist.leaderVersion
}

// RunService start the TCP listener to handle incoming requests for the service
func (leaderlist *Leaderlist) RunService() (bool, error) {

	return leaderlist.tcpListen(net.JoinHostPort(leaderlist.serviceIp, leaderlist.servicePort),
		true)
}

// RunClient start the TCP listener to handle incoming replies from the service
func (leaderlist *Leaderlist) RunClient() (bool, error) {

	return leaderlist.tcpListen(net.JoinHostPort(leaderlist.member.Ip, leaderlist.member.Port),
		false)
}

// SyncService sends a request to synchronize the leaderlist, i.e. returns SERVICE and WORKING members only
func (leaderlist *Leaderlist) SyncService() error {

	leaderlist.message.MsgType = leadlist.Message_LEADER_SYNC_REQUEST

	go leaderlist.tcpSend(leaderlist.message, net.JoinHostPort(leaderlist.serviceIp, leaderlist.servicePort))

	return nil
}

// PingMember sends a ping like request to a named member of the list
func (leaderlist *Leaderlist) PingMember(member string) error {
	for _, m := range leaderlist.list {

		if m.Name == member {
			leaderlist.message.MsgType = leadlist.Message_PING_REQUEST
			return leaderlist.tcpSend(leaderlist.message, net.JoinHostPort(m.Ip, m.Port))
		}
	}
	return &MemberNotFound{member}
}
