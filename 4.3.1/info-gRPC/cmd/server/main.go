package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"bitbucket.org/stefanhans/go-thesis/4.3.1/info-gRPC/info"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	// Create and register server
	var infos infoServer
	srv := grpc.NewServer()
	info.RegisterInfosServer(srv, infos)

	// Create listener
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("could not listen to :8888: \v", err)
	}
	// Serve messages via listener
	log.Fatal(srv.Serve(l))
}

// Receiver for implementing the server service interface
type infoServer struct{}

// Server's Write implementation
func (s infoServer) Write(ctx context.Context, info *info.Info) (*info.Info, error) {

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
func (s infoServer) Read(ctx context.Context, void *info.Void) (*info.InfoList, error) {

	// Read file
	b, err := ioutil.ReadFile("storage")
	if err != nil {
		return nil, fmt.Errorf("could not read storage: %v", err)
	}

	// Iterate over read bytes
	var infos info.InfoList
	for {
		if len(b) == 0 {
			// Return result
			return &infos, nil
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
		var info info.Info
		if err := proto.Unmarshal(b[:length], &info); err != nil {
			return nil, fmt.Errorf("could not read info: %v", err)
		}
		b = b[length:]
		infos.Infos = append(infos.Infos, &info)
	}
}
