package leader

import (
	"fmt"
	"net"
	"os"
	"strconv"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

func NewLeaderlist(name string,
	serviceIp  string,
	servicePort int,
	memberName string,
	memberIp string,
	memberPort int,
	memberStatus leadlist.Leader_LeaderStatus) (*Leaderlist, error) {

	// Resolve IP string of service and update accordingly
	addr, err := net.ResolveIPAddr("ip", serviceIp)
	if err != nil {
		fmt.Printf("no valid ip address of service %q for publishing service: %v\n", serviceIp, err.Error())
		os.Exit(1)
	}
	serviceIp = addr.String()

	// Resolve IP string of member and update accordingly
	addr, err = net.ResolveIPAddr("ip", memberIp)
	if err != nil {
		fmt.Printf("no valid ip address of member %q for publishing service: %v\n", memberIp, err.Error())
		os.Exit(1)
	}
	memberIp = addr.String()

	//fmt.Println("CHANGED")
	member := &leadlist.Leader{
		Name:   memberName,
		Ip:     memberIp,
		Port:   fmt.Sprint(memberPort),
		Status: memberStatus,
	}
	var list []*leadlist.Leader
	list = append(list, member)

	actionMap := make(map[leadlist.Message_MessageType]func(*leadlist.Message, net.Addr) error)

	leader := &Leaderlist{
		name:    name,
		leaderVersion: 0,
		serviceIp: serviceIp,
		servicePort: servicePort,
		member: member,
		List:    list,
		Message: &leadlist.Message{
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

func (leaderlist *Leaderlist) String() string {
	out := "leadergroup.Leaderlist:\n"
	out += fmt.Sprintf("\tName: %q\n", leaderlist.name)
	out += fmt.Sprintf("\tService: ip:%q port:%d\n",
		leaderlist.serviceIp, leaderlist.servicePort)
	for i, l := range leaderlist.List {
		out += fmt.Sprintf("\tList[%d]: name:%q ip:%q port:%q status:%v\n",
			i, l.Name, l.Ip, l.Port, l.Status)
	}
	out += fmt.Sprintf("\tMessage: %v\n", leaderlist.Message)
	out += fmt.Sprintf("\tActionMap: %v\n", leaderlist.actionMap)

	return out
}

func (leaderlist *Leaderlist) SetServiceIp(ip string) {
	leaderlist.serviceIp = ip
}

func (leaderlist *Leaderlist) ServiceIp() string {

	return leaderlist.serviceIp
}

func (leaderlist *Leaderlist) SetServicePort(port string) error {
	p, err := strconv.Atoi(port)
	// Todo check valid port
	if err != nil {
		return err
	}
	leaderlist.servicePort = p
	return nil
}

func (leaderlist *Leaderlist) ServicePort() int {

	return leaderlist.servicePort
}

func (leaderlist *Leaderlist) SetMemberStatus(memberStatus leadlist.Leader_LeaderStatus) {
	leaderlist.member.Status = memberStatus
}

func (leaderlist *Leaderlist) MemberStatus() leadlist.Leader_LeaderStatus {

	return leaderlist.member.Status
}



func (leaderlist *Leaderlist) LeaderAddress() string {

	for _, m := range leaderlist.List {
		if m.Status == leadlist.Leader_WORKING {
			return fmt.Sprintf("%s:%s", m.Ip, m.Port)
		}
	}
	return ""
}
func (leaderlist *Leaderlist) LeaderVersion() int {

	return leaderlist.leaderVersion
}

func (leaderlist *Leaderlist) RunService() (bool, error) {

	return leaderlist.TcpListen(fmt.Sprintf("%s:%d", leaderlist.serviceIp, leaderlist.servicePort),
		true)
}

func (leaderlist *Leaderlist) RunClient() (bool, error) {

	return leaderlist.TcpListen(fmt.Sprintf("%s:%s", leaderlist.member.Ip, leaderlist.member.Port),
		false)
}

func (leaderlist *Leaderlist) SyncService() error {

	leaderlist.Message.MsgType = leadlist.Message_LEADER_SYNC_REQUEST

	go leaderlist.TcpSend(leaderlist.Message, fmt.Sprintf("%s:%d", leaderlist.serviceIp, leaderlist.servicePort))

	return nil
}

func (leaderlist *Leaderlist) PingMember(member string) error {
	for _, m := range leaderlist.List {

		if m.Name == member {
			leaderlist.Message.MsgType = leadlist.Message_PING_REQUEST
			return leaderlist.TcpSend(leaderlist.Message, fmt.Sprintf("%s:%s", m.Ip, m.Port))
		}
	}
	// todo error member not found
	return nil
}
