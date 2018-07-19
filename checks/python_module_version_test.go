// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package checks

import (
	"os/exec"
	"testing"
)

func TestPythonModuleVersion(t *testing.T) {
	// First make sure python is in the PATH
	// Don't run these tests unless you are on a system with python installed
	_, err := exec.LookPath("python")
	if err != nil {
		return
	}

	// Test a module exists, don't specify a version
	err = PythonModuleVersion{Module: "yaml"}.Check()
	if err != nil {
		t.Error(err)
		return
	}

	// A python module you certainly don't want on your system
	// Don't specify a version
	err = PythonModuleVersion{Module: "i-am-not-a-python-module-at-all"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expected behavior here.", err)
		return
	}

	// Test a module exists and that the version is correct
	err = PythonModuleVersion{Module: "yaml", Version: "3.13"}.Check()
	if err != nil {
		t.Error(err)
		return
	}

	// Test a module exists and that the version is correct
	err = PythonModuleVersion{Module: "yaml", Version: "3.00"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expected behavior here.", err)
		return
	} // Test a module exists, but the version is incorrect

	checker, err := PythonModuleVersion{}.FromArgs(Args{"module": "yaml"})
	if err != nil {
		t.Error(err)
		return
	}

	checker, err = PythonModuleVersion{}.FromArgs(Args{"module": "yaml", "version": "3.13"})
	if err != nil {
		t.Error(err)
		return
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
		return
	}
}
