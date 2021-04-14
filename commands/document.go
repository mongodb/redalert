// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/mongodb/redalert/reports"
	"github.com/spf13/cobra"
)

// Document will list the installed packages in the system.
var Document = &cobra.Command{
	Use:   "document",
	Short: "Document the current image",
	Run: func(cmd *cobra.Command, args []string) {
		os := runtime.GOOS
		details := make(map[string]interface{})

		pacakgeDetails, err := reports.GetPackagesDetails(os)
		if err != nil {
			fmt.Println("ERR: " + err.Error())
		}

		toolchainDetails := reports.GetToolchainDetails()
		details["toolchains"] = toolchainDetails
		details["packages"] = pacakgeDetails

		jsonString, _ := json.Marshal(details)
		fmt.Println(string(jsonString))
		commandsParsed := parseCommandOuput(string(commandRes), systemtype)

		formattedPackages, err := formatPacakges(commandsParsed, "json")

		if err != nil {
			fmt.Println("Could not format the pacakges")
			return
		}

		fmt.Println(formattedPackages)
	},
}
