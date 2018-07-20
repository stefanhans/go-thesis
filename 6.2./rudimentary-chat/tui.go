package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var (
	clientGui *gocui.Gui
)

// Run the text-based UI
func runTUI() error {
	var err error

	// Create the TUI
	clientGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return fmt.Errorf("could not create tui: %v\n", err)
	}
	defer clientGui.Close()

	// Set function to manage all views and keybindings
	clientGui.SetManagerFunc(layout)

	// Bind keys with functions
	clientGui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	clientGui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send)

	// Start main event loop of the TUI
	return clientGui.MainLoop()
}

// Content to be displayed in the TUI
func layout(g *gocui.Gui) error {

	// Get size of the terminal
	maxX, maxY := clientGui.Size()

	// Creates view "messages"
	messages, err := clientGui.SetView("messages", 0, 0, maxX-1, maxY-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Autoscroll = true
		messages.Wrap = true
	}

	// Creates view "input"
	input, err := clientGui.SetView("input", 0, maxY-4, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Wrap = true
		input.Editable = true
	}

	// Set view "input" as the current view with focus and cursor
	_, err = clientGui.SetCurrentView("input")
	if err != nil {
		return err
	}

	// Show cursor
	clientGui.Cursor = true

	return nil
}

// Quit the TUI
func quit(g *gocui.Gui, v *gocui.View) error {

	// Unsubscribe via PublisherClient
	Unsubscribe(memberName)

	return gocui.ErrQuit
}

// Send content from the view "input"
func send(g *gocui.Gui, inputView *gocui.View) error {

	// Get the "messages"
	m, err := clientGui.View("messages")
	if err != nil {
		log.Fatal(err)
	} else {

		// Send "input" to Publisher
		Publish(memberName, strings.Trim(inputView.Buffer(), "\n"))

		// Write "input" to "messages"
		m.Write([]byte(fmt.Sprintf("%s: %s", memberName, inputView.Buffer())))
	}

	// Clear the "input" and reset the cursor
	inputView.Clear()
	if err := inputView.SetCursor(0, 0); err != nil {
		log.Fatal(err)
	}
	return nil
}

// Display text in "messages"
func displayText(txt string) error {

	messagesView, _ := clientGui.View("messages")
	clientGui.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(messagesView, txt)
		return nil
	})
	return nil
}
