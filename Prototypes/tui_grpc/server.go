package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/tui_grpc/chat-message-pb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// Receiver for implementing the server service interface
type chatMessageServer struct{}

// Server's Write implementation
func (s chatMessageServer) Write(ctx context.Context, info *chatmessage.ChatMessage) (*chatmessage.ChatMessage, error) {

	// Marshall message
	b, err := proto.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("could not encode info: %v", err)
	}

	// Open file
	f, err := os.OpenFile("storage", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open storage: %v", err)
	}

	// Encode message and write to file
	if err := binary.Write(f, binary.LittleEndian, int64(len(b))); err != nil {
		return nil, fmt.Errorf("could not encode length of message: %v", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write info to file: %v", err)
	}

	// Close file
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file storage: %v", err)
	}
	return info, nil
}

// Server's Read implementation
func (s chatMessageServer) Read(ctx context.Context, void *chatmessage.Void) (*chatmessage.ChatMessageList, error) {

	// Read file
	b, err := ioutil.ReadFile("storage")
	if err != nil {
		return nil, fmt.Errorf("could not read storage: %v", err)
	}

	// Iterate over read bytes
	var chatmessages chatmessage.ChatMessageList
	for {
		if len(b) == 0 {
			// Return result
			return &chatmessages, nil
		} else if len(b) < 8 {
			return nil, fmt.Errorf("remaining odd %d bytes", len(b))
		}

		// Decode message
		var length int64
		if err := binary.Read(bytes.NewReader(b[:8]), binary.LittleEndian, &length); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[8:]

		// Unmarshall message and append it
		var info chatmessage.ChatMessage
		if err := proto.Unmarshal(b[:length], &info); err != nil {
			return nil, fmt.Errorf("could not read info: %v", err)
		}
		b = b[length:]
		chatmessages.Chatmessages = append(chatmessages.Chatmessages, &info)
	}
}
