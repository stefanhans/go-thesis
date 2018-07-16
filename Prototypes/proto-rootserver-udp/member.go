package main

import (
	"fmt"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-udp/member-group"
	"github.com/golang/protobuf/proto"
)

func handleMemberRequest(data []byte, remoteAddr net.Addr)  ([]byte, error) {

	var newMember membergroup.Member
	if err := proto.Unmarshal(data, &newMember); err != nil {
		return nil, fmt.Errorf("could not unmarshall member request: %v", err)
	}
	fmt.Printf("Member (%v byte) received from %s: %v\n", newMember.XXX_Size(), remoteAddr, newMember)

	mbStorage = append(mbStorage, &newMember)

	return []byte("ToDo: implement reply"), nil
}
