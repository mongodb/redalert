package checks

import "testing"

type checkerTest struct {
	Name        string
	Args        Args
	ShouldError bool
	Error       string
}

type checkerTests []checkerTest

func runCheckerTests(t *testing.T, tests checkerTests, argable ArgFunc) {
	for _, test := range tests {
		checker, err := argable(test.Args)
		if err != nil && test.ShouldError {
			if test.Error != "" && err.Error() != test.Error {
				t.Errorf("%s: Expected error message: %s Got: %s", test.Name, test.Error, err.Error())
			}

			t.Logf("%s: Error creating checker: %s", test.Name, err)
			continue
		} else if err != nil {
			t.Errorf("%s: Unexpected error: %s", test.Name, err)
			continue
		}

		err = checker.Check()
		if err != nil && test.ShouldError {
			if test.Error != "" && err.Error() != test.Error {
				t.Errorf("%s: Expected error message: %s Got: %s", test.Name, test.Error, err.Error())
			}

			t.Logf("%s: Error running checker: %s", test.Name, err)
		} else if err != nil {
			t.Errorf("%s: Check error: %s", test.Name, err)
		} else if err == nil && test.ShouldError {
			t.Errorf("%s: Didn't get an error and should have.", test.Name)
		}
	}
}
