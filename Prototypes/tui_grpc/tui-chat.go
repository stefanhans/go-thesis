package main

import (
	"log"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/tui-protobuffer/chat-message-pb"
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

		messages.Write([]byte(initialMessages))
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

	file, _ := openStorageToWrite("storage")

	if m, err := g.View("messages"); err != nil {
		log.Fatal(err)
	} else {
		writeMessage(&chatmessage.Chatmessage{
			Text: m.Buffer(),
			From: "",
		}, file)
	}

	// Close file
	closeStorage(file)

	return gocui.ErrQuit
}

// Send content from the bottom window to the top window
func send(g *gocui.Gui, v *gocui.View) error {
	msg := chatmessage.Chatmessage{
		Text: v.Buffer(),
		From: nickname,
	}

	// Get the top window view and write the buffer of the bottom window view to it
	if m, err := g.View("messages"); err != nil {
		log.Fatal(err)
	} else {
		if len(msg.Text) != 0 {
			m.Write([]byte(" " + msg.From + ": " + msg.Text))
		}
	}

	// Clear the bottom window and reset the cursor
	v.Clear()
	if err := v.SetCursor(0, 0); err != nil {
		log.Fatal(err)
	}

	// Send message to other clients

	return nil
}

//func initMessages(g *gocui.Gui, text string) {
//
//
//	// Get the top window view and write the buffer of the bottom window view to it
//	if m, err := gui.View("messages"); err != nil {
//		log.Fatal(err)
//	} else {
//		m.Write([]byte(text))
//	}
//}
