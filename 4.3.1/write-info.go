package main

import (
	"encoding/binary"
	"fmt"
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
}
