// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chasinglogic/redalert/testfile"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type multiFlags []string

func (sf *multiFlags) String() string {
	return strings.Join(*sf, ", ")
}

func (sf *multiFlags) Set(value string) error {
	*sf = append(*sf, value)
	return nil
}

func (sf *multiFlags) Type() string {
	return "string"
}

var (
	fileFlag string
	suites   multiFlags
	tests    multiFlags
	output   string
)

func init() {
	Run.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to test file to run. Note: this overrides default test file detection.")
	Run.Flags().StringVarP(&output, "output", "o", "", "Output format to use. Valid values: text, json, csv. Default: text")
	Run.Flags().Var(&suites, "suite", "Suite or alias name to run, can be passed multiple times.")
	Run.Flags().Var(&tests, "test", "Specific test name to run, can be passed multiple times.")
}

// Run will simply run the tests. It takes zero to one arguments. If no
// arguments are given it looks for a tests.(yml|yaml) or redalert.(yml|yaml) in
// the local directory, in `$HOME/.redalert/` (`%APPDATA%\redalert` instead of
// `$HOME/.redalert` on windows.) and finally in `/etc/redalert/` (`C:\redalert`
// on windows).
//
// It takes the following flags:
//
//  - `--quiet` only report test failures.
//  - `--jobs $VALUE` specify the number of parallel tests to run. Default to # of cores.
//  - `--test $TEST_NAME` specify a test by name to run, can be provided multiple times.
//  - `--suite $SUITE_NAME` which test suite to run, an [alias](#aliases) can be provided as the suite name.
//    Can be provided multiple times. There is a default "all" alias which matches
//    all tests and is used if `--suite` is not provided.
//  - `--file $FILE_PATH` in lieu of searching in directories you can specify a file using this flag.
//    If provided the aliases.(yml|yaml) file will not be looked for but aliases in the file will still be loaded.
//  - `--output $FORMAT` specify the output format. Valid values: text, json, csv. Default: text
//
// Run will attempt to load an additional file `aliases.(yml|yaml)` which
// specifies the available [aliases](#aliases). If not found aliases will be
// looked for as a key in the `tests.(yml|yaml)` according to the
// [aliases](#aliases) section.
//
// Run will run tests in parallel.
var Run = &cobra.Command{
	Use:   "run",
	Short: "Run tests against this system.",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var testsFile testfile.TestFile

		if fileFlag == "" {
			testsFile, err = loadTestFile(findTestFile())
		} else {
			testsFile, err = loadTestFile(fileFlag)
		}

		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		results := map[string]error{}

		for _, suite := range suites {
			tests := testsFile.TestsToRun(suite)
			loadedChecks, err := testfile.LoadChecks(tests)
			if err != nil {
				fmt.Println("ERROR: Unable to load checks:", err)
				continue
			}

			checksToRun := make(chan testfile.CheckToRun, len(loadedChecks))
			checkResults := make(chan checkResult, len(loadedChecks))

			for i := 0; i < runtime.NumCPU(); i++ {
				go runCheck(checksToRun, checkResults)
			}

			for _, c := range loadedChecks {
				checksToRun <- c
			}

			for i := 0; i < len(loadedChecks); i++ {
				result := <-checkResults
				results[result.Name] = result.Result
			}

			close(checksToRun)
			close(checkResults)
		}

		exitCode := 0

		for name, result := range results {
			if result == nil {
				fmt.Printf("%s: SUCCESS\n", name)
				continue
			}

			fmt.Printf("%s: FAILURE\n%s\n", name, result)
			exitCode = 1
		}

		os.Exit(exitCode)
	},
}

type checkResult struct {
	Name   string
	Result error
}

func runCheck(incomingChecks chan testfile.CheckToRun, results chan checkResult) {
	for check := range incomingChecks {
		results <- checkResult{
			Name:   check.Name,
			Result: check.Check(),
		}
	}
}

func findTestFile() string {
	possiblePaths := []string{
		"tests.yml",
		"tests.yaml",
		"redalert.yml",
		"redalert.yaml",
		filepath.Join(os.Getenv("HOME"), ".redalert", "tests.yml"),
		filepath.Join(os.Getenv("HOME"), ".redalert", "tests.yaml"),
		filepath.Join(os.Getenv("APPDATA"), "redalert", "tests.yml"),
		filepath.Join(os.Getenv("APPDATA"), "redalert", "tests.yaml"),
		filepath.Join(os.Getenv("HOME"), ".redalert", "redalert.yml"),
		filepath.Join(os.Getenv("HOME"), ".redalert", "redalert.yaml"),
		filepath.Join(os.Getenv("APPDATA"), "redalert", "redalert.yml"),
		filepath.Join(os.Getenv("APPDATA"), "redalert", "redalert.yaml"),
		"/etc/redalert/tests.yml",
		"/etc/redalert/tests.yaml",
		"/etc/redalert/redalert.yml",
		"/etc/redalert/redalert.yaml",
		"C:\\redalert\\tests.yml",
		"C:\\redalert\\tests.yaml",
		"C:\\redalert\\redalert.yml",
		"C:\\redalert\\redalert.yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path
		}
	}

	return ""
}

func loadTestFile(path string) (testfile.TestFile, error) {
	if path == "" {
		return testfile.TestFile{}, errors.New("no tests files found")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return testfile.TestFile{}, err
	}

	var tf testfile.TestFile

	err = yaml.Unmarshal(content, &tf)
	if err != nil {
		return testfile.TestFile{}, err
	}

	dir := filepath.Dir(path)
	possibleAliasFiles := []string{
		filepath.Join(dir, "aliases.yaml"),
		filepath.Join(dir, "aliases.yml"),
	}

	var aliasContent []byte

	for _, file := range possibleAliasFiles {
		aliasContent, err = ioutil.ReadFile(file)
		if err != nil && !os.IsNotExist(err) {
			return tf, err
		} else if err == nil {
			break
		}
	}

	if len(aliasContent) == 0 {
		return tf, nil
	}

	var other testfile.TestFile

	err = yaml.Unmarshal(aliasContent, &other)
	if err != nil {
		return tf, err
	}

	return tf.Join(other), nil
}
