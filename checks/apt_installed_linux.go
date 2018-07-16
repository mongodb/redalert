package checks

import (
	"fmt"
	"os/exec"
)

// AptInstalled checks if an apt/deb package is installed on the system
//
// Type:
//   - apt-installed
//
// Support Platforms:
//   - Linux
//
// Arguments:
//   name (required): A string value that represents the deb package
type AptInstalled struct {
	Name string
}

// Check if a deb package is installed on the system
func (ai AptInstalled) Check() error {
	out, err := exec.Command("dpkg", "-l", ai.Name).Output()
	if err != nil {
		return fmt.Errorf("%s isn't installed and should be: %s", ai.Name, err)
	}

	if len(out) <= 0 {
		return fmt.Errorf("%s isn't installed and should be", ai.Name)
	}

	return nil
}

// FromArgs will populate the AptInstalled struct with the args given in the tests YAML
// config
func (ai AptInstalled) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "name"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &ai); err != nil {
		return nil, err
	}

	return ai, nil
}
