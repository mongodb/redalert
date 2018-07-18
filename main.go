package main

import (
	"fmt"

	"github.com/chasinglogic/redalert/commands"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("%v, commit %v, built at %v\n", version, commit, date)
	commands.Root.Execute()
}
