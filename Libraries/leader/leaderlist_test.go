package leader

import (
	"log"
	"testing"

	"bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
	"time"
	"fmt"
)

//func TestServiceLeaderlist(t *testing.T) {
//
//	testleaders, err := NewLeaderlist("testlist", "localhost", 22365, "alice",
//		"localhost", 12345, leaderlist.Leader_CANDIDATE)
//	if err != nil {
//		log.Fatalf("could not create new leaderlist: %v", err)
//	}
//	if testleaders.String() != "" {
//
//		t.Errorf("test error!!!: %v", testleaders)
//	}
//}



func TestServiceAndLeader(t *testing.T) {

	service, err := NewLeaderlist("testlist", "localhost", 22365, "service",
		"localhost", 12345, leaderlist.Leader_SERVICE)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	if service.LeaderVersion() != 0 {
		t.Errorf("Unexpected initial leader version: %v", service.LeaderVersion())
	}

	service.SetMemberStatus(leaderlist.Leader_SERVICE)
	_, err = service.RunService()
	if err != nil {
		t.Fatalf("could not run service: %v", err)
	}

	if service.LeaderVersion() != 0 {
		t.Errorf("Unexpected first leader version: %v", service.LeaderVersion())
	}

	leader, err := NewLeaderlist("testlist", "localhost", 22365, "alice",
		"localhost", 12346, leaderlist.Leader_CANDIDATE)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}

	if leader.LeaderVersion() != 0 {
		t.Errorf("Unexpected initial leader version: %v", leader.LeaderVersion())
	}

	leader.SetMemberStatus(leaderlist.Leader_CANDIDATE)
	_, err = leader.RunClient()
	if err != nil {
		t.Fatalf("could not run leader: %v", err)
	}

	err = leader.SyncService()
	if err != nil {
		log.Fatalf("could not sync leader and service: %v", err)
	}


	fmt.Printf("start waiting\n")

	for leader.LeaderVersion() <= 0 {
		time.Sleep(time.Millisecond*10)
	}

	if leader.LeaderVersion() != service.leaderVersion {

		t.Errorf("Unsyncronized second leader version: service %v != leader %v", service.LeaderVersion(), leader.leaderVersion)
	}
}