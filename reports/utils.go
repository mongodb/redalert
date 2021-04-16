package reports

func isSupported(packageManager string) bool {
	if _, ok := externalCommands[packageManager]; !ok {
		return false
	}

	return true
}
