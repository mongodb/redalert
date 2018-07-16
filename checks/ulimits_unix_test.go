package checks

import (
	"syscall"
	"testing"
)

func TestEveryType(t *testing.T) {
	for limitName, _ := range limitsByName {
		err := UlimitChecker{IsHard: false, Item: limitName, Limit: 0}.Check()
		if err != nil {
			t.Error(err)
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

	err = UlimitChecker{IsHard: true, Item: "files", Limit: rLimit.Max - 1}.Check()

	if err != nil {
		t.Error(err)
	}

	err = UlimitChecker{IsHard: false, Item: "files", Limit: rLimit.Cur - 1}.Check()

	if err != nil {
		t.Error(err)
	}

}
func TestArgBuilding(t *testing.T) {
	arged, err := UlimitChecker{}.FromArgs(map[string]interface{}{"item": "as", "type": "soft", "limit": 1024})
	if err != nil {
		t.Error(err)
	}

	checker, ok := arged.(UlimitChecker)
	if !ok {
		t.Error("Expected a ulimit checker")
	}

	if checker.Item != "as" {
		t.Errorf("Item (%s) not set correctly from args (as)", checker.Item)
	}

	if checker.IsHard != false {
		t.Errorf("IsHard (%t) not set correctly from args (soft)", checker.IsHard)
	}

	if checker.Limit != 1024 {
		t.Errorf("Limit (%d) not set correctly from args (1024)", checker.Limit)
	}

	arged, err = UlimitChecker{}.FromArgs(map[string]interface{}{"item": "files", "type": "hard", "limit": -1})
	if err != nil {
		t.Error(err)
	}

	checker, ok = arged.(UlimitChecker)
	if !ok {
		t.Error("Expected a ulimit checker")
	}

	if checker.Item != "files" {
		t.Errorf("Item (%s) not set correctly from args (files)", checker.Item)
	}

	if checker.IsHard != true {
		t.Errorf("IsHard (%t) not set correctly from args (hard)", checker.IsHard)
	}

	if checker.Limit != syscall.RLIM_INFINITY {
		t.Errorf("Limit (%d) not set correctly from args (syscall.RLIM_INFINITY)", checker.Limit)
	}

	arged, err = UlimitChecker{}.FromArgs(map[string]interface{}{"item": "core", "limit": 1024})
	if err != nil {
		t.Error(err)
	}

	checker, ok = arged.(UlimitChecker)
	if !ok {
		t.Error("Expected a ulimit checker")
	}

	if checker.Item != "core" {
		t.Errorf("Item (%s) not set correctly from args (core)", checker.Item)
	}

	if checker.IsHard != false {
		t.Errorf("IsHard (%t) not set correctly from args (none)", checker.IsHard)
	}

	if checker.Limit != 1024 {
		t.Errorf("Limit (%d) not set correctly from args (1024)", checker.Limit)
	}

}

// The stack size limit is a good limit for this test because it commonly
// has a non-inifinity hard limit (65532 in sample MacOS/RHEL 6 shells).
func TestFailure(t *testing.T) {
	checker := UlimitChecker{IsHard: false, Item: "stack", Limit: syscall.RLIM_INFINITY}
	err := checker.Check()
	if err == nil {
		t.Errorf("Didn't get err when expected one: %+v", checker)
	}
}

func TestLegacyTypes(t *testing.T) {
	checker, _ := availableChecks["open-files"](map[string]interface{}{"value": 1024})
	equivalentChecker := UlimitChecker{IsHard: true, Item: "nofile", Limit: 1024, Type: "hard"}

	if checker != equivalentChecker {
		t.Errorf("Legacy conversion failed (%+v != %+v)", checker, equivalentChecker)
	}

	checker, _ = availableChecks["address-size"](map[string]interface{}{"value": 1024})
	equivalentChecker = UlimitChecker{IsHard: true, Item: "as", Limit: 1024, Type: "hard"}

	if checker != equivalentChecker {
		t.Errorf("Legacy conversion failed (%+v != %+v)", checker, equivalentChecker)
	}

}
