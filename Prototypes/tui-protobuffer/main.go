package main

import (
	"log"

	"fmt"
	"github.com/jroimartin/gocui"
)

var (
	msg             string
	nickname        string
	initialMessages string
)

func init() {
	initialMessages = "init()"
}

func main() {

	nickname = "me"

	// Open file for reading info
	byteArray, _ := OpenStorageForRead("storage")

	fmt.Printf("1: %q\n", initialMessages)

	initialMessages = readMessage(byteArray)
	fmt.Printf("2: %q\n", initialMessages)

	// Create the terminal GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	// Set function to manage all views and keybindings
	g.SetManagerFunc(layout)

	// Bind keys with functions
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send)

	// Start main event loop of the GUI
	g.MainLoop()
}
