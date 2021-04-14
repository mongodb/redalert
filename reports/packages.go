package reports

import (
	"fmt"
	"os/exec"
)

type externalCommand []string

var externalCommands = map[string]externalCommand{"debian": []string{"dpkg", "-l"}, "macos": []string{"pkgutil", "--pkgs"}}

func Packages(systemtype string) (string, error) {

	if _, ok := externalCommands[systemtype]; !ok {
		return "", fmt.Errorf("system type not found: " + systemtype)
	}
	externalCommand := externalCommands[systemtype]

	command := exec.Command(externalCommand[0], externalCommand[1:]...)

	commandRes, err := command.CombinedOutput()

	if err != nil {
		return "", err
	}

	fmt.Println(string(commandRes))
	return string(commandRes), nil
}
