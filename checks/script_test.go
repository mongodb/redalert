// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"os/exec"
	"testing"
)

func TestRunUnixShellScript(t *testing.T) {
	// Check if sh in path since it's required
	_, err := exec.LookPath("sh")
	if err != nil {
		return
	}

	// Check if bash in path since it's required
	_, err = exec.LookPath("bash")
	if err != nil {
		return
	}

	tests := checkerTests{
		{
			Name: "exit 0 check for exit code",
			Args: Args{
				"source": "exit 0",
			},
		},
		{
			Name: "exit 1 check for exit code",
			Args: Args{
				"source": "exit 1",
			},
			ShouldError: true,
		},
		{
			Name: "run bad command",
			Args: Args{
				"source": "ls -l /fake_dir",
			},
			ShouldError: true,
		},
		{
			Name: "echo 123 checks for 123",
			Args: Args{
				"output": "123",
				"source": "echo 123",
			},
		},
		{
			Name: "echo 123 checks for 111",
			Args: Args{
				"output": "111",
				"source": "echo 123",
			},
			ShouldError: true,
		},
		{
			Name: "echo 123 with good interpreter",
			Args: Args{
				"source":      "echo 123",
				"interpreter": "/bin/sh",
			},
		},
		{
			Name: "echo 123 with bad interpreter",
			Args: Args{
				"source":      "echo 123",
				"interpreter": "/bin/shhh",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["run-bash-script"])
}

func TestRunPythonScript(t *testing.T) {
	// Check for python before running these tests
	if _, err := exec.LookPath("python"); err != nil {
		return
	}

	tests := checkerTests{
		{
			Name: "python print 123 expecting 123",
			Args: Args{
				"source": "print('123')",
				"output": "123",
			},
		},
		{
			Name: "python print 123 expecting 111",
			Args: Args{
				"source": "print('123')",
				"output": "111",
			},
			ShouldError: true,
		},
		{
			Name: "python exit 0 check for exit code",
			Args: Args{
				"source": "exit(0)",
			},
		},
		{
			Name: "python exit 1 check for exit code",
			Args: Args{
				"source": "exit(1)",
			},
			ShouldError: true,
		},
		{
			Name: "python import module check exit code",
			Args: Args{
				"source": `import datetime
print(datetime.date.today())`,
			},
		},
		{
			Name: "python import bad module check exit code",
			Args: Args{
				"source": "import datetimes",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["run-python-script"])
}
