// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package checks

import (
	"os/exec"
	"testing"
)

func TestGemInstalled(t *testing.T) {
	// First make sure gem is in the PATH
	// Don't run these tests unless you are on a system with gem installed
	_, err := exec.LookPath("gem")
	if err != nil {
		return
	}

	// A gem that should exist everywhere
	err = GemInstalled{Name: "rake"}.Check()
	if err != nil {
		t.Error(err)
	}

	// A gem you certainly don't want on your system
	err = GemInstalled{Name: "i-am-not-a-gem-at-all"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expected behavior here.")
	}

	checker, err := GemInstalledFromArgs(Args{"name": "rake"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
