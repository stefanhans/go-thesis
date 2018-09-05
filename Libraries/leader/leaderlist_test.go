package leader

import (
	"net"
	"testing"
	"time"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

func TestUpdateRemoteIpAddr(t *testing.T) {
	name := "TestUpdateRemoteIpAddr"
	serviceIp := "127.0.0.1"
	servicePort := "22369"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "22369"
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist: %v", err)
	}

	if service.LeaderVersion() != 0 {
		t.Errorf("Unexpected initial leader version: %v", service.LeaderVersion())
	}

	_, err = service.RunService()
	if err != nil {
		t.Fatalf("could not run service: %v", err)
	}

	if service.LeaderVersion() != 0 {
		t.Errorf("Unexpected first leader version: %v", service.LeaderVersion())
	}

	memberName = "alice"
	memberPort = "12349"
	memberStatus = leadlist.Leader_CANDIDATE
	leader, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist: %v", err)
	}

	if leader.LeaderVersion() != 0 {
		t.Errorf("Unexpected initial leader version: %v", leader.LeaderVersion())
	}

	_, err = leader.RunClient()
	if err != nil {
		t.Fatalf("could not run leader: %v", err)
	}

	lastVersion := leader.LeaderVersion()

	err = leader.SyncService()
	if err != nil {
		t.Fatalf("could not sync leader and service: %v", err)
	}

	for leader.LeaderVersion() <= lastVersion {
		time.Sleep(time.Millisecond * 10)
	}

	// Update only leader
	leader.message.Sender.Ip = "needsUpdate"
	for i, member := range leader.list {
		if member.Name == leader.message.Sender.Name {
			leader.list[i] = leader.message.Sender
		}
	}
	err = leader.SyncService()
	if err != nil {
		t.Fatalf("could not sync leader and service: %v", err)
	}
	for leader.LeaderVersion() <= lastVersion {
		time.Sleep(time.Millisecond * 10)
	}

	// Update leader and service
	leader.message.Sender.Ip = "needsUpdate"
	for i, member := range leader.list {
		if member.Name == leader.message.Sender.Name {
			leader.list[i] = leader.message.Sender
		}
	}

	for i, member := range service.list {
		if member.Name == leader.message.Sender.Name {
			service.list[i] = leader.message.Sender
		}
	}

	lastVersion = leader.LeaderVersion()
	err = leader.SyncService()
	if err != nil {
		t.Fatalf("could not sync leader and service: %v", err)
	}
	for leader.LeaderVersion() <= lastVersion {
		time.Sleep(time.Millisecond * 10)
	}

	for _, member := range leader.list {
		if member.Name == memberName {
			if member.Ip == "needsUpdate" {
				t.Errorf("Unexpected IP: %v", leader.message.Sender.Ip)
			}
		}
	}

	if leader.message.Sender.Ip == "needsUpdate" {
		t.Errorf("Unexpected sender IP: %v", leader.message.Sender.Ip)
	}

	leader.message.MsgType = leadlist.Message_INVALID
	err = leader.tcpSend(leader.message, net.JoinHostPort(serviceIp, servicePort))
	if err != nil {
		t.Errorf("Unexpected error from tcpSend(): %v", err)
	}
}

func TestTcpSend(t *testing.T) {
	name := "TestTcpSend"
	serviceIp := "127.0.0.1"
	servicePort := "22369"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "22369"
	memberStatus := leadlist.Leader_NOTFOUND

	sender, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist: %v", err)
	}

	err = sender.tcpSend(sender.message, "bla:foo")
	if err == nil {
		t.Errorf("Unexpected error from tcpSend(): %v", err)
	}
}

