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
availableChecks["run-script"] = func(args map[string]interface{}) (Checker, error) {
		return RunScript{}.FromArgs(args)
	}

	availableChecks["run-bash-script"] = func(args map[string]interface{}) (Checker, error) {
		args["interpreter"] = "/bin/bash"
		return RunScript{}.FromArgs(args)
	}

	availableChecks["run-python-script"] = func(args map[string]interface{}) (Checker, error) {
		args["interpreter"] = "python"
		return RunScript{}.FromArgs(args)
	}

	availableChecks["run-python2-script"] = func(args map[string]interface{}) (Checker, error) {
		args["interpreter"] = "python2"
		return RunScript{}.FromArgs(args)
	}

	availableChecks["run-python3-script"] = func(args map[string]interface{}) (Checker, error) {
		args["interpreter"] = "python3"
		return RunScript{}.FromArgs(args)
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

	if rs.Output != "" && strings.TrimRight(string(out), "\n") != rs.Output {
		return fmt.Errorf("Output doesn't match. Expected '%s'\nActual output: %s", rs.Output, string(out))
	}

	return nil
}

// FromArgs will populate the RunScript with the args given in the tests YAML
// config
func (rs RunScript) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "source"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &rs); err != nil {
		return nil, err
	}

	if _, interpreterGiven := args["interpreter"]; rs.Interpreter == "" && !interpreterGiven {
		rs.Interpreter = "/bin/bash"
	}

	if _, outputGiven := args["output"]; rs.Output == "" && !outputGiven {
		rs.Output = ""
	}

	return rs, nil
}
