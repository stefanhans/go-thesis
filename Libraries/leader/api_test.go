package leader

import (
	"fmt"
	"log"
	"net"
	"strings"
	"testing"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"time"
)

func TestSimpleNewLeaderlist(t *testing.T) {
	name := "testlist"
	serviceIp := "127.0.0.1"
	servicePort := "22365"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_UNKNOWN

	leaderlist, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// name
	t.Run("name", func(t *testing.T) {
		if leaderlist.name != name {
			t.Errorf("Unexpected name: %v", leaderlist.name)
		}
	})

	// serviceIp
	t.Run("serviceIp", func(t *testing.T) {
		if leaderlist.serviceIp != serviceIp {
			t.Errorf("Unexpected serviceIp: %v", leaderlist.serviceIp)
		}
	})

	// servicePort
	t.Run("servicePort", func(t *testing.T) {

		if leaderlist.servicePort != servicePort {
			t.Errorf("Unexpected servicePort: %v", leaderlist.servicePort)
		}
	})

	// memberName
	t.Run("memberName", func(t *testing.T) {

		if leaderlist.member.Name != memberName {
			t.Errorf("Unexpected memberName: %v", leaderlist.member.Name)
		}
	})

	// memberIp
	t.Run("memberIp", func(t *testing.T) {

		if leaderlist.member.Ip != memberIp {
			t.Errorf("Unexpected memberIp: %v", leaderlist.member.Ip)
		}
	})

	// memberPort
	t.Run("memberPort", func(t *testing.T) {

		if leaderlist.member.Port != memberPort {
			t.Errorf("Unexpected memberPort: %v", leaderlist.member.Port)
		}
	})

	// memberStatus
	t.Run("memberStatus", func(t *testing.T) {

		if leaderlist.member.Status != memberStatus {
			t.Errorf("Unexpected serviceIp: %v", leaderlist.member.Status)
		}
	})
}

