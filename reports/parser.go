package reports

import (
	"strings"
)

func parseDebianCommandOutput(commandOutput string) []Package {
	rows := strings.Split(commandOutput, "\n")
	var packages = []Package{}

	for _, row := range rows {
		if row == "" {
			continue
		}

		rowElements := strings.Split(row, ";")
		p := Package{Name: rowElements[0]}
		if len(rowElements) > 1 {
			p.Version = rowElements[1]
		}

		packages = append(packages, p)
	}

	return packages
}

func parseCommandOuput(commandOutput string, packageManager string) []Package {
	if packageManager == "dpkg" {
		return parseDebianCommandOutput(commandOutput)
	}

	return []Package{}
}
