package main

import (
	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
)

const (

	// Publishing service on a commonly known address
	serverIp          string = "192.168.1.126"
	serverPort        string = "22365"
	publishingService string = serverIp + ":" + serverPort

	// Switch debugging
	debug bool = true
)

var (

	// Application identity set by command args
	//memberName        string
	//memberIp          string
	//memberPort        string
	displayingService string
	selfMember        *chatgroup.Member

	// Publisher storage for member of chat group
	cgMember []*chatgroup.Member
)
