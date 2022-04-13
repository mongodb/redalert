// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

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
			Name: "check open files is equal",
			Args: Args{
				"type":  "soft",
				"item":  "nofile",
				"value": openFilesValue,
			},
		},
		{
			Name: "check open files is wrong value",
			Args: Args{
				"type":         "hard",
				"item":         "nofile",
				"value":        openFilesValue - 1,
				"greater_than": true,
			},
		},
		{
			Name: "stack should fail",
			Args: Args{
				"value": openFilesValue - 1,
				"item":  "nofile",
				"type":  "soft",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["ulimit"])
}

func TestEveryType(t *testing.T) {
	for limitName := range limitsByName {
		err := UlimitChecker{IsHard: false, Item: limitName, Value: 0, GreaterThan: true}.Check()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestArgBuilding(t *testing.T) {
	tests := []struct {
		Args     Args
		Expected UlimitChecker
	}{
		{
			Args: Args{
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
			Args: Args{
				"item":  "nofile",
				"type":  "hard",
				"value": -1,
			},
			Expected: UlimitChecker{
				IsHard: true,
				Item:   "nofile",
				Type:   "hard",
				Value:  int(syscall.RLIM_INFINITY),
			},
		},
	}

	for _, test := range tests {
		arged, err := UlimitCheckerFromArgs(test.Args)
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

	err = UlimitChecker{IsHard: true, Item: "nofile", GreaterThan: true, Value: int(rLimit.Max - 1)}.Check()

	if err != nil {
		t.Error(err)
	}

	err = UlimitChecker{IsHard: false, Item: "nofile", GreaterThan: true, Value: int(rLimit.Cur - 1)}.Check()

	if err != nil {
		t.Error(err)
	}

}

func TestLegacyTypes(t *testing.T) {
	checker, _ := availableChecks["open-files"](Args{"value": 1024})
	equivalentChecker := UlimitChecker{IsHard: true, Item: "nofile", Value: 1024, Type: "hard"}

	if checker != equivalentChecker {
		t.Errorf("Legacy conversion failed (%+v != %+v)", checker, equivalentChecker)
	}

	checker, _ = availableChecks["address-size"](Args{"value": 1024})
	equivalentChecker = UlimitChecker{IsHard: true, Item: "as", Value: 1024, Type: "hard"}

	if checker != equivalentChecker {
		t.Errorf("Legacy conversion failed (%+v != %+v)", checker, equivalentChecker)
	}

}
