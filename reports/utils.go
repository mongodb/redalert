package reports

func IsValidSystemType(systemType string) bool {
	if _, ok := externalCommands[systemType]; !ok {
		return false
	}

	return true
}
