package main

import (
	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
)

const (

	// Publishing service on a commonly known address
	serverIp          string = "localhost"
	serverPort        string = "22365"
	publishingService string = serverIp + ":" + serverPort

	// Switch debugging
	debug bool = true
)

var (

	// Application identity set by command args
	memberName        string
	memberIp          string
	memberPort        string
	displayingService string

	cgMember []*chatgroup.Member
)
