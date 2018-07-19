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
