// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package checks

import (
	"fmt"
	"log"
	"os/exec"
)

func init() {
	availableChecks["yum-installed"] = YumInstalledFromArgs
}

// YumInstalled checks if an rpm package is installed on the system
//
// Type:
//   - yum-installed
//
// Support Platforms:
//   - Linux
//
// Arguments:
//   package (required): A string value that represents the rpm package
type YumInstalled struct {
	Packages []string
}

// Check if an rpm is installed on the system
func (yi YumInstalled) Check() error {
	for _, p := range yi.Packages {
		out, err := exec.Command("yum", "list", "installed", p).Output()
		if err != nil {
			log.Fatal(err)
		}

		if len(out) <= 0 {
			return fmt.Errorf("%s isn't installed and should be", p)
		}

	}

	return nil
}

// YumInstalledFromArgs will populate the YumInstalled struct with the args
// given in the tests YAML config
func YumInstalledFromArgs(args Args) (Checker, error) {
	yi := YumInstalled{}

	if err := requiredArgs(args, "packages"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &yi); err != nil {
		return nil, err
	}

	return yi, nil
}
