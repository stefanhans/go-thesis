package main

import (
	"fmt"
	"net"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/proto-rootserver-tcp/member-group"
	"github.com/golang/protobuf/proto"
)

func handleMemberRequest(data []byte, remoteAddr net.Addr) error {

	fmt.Printf("remoteAddr: %v\n", remoteAddr)

	var newMember membergroup.Member
	if err := proto.Unmarshal(data, &newMember); err != nil {
		return fmt.Errorf("could not unmarshall member request: %v", err)
	}
	fmt.Printf("Member received (%v byte): %v\n", newMember.XXX_Size(), newMember)

	mbStorage = append(mbStorage, &newMember)

	return nil
}
