package checks

import (
	"fmt"
	"os"
)

// FileChecker checks if a file exists or does not
type FileChecker struct {
	Exists bool
	Name   string
}

// Check if a file exists or does not
func (fc FileChecker) Check() error {
	_, err := os.Stat(fc.Name)

	switch os.IsNotExist(err) {
	case true && fc.Exists:
		return fmt.Errorf("%s doesn't exist and should", fc.Name)
	case false && !fc.Exists:
		return fmt.Errorf("%s does exist and shouldn't", fc.Name)
	default:
		return nil
	}
}

// FromArgs will populate the FileChecker with the args given in the tests YAML
// config
func (fc FileChecker) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "name"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &fc); err != nil {
		return nil, err
	}

	if _, existsGiven := args["exists"]; !fc.Exists && !existsGiven {
		fc.Exists = true
	}

	return fc, nil
}
