// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import "github.com/spf13/cobra"

func init() {
	Root.AddCommand(Run)
	Root.AddCommand(Document)
}

// Root CLI command. This should have no functionality.
var Root = &cobra.Command{
	Use:   "redalert",
	Short: "Validate system state.",
}
