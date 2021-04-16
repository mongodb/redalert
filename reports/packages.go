package reports

import (
	"errors"
	"os/exec"
)

type externalCommand []string

var externalCommands = map[string]externalCommand{
	"dpkg":   []string{"dpkg-query", "-W", "-f=${binary:Package};${Version}|"},
	"rpm":    []string{"rpm", "-qa", "--queryformat", "%{NAME};%{VERSION}|"},
	"zypper": []string{"bash", "-c", "zypper search  --installed-only -s | awk 'NR>5 {printf $3\";\"$7\"|\"}'"},
}

func GetPackagesDetails(packageManager string) (map[string]string, error) {
	if !isSupported(packageManager) {
		return nil, errors.New("Package manager is not supported: " + packageManager)
	}
	externalCommand := externalCommands[packageManager]
	command := exec.Command(externalCommand[0], externalCommand[1:]...)
	packageDetails, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}

	packageDetailsarsed := parseCommandOuput(string(packageDetails), packageManager)
	return packageDetailsarsed, nil
}