func TestResolveIpNewLeaderlist(t *testing.T) {
	name := "testlist"
	serviceIp := "localhost"
	servicePort := "22365"
	memberName := "service"
	memberIp := "localhost"
	memberPort := "12345"
	memberStatus := leadlist.Leader_UNKNOWN

	leaderlist, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// serviceIp
	t.Run("serviceIp", func(t *testing.T) {
		addr, _ := net.ResolveIPAddr("ip", serviceIp)
		if leaderlist.serviceIp != addr.String() {
			t.Errorf("Unexpected serviceIp: %v", leaderlist.serviceIp)
		}
	})

	// memberIp
	t.Run("memberIp", func(t *testing.T) {
		addr, _ := net.ResolveIPAddr("ip", memberIp)

		if leaderlist.member.Ip != addr.String() {
			t.Errorf("Unexpected memberIp: %v", leaderlist.member.Ip)
		}
	})

	// error test serviceIp
	serviceIp = "invalid"
	t.Run("invalid serviceIp", func(t *testing.T) {
		_, err = NewLeaderlist(name, serviceIp, servicePort, memberName,
			memberIp, memberPort, memberStatus)
		if !strings.HasSuffix(err.Error(), "no such host\n") {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	// error test memberIp
	serviceIp = "localhost"
	memberIp = "invalid"
	t.Run("invalid memberIp", func(t *testing.T) {
		_, err = NewLeaderlist(name, serviceIp, servicePort, memberName,
			memberIp, memberPort, memberStatus)
		if !strings.HasSuffix(err.Error(), "no such host\n") {
			t.Errorf("Unexpected error: %q", err)
		}
	})
}

func TestString(t *testing.T) {
	name := "testlist"
	serviceIp := "127.0.0.1"
	servicePort := "22365"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_UNKNOWN

	leaderlist, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// String
	t.Run("String", func(t *testing.T) {
		if !strings.Contains(leaderlist.String(), name) {
			t.Errorf("No name %q in String(): %v", name, leaderlist.String())
		}
		if !strings.Contains(leaderlist.String(), serviceIp) {
			t.Errorf("No serviceIp %q in String(): %v", serviceIp, leaderlist.String())
		}
		if !strings.Contains(leaderlist.String(), servicePort) {
			t.Errorf("No servicePort %q in String(): %v", servicePort, leaderlist.String())
		}
		if !strings.Contains(leaderlist.String(), memberName) {
			t.Errorf("No memberName %q in String(): %v", memberName, leaderlist.String())
		}
		if !strings.Contains(leaderlist.String(), memberIp) {
			t.Errorf("No memberIp %q in String(): %v", memberIp, leaderlist.String())
		}
		if !strings.Contains(leaderlist.String(), memberPort) {
			t.Errorf("No memberPort %q in String(): %v", memberPort, leaderlist.String())
		}
		if !strings.Contains(leaderlist.String(), fmt.Sprintf("%v", memberStatus)) {
			t.Errorf("No memberStatus %v in String(): %v", memberStatus, leaderlist.String())
		}
	})
}

func TestSetterGetter(t *testing.T) {
	name := "testlist"
	serviceIp := "127.0.0.1"
	servicePort := "22365"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_UNKNOWN

	leaderlist, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// ServiceIp
	t.Run("ServiceIp", func(t *testing.T) {
		leaderlist.SetServiceIp("")
		if leaderlist.ServiceIp() != "" {
			t.Errorf("Unexpected result from ServiceIp(): %v", leaderlist.ServiceIp())
		}
	})

	// ServicePort
	var porttests = []struct {
		port  string
		valid bool
	}{
		{"22365", true},
		{"1024", true},
		{"65535", true},
		{"invalid", false},
		{"-1", false},
		{"0", false},
		{"1023", false},
		{"65536", false},
		{"1234567890", false},
	}
	for _, pt := range porttests {
		t.Run("ServicePort", func(t *testing.T) {
			leaderlist.SetServicePort(pt.port)
			if (leaderlist.ServicePort() == pt.port) != pt.valid {
				t.Errorf("Unexpected validity %v from ServicePort(%q): %v", pt.valid, pt.port, leaderlist.ServicePort())
			}
		})
	}

	// MemberStatus
	t.Run("MemberStatus", func(t *testing.T) {
		leaderlist.SetMemberStatus(leadlist.Leader_CANDIDATE)
		if leaderlist.MemberStatus() != leadlist.Leader_CANDIDATE {
			t.Errorf("Unexpected result from MemberStatus(): %v", leaderlist.MemberStatus())
		}
	})

	// LeaderAddress
	t.Run("No LeaderAddress", func(t *testing.T) {
		if leaderlist.LeaderAddress() != "" {
			t.Errorf("Unexpected result from LeaderAddress(): %v", leaderlist.LeaderAddress())
		}
	})

	// LeaderVersion
	t.Run("Initial LeaderVersion", func(t *testing.T) {
		if leaderlist.LeaderVersion() != 0 {
			t.Errorf("Unexpected result from LeaderVersion(): %v", leaderlist.LeaderVersion())
		}
	})
}

func TestService(t *testing.T) {
	name := "testlist"
	serviceIp := "127.0.0.1"
	servicePort := "22365"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// First RunService
	t.Run("First RunService", func(t *testing.T) {
		started, err := service.RunService()
		if (!started) || err != nil {
			t.Errorf("Unexpected result from RunService(): %v %v", started, err)
		}
	})

	// Second RunService
	t.Run("Second RunService", func(t *testing.T) {
		started, err := service.RunService()
		if started || err != nil {
			t.Errorf("Unexpected result from RunService(): %v %v", started, err)
		}
	})

	// Invalid RunService
	t.Run("Invalid RunService", func(t *testing.T) {
		service.SetServiceIp("invalid")
		started, err := service.RunService()
		if started || err == nil {
			t.Errorf("Unexpected result from RunService(): %v %v", started, err)
		}
	})
}

func TestLeaderAddress(t *testing.T) {
	name := "testlist"
	serviceIp := "127.0.0.1"
	servicePort := "22366"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// First RunService
	t.Run("First RunService", func(t *testing.T) {
		started, err := service.RunService()
		if (!started) || err != nil {
			t.Errorf("Unexpected result from RunService(): %v %v", started, err)
		}
	})
}

func TestSyncService(t *testing.T) {
	name := "TestSyncService"
	serviceIp := "127.0.0.1"
	servicePort := "22367"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
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
	memberPort = "12345"
	memberStatus = leadlist.Leader_CANDIDATE
	leader, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	if leader.LeaderVersion() != 0 {
		t.Errorf("Unexpected initial leader version: %v", leader.LeaderVersion())
	}

	_, err = leader.RunClient()
	if err != nil {
		t.Fatalf("could not run leader: %v", err)
	}

	err = leader.SyncService()
	if err != nil {
		log.Fatalf("could not sync leader and service: %v", err)
	}

	for leader.LeaderVersion() <= 0 {
		time.Sleep(time.Millisecond * 10)
	}

	if leader.LeaderVersion() != service.leaderVersion {
		t.Errorf("Unsyncronized second leader version: service %v != leader %v", service.LeaderVersion(), leader.leaderVersion)
	}

	if leader.LeaderAddress() != fmt.Sprintf("%v:%v", leader.member.Ip, leader.member.Port) {
		t.Errorf("Unexpected LeaderAddress(): %v", leader.LeaderAddress())
	}
}

func TestPingMember(t *testing.T) {
	name := "TestPingMember"
	serviceIp := "127.0.0.1"
	servicePort := "22368"
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := "12345"
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
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
	memberPort = "12346"
	memberStatus = leadlist.Leader_CANDIDATE
	leader, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	if leader.LeaderVersion() != 0 {
		t.Errorf("Unexpected initial leader version: %v", leader.LeaderVersion())
	}

	_, err = leader.RunClient()
	if err != nil {
		t.Fatalf("could not run leader: %v", err)
	}

	err = leader.SyncService()
	if err != nil {
		log.Fatalf("could not sync leader and service: %v", err)
	}

	for leader.LeaderVersion() <= 0 {
		time.Sleep(time.Millisecond * 10)
	}

	err = service.PingMember("alice")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = service.PingMember("invalid")
	switch err := err.(type) {
	case *MemberNotFound:
		break
	default:
		t.Errorf("Not the expected error type 'MemberNotFound': %v", err)
	}

	// For coverage
	_ = err.Error()
}

//t.Run("ServiceVersion", func(t *testing.T) {
//	if service.LeaderVersion() != 0 {
//		t.Errorf("Unexpected initial leader version: %v", service.LeaderVersion())
//	}
//})
//
//t.Run("leader", func(t *testing.T) {
//	leader, err := NewLeaderlist("testlist", "localhost", 22365, "alice",
//		"localhost", 12346, leadlist.Leader_CANDIDATE)
//	if err != nil {
//		log.Fatalf("could not create new leaderlist: %v", err)
//	}
//	t.Run("LeaderVersion", func(t *testing.T) {
//		if leader.LeaderVersion() != 0 {
//			t.Errorf("Unexpected initial leader version: %v", leader.LeaderVersion())
//		}
//	})
//})

//t.Run("group", func(t *testing.T) {
//	t.Run("Test1", func(t *testing.T) {
//
//	})
//	t.Run("Test2", parallelTest2)
//	t.Run("Test3", parallelTest3)
//})
