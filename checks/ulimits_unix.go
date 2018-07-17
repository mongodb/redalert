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
//   - value (required): Numerical value representing the minimum value to be tested
//   - type: "hard" or "soft" with a default of "hard"
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
	Item   string
	Value  int
	IsHard bool
	Type   string
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
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(limitsByName[uc.Item], &rLimit)

	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}

	var LimitToCheck uint64
	if uc.IsHard {
		uc.Type = "hard"
		LimitToCheck = rLimit.Max
	} else {
		uc.Type = "soft"
		LimitToCheck = rLimit.Cur
	}

	if uint64(uc.Value) == syscall.RLIM_INFINITY && LimitToCheck != syscall.RLIM_INFINITY {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is lower than required (unlimited)", uc.Type, LimitToCheck, uc.Item)
	} else if LimitToCheck < uint64(uc.Value) {
		return fmt.Errorf("Process %s ulimit (%d) of type \"%s\" is lower than required (%d)", uc.Type, LimitToCheck, uc.Item, uc.Value)
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
