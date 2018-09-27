// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"os/exec"
	"testing"
)

func TestYumInstalled(t *testing.T) {

	// First make sure rpm is in the PATH
	// Don't run these tests unless you are on a system with rpm
	_, err := exec.LookPath("rpm")
	if err != nil {
		return
	}

	tests := checkerTests{
		{
			Name: "kernel package should be installed",
			Args: Args{
				"package": "kernel",
			},
		},
		{
			Name: "NOT_A_YUM_PACKAGE should not be installed",
			Args: Args{
				"package": "NOT_A_YUM_PACKAGE",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["yum-installed"])
}
