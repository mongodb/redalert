// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import "testing"

type checkerTest struct {
	Name        string
	Args        Args
	ShouldError bool
	Error       string
}

type checkerTests []checkerTest

func runCheckerTests(parentT *testing.T, tests checkerTests, argable ArgFunc) {
	for _, test := range tests {
		parentT.Run(test.Name, func(t *testing.T) {
			checker, err := argable(test.Args)
			if err != nil && test.ShouldError {
				if test.Error != "" && err.Error() != test.Error {
					t.Errorf("%s: Expected error message: %s Got: %s", test.Name, test.Error, err.Error())
				}

				t.Logf("%s: Error creating checker: %s", test.Name, err)
				return
			} else if err != nil {
				t.Errorf("%s: Unexpected error: %s", test.Name, err)
				return
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
		})

	}
}
