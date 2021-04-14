// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"encoding/json"
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
		if !reports.IsValidSystemType(systemtype) {
			fmt.Errorf("System type not supported: " + systemtype)
			return
		}

		details := make(map[string]interface{})

		pacakgeDetails, err := reports.GetPackagesDetails(systemtype)
		if err != nil {
			fmt.Println("ERR: " + err.Error())
		}

		toolchainDetails := reports.GetToolchainDetails()
		details["toolchains"] = toolchainDetails
		details["packages"] = pacakgeDetails

		jsonString, _ := json.Marshal(details)
		fmt.Println(string(jsonString))
	},
}

func init() {
	Document.Flags().StringVarP(&systemtype, "type", "t", "", "type of the system, valid values: debian, rpm")
}
