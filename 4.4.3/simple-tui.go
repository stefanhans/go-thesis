package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

// Content to be displayed in the GUI
func layout(g *gocui.Gui) error {

	// Get size of the terminal
	maxX, maxY := g.Size()

	// Creates view "messages"
	if messages, err := g.SetView("messages", 0, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Autoscroll = true
		messages.Wrap = true
	}

	// Creates view "input"
	if input, err := g.SetView("input", 0, maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Wrap = true
		input.Editable = true
	}

	// Set view "input" as the current view with focus and cursor
	if _, err := g.SetCurrentView("input"); err != nil {
		return err
	}

	// Show cursor
	g.Cursor = true

	return nil
}

// Quit the GUI
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Send content from the bottom window to the top window
func send(g *gocui.Gui, v *gocui.View) error {

	// Get the top window view and write the buffer of the bottom window view to it
	if m, err := g.View("messages"); err != nil {
		log.Fatal(err)
	} else {
		m.Write([]byte(v.Buffer()))
	}

	// Clear the bottom window and reset the cursor
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {

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
