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
	"context"
	"fmt"

	"bitbucket.org/stefanhans/go-thesis/4.3.1/info-gRPC/info"
	"github.com/spf13/cobra"
)

// *** CUSTOMIZED ***
// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads all info from server and prints it",
	RunE: func(cmd *cobra.Command, args []string) error {
		return read(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}

// *** ADDED ***
// Read wrapper function
func read(ctx context.Context) error {

	// Read from gRPC client
	l, err := client.Read(ctx, &info.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch info: %v", err)
	}

	// Print messages
	for _, t := range l.Infos {
		fmt.Printf("%s: %s\n", t.From, t.Text)
	}
	return nil
}

