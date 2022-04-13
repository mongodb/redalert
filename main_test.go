// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mongodb/redalert/checks"
	"github.com/mongodb/redalert/testfile"
	yaml "gopkg.in/yaml.v2"
)

func TestRunCommand(t *testing.T) {
	dir, err := os.MkdirTemp("", "redalert")
	if err != nil {
		t.Fatalf("create temp: %s", err)
	}
	defer os.RemoveAll(dir)

	testFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), os.ModePerm); err != nil {
		t.Fatalf("write test data: %s", err)
	}

	tf := testfile.TestFile{
		Tests: []testfile.Test{
			{
				Name:   "should run",
				Suites: []string{"any"},
				Type:   "file-exists",
				Args: checks.Args{
					"name": testFile,
				},
			},
			{
				Name:   "should also run",
				Suites: []string{"any"},
				Type:   "file-exists",
				Args: checks.Args{
					"name": testFile,
				},
			},
			{
				Name:   "Should not run",
				Suites: []string{"NOPE"},
				Type:   "file-exists",
				Args: checks.Args{
					"name": "NOPE.txt",
				},
			},
		},
	}

	content, err := yaml.Marshal(tf)
	if err != nil {
		t.Fatalf("yaml marshal: %s", err)
	}

	testYaml := filepath.Join(dir, "tests.yml")
	if err := os.WriteFile(testYaml, content, os.ModePerm); err != nil {
		t.Fatalf("write tests.yml: %s", err)
	}

	bin := filepath.Join(dir, "bin")

	if err := compileBinary(bin); err != nil {
		t.Fatalf("compile binary: %s", err)
	}

	cmd := exec.Command(bin, "run", "--suite", "any", "--file", testYaml)

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run redalert: %s: %s", err, out)
	}
}

func compileBinary(path string) error {
	fmt.Println(os.Getwd())
	cmd := exec.Command("go", "build", "-o", path)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build go binary: %s: %w", out, err)
	}

	return nil
}
