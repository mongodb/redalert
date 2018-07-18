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

	err = YumInstalled{Package: "kernel"}.Check()
	if err != nil {
		t.Error(err)
	}

	// This should fail
	err = YumInstalled{Package: "DonaldTrump"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expecte behavior here.")
	}

	checker, err := YumInstalled{}.FromArgs(map[string]interface{}{"name": "kernel"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
