// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"fmt"

	"github.com/mongodb/redalert/reports"
	"github.com/spf13/cobra"
)

var (
	systemtype string
)

// Document will list the installed packages in current machinees zero to one arguments.
//
// It takes the following flags:
//
//  - `--type` tye type of the system: supported values are macos, debian
var Document = &cobra.Command{
	Use:   "document",
	Short: "Document the current image",
	Run: func(cmd *cobra.Command, args []string) {
<<<<<<< HEAD
		if len(systemtype) > 0 {
			commandRes, err := reports.Packages(systemtype)
			if err != nil {
				fmt.Println("ERR: " + err.Error())
				return
			}
			fmt.Println(string(commandRes))
		}

=======
		fmt.Println("Hello from the ducment command")
		if len(args) > 0 {
			fmt.Printf("Args: %v\n", args)
		}
>>>>>>> 84d7906... feat: test args to document command
	},
}

func init() {
	Document.Flags().StringVarP(&systemtype, "type", "t", "", "type of the system, valid values: debian, rpm")
}
