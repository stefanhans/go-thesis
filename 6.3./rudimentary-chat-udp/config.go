package main

import (
	"net"
	"os"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-udp/chat-group"
)

const (

	// Publishing service on a commonly known address
	//serverIp          string = "192.168.1.126"

	serverIp          string = "127.0.0.1"
	serverPort        string = "22365"
	publishingService string = serverIp + ":" + serverPort

	// The maximum safe UDP payload is 508 bytes.
	// This is a packet size of 576 (IPv4 minimum reassembly buffer size),
	// minus the maximum 60-byte IP header and the 8-byte UDP header.
	bufferSize = 508

	// Switch debugging
	debug bool = true
)

var (

	// Application identity set by command args
	displayingService string
	selfMember        *chatgroup.Member

	// Publisher storage for member of chat group
	// todo refactor chatgroup.memberlist instead of []*chatgroup.Member
	cgMember []*chatgroup.Member

	//
	logfilename string
	logfile     *os.File
)

var requestActionMap = map[chatgroup.Message_MessageType]func(*chatgroup.Message, net.Addr) error{
	chatgroup.Message_SUBSCRIBE_REQUEST:   handleSubscribeRequest,
	chatgroup.Message_UNSUBSCRIBE_REQUEST: handleUnsubscribeRequest,
	chatgroup.Message_PUBLISH_REQUEST:     handlePublishRequest,
}

var replyActionMap = map[chatgroup.Message_MessageType]func(*chatgroup.Message) error{
	chatgroup.Message_SUBSCRIBE_REPLY:   handleSubscribeReply,
	chatgroup.Message_UNSUBSCRIBE_REPLY: handleUnsubscribeReply,
	chatgroup.Message_PUBLISH_REPLY:     handlePublishReply,
}
