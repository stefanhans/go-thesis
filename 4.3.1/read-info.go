package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"bitbucket.org/stefanhans/go-thesis/4.3.1/info-pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	// Open file for reading info
	filename := "storage"
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
