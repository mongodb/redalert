// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package reports

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadRevision(fn string) (string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string

	// We expect the Revision to be on the first line
	line, err = bf.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")
	if err != nil {
		return "", err
	}

	if line == "" {
		return "", fmt.Errorf("empty")
	}
	return line, nil
}

func GetToolchainDetails() map[string]string {
	toolchains := make(map[string]string)
	toolchains["mongodb"] = "/tmp/revision.txt"
	toolchains["python"] = "/tmp/rev_python.txt"

	toolchainDetails := make(map[string]string)
	for name, path := range toolchains {
		revision, err := ReadRevision(path)
		if err != nil {
			fmt.Println(err)
			toolchainDetails[name] = "none"
		} else {
			toolchainDetails[name] = revision
		}
	}

	return toolchainDetails
}
