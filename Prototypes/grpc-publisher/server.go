package main

import (
	"bitbucket.org/stefanhans/go-thesis/Prototypes/grpc-publisher/memberlist"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Create and register server
	var members memberlist.MembersServer
	srv := grpc.NewServer()
	memberlist.RegisterMembersServer(srv, members)

	// Create listener
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("could not listen to :8888: \v", err)
	}
	// Serve messages via listener
	log.Fatal(srv.Serve(l))
}

// Receiver for implementing the server service interface
type memberServer struct{}

// Server's Subscribe implementation
//func (ms memberServer) Subscribe(ctx context.Context, member *memberlist.Member) (*memberlist.Void, error) {
//	memberStore = append(memberStore, member)
//	return &memberlist.Void{}, nil
//}

// Server's List implementation
func (ms memberServer) List(ctx context.Context, void *memberlist.Void) (*memberlist.MemberList, error) {

	var members memberlist.MemberList
	//mb := memberlist.Member{ Name: "me", Ip: "ip", Port: "port", Leader: true}
	//members.List = append(members.List, &mb)
	////for _, member := range memberStore {
	////
	////	members.List = append(members.List, member)
	////}
	return &members, nil
}
