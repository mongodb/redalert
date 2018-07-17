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
	err = PythonModuleVersion{Name: "yaml"}.Check()
	if err != nil {
		t.Error(err)
	}

	// A python module you certainly don't want on your system
	// Don't specify a version
	err = PythonModuleVersion{Name: "i-am-not-a-python-module-at-all"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expected behavior here.", err)
	}

	// Test a module exists and that the version is correct
	err = PythonModuleVersion{Name: "yaml", Version: "3.13"}.Check()
	if err != nil {
		t.Error(err)
	}

	// Test a module exists and that the version is correct
	err = PythonModuleVersion{Name: "yaml", Version: "3.99"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expected behavior here.", err)
	} // Test a module exists, but the version is incorrect

	checker, err := PythonModuleVersion{}.FromArgs(map[string]interface{}{"name": "yaml"})
	if err != nil {
		t.Error(err)
	}

	checker, err = PythonModuleVersion{}.FromArgs(map[string]interface{}{"name": "yaml", "version": "3.13"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
