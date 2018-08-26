package leader

import (
	"log"
	"testing"

	"bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"
)

var (
	serviceLeaderlist = "testlist"
	serviceMember     = "service"
	serviceIp         = "127.0.0.1"
	servicePort       = 22365

	testleaders *Leaderlist

	err error
)

func TestServiceLeaderlist(t *testing.T) {
	config := &Config{
		Leaderlist: serviceLeaderlist,
		Member:     serviceMember,
		Ip:         serviceIp,
		Port:       servicePort,
		Status:     leaderlist.Leader_SERVICE,
		Transport:  "TCP",
	}

	testleaders, err = NewLeaderlist(config)
	if err != nil {
		log.Fatalf("could not create NewLeaderlist: %v", err)
	}
	if testleaders.String() != "" {

		t.Errorf("test error!!!: %v", testleaders)
	}
}
