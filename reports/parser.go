package reports

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func parseDebianCommandOutput(commandOutput string) []Package {
	//Overriding the param for testing the pre-pupulated debian output
	debex, err := ioutil.ReadFile("/tmp/deb.out")
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

	//Only debian is supported as of now
	return parseDebianCommandOutput(commandOutput)

	return []Package{}
}
