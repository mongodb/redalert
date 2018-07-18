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
