package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"bitbucket.org/stefanhans/go-thesis/4.3.1/info-pb"
	"github.com/golang/protobuf/proto"
)

func main() {

	// Declare array with protobuffer messages
	infos := []info_pb.Info{
		info_pb.Info{
			Text: "I am a painter",
			From: "Marc Chagall",
		},
		info_pb.Info{
			Text: "I am a writer",
			From: "Edgar Allan Poe",
		},
	}

	// Open file for appending info
	filename := "storage"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("could not open %s: %v", filename, err)
	}

	// Range over protobuffers
	for _, v := range infos {

		// Marshal into binary format
		byteArray, err := proto.Marshal(&v)
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

	// Close file
	if err := file.Close(); err != nil {
		fmt.Errorf("could not close file %s: %v", filename, err)
	}

	// Open file for reading info
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("could not read %s: %v", filename, err)
	}

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
		var info info_pb.Info
		if err := proto.Unmarshal(byteArray[:length], &info); err != nil {
			fmt.Errorf("could not read info: %v", err)
		}
		byteArray = byteArray[length:]

		fmt.Printf("%s: %q\n", info.From, info.Text)
	}
}
