// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
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

	// This should fail
	err = YumInstalled{Package: "DonaldTrump"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expecte behavior here.")
	}

	checker, err := YumInstalled{}.FromArgs(Args{"name": "kernel"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
