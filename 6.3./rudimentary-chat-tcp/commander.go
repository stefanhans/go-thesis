package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
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

	if selfMember.Leader {

		text := ""
		for i, member := range cgMember {
			text += fmt.Sprintf("<CMD_LIST>: %v: %s\n", i, member)
		}
		displayText(strings.Trim(text, "\n"))
	} else {
		displayText("<CMD_LIST>Error: Only publishing service can tell!")
	}
}

func showLogfile(arguments []string) {

	displayText(fmt.Sprintf("<CMD_LOGFILE>: %v", logfilename))
}

func publisher(arguments []string) {

	displayText(fmt.Sprintf("<CMD_PUBLISHER>: %v", publishingService))
}

func self(arguments []string) {

	displayText(fmt.Sprintf("<CMD_SELF>: %v", selfMember))
}
