package main

import (
	"fmt"
	"log"
	"strings"

	"bitbucket.org/stefanhans/go-thesis/6.3./rudimentary-chat-tcp/chat-group"
	"sort"
)

var (
	cmdUsage map[string]string
	keys     []string
)

func commandUsageInit() {
	cmdUsage = make(map[string]string)

	cmdUsage["clear"] = "\\clear"
	cmdUsage["list"] = "\\list"
	cmdUsage["publisher"] = "\\publisher [details]"
	cmdUsage["self"] = "\\self"
	cmdUsage["subscribe"] = "\\subscribe <name> <ip> <port>"
	cmdUsage["unsubscribe"] = "\\unsubscribe <name>"

	// To store the keys in sorted order
	for key := range cmdUsage {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	log.Printf("commandUsageInit: keys: %v\n", keys)
}

func executeCommand(commandline string) {

	commandFields := strings.Fields(strings.Trim(commandline, "\\"))

	if len(commandFields) > 0 {
		log.Printf("Command: %q\n", commandFields[0])
		log.Printf("Arguments (%v): %v\n", len(commandFields[1:]), commandFields[1:])

		switch commandFields[0] {

		case "clear":
			log.Printf("CMD_CLEAR\n")
			clear(commandFields[1:])

		case "list":
			log.Printf("CMD_LIST\n")
			list(commandFields[1:])

		case "publisher":
			log.Printf("CMD_PUBLISHER\n")
			publisher(commandFields[1:])

		case "self":
			log.Printf("CMD_SELF\n")
			self(commandFields[1:])

		case "subscribe":
			log.Printf("CMD_SUBSCRIBE\n")
			subscribe(commandFields[1:])

		case "unsubscribe":
			log.Printf("CMD_UNSUBSCRIBE\n")
			unsubscribe(commandFields[1:])

		default:
			for _, value := range cmdUsage {

				displayText(fmt.Sprintf("<CMD USAGE>: %s", value))
			}
		}

	} else {
		for _, key := range keys {
			displayText(fmt.Sprintf("<CMD USAGE>: %s", cmdUsage[key]))
		}
	}
}

func clear(arguments []string) error {

	messagesView, _ := clientGui.View("messages")
	messagesView.Clear()
	return nil
}

func list(arguments []string) error {

	err := List()
	if err != nil {
		log.Printf("could not request \\list: %v\n", err)
	}
	return nil
}

func publisher(arguments []string) error {
	if len(arguments) > 0 {

		switch arguments[0] {
		case "details":

			displayText(fmt.Sprintf("<CMD_PUBLISHER>: %v", selfMember))
		default:

			displayText(fmt.Sprintf("<CMD_PUBLISHER>: Usage: %s", cmdUsage["publisher"]))
		}

	} else {

		displayText(fmt.Sprintf("<CMD_PUBLISHER>: %v", displayingService))
	}
	return nil
}

func self(arguments []string) error {
	displayText(fmt.Sprintf("<CMD_SELF>: %v", selfMember))

	return nil
}

func subscribe(arguments []string) error {
	if len(arguments) >= 3 {
		newMember := &chatgroup.Message{
			MsgType: chatgroup.Message_SUBSCRIBE,
			Sender:  &chatgroup.Member{Name: arguments[0], Ip: arguments[1], Port: arguments[2]}}

		sendPublisherRequest(newMember)

		// Append subscription message in "messages" view
		displayText(fmt.Sprintf("<CMD_SUBSCRIBE>: %s (%s:%s) has joined", newMember.Sender.Name, newMember.Sender.Ip, newMember.Sender.Port))
	} else {
		displayText(fmt.Sprintf("<CMD_SUBSCRIBE>: Usage: %s", cmdUsage["subscribe"]))
	}
	return nil
}

func unsubscribe(arguments []string) error {
	if len(arguments) > 0 {
		leavingMember := &chatgroup.Message{
			MsgType: chatgroup.Message_UNSUBSCRIBE,
			Sender:  &chatgroup.Member{Name: arguments[0]}}

		sendPublisherRequest(leavingMember)

		// Append subscription message in "messages" view
		displayText(fmt.Sprintf("<CMD_UNSUBSCRIBE>: %s (%s:%s) has left", leavingMember.Sender.Name))
	} else {
		displayText(fmt.Sprintf("<CMD_UNSUBSCRIBE>: Usage: %s", cmdUsage["unsubscribe"]))
	}
	return nil
}
