// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"bitbucket.org/stefanhans/go-thesis/4.3.1/info-gRPC/info"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// *** CUSTOMIZED ***
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "info",
	Short: "Connects to the info server to read or write",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// *** ADDED ***
var client info.InfosClient

// *** CUSTOMIZED ***
func init() {
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}
	client = info.NewInfosClient(conn)
}
