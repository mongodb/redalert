package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chasinglogic/redalert/testfile"
	yaml "gopkg.in/yaml.v2"
)

func TestRunCommand(t *testing.T) {
	err := ioutil.WriteFile("run_test.txt", []byte("test"), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("run_test.txt")

	tf := testfile.TestFile{
		Tests: []testfile.Test{
			{
				Name:   "should run",
				Suites: []string{"any"},
				Type:   "file-exists",
				Args: Args{
					"name": "run_test.txt",
				},
			},
			{
				Name:   "should also run",
				Suites: []string{"any"},
				Type:   "file-exists",
				Args: Args{
					"name": "run_test.txt",
				},
			},
			{
				Name:   "Should not run",
				Suites: []string{"NOPE"},
				Type:   "file-exists",
				Args: Args{
					"name": "NOPE.txt",
				},
			},
		},
	}

	content, err := yaml.Marshal(tf)
	err = ioutil.WriteFile("tests.yml", content, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("tests.yml")

	Root.SetArgs([]string{"run", "--suite", "any"})
	err = Root.Execute()
	if err != nil {
		t.Errorf("Error running command: %s", err)
	}
}
