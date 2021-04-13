// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/redalert/reports"
	"github.com/spf13/cobra"
)

<<<<<<< HEAD
var (
	systemtype string
)

// Document will list the installed packages in current machinees zero to one arguments.
//
// It takes the following flags:
//
//  - `--type` tye type of the system: supported values are macos, debian
=======
type details struct {
	Type     string
	Versions []versions
}

type versions struct {
	Name    string
	Version string
}

func ReadRevision(fn string) (string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string

	// We expect the Revision to be on the first line
	line, err = bf.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	if err != nil {
		return "", err
	}

	if line == "" {
		return "", fmt.Errorf("empty")
	}
	return line, nil
}

>>>>>>> 61f90a4... feat: get toolchain revisions from preset configs
var Document = &cobra.Command{
	Use:   "document",
	Short: "Document the current image",
	Run: func(cmd *cobra.Command, args []string) {
<<<<<<< HEAD
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
=======
		fmt.Println("Hello from the document command")
		if len(args) > 0 {
			fmt.Printf("Args: %v\n", args)
		}
		toolchains := make(map[string]string)

		toolchains["mongodb"] = "/tmp/revision.txt"
		toolchains["python"] = "/tmp/rev_python.txt"

		details := make(map[string]map[string]string)
		toolchainDetails := make(map[string]string)
		for name, path := range toolchains {
			revision, err := ReadRevision(path)
			if err != nil {
				fmt.Println(err)
				toolchainDetails[name] = "none"
			} else {
				toolchainDetails[name] = revision
			}
			//fmt.Println(toolchainDetails)
		}
		details["toolchains"] = toolchainDetails
		//fmt.Println(details)

		jsonString, _ := json.Marshal(details)
		fmt.Println(string(jsonString))
>>>>>>> 61f90a4... feat: get toolchain revisions from preset configs
	},
}

func init() {
	Document.Flags().StringVarP(&systemtype, "type", "t", "", "type of the system, valid values: debian, rpm")
}
