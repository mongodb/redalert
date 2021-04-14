package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func parseDebianCommandOutput(commandOutput string) []Package {
	debex, err := ioutil.ReadFile("/Users/tural.farhadov/work/redalert/deb.out")
	if err != nil {
		fmt.Print(err.Error())
	}
	commandOutput = string(debex)

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

func parseCommandOuput(commandOutput string, systemtype string) []Package {
	if systemtype == "debian" {
		return parseDebianCommandOutput(commandOutput)
	}

	//Testing
	return parseDebianCommandOutput(commandOutput)

	return []Package{}
}

func formatPacakges(packages []Package, format string) (string, error) {
	if format == "json" {
		return formatPackagesIntoJson(packages), nil
	}

	return "", errors.New("format not supported: " + format)
}

func formatPackagesIntoJson(packages []Package) string {
	res, err := json.Marshal(packages)
	if err != nil {
		fmt.Println("Could not marshal")
		return ""
	}

	jsonresult := fmt.Sprintf("{\"packages\": %s}", res)

	return jsonresult
}
