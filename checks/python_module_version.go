package checks

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/blang/semver"
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
//   module (required): A string value that represents the python module
//   version: An optional version number to check
//            Leave version blank to just verify module is present
//   statement: Optional python statement, the result will be passed to print()
//              Defaults to module.__version__
//   relationship: Optional comparison operator for the version provided. Valid
//                 values are lt, lte, gt, gte, eq. Defaults to gte (greater than or equal to)
type PythonModuleVersion struct {
	Module       string
	Version      string
	Relationship string
	Statement    string
}

func (pmv PythonModuleVersion) makeStringSemverCompatible(s string) string {
	switch len(strings.Split(s, ".")) {
	case 2:
		return s + ".0"
	case 1:
		return s + ".0.0"
	default:
		return s
	}
}

// Check if a python module is installed on the system and verify version if
// the Version string is set
func (pmv PythonModuleVersion) Check() error {
	if pmv.Statement == "" {
		pmv.Statement = pmv.Module + ".__version__"
	}

	pyCommand := "import " + pmv.Module + "; print(" + pmv.Statement + ")"
	out, err := exec.Command("python", "-c", pyCommand).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s isn't installed and should be: %s, %s", pmv.Module, err, string(out))
	}

	// Don't check semver if not provided
	if pmv.Version == "" {
		return nil
	}

	pmv.Version = pmv.makeStringSemverCompatible(pmv.Version)

	// strip the newline added by python's print()
	strippedOutput := strings.TrimRight(string(out), "\n")
	strippedOutput = pmv.makeStringSemverCompatible(strippedOutput)

	installedVersion, err := semver.Parse(strippedOutput)
	if err != nil {
		return fmt.Errorf("Unable to parse semver from python output: %s: %s", strippedOutput, err)
	}

	requestedVersion, err := semver.Parse(pmv.Version)
	if err != nil {
		return fmt.Errorf("Unable to parse semver from args: %s", err)
	}

	switch pmv.Relationship {
	case "eq":
		if !requestedVersion.EQ(installedVersion) {
			return fmt.Errorf("%s is not equal to %s", installedVersion, requestedVersion)
		}
	case "lt":
		if !installedVersion.LT(requestedVersion) {
			return fmt.Errorf("%s is not less than %s", installedVersion, requestedVersion)
		}
	case "lte":
		if !installedVersion.LTE(requestedVersion) {
			return fmt.Errorf("%s is not less than or equal to %s", installedVersion, requestedVersion)
		}
	case "gt":
		if !requestedVersion.GT(installedVersion) {
			return fmt.Errorf("%s is not greater than %s", installedVersion, requestedVersion)
		}
	default:
		if !requestedVersion.GTE(installedVersion) {
			return fmt.Errorf("%s is not greater than or equal to %s", installedVersion, requestedVersion)
		}
	}
	return nil
}

// FromArgs will populate the PythonModuleVersion struct with the args given in the tests YAML
// config
func (pmv PythonModuleVersion) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "module"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &pmv); err != nil {
		return nil, err
	}

	return pmv, nil
}
