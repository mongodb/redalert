// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/blang/semver"
)

func init() {
	availableChecks["python-module-version"] = PipInstalledFromArgs
}

// PipInstalled checks if python module is installed on the system
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
type PipInstalled struct {
	Python       string
	Module       string
	Version      string
	Relationship string
	Statement    string
}

func (pmv PipInstalled) makeStringSemverCompatible(s string) string {
	split := strings.Split(s, ".")
	switch len(split) {
	case 2:
		split = append(split, "0")
	case 1:
		split = append(split, "0", "0")
	}

	for i := range split {
		// Cleanup whitesapce
		split[i] = strings.TrimSpace(split[i])

		// Check for 0 padded numbers
		if !strings.HasPrefix(split[i], "0") {
			continue
		}

		zeroPadded := []rune(split[i])
		for x := range zeroPadded {
			if x == len(zeroPadded)-1 {
				split[i] = "0"
				break
			}

			if zeroPadded[x] != '0' {
				split[i] = string(zeroPadded[x:])
				break
			}
		}
	}

	return strings.TrimSpace(strings.Join(split, "."))
}

// Check if a python module is installed on the system and verify version if
// the Version string is set
func (pmv PipInstalled) Check() error {
	if pmv.Statement == "" {
		pmv.Statement = pmv.Module + ".__version__"
	}

	pyCommand := "import " + pmv.Module + "; print(" + pmv.Statement + ")"
	out, err := exec.Command(pmv.Python, "-c", pyCommand).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s isn't installed and should be: %s, %s", pmv.Module, err, string(out))
	}

	// Don't check semver if not provided
	if pmv.Version == "" {
		return nil
	}

	pmv.Version = pmv.makeStringSemverCompatible(pmv.Version)

	strippedOutput := pmv.makeStringSemverCompatible(string(out))

	installedVersion, err := semver.Parse(strippedOutput)
	if err != nil {
		return fmt.Errorf("Unable to parse semver from python output: %s: %s", strippedOutput, err)
	}

	strippedInput := pmv.makeStringSemverCompatible(pmv.Version)
	requestedVersion, err := semver.Parse(strippedInput)
	if err != nil {
		return fmt.Errorf("Unable to parse semver from args: %s", err)
	}

	switch pmv.Relationship {
	case "eq":
		if !installedVersion.EQ(requestedVersion) {
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
		if !installedVersion.GT(requestedVersion) {
			return fmt.Errorf("%s is not greater than %s", installedVersion, requestedVersion)
		}
	default:
		if !installedVersion.GTE(requestedVersion) {
			return fmt.Errorf("%s is not greater than or equal to %s", installedVersion, requestedVersion)
		}
	}
	return nil
}

// FromArgs will populate the PipInstalled struct with the args given in the tests YAML
// config
func PipInstalledFromArgs(args Args) (Checker, error) {
	pmv := PipInstalled{}

	if err := requiredArgs(args, "module"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &pmv); err != nil {
		return nil, err
	}

	return pmv, nil
}
