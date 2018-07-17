package checks

import (
	"syscall"
	"testing"
)

func TestUlimitsChecks(t *testing.T) {
	// The open files limit is a good limit for this test because it commonly
	// has different soft/hard limits (4864/unlimited in sample MacOS shell,
	// 1024/4096 in sample RHEL 6 shell).
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		t.Error(err)
	}

	openFilesValue := int(rLimit.Cur)

	tests := checkerTests{
		{
			Name: "check open files",
			Args: map[string]interface{}{
				"type":  "hard",
				"item":  "files",
				"value": openFilesValue,
			},
		},
		{
			Name: "check open files is wrong value",
			Args: map[string]interface{}{
				"type":  "hard",
				"item":  "files",
				"value": openFilesValue - 1,
			},
		},
		{
			Name: "stack should fail",
			Args: map[string]interface{}{
				"value": int(syscall.RLIM_INFINITY),
				"item":  "stack",
				"type":  "soft",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["ulimit"])
}

func TestEveryType(t *testing.T) {
	for limitName := range limitsByName {
		err := UlimitChecker{IsHard: false, Item: limitName, Value: 0}.Check()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestArgBuilding(t *testing.T) {
	tests := []struct {
		Args     map[string]interface{}
		Expected UlimitChecker
	}{
		{
			Args: map[string]interface{}{
				"item":  "as",
				"type":  "soft",
				"value": 1024,
			},
			Expected: UlimitChecker{
				IsHard: false,
				Item:   "as",
				Value:  1024,
				Type:   "soft",
			},
		},
		{
			Args: map[string]interface{}{
				"item":  "files",
				"type":  "hard",
				"value": -1,
			},
			Expected: UlimitChecker{
				IsHard: true,
				Item:   "files",
				Type:   "hard",
				Value:  int(syscall.RLIM_INFINITY),
			},
		},
	}

	for _, test := range tests {
		arged, err := UlimitChecker{}.FromArgs(test.Args)
		if err != nil {
			t.Error(err)
		}

		checker, ok := arged.(UlimitChecker)
		if !ok {
			t.Error("Expected a ulimit checker")
		}

		if checker != test.Expected {
			t.Errorf("Expected: %v Got: %v", test.Expected, checker)
		}
	}
}

// The open files limit is a good limit for this test because it commonly
// has different soft/hard limits (4864/unlimited in sample MacOS shell,
// 1024/4096 in sample RHEL 6 shell).
func TestSoftHard(t *testing.T) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		t.Error(err)
	}

	err = UlimitChecker{IsHard: true, Item: "files", Value: int(rLimit.Max - 1)}.Check()

	if err != nil {
		t.Error(err)
	}

	err = UlimitChecker{IsHard: false, Item: "files", Value: int(rLimit.Cur - 1)}.Check()

	if err != nil {
		t.Error(err)
	}

}

func TestLegacyTypes(t *testing.T) {
	checker, _ := availableChecks["open-files"](map[string]interface{}{"value": 1024})
	equivalentChecker := UlimitChecker{IsHard: true, Item: "nofile", Value: 1024, Type: "hard"}

	if checker != equivalentChecker {
		t.Errorf("Legacy conversion failed (%+v != %+v)", checker, equivalentChecker)
	}

	checker, _ = availableChecks["address-size"](map[string]interface{}{"value": 1024})
	equivalentChecker = UlimitChecker{IsHard: true, Item: "as", Value: 1024, Type: "hard"}

	if checker != equivalentChecker {
		t.Errorf("Legacy conversion failed (%+v != %+v)", checker, equivalentChecker)
	}

}
