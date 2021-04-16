package reports

import (
	"strings"
)

func parseLinuxOutput(commandOutput string) map[string]string {
	rows := strings.Split(commandOutput, "|")
	packages := make(map[string]string)

	for _, row := range rows {
		if row == "" {
			continue
		}

		rowElements := strings.Split(row, ";")
		if len(rowElements) > 1 {
			packages[rowElements[0]] = rowElements[1]
		}

	}

	return packages
}

func parseCommandOuput(commandOutput string, packageManager string) map[string]string {
	switch packageManager {
	case
		"dpkg",
		"rpm",
		"zypper":
		return parseLinuxOutput(commandOutput)
	}

	return nil
}
