// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package reports

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetRevision(t *testing.T) {
	err := ioutil.WriteFile("rev_python.txt", []byte("123456789"), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("rev_python.txt")

	toolchains := make(map[string]string)
	toolchains["python"] = "rev_python.txt"

	toolchainDetails := GetToolchainDetails(toolchains)
	if toolchainDetails["python"] != "123456789" {
		t.Error("failed")
	}
}
