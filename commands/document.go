// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mongodb/redalert/reports"
	"github.com/mongodb/redalert/testfile"
	"github.com/spf13/cobra"
)

var (
	packageManager string
)

// Document will list the installed packages in the system.
var Document = &cobra.Command{
	Use:   "document",
	Short: "Document the current image",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Package manager: " + packageManager)
		details := make(map[string]interface{})

		pacakgeDetails, err := reports.GetPackagesDetails(packageManager)
		if err != nil {
			fmt.Println("ERR: " + err.Error())
		}

		// Load reports conf file
		var testsFile testfile.TestFile

		if fileFlag == "" {
			testsFile, err = loadTestFile(findReportFile())
		} else {
			testsFile, err = loadTestFile(fileFlag)
		}

		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		toolchains := make(map[string]string)
		// Find reports for this suite
		for _, suite := range suites {
			reports := testsFile.TestsToRun(suite)
			for _, report := range reports {
				if report.Type == "toolchains" {
					toolchains[report.Name] = report.Args["path"].(string)
				}
			}
		}

		toolchainDetails := reports.GetToolchainDetails(toolchains)
		details["toolchains"] = toolchainDetails
		details["packages"] = pacakgeDetails

		jsonString, _ := json.Marshal(details)
		fmt.Println(string(jsonString))
	},
}

func init() {
	Document.Flags().StringVarP(&packageManager, "pkg-mngr", "p", "", "Package manager to be used to list the installed pacakges")
	Document.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to test file to run. Note: this overrides default test file detection.")
	Document.Flags().Var(&suites, "suite", "Suite or alias name to run, can be passed multiple times.")
}

func findReportFile() string {
	possiblePaths := []string{
		"toolchains.yml",
		"toolchains.yaml",
		"reports.yml",
		"reports.yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path
		}
	}

	return ""
}
