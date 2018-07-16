package main

import (
	"fmt"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-tcp/member-group"
	"github.com/golang/protobuf/proto"
)

func handleMemberListRequest(data []byte) error {

	var receivedMembers membergroup.MemberList
	if err := proto.Unmarshal(data, &receivedMembers); err != nil {
		return fmt.Errorf("could not unmarshall memberlist request: %v", err)
	}
	fmt.Printf("MemberList received (%v byte): %v\n", receivedMembers.XXX_Size(), receivedMembers)

	return nil
}
