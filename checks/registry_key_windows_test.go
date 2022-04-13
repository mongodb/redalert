// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import "testing"

func TestRegistryKeyChecker(t *testing.T) {
	tests := checkerTests{
		{
			Name: "check dword value",
			Args: Args{
				"root":       "HKLM",
				"path":       "SYSTEM\\CurrentControlSet\\services\\LanmanServer\\Parameters",
				"key":        "Size",
				"value_type": "DWORD",
				"value":      3,
			},
		},
		{
			Name: "check path exists",
			Args: Args{
				"root": "HKEY_LOCAL_MACHINE",
				"path": "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion",
			},
		},
		{
			Name: "check path doesn't exist",
			Args: Args{
				"root": "HKEY_LOCAL_MACHINE",
				"path": "SOFTWARE\\Microsoft\\Windows NT\\NotARealPath",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["registry"])
}
