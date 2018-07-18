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

	err = AptInstalled{Package: "linux-base"}.Check()
	if err != nil {
		t.Error(err)
	}

	// This should fail
	err = AptInstalled{Package: "DonaldTrump"}.Check()
	if err == nil {
		t.Error("Got no error, which is the expected behavior here.")
	}

	checker, err := AptInstalled{}.FromArgs(map[string]interface{}{"name": "linux-base"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
