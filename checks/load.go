package checks

import (
	"fmt"

	"github.com/chasinglogic/redalert/testfile"
)

var availableChecks = map[string]ArgableFunc{}

// LoadCheck will return the appropriate Checker based on the test type name.
// As documented on the various checkers
func LoadCheck(name string, args map[string]interface{}) (Checker, error) {
	if argFunc, exists := availableChecks[name]; exists {
		checker, err := argFunc(args)
		if err != nil {
			return checker, fmt.Errorf("Error loading check %s: %s", name, err)
		}

		return checker, nil
	}

	return nil, fmt.Errorf("%s is not a known check type", name)
}

// CheckToRun keeps the name and actual check object together for easy
// reporting to the user.
type CheckToRun struct {
	Name    string
	Checker Checker
}

// Check makes CheckToRun a Checker
func (ctr CheckToRun) Check() error {
	return ctr.Checker.Check()
}

// LoadChecks takes a slice of tesfile.Tests and returns a slice of Checks to run
func LoadChecks(tests []testfile.Test) ([]CheckToRun, error) {
	checks := make([]CheckToRun, len(tests))

	var err error

	for i, test := range tests {
		checks[i].Name = test.Name
		checks[i].Checker, err = LoadCheck(test.Type, test.Args)
		if err != nil {
			return nil, err
		}
	}

	return checks, nil
}
