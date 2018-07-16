package commands

import "github.com/spf13/cobra"

func init() {
	Root.AddCommand(Run)
}

// Root CLI command. This should have no functionality.
var Root = &cobra.Command{
	Use:   "redalert",
	Short: "Validate system state.",
}
