package reports

import (
	"errors"
	"fmt"
	"os/exec"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type externalCommand []string

var externalCommands = map[string]externalCommand{
	"dpkg": []string{"dpkg-query", "-W", "-f='${binary:Package};${Version}\n'"}}

func GetPackagesDetails(packageManager string) ([]Package, error) {
	if !isSupported(packageManager) {
		return nil, errors.New("Package manager is not supported: " + packageManager)
	}
	externalCommand := externalCommands[packageManager]
	command := exec.Command(externalCommand[0], externalCommand[1:]...)
	fmt.Println(command.String())
	packageDetails, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}

	packageDetailsarsed := parseCommandOuput(string(packageDetails), packageManager)
	return packageDetailsarsed, nil
}