func TestHandleLeaderSyncRequest(t *testing.T) {

	// Establish the service
	name := "TestHandleLeaderSyncRequest"
	serviceIp := "127.0.0.1"
	servicePort := "22370"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "22370"
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist: %v", err)
	}

	started, err := service.RunService()
	if (!started) || err != nil {
		t.Errorf("Unexpected result from RunService(): %v %v", started, err)
	}

	// *******************************************
	// TESTCASE 0: no Leader_WORKING member
	// *******************************************
	service.list = []*leadlist.Leader{
		&leadlist.Leader{Name: "member1", Status: leadlist.Leader_CANDIDATE},
		&leadlist.Leader{Name: "member2", Status: leadlist.Leader_CANDIDATE},
		&leadlist.Leader{Name: "member3", Status: leadlist.Leader_CANDIDATE}}

	memberName = "requestor1"
	memberPort = "22354"
	memberStatus = leadlist.Leader_CANDIDATE

	requestor0, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist for candidate1: %v", err)
	}

	started, err = requestor0.RunClient()
	if (!started) || err != nil {
		t.Fatalf("Unexpected result from %q.RunClient(): %v %v", memberName, started, err)
	}

	err = requestor0.SyncService()
	if err != nil {
		t.Fatalf("Unexpected error from %q.SyncService(): %v", memberName, err)
	}
	for requestor0.LeaderVersion() <= 0 {
		time.Sleep(time.Millisecond * 10)
	}

	for _, member := range requestor0.list {
		if member.Name == "requestor0" {
			if member.Status != leadlist.Leader_WORKING {
				t.Errorf("Unexpected leader, i.e. status of %v is %v and not WORKING\n",
					member.Name, member.Status)
			}
		}
	}

	// *******************************************
	// TESTCASE 1: requestor has Leader_WORKING
	// *******************************************
	service.list = []*leadlist.Leader{
		&leadlist.Leader{Name: "member1", Status: leadlist.Leader_CANDIDATE},
		&leadlist.Leader{Name: "member2", Status: leadlist.Leader_CANDIDATE},
		&leadlist.Leader{Name: "member3", Status: leadlist.Leader_CANDIDATE}}

	memberName = "requestor1"
	memberPort = "22355"
	memberStatus = leadlist.Leader_WORKING

	requestor1, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist for candidate1: %v", err)
	}

	started, err = requestor1.RunClient()
	if (!started) || err != nil {
		t.Fatalf("Unexpected result from %q.RunClient(): %v %v", memberName, started, err)
	}

	err = requestor1.SyncService()
	if err != nil {
		t.Fatalf("Unexpected error from %q.SyncService(): %v", memberName, err)
	}
	for requestor1.LeaderVersion() <= 0 {
		time.Sleep(time.Millisecond * 10)
	}

	for _, member := range requestor1.list {
		if member.Name == "requestor1" {
			if member.Status != leadlist.Leader_WORKING {
				t.Errorf("Unexpected leader, i.e. status of %v is %v and not WORKING\n",
					member.Name, member.Status)
			}
		}
	}

	// *******************************************
	// TESTCASE 2: only one member has Leader_WORKING
	// *******************************************
	service.list = []*leadlist.Leader{
		&leadlist.Leader{Name: "member1", Status: leadlist.Leader_CANDIDATE},
		&leadlist.Leader{Name: "member2", Status: leadlist.Leader_CANDIDATE},
		&leadlist.Leader{Name: "member3", Status: leadlist.Leader_WORKING}}

	memberName = "requestor2"
	memberPort = "22356"
	memberStatus = leadlist.Leader_CANDIDATE

	requestor2, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist for %q: %v", memberName, err)
	}

	started, err = requestor2.RunClient()
	if (!started) || err != nil {
		t.Fatalf("Unexpected result from %q.RunClient(): %v %v", memberName, started, err)
	}

	err = requestor2.SyncService()
	if err != nil {
		t.Fatalf("Unexpected error from %q.SyncService(): %v", memberName, err)
	}
	for requestor2.LeaderVersion() <= 0 {
		time.Sleep(time.Millisecond * 10)
	}

	for _, member := range requestor2.list {
		if member.Name == "member3" {
			if member.Status != leadlist.Leader_WORKING {
				t.Errorf("Unexpected leader, i.e. status of %v is %v and not WORKING\n",
					member.Name, member.Status)
			}
		}
	}
}

func TestHandlePingRequest(t *testing.T) {
	name := "TestHandlePingRequest"
	serviceIp := "localhost"
	servicePort := "22371"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "22371"
	memberStatus := leadlist.Leader_UNKNOWN

	sender, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		t.Fatalf("could not create new leaderlist: %v", err)
	}

	addr, _ := net.ResolveIPAddr("tcp", "127.10.10.10")
	err = sender.handlePingRequest(sender.message, addr)
	if err == nil {
		t.Errorf("Unexpected error from tcpSend(): %v", err)
	}
}
