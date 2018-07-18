package testfile

import "fmt"

// Aliases is used to map one suite name to multiple "suites"
type Aliases map[string][]string

// Matches will return all aliases which match the given suite
func (a Aliases) Matches(suite string) []string {
	matches := []string{}

	// TODO: Optimize this.
	for alias, suites := range a {
		for _, s := range suites {
			if s == suite {
				matches = append(matches, alias)
				break
			}
		}
	}

	return matches
}

// Test is a single check to run against a system.
type Test struct {
	Name   string
	Type   string
	Suites []string
	Args   map[string]interface{}
}

// Matches will return a boolean indicating whether this test should be run
// based on the given suiteNames
func (t Test) Matches(suiteNames []string) bool {
	// TODO: Optimize this.
	for _, s := range t.Suites {
		for _, name := range suiteNames {
			if s == name {
				return true
			}
		}
	}

	return false
}

// TestFile is the YAML file that contains tests and aliases to run
type TestFile struct {
	Aliases Aliases
	Tests   []Test
}

// Validate will ensure no test names are duplicated
func (tf TestFile) Validate() error {
	// Verify no names were duplicated in other.Tests
	var uniqueTests = make(map[string]struct{}, len(tf.Tests))

	for _, test := range tf.Tests {
		if _, ok := uniqueTests[test.Name]; ok {
			return fmt.Errorf("%s is duplicated, cannot have multiple tests with the same name", test.Name)
		}

		uniqueTests[test.Name] = struct{}{}
	}

	return nil
}

// Join returns the union of this testfile and another.  It is up to the caller
// to run Validate, since it's possible that multiple Joins could occur and
// calling Validate for each one would be expensive.
func (tf TestFile) Join(other TestFile) TestFile {
	if tf.Aliases == nil && other.Aliases != nil {
		tf.Aliases = other.Aliases
	} else if tf.Aliases != nil && other.Aliases != nil {
		for k, v := range other.Aliases {
			if _, exists := tf.Aliases[k]; exists {
				tf.Aliases[k] = append(tf.Aliases[k], v...)
			} else {
				tf.Aliases[k] = v
			}
		}
	} else {
		tf.Aliases = other.Aliases
	}

	tf.Tests = append(tf.Tests, other.Tests...)

	return tf
}

// TestsToRun will return a slice of Tests which shuld be run for the given
// suite.
func (tf TestFile) TestsToRun(suite string) []Test {
	testsToRun := []Test{}

	aliasedSuites := []string{suite}
	if aliases, ok := tf.Aliases[suite]; ok {
		aliasedSuites = aliases
	}

	for _, test := range tf.Tests {
		if test.Matches(aliasedSuites) {
			testsToRun = append(testsToRun, test)
		}
	}

	return testsToRun
}
