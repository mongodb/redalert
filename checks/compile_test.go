// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.

package checks

import (
	"testing"
)

func TestCompileGcc(t *testing.T) {

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
