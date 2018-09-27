// Copyright 2018 MongoDB Inc. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"os/exec"
	"testing"
)

func TestCompileGcc(t *testing.T) {
	// First make sure gcc is in the PATH
	// Don't run these tests unless you are on a system with gcc installed
	_, err := exec.LookPath("gcc")
	if err != nil {
		return
	}

	tests := checkerTests{
		{
			Name: "gcc compile hello world",
			Args: Args{
				"source": `#include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
			},
		},
		{
			Name: "gcc compile bad hello world",
			Args: Args{
				"source": `include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
			},
			ShouldError: true,
		},
		{
			Name: "gcc compile hello world with args",
			Args: Args{
				"source": `#include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
				"cflags": "-g",
			},
		},
		{
			Name: "gcc compile hello world with bad args",
			Args: Args{
				"source": `include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
				"cflags": "-fake",
			},
			ShouldError: true,
		},
		{
			Name: "gcc compile hello world with cflags_command",
			Args: Args{
				"source": `#include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
				"cflags_command": "/bin/bash -c 'echo -g'",
			},
		},
		{
			Name: "gcc compile hello world with bad cflags_command",
			Args: Args{
				"source": `#include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
				"cflags_command": "/bin/false",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["compile-gcc"])
}
