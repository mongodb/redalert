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

var Document = &cobra.Command{
	Use:   "document",
	Short: "Document the current image",
	Run: func(cmd *cobra.Command, args []string) {
		if _, ok := externalCommands[systemtype]; !ok {
			fmt.Println("system type not found: " + systemtype)
			return
		}
		externalCommand := externalCommands[systemtype]

		res := exec.Command(externalCommand[0], externalCommand[1:]...)

		stdRes, err := res.CombinedOutput()

		if err != nil {
			fmt.Println("ERR: " + err.Error())
			return
		}

		fmt.Println(string(stdRes))
	},
}

func init() {
	Document.Flags().StringVarP(&systemtype, "type", "t", "", "type of the system, valid values: debian, rpm")
}
