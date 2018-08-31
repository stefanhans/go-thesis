package leader

import (
	"fmt"
	"log"
	"net"
	"testing"

	leadlist "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

func TestSimpleNewLeaderlist(t *testing.T) {
	name := "testlist"
	serviceIp := "127.0.0.1"
	servicePort := 22365
	memberName := "service"
	memberIp := "127.0.0.1"
	memberPort := 12345
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// name
	t.Run("name", func(t *testing.T) {
		if service.name != name {
			t.Errorf("Unexpected name: %v", service.name)
		}
	})

	// serviceIp
	t.Run("serviceIp", func(t *testing.T) {
		if service.serviceIp != serviceIp {
			t.Errorf("Unexpected serviceIp: %v", service.serviceIp)
		}
	})

	// servicePort
	t.Run("servicePort", func(t *testing.T) {

		if service.servicePort != servicePort {
			t.Errorf("Unexpected servicePort: %v", service.servicePort)
		}
	})

	// memberName
	t.Run("memberName", func(t *testing.T) {

		if service.member.Name != memberName {
			t.Errorf("Unexpected memberName: %v", service.member.Name)
		}
	})

	// memberIp
	t.Run("memberIp", func(t *testing.T) {

		if service.member.Ip != memberIp {
			t.Errorf("Unexpected memberIp: %v", service.member.Ip)
		}
	})

	// memberPort
	t.Run("memberPort", func(t *testing.T) {

		if service.member.Port != fmt.Sprintf("%d", memberPort) {
			t.Errorf("Unexpected memberPort: %v", service.member.Port)
		}
	})

	// memberStatus
	t.Run("memberStatus", func(t *testing.T) {

		if service.member.Status != memberStatus {
			t.Errorf("Unexpected serviceIp: %v", service.member.Status)
		}
	})

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

}

func TestResolveIpNewLeaderlist(t *testing.T) {
	name := "testlist"
	serviceIp := "localhost"
	servicePort := 22365
	memberName := "service"
	memberIp := "localhost"
	memberPort := 12345
	memberStatus := leadlist.Leader_SERVICE

	service, err := NewLeaderlist(name, serviceIp, servicePort, memberName,
		memberIp, memberPort, memberStatus)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	// serviceIp
	t.Run("serviceIp", func(t *testing.T) {
		addr, _ := net.ResolveIPAddr("ip", serviceIp)
		if service.serviceIp != addr.String() {
			t.Errorf("Unexpected serviceIp: %v", service.serviceIp)
		}
	})

	// memberIp
	t.Run("memberIp", func(t *testing.T) {
		addr, _ := net.ResolveIPAddr("ip", memberIp)

		if service.member.Ip != addr.String() {
			t.Errorf("Unexpected memberIp: %v", service.member.Ip)
		}
	})
}
