package checks

import (
	"fmt"
)

var availableChecks = map[string]ArgableFunc{}

// LoadCheck will return the appropriate Checker based on the test type name.
// As documented on the various checkers
func LoadCheck(name string, args Args) (Checker, error) {
	if argFunc, exists := availableChecks[name]; exists {
		checker, err := argFunc(args)
		if err != nil {
			return checker, fmt.Errorf("Error loading check %s: %s", name, err)
		}

		return checker, nil
	}

	return nil, fmt.Errorf("%s is not a known check type", name)
}
