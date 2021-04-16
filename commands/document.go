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
		details := make(map[string]interface{})

		// Load reports conf file
		var testsFile testfile.TestFile
		err := fmt.Errorf("")
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
			reportsToRun := testsFile.TestsToRun(suite)
			for _, report := range reportsToRun {
				if report.Type == "toolchains" {
					toolchains[report.Name] = report.Args["path"].(string)
				}
				if report.Type == "packages" {
					fmt.Println(report.Args["pkg-mgr"])
					packageDetails, err := reports.GetPackagesDetails(report.Args["pkg-mgr"].(string))
					if err != nil {
						fmt.Println("ERR: " + err.Error())
					}
					details["packages"] = packageDetails
				}
			}
		}

		toolchainDetails := reports.GetToolchainDetails(toolchains)
		if len(toolchainDetails) > 0 {
			details["toolchains"] = toolchainDetails
		}
		//details["packages"] = packageDetails

		jsonString, _ := json.Marshal(details)
		fmt.Println(string(jsonString))
	},
}

func init() {
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
