package leader

import "bitbucket.org/stefanhans/go-thesis/Libraries/leader/leaderlist"

type Config struct {
	Leaderlist string

	Member string
	Ip     string
	Port   int
	Status leaderlist.Leader_LeaderStatus

	// TCP or UDP
	Transport string

	// Todo: Implement CanServiceBeLeader
	CanServiceBeLeader bool
}
