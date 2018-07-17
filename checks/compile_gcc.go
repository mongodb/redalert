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
//   interpreter: path to bash. Default is 'gcc' from the PATH
//				  but it will also check if gcc exists in toolchain folders.
type CompileGcc struct {
	Source      string
	Interpreter string
	Cflags      string
}

// Check Runs a gcc command and checks the return code
func (cg CompileGcc) Check() error {

	tmpfolder, err := ioutil.TempDir("", "compileGcc_")
	if err != nil {
		return fmt.Errorf("Problem creating a tmpdir: %s", err)
	}

	srcfileName := tmpfolder + "/src.c"
	outfileName := tmpfolder + "/out.o"

	srcfile, err := os.Create(srcfileName)
	defer os.RemoveAll(tmpfolder)

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

	cmd := exec.Command(cg.Interpreter, argv...)
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

	if _, interpreterGiven := args["interpreter"]; cg.Interpreter == "" && !interpreterGiven {

		cg.Interpreter = "gcc"

		paths := []string{
			"/opt/mongodbtoolchain/v2/bin/gcc",
			"/opt/mongodbtoolchain/v1/bin/gcc",
			"/opt/mongodbtoolchain/bin/gcc",
			"/usr/bin/gcc",
			"/usr/local/bin/gcc",
		}

		for _, path := range paths {
			if _, err := os.Stat(path); !os.IsNotExist(err) {
				cg.Interpreter = path
				break
			}
		}
	}

	if _, cflagsGiven := args["cflags"]; cg.Cflags == "" && !cflagsGiven {
		cg.Cflags = ""
	}
	return cg, nil
}
