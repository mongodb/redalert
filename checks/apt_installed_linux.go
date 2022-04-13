// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"fmt"
	"os/exec"
)

func init() {
	availableChecks["apt-installed"] = AptInstalledFromArgs
	// alias dpkg-installed to apt-installed for backwards compatibility
	availableChecks["dpkg-installed"] = availableChecks["apt-installed"]
}

// AptInstalled checks if an apt/deb package is installed on the system
//
// Type:
//   - apt-installed
//
// Support Platforms:
//   - Linux
//
// Arguments:
//   package (required): A string value that represents the deb package
type AptInstalled struct {
	Package string
}

// Check if a deb package is installed on the system
func (ai AptInstalled) Check() error {
	out, err := exec.Command("dpkg", "-l", ai.Package).Output()
	if err != nil {
		return fmt.Errorf("%s isn't installed and should be: %s", ai.Package, err)
	}

	if len(out) <= 0 {
		return fmt.Errorf("%s isn't installed and should be", ai.Package)
	}

	return nil
}

// AptInstalledFromArgs will populate the AptInstalled struct with the args
// given in the tests YAML config
func AptInstalledFromArgs(args Args) (Checker, error) {
	ai := AptInstalled{}

	if err := requiredArgs(args, "package"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &ai); err != nil {
		return nil, err
	}

	return ai, nil
}
