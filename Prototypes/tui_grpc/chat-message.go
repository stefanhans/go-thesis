package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/stefanhans/go-thesis/Prototypes/tui-protobuffer/chat-message-pb"
	"github.com/golang/protobuf/proto"
)

// Declare protobuffer message
func declareExample() chatmessage.Chatmessage {
	return chatmessage.Chatmessage{
		Text: msg,
		From: nickname,
	}
}

// Open file for appending info
func openStorageToWrite(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("could not open %s: %v", filename, err)
	}
	return file, nil
}

func OpenStorageForRead(filename string) ([]byte, error) {

	// Open file for reading info
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("could not read %s: %v", filename, err)
	}
	return byteArray, nil
}

func writeMessage(chatMessage *chatmessage.Chatmessage, file *os.File) {
	// Marshal into binary format
	byteArray, err := proto.Marshal(chatMessage)
	if err != nil {
		fmt.Errorf("could not encode info: %v", err)
		os.Exit(1)
	}

	// Write binary representation
	if err := binary.Write(file, binary.LittleEndian, int64(len(byteArray))); err != nil {
		fmt.Errorf("could not encode length of message: %v", err)
	}

	// Write to file
	_, err = file.Write(byteArray)
	if err != nil {
		fmt.Errorf("could not write info to file: %v", err)
	}
}

func readMessage(byteArray []byte) string {
	out := "\n"
	for {
		// Check length of remaining bytes
		if len(byteArray) == 0 {
			break
		} else if len(byteArray) < 8 {
			fmt.Errorf("remaining odd %d bytes, what to do?", len(byteArray))
		}

		// Decode binary data and shift array forward
		var length int64
		if err := binary.Read(bytes.NewReader(byteArray[:8]), binary.LittleEndian, &length); err != nil {
			fmt.Errorf("could not decode message length: %v", err)
		}
		byteArray = byteArray[8:]

		// Unmarshall info
		var chatmessage chatmessage.Chatmessage
		if err := proto.Unmarshal(byteArray[:length], &chatmessage); err != nil {
			fmt.Errorf("could not read info: %v", err)
		}
		byteArray = byteArray[length:]

		//out = fmt.Sprint("\n  %q %q", chatmessage.From, chatmessage.Text)
		if len(chatmessage.Text) != 0 {
			out += fmt.Sprintf("%s", chatmessage.Text)
		}
	}
	return out
}

// Close file
func closeStorage(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Errorf("could not close file %s: %v", file.Name(), err)
	}
}
