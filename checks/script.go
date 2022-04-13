// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func init() {
	availableChecks["run-script"] = RunScriptFromArgs

	availableChecks["run-bash-script"] = func(args Args) (Checker, error) {
		args["interpreter"] = "/bin/bash"
		return RunScriptFromArgs(args)
	}

	availableChecks["run-python-script"] = func(args Args) (Checker, error) {
		args["interpreter"] = "python"
		return RunScriptFromArgs(args)
	}

	availableChecks["run-python2-script"] = func(args Args) (Checker, error) {
		args["interpreter"] = "python2"
		return RunScriptFromArgs(args)
	}

	availableChecks["run-python3-script"] = func(args Args) (Checker, error) {
		args["interpreter"] = "python3"
		return RunScriptFromArgs(args)
	}
}

// RunScript runs a bash script and optionally checks output.
//
// Type:
//	 - run-bash-script
//
// Support Platforms:
//   - Mac
//   - Linux
//   - Windows
//
// Arguments:
//   source (required): The source code of the script.
//   output: string to which the output of the script will be compared to determine a successful run. If omitted, only checks that the script exits with returncode 0.
//   interpreter: path to bash. Default is '/bin/bash'.
//
// Notes:
type RunScript struct {
	Source      string
	Output      string
	Interpreter string
}

// Check Runs a script and checks the return code or output
func (rs RunScript) Check() error {
	tmpfile, err := ioutil.TempFile(os.TempDir(), "testScript_")
	if err != nil {
		return fmt.Errorf("Problem creating a tmpfile: %s", err)
	}

	if runtime.GOOS == "windows" {
		rs.Source = strings.Replace(rs.Source, "\n", "\r\n", -1)
	}

	content := []byte(rs.Source)

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		return fmt.Errorf("Problem writing to a tmpfile: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("Problem closing a tmpfile: %s", err)
	}

	cmd := exec.Command(rs.Interpreter, tmpfile.Name())
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Problem running the script: %s: %s", err.Error(), string(out))
	}

	if rs.Output == "" {
		return nil
	}

	var trimmed string
	if runtime.GOOS == "windows" {
		trimmed = strings.TrimRight(string(out), "\r\n")
	} else {
		trimmed = strings.TrimRight(string(out), "\n")
	}

	if trimmed != rs.Output {
		return fmt.Errorf("Output doesn't match. Expected '%s'\nActual output: %s", rs.Output, string(out))
	}

	return nil
}

// RunScriptFromArgs will populate the RunScript with the args given in the
// tests YAML config
func RunScriptFromArgs(args Args) (Checker, error) {
	rs := RunScript{}

	if err := requiredArgs(args, "source"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &rs); err != nil {
		return nil, err
	}

	if _, interpreterGiven := args["interpreter"]; rs.Interpreter == "" && !interpreterGiven {
		rs.Interpreter = "bash"
	}

	if _, outputGiven := args["output"]; rs.Output == "" && !outputGiven {
		rs.Output = ""
	}

	return rs, nil
}
