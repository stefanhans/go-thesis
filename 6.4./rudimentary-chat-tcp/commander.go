package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/6.4./rudimentary-chat-tcp/chat-group"
	// "github.com/golang/protobuf/proto"
)

var (
	cmdUsage map[string]string
	keys     []string
)

func commandUsageInit() {
	cmdUsage = make(map[string]string)

	cmdUsage["list"] = "\\list"
	cmdUsage["logfile"] = "\\logfile"
	cmdUsage["publisher"] = "\\publisher"
	cmdUsage["self"] = "\\self"
	cmdUsage["sync"] = "\\sync"
	cmdUsage["echo"] = "\\echo"

	// To store the keys in sorted order
	for key := range cmdUsage {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	log.Printf("commandUsageInit: keys: %v\n", keys)
}

// Execute a command specified by the argument string
func executeCommand(commandline string) {

	// Trim prefix and split string by white spaces
	commandFields := strings.Fields(strings.Trim(commandline, "\\"))

	// Check for empty string without prefix
	if len(commandFields) > 0 {
		log.Printf("Command: %q\n", commandFields[0])
		log.Printf("Arguments (%v): %v\n", len(commandFields[1:]), commandFields[1:])

		// Switch according to the first word and call appropriate function with the rest as arguments
		switch commandFields[0] {

		case "list":
			log.Printf("CMD_LIST\n")
			list(commandFields[1:])

		case "logfile":
			log.Printf("CMD_LOGFILE\n")
			showLogfile(commandFields[1:])

		case "publisher":
			log.Printf("CMD_PUBLISHER\n")
			publisher(commandFields[1:])

		case "self":
			log.Printf("CMD_SELF\n")
			self(commandFields[1:])

		case "sync":
			log.Printf("CMD_SYNC\n")
			syncMemberlist(commandFields[1:])

		case "echo":
			log.Printf("CMD_ECHO\n")
			echoFromMemberlist(commandFields[1:])

		default:
			usage()
		}

	} else {
		usage()
	}
}

// Display the usage of all available commands
func usage() {
	// todo: order not deterministic bug
	for _, key := range keys {
		displayText(fmt.Sprintf("<CMD USAGE>: %s", cmdUsage[key]))
	}
}

func list(arguments []string) {

	if len(selfMemberList) == 0 {
		displayText(fmt.Sprintf("<CMD_LIST>: no member found - call \\sync"))
		return
	}

	text := ""
	for i, member := range selfMemberList {
		text += fmt.Sprintf("<CMD_LIST>: %v: %s\n", i, member)
	}
	displayText(strings.Trim(text, "\n"))
}

func showLogfile(arguments []string) {

	displayText(fmt.Sprintf("<CMD_LOGFILE>: %v", logfilename))
}

func publisher(arguments []string) {

	displayText(fmt.Sprintf("<CMD_PUBLISHER>: %v", leaderService))
}

func self(arguments []string) {

	displayText(fmt.Sprintf("<CMD_SELF>: %v", selfMember))
}

func syncMemberlist(arguments []string) {

	if selfMember.Leader {
		displayText(fmt.Sprintf("<CMD_SYNC>: requestor provides publishing service itself - call \\list"))
		return
	}

	message := &chatgroup.Message{
		MsgType:    chatgroup.Message_MEMBERLIST_REQUEST,
		Sender:     selfMember,
		MemberList: &chatgroup.MemberList{},
	}

	err := sendPublisherRequest(message)
	if err != nil {
		displayText(fmt.Sprintf("<CMD_SYNC_ERR>: %v", err))
		return
	}
	displayText(fmt.Sprintf("<CMD_SYNC_REQUEST>: request sent to %v", leaderService))
}

func echoFromMemberlist(arguments []string) {


	message := &chatgroup.Message{
		MsgType:    chatgroup.Message_MEMBERLIST_REQUEST,
		Sender:     selfMember,
		MemberList: &chatgroup.MemberList{},
	}


	for i, member := range selfMemberList {

		// Send message to requestor
		err := sendMessage(message, member.Ip+":"+member.Port)
		if err != nil {
			displayText(fmt.Sprintf("<CMD_ECHO_ERR> %d: failed send request to %v: %v", i, member, err))
		} else {
			displayText(fmt.Sprintf("<CMD_ECHO> %d: request %q sent to %v",
				i, strings.Join(arguments, " "), member))
		}
	}

}
