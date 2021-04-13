// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var externalCommands = map[string]externalCommand{"debian": []string{"dpkg", "-l"}, "macos": []string{"pkgutil", "--pkgs"}}

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
		if _, ok := externalCommands[systemtype]; !ok {
			fmt.Println("system type not found: " + systemtype)
			return
		}
		externalCommand := externalCommands[systemtype]

		command := exec.Command(externalCommand[0], externalCommand[1:]...)

		commandRes, err := command.CombinedOutput()

		if err != nil {
			fmt.Println("ERR: " + err.Error())
			return
		}

		fmt.Println(string(commandRes))
	},
}

func init() {
	Document.Flags().StringVarP(&systemtype, "type", "t", "", "type of the system, valid values: debian, rpm")
}
