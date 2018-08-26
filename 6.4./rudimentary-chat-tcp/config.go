package main

import (
	"os"

	"bitbucket.org/stefanhans/go-thesis/6.4./rudimentary-chat-tcp/chat-group"
)

const (

	// Publishing service on a commonly known address
	//serverIp          string = "192.168.1.126"

	serverIp      string = "localhost"
	serverPort    string = "22365"
	leaderService string = serverIp + ":" + serverPort

	// Switch debugging
	debug bool = true
)

var (

	// Application identity set by command args
	displayingService string

	// Client as member
	selfMember *chatgroup.Member

	// Storage for member list of chat group
	selfMemberList []*chatgroup.Member

	//
	logfilename string
	logfile     *os.File
)
