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

// Content to be displayed in the GUI
func layout(g *gocui.Gui) error {

	// Get size of the terminal
	maxX, maxY := clientGui.Size()

	// Creates view "messages"
	if messages, err := clientGui.SetView("messages", 0, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Autoscroll = true
		messages.Wrap = true
	}

	// Creates view "input"
	if input, err := clientGui.SetView("input", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Wrap = true
		input.Editable = true
	}

	// Set view "input" as the current view with focus and cursor
	if _, err := clientGui.SetCurrentView("input"); err != nil {
		return err
	}

	// Show cursor
	clientGui.Cursor = true

	return nil
}

// Quit the GUI
func quit(g *gocui.Gui, v *gocui.View) error {
	unsubscribeClient(memberName)
	return gocui.ErrQuit
}

// Send content from the bottom window to the top window
func send(g *gocui.Gui, v *gocui.View) error {

	// Get the top window view and write the buffer of the bottom window view to it
	if m, err := clientGui.View("messages"); err != nil {
		log.Fatal(err)
	} else {
		sendMessage(memberName, strings.Trim(v.Buffer(), "\n"))
		m.Write([]byte(fmt.Sprintf("%s: %s", memberName, v.Buffer())))
	}

	// Clear the bottom window and reset the cursor
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		log.Fatal(err)
	}

	return nil
}

func startTui() error {
	var err error

	// Create the terminal GUI
	clientGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return fmt.Errorf("could start tui: %v\n", err)
	}
	defer clientGui.Close()

	// Set function to manage all views and keybindings
	clientGui.SetManagerFunc(layout)

	// Bind keys with functions
	clientGui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	clientGui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send)

	// Start main event loop of the GUI
	return clientGui.MainLoop()
}

func displayText(txt string) error {

	messagesView, _ := clientGui.View("messages")
	clientGui.Update(func(g *gocui.Gui) error {
		fmt.Fprintln(messagesView, txt)
		return nil
	})
	return nil
}
