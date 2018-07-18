// +build !windows

package checks

import (
	"fmt"
	"syscall"
)

func init() {
	availableChecks["ulimit"] = func(args map[string]interface{}) (Checker, error) {
		return UlimitChecker{}.FromArgs(args)
	}

	// Legacy greenbay types
	availableChecks["open-files"] = func(args map[string]interface{}) (Checker, error) {
		args["item"] = "nofile"
		args["type"] = "hard"
		return UlimitChecker{}.FromArgs(args)
	}

	availableChecks["address-size"] = func(args map[string]interface{}) (Checker, error) {
		args["item"] = "as"
		args["type"] = "hard"
		return UlimitChecker{}.FromArgs(args)
	}
}

// UlimitChecker checks if current process resource limits are above a given minimum
//
// Type:
//   - ulimit
//
// Supported Platforms:
//   - MacOS
//   - Linux
//
// Arguments:
//   - item (required): A string value representing the type of limit to check
//   - value (required): Numerical value representing the value to be tested
//   - type: "hard" or "soft" with a default of "hard"
//   - greater_than: If provided will verify that the limit is greater than or
//                   equal to value instead of strictly equal to
//
// Notes:
//   - "item" strings are from http://www.linux-pam.org/Linux-PAM-html/sag-pam_limits.html
//     - The following values are supported:
//       - core
//       - data
//       - fsize
//       - nofile
//       - stack
//       - cpu
//       - as
//   - "value" can be '-1' to represent that the resource limit should be unlimited
type UlimitChecker struct {
	Item        string
	Value       int
	IsHard      bool
	GreaterThan bool `mapstructure:"greater_than"`
	Type        string
}

// Map symbolic limit names to rlimit constants
var limitsByName = map[string]int{
	"core":   syscall.RLIMIT_CORE,
	"data":   syscall.RLIMIT_DATA,
	"fsize":  syscall.RLIMIT_FSIZE,
	"nofile": syscall.RLIMIT_NOFILE,
	"stack":  syscall.RLIMIT_STACK,
	"cpu":    syscall.RLIMIT_CPU,
	"as":     syscall.RLIMIT_AS,
}

// Check if a ulimit is high enough
func (uc UlimitChecker) Check() error {
	limit, ok := limitsByName[uc.Item]
	if !ok {
		return fmt.Errorf("%s is not a valid limit name", uc.Item)
	}

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(limit, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}

	var LimitToCheck int
	if uc.IsHard {
		LimitToCheck = int(rLimit.Max)
	} else {
		LimitToCheck = int(rLimit.Cur)
	}

	if uc.Value == int(syscall.RLIM_INFINITY) && LimitToCheck != int(syscall.RLIM_INFINITY) {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is lower than required (unlimited)", uc.Type, LimitToCheck, uc.Item)
	} else if uc.GreaterThan && !(LimitToCheck < 0 || LimitToCheck >= uc.Value) {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is lower than required (%d)", uc.Type, LimitToCheck, uc.Item, uc.Value)
	} else if !uc.GreaterThan && LimitToCheck != uc.Value {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is not equal to %d", uc.Type, LimitToCheck, uc.Item, uc.Value)
	}

	return nil
}

// FromArgs will populate the UlimitChecker with the args given in the tests YAML
// config
func (uc UlimitChecker) FromArgs(args map[string]interface{}) (Checker, error) {
	if err := requiredArgs(args, "item", "value"); err != nil {
		return nil, err
	}

	if err := decodeFromArgs(args, &uc); err != nil {
		return nil, err
	}

	if uc.Value == -1 {
		uc.Value = int(syscall.RLIM_INFINITY)
	} else if uc.Value < 0 {
		return nil, fmt.Errorf("negative values other than -1 are invalid, got: %d", uc.Value)
	}

	uc.IsHard = !(uc.Type == "soft" || uc.Type == "")
	return uc, nil
}
