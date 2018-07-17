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
	availableChecks["compile-gcc"] = func(args map[string]interface{}) (Checker, error) {
		return CompileGcc{}.FromArgs(args)
	}
}

// CompileGcc runs gcc compile.
//
// Type:
//	 - compile-gcc
//
// Support Platforms:
//   - Mac
//   - Linux
//   - Windows
//
// Arguments:
//   source (required): The source code of the script.
//   compiler: path to the compiler. Default is 'gcc' from the PATH
//	 cflags: compiles flags, string, e.g "-lss -lsasl2"
//	 cflags_command: command to get clags, e.g. "net-snmp-config --agent-libs"
type CompileGcc struct {
	Source        string
	Compiler      string
	Cflags        string
	CflagsCommand string `mapstructure:"cflags_command"`
}

// Check Runs a gcc command and checks the return code
func (cg CompileGcc) Check() error {

	tmpfolder, err := ioutil.TempDir("", "compileGcc_")
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
		return fmt.Errorf("Problem closing a tmpfile: %s", err)
	}

	argv := []string{"-Werror", "-o", outfileName, srcfileName}
	argv = append(argv, cg.Cflags)

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

		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Problem running the script: %s: %s", err.Error(), string(out))
		}

		flags := strings.TrimRight(string(out), "\n")
		argv = append(argv, strings.Split(flags, " ")...)
	}

	cmd := exec.Command(cg.Compiler, argv...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Problem running the gcc compile: %s: %s", err.Error(), string(out))
	}

	return nil
}

// FromArgs will populate the CompileGcc with the args given in the tests YAML
// config
func (cg CompileGcc) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "source"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &cg); err != nil {
		return nil, err
	}

	if _, compilerGiven := args["compiler"]; cg.Compiler == "" && !compilerGiven {
		cg.Compiler = "gcc"
	}

	if _, cflagsGiven := args["cflags"]; cg.Cflags == "" && !cflagsGiven {
		cg.Cflags = ""
	}

	if _, cflagscommandGiven := args["cflags_command"]; cg.CflagsCommand == "" && !cflagscommandGiven {
		cg.CflagsCommand = ""
	}

	return cg, nil
}
