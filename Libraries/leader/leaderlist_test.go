package leader

import (
	"log"
	"testing"

	"bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

func TestServiceLeaderlist(t *testing.T) {

	testleaders, err := NewLeaderlist("testlist", "localhost", 22365, "alice",
		"localhost", 12345, leaderlist.Leader_CANDIDATE)
	if err != nil {
		log.Fatalf("could not create new leaderlist: %v", err)
	}
	if testleaders.String() != "" {

		t.Errorf("test error!!!: %v", testleaders)
	}
}
