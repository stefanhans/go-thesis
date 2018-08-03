package main

import (
	"net"
	"os"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
)

const (

	// Publishing service on a commonly known address
	//serverIp          string = "192.168.1.126"

	serverIp          string = "localhost"
	serverPort        string = "22365"
	publishingService string = serverIp + ":" + serverPort

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
