// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/mongodb/redalert/commands"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	commands.Root.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("redalert-%s-%s-%s\n", version, commit, date)
	},
}

func main() {
	commands.Root.Execute()
}
