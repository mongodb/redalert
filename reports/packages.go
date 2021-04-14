package reports

import (
	"errors"
	"os/exec"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type externalCommand []string

var externalCommands = map[string]externalCommand{
	"linux": []string{"dpkg-query", "-W -f='${binary:Package};${Version}\n'"}}

func GetPackagesDetails(systemtype string) ([]Package, error) {
	if !isValidSystemType(systemtype) {
		return nil, errors.New("System type not supported: " + systemtype)
	}
	externalCommand := externalCommands[systemtype]
	command := exec.Command(externalCommand[0], externalCommand[1:]...)
	packageDetails, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}

	packageDetailsarsed := parseCommandOuput(string(packageDetails), systemtype)
	return packageDetailsarsed, nil
}
