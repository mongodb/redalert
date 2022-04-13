// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"fmt"
)

var availableChecks = map[string]ArgFunc{}

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
