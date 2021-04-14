package reports

func isValidSystemType(systemType string) bool {
	if _, ok := externalCommands[systemType]; !ok {
		return false
	}

	return true
}
