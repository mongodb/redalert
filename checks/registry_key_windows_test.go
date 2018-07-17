package checks

func TestRegistryKeyChecker(t *testing.T) {
	tests := checkerTests{
		{
			Name: "check dword value",
			Args: map[string]interface{}{
				"root": "HKLM",
				"path": "SYSTEM\\CurrentControlSet\\services\\LanmanServer\\Parameters",
				"key": "IRPStackSize",
				"value_type": "DWORD",
				"value": "3",
			},
		},
	}

	runCheckerTests(t, tests, availableChecks["registry"])
}
