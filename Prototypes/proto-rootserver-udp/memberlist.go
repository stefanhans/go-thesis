package main

import (
	"fmt"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-udp/member-group"
	"github.com/golang/protobuf/proto"
	"net"
)

func handleMemberListRequest(data []byte, remoteAddr net.Addr)  ([]byte, error) {

	var receivedMembers membergroup.MemberList
	if err := proto.Unmarshal(data, &receivedMembers); err != nil {
		return nil, fmt.Errorf("could not unmarshall memberlist request: %v", err)
	}
	fmt.Printf("MemberList (%v byte) received from %s: %v\n", receivedMembers.XXX_Size(), remoteAddr, receivedMembers)


	return []byte("ToDo: implement reply"), nil
}
