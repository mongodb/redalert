// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package main

import (
	"fmt"

	"github.com/chasinglogic/redalert/commands"
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
