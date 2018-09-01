package leader

import (
	"testing"
	"time"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

func TestUpdateRemoteIpAddr(t *testing.T) {
	name := "TestUpdateRemoteIpAddr"
	serviceIp := "127.0.0.1"
	servicePort := 22369
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := 22369
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
	memberPort = 12349
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
	leader.Message.Sender.Ip = "needsUpdate"
	for i, member := range leader.List {
		if member.Name == leader.Message.Sender.Name {
			leader.List[i] = leader.Message.Sender
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
	leader.Message.Sender.Ip = "needsUpdate"
	for i, member := range leader.List {
		if member.Name == leader.Message.Sender.Name {
			leader.List[i] = leader.Message.Sender
		}
	}
	for i, member := range service.List {
		if member.Name == leader.Message.Sender.Name {
			service.List[i] = leader.Message.Sender
		}
	}
	err = leader.SyncService()
	if err != nil {
		t.Fatalf("could not sync leader and service: %v", err)
	}
	for leader.LeaderVersion() <= lastVersion {
		time.Sleep(time.Millisecond * 10)
	}
	// todo Fix the following
	//
	//for _, member := range leader.List {
	//	if member.Name == leader.Message.Sender.Name {
	//		if member.Ip == "needsUpdate" {
	//			t.Errorf("Unexpected IP: %v", leader.Message.Sender.Ip)
	//			t.Errorf("leader: %v", leader)
	//		}
	//	}
	//}
	//
	//if leader.Message.Sender.Ip == "needsUpdate" {
	//	t.Errorf("Unexpected sender IP: %v", leader.Message.Sender.Ip)
	//	t.Errorf("leader: %v", leader)
	//	t.Errorf("service: %v", service) // !!!!!!
	//}
}
