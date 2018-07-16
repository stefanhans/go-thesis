package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-rootserver/member-group"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var memberlist membergroup.MemberList

func main() {
	// Create server
	srv := grpc.NewServer()

	// Register server
	var members memberServer
	membergroup.RegisterMembersServer(srv, members)

	// Create listener
	l, err := net.Listen("tcp", ":22365")
	if err != nil {
		log.Fatal("could not listen to :22365: \v", err)
	}
	fmt.Printf("Rootserver does listen on localhost:22365\n")

	// Serve messages via listener
	log.Fatal(srv.Serve(l))
}

// Receiver for implementing the server service interface
type memberServer struct{}

// Server's Subscribe implementation
func (s memberServer) Register(ctx context.Context, requester *membergroup.Member) (*membergroup.Member, error) {
	fmt.Printf("Request for registration: %v\n", requester)

	exists := false

	for _, registeredMember := range memberlist.Member {
		if registeredMember.Name == requester.Name {
			// TODO Update registeredMember
			exists = true
		}
		if exists {
			break
		}
	}
	if !exists {
		memberlist.Member = append(memberlist.Member, requester)
		fmt.Printf("Registration: %v\n", requester)
	}

	fmt.Printf("Registered Leaders and Candidates:\n")
	for i, registeredMember := range memberlist.Member {
		fmt.Printf("%v\tRegistered: %v\n", i+1, registeredMember)
	}
	return requester, nil
}

// Server's List implementation
func (s memberServer) List(ctx context.Context, void *membergroup.Void) (*membergroup.MemberList, error) {

	return &memberlist, nil
}
