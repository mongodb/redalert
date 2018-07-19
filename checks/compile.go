// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package checks

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/shlex"
)

func init() {
	availableChecks["compile"] = func(args map[string]interface{}) (Checker, error) {
		return CompileChecker{}.FromArgs(args)
	}
}

// CompileChecker runs gcc compile.
//
// Type:
//	 - compile
//	 - compile-gcc
//	 - compile-visual-studio
//
// Support Platforms:
//   - Mac
//   - Linux
//   - Windows
//
// Arguments:
//   source (required): The source code of the script.
//   compiler: path to the compiler. Default is 'gcc' from the PATH on Linux and
//             the most recent Visual Studio on Windows
//   cflags: compiles flags, string, e.g "-lss -lsasl2"
//   cflags_command: command to get clags, e.g. "net-snmp-config --agent-libs"
type CompileChecker struct {
	Source        string
	Compiler      string
	Cflags        string
	CflagsCommand string `mapstructure:"cflags_command"`
	Run           bool
}

// Check is implmented in the platform specific files compile_windows.go and
// compile_unix.go

// Check Runs a gcc command and checks the return code
func (cg CompileChecker) Check() error {
	tmpfolder, err := ioutil.TempDir("", "CompileChecker_")
	if err != nil {
		return fmt.Errorf("Problem creating a tmpdir: %s", err)
	}
	defer os.RemoveAll(tmpfolder)

	srcfileName := filepath.Join(tmpfolder, "src.c")
	outfileName := filepath.Join(tmpfolder, "out.o")

	srcfile, err := os.Create(srcfileName)
	if err != nil {
		return fmt.Errorf("Problem creating a srcfile: %s", err)
	}

	if runtime.GOOS == "windows" {
		cg.Source = strings.Replace(cg.Source, "\n", "\r\n", -1)
	}

	content := []byte(cg.Source)

	if _, err := srcfile.Write(content); err != nil {
		return fmt.Errorf("Problem writing to a tmpfile: %s", err)
	}

	if err := srcfile.Close(); err != nil {
		return fmt.Errorf("Problem closing the tmpfile: %s", err)
	}

	argv := []string{"-Werror", "-o", outfileName}

	if cg.Cflags != "" {
		flags, err := shlex.Split(cg.Cflags)
		if err != nil {
			return fmt.Errorf("Unable to parse cflags: %s", err)
		}

		argv = append(argv, flags...)
	}

	if cg.CflagsCommand != "" {
		fields, err := shlex.Split(cg.CflagsCommand)
		if err != nil {
			return fmt.Errorf("Problem parsing cflags_command: %s", err)
		}

		var cmd *exec.Cmd

		switch len(fields) {
		case 0:
			return fmt.Errorf("Got unexecutable cflags command: %s", cg.CflagsCommand)
		case 1:
			cmd = exec.Command(fields[0])
		default:
			cmd = exec.Command(fields[0], fields[1:]...)
		}

		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("Problem running the script: %s: %s", err.Error(), string(out))
		}

		flags := strings.TrimRight(string(out), "\n")
		flagArgs, err := shlex.Split(flags)
		if err != nil {
			return fmt.Errorf("Problem parsing clfags command output: %s: %s", err, flags)
		}

		argv = append(argv, flagArgs...)
	}

	argv = append(argv, srcfileName)
	cmd := exec.Command(cg.Compiler, argv...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Problem running the gcc compile: %s: %s", err, string(out))
	}

	if cg.Run {
		cmd = exec.Command(outfileName)
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Problem running compiled executable: %s: %s", err, string(out))
		}
	}

	return nil
}

// FromArgs will populate the CompileChecker with the args given in the tests YAML
// config
func (cg CompileChecker) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "source"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &cg); err != nil {
		return nil, err
	}

	if _, compilerGiven := args["compiler"]; cg.Compiler == "" && !compilerGiven {
		cg.Compiler = "gcc"
	}

	return cg, nil
}
