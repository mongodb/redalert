package checks

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAndChecker(t *testing.T) {
	err := ioutil.WriteFile("test_exists.txt", []byte("test"), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("test_exists.txt")

	tests := checkerTests{
		{
			Name: "should succeed when two tests succeed",
			Args: Args{
				"checks": []Args{
					{
						"type": "file-exists",
						"args": Args{
							"name": "test_exists.txt",
						},
					},
					{
						"type": "file-does-not-exist",
						"args": Args{
							"name": "totally_not_a_file",
						},
					},
				},
			},
		},
		{
			Name:        "should fail when one test fails",
			Error:       "totally_not_a_file doesn't exist and should",
			ShouldError: true,
			Args: Args{
				"checks": []Args{
					{
						"type": "file-exists",
						"args": Args{
							"name": "test_exists.txt",
						},
					},
					{
						"type": "file-exists",
						"args": Args{
							"name": "totally_not_a_file",
						},
					},
				},
			},
		},
	}

	runCheckerTests(t, tests, availableChecks["and"])
}

func TestOrChecker(t *testing.T) {
	err := ioutil.WriteFile("test_exists.txt", []byte("test"), os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("test_exists.txt")

	tests := checkerTests{
		{
			Name: "should succeed when both tests succeed",
			Args: Args{
				"checks": []Args{
					{
						"type": "file-exists",
						"args": Args{
							"name": "test_exists.txt",
						},
					},
					{
						"type": "file-does-not-exist",
						"args": Args{
							"name": "totally_not_a_file",
						},
					},
				},
			},
		},
		{
			Name: "should succeed when one test fails but other passes",
			Args: Args{
				"checks": []Args{
					{
						"type": "file-exists",
						"args": Args{
							"name": "totally_not_a_file",
						},
					},
					{
						"type": "file-exists",
						"args": Args{
							"name": "test_exists.txt",
						},
					},
				},
			},
		},
		{
			Name:        "should fail when all fail",
			ShouldError: true,
			Args: Args{
				"checks": []Args{
					{
						"type": "file-exists",
						"args": Args{
							"name": "totally_not_a_file",
						},
					},
					{
						"type": "file-exists",
						"args": Args{
							"name": "totally_not_a_file",
						},
					},
				},
			},
		},
	}

	runCheckerTests(t, tests, availableChecks["or"])
}
