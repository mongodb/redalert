// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Document = &cobra.Command{
	Use:   "document",
	Short: "Document the current image",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from the ducment command")
	},
}
