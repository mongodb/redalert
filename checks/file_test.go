package checks

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {
	err := ioutil.WriteFile("test_exists.txt", []byte("test"), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("test_exists.txt")

	err = FileChecker{Exists: true, Name: "test_exists.txt"}.Check()
	if err != nil {
		t.Error(err)
	}

	checker, err := FileChecker{}.FromArgs(map[string]interface{}{"name": "test_exists.txt"})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}

func TestFileNotExists(t *testing.T) {
	err := FileChecker{Exists: false, Name: "test_not_exists.txt"}.Check()
	if err != nil {
		t.Error(err)
	}

	checker, err := FileChecker{}.FromArgs(map[string]interface{}{"name": "test_not_exists.txt", "exists": false})
	if err != nil {
		t.Error(err)
	}

	err = checker.Check()
	if err != nil {
		t.Error(err)
	}
}
