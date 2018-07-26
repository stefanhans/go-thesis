package main

import (
	"fmt"
	"log"
	"strings"
)

var (
	cmdUsage map[string]string
)

func commandUsageInit() {
	cmdUsage = make(map[string]string)

	cmdUsage["clear"] = "\\clear"
	cmdUsage["list"] = "\\list"
	cmdUsage["publisher"] = "\\publisher [details]"
	cmdUsage["self"] = "\\self"
}

func executeCommand(commandline string) {

	commandFields := strings.Fields(strings.Trim(commandline, "\\"))

	if len(commandFields) > 0 {
		log.Printf("Command: %q\n", commandFields[0])
		log.Printf("Arguments (%v): %v\n", len(commandFields[1:]), commandFields[1:])

		switch commandFields[0] {

		case "clear":
			clear(commandFields[1:])

		case "list":
			list(commandFields[1:])

		case "publisher":
			publisher(commandFields[1:])

		case "self":
			self(commandFields[1:])

		default:
			for _, value := range cmdUsage {

				displayText(fmt.Sprintf("<CMD USAGE>: %s", value))
			}
		}

	} else {
		for _, value := range cmdUsage {

			displayText(fmt.Sprintf("<CMD USAGE>: %s", value))
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

			displayText(fmt.Sprintf("<PUBLISHER>: %v", selfMember))
		default:

			displayText(fmt.Sprintf("<PUBLISHER>: Usage: %s", cmdUsage["publisher"]))
		}

	} else {

		displayText(fmt.Sprintf("<PUBLISHER>: %v", displayingService))
	}
	return nil
}

func self(arguments []string) error {
	displayText(fmt.Sprintf("<SELF>: %v", selfMember))

	return nil
}
