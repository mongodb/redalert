// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package checks

import (
	"os/exec"
	"testing"
)

func TestAptInstalled(t *testing.T) {
	// First make sure dpgk is in the PATH
	// Don't run these tests unless you are on a system with dpkg
	_, err := exec.LookPath("dpkg")
	if err != nil {
		return
	}

	tests := checkerTests{
		{
			Name: "linux-base should be installed",
			Args: Args{
				"package": "linux-base",
			},
		},
		{
			Name: "NOT_A_PACKAGE should not be installed",
			Args: Args{
				"package": "NOT_A_PACKAGE",
			},
			ShouldError: true,
		}
	}

	checker, err := AptInstalled{}.FromArgs(Args{"name": "linux-base"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
