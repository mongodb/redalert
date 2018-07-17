package checks

import (
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	availableChecks["gem-installed"] = func(args map[string]interface{}) (Checker, error) {
		return GemInstalled{}.FromArgs(args)
	}
}

// GemInstalled checks if ruby gem is installed on the system
//
// Type:
//   - gem-installed
//
// Support Platforms:
//   - Linux
//   - Windows
//
// Argument:
//   name (required): A string value that represents the gem name
type GemInstalled struct {
	Name string
}

// Check if a deb package is installed on the system
func (gi GemInstalled) Check() error {
	out, err := exec.Command("gem", "list", gi.Name).Output()
	if err != nil {
		return fmt.Errorf("%s isn't installed and should be: %s", gi.Name, err)
	}

	// No output
	if len(out) <= 0 {
		return fmt.Errorf("%s isn't installed and should be", gi.Name)
	}

	// Split the gemlist into an array and check for gem package explicitly
	gemList := strings.Split(string(out), "\n")
	for _, gem := range gemList {
		gemName := strings.Split(gem, " ")

		// Skip erroneous output
		if len(gemName) == 0 {
			continue
		}

		if gi.Name == gemName[0] {
			return nil
		}
	}
	return fmt.Errorf("The %s gem was not found on the system", gi.Name)
}

// FromArgs will populate the GemInstalled struct with the args given in the tests YAML
// config
func (gi GemInstalled) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "name"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &gi); err != nil {
		return nil, err
	}

	return gi, nil
}
