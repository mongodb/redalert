package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mongodb/redalert/checks"
	"github.com/mongodb/redalert/testfile"
	yaml "gopkg.in/yaml.v2"
)

func TestDocumentCommand(t *testing.T) {

	err := ioutil.WriteFile("toolchain_test.txt", []byte("revision12345"), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("toolchain_test.txt")

	tf := testfile.TestFile{
		Tests: []testfile.Test{
			{
				Name:   "toolchains",
				Suites: []string{"any"},
				Type:   "toolchains",
				Args: checks.Args{
					"path": "toolchain_test.txt",
				},
			},
		},
	}

	content, err := yaml.Marshal(tf)
	err = ioutil.WriteFile("reports.yml", content, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("reports.yml")

	Root.SetArgs([]string{"document", "--suite", "any"})
	err = Root.Execute()
	if err != nil {
		t.Errorf("Error running command: %s", err)
	}

}
