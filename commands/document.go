// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/mongodb/redalert/reports"
	"github.com/spf13/cobra"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var externalCommands = map[string]externalCommand{
	"debian": []string{"dpkg-query", "-W -f='${binary:Package};${Version}\n'"},
	"macos":  []string{"pkgutil", "--pkgs"}}

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
		if len(systemtype) > 0 {
			commandRes, err := reports.GetPackagesDetails(systemtype)
			if err != nil {
				fmt.Println("ERR: " + err.Error())
			}
			fmt.Println(string(commandRes))
		}

		details := make(map[string]map[string]string)
		externalCommand := externalCommands[systemtype]
		command := exec.Command(externalCommand[0], externalCommand[1:]...)
		commandRes, err := command.CombinedOutput()

		toolchainDetails := reports.GetToolchainDetails()
		details["toolchains"] = toolchainDetails
		//fmt.Println(details)

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

func init() {
	Document.Flags().StringVarP(&systemtype, "type", "t", "", "type of the system, valid values: debian, rpm")
}
