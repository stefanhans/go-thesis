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

// Create the text-based UI.
// Afterwards we can use clientGui.Update or displayText, respectively,
// to display change views as soon as TUI is running
func initTUI() error {
	var err error

	// Create the TUI
	clientGui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return fmt.Errorf("could not create tui: %v\n", err)
	}
	defer clientGui.Close()

	return err
}

// Configure and run the text-based UI
func runTUI() error {

	// Set function to manage all views and keybindings
	clientGui.SetManagerFunc(layout)

	// Bind keys with functions
	clientGui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	clientGui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, send)
	// todo fix up/down or remove it
	clientGui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, up)
	clientGui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, down)

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
	Unsubscribe(selfMember.Name)

	return gocui.ErrQuit
}

// Send content from the view "input"
func send(g *gocui.Gui, inputView *gocui.View) error {

	input := strings.Trim(inputView.Buffer(), "\n")

	// Distinguish between command and chat mode by '\'-prefix
	if strings.HasPrefix(input, "\\") {

		// Interpret "input" as command
		executeCommand(input)

	} else {
		// Send "input" to publish
		Publish(input)
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

	// Update the "messages" view as soon as possible
	clientGui.Update(func(g *gocui.Gui) error {
		messagesView, err := clientGui.View("messages")
		if err != nil {
			return fmt.Errorf("could not display text: %v\n", err)
		}
		fmt.Fprintln(messagesView, txt)
		return nil
	})
	return nil
}

func up(g *gocui.Gui, inputView *gocui.View) error {

	messagesView, err := clientGui.View("messages")
	if err != nil {
		return fmt.Errorf("could not move messages down: %v\n", err)
	}
	x, y := messagesView.Origin()
	messagesView.SetOrigin(x, y-1)
	x2, y2 := messagesView.Origin()

	log.Printf("Moved from %v/%v to %v/%v\n", x, y, x2, y2)
	return nil
}

func down(g *gocui.Gui, inputView *gocui.View) error {

	messagesView, err := clientGui.View("messages")
	if err != nil {
		return fmt.Errorf("could not move messages down: %v\n", err)
	}
	x, y := messagesView.Origin()
	messagesView.SetOrigin(x, y+1)
	x2, y2 := messagesView.Origin()

	log.Printf("Moved from %v/%v to %v/%v\n", x, y, x2, y2)
	return nil
}
