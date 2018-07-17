package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	availableChecks["python-module-version"] = func(args map[string]interface{}) (Checker, error) {
		return PythonModuleVersion{}.FromArgs(args)
	}
}

// PythonModuleVersion checks if python module is installed on the system
// And verifies its version.
//
// Type:
//   - python-module-version
//
// Supported Platforms:
//   - Linux
//   - Windows
//
// Argument:
//   name (required): A string value that represents the python module
//   version: An optional version number to check
//            Leave version blank to just verify module is present
type PythonModuleVersion struct {
	Name    string
	Version string
}

// Check if a python module is installed on the system and verify version if
// the Version string is set
func (pmv PythonModuleVersion) Check() error {
	pyCommand := "import " + pmv.Name + "; print(" + pmv.Name + ".__version__)"
	out, err := exec.Command("python", "-c", pyCommand).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s isn't installed and should be: %s, %s", pmv.Name, err, string(out))
	}

	if len(pmv.Version) >= 1 {
		if string(pmv.Version) != strings.TrimRight(string(out), "\n") {
			return fmt.Errorf("%s has version %s and it should be version: %s", pmv.Name, out, pmv.Version)
		}
	}

	return nil
}

// FromArgs will populate the PythonModuleVersion struct with the args given in the tests YAML
// config
func (pmv PythonModuleVersion) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "name"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &pmv); err != nil {
		return nil, err
	}

	return pmv, nil
}
