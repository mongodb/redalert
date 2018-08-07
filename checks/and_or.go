package checks

import (
	"errors"
	"fmt"
)

func init() {
	availableChecks["and"] = andFromArgs
	availableChecks["or"] = orFromArgs
}

// metaCheckArgs is used for any checker which is "meta" meaning that it's args
// are a list of other checks to run.
type metaCheckArgs struct {
	Type string
	Args Args
}

// Load will load the correct checker implementation for the Type of this checker.
func (mc metaCheckArgs) Load() (Checker, error) {
	argFunc, ok := availableChecks[mc.Type]
	if !ok {
		return nil, fmt.Errorf("%s is not an available check type", mc.Type)
	}

	return argFunc(mc.Args)
}

// AndCheck performs multiple checks in order verifying that all are true
// Useful for verifying a set of tests are performed serially
//
// Type:
//     - and
//
// Supported Platforms:
//     - Linux
//     - MacOS
//     - Windows
//
// Arguments:
//  checks (required): A list of yaml maps that are the same as regular tests
//                     an example can be found below.
//
// Notes:
//   An example usage if you wanted to and multiple file checks would be:
//     - name: bashrc and bash_profile exist
//       type: and
//       args:
//         checks:
//           - type: file-exists
//             args:
//               name: .bashrc
//           - type: file-exists
//             args:
//               name: .bash_profile
type AndCheck struct {
	Checks []metaCheckArgs
}

// Check will run all checks serially finding the first failure if there is one
// if no failures found it returns nil for a success
func (and AndCheck) Check() error {
	for _, check := range and.Checks {
		checker, err := check.Load()
		if err != nil {
			return err
		}

		if err := checker.Check(); err != nil {
			return err
		}
	}

	return nil
}

// OrCheck performs multiple checks in order verifying that at least one is true
// Useful for verifying a set of tests are performed serially
//
// Type:
//     - or
//
// Supported Platforms:
//     - Linux
//     - MacOS
//     - Windows
//
// Arguments:
//  checks (required): A list of yaml maps that are the same as regular tests
//                     an example can be found below.
//
// Notes:
//   An example usage if you wanted to and multiple file checks would be:
//     - name: totally_not_a_file or bash_profile exist
//       type: or
//       args:
//         checks:
//           - type: file-exists
//             args:
//               name: totally_not_a_file
//           - type: file-exists
//             args:
//               name: .bash_profile
type OrCheck struct {
	Checks []metaCheckArgs
}

// Check will run all checks finding the first success and returning it
// otherwise will return an error that is a concatenation of all failures found
func (or OrCheck) Check() error {
	errorMessages := ""

	for _, check := range or.Checks {
		checker, err := check.Load()
		if err != nil {
			return err
		}

		err = checker.Check()
		if err == nil {
			return nil
		}

		errorMessages += err.Error() + "\n"
	}

	return errors.New(errorMessages)
}

func andFromArgs(args Args) (Checker, error) {
	and := AndCheck{}

	if err := requiredArgs(args, "checks"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &and); err != nil {
		return nil, err
	}

	return and, nil
}

func orFromArgs(args Args) (Checker, error) {
	or := OrCheck{}

	if err := requiredArgs(args, "checks"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &or); err != nil {
		return nil, err
	}

	return or, nil
}
