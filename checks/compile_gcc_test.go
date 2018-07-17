package checks

import (
	"testing"
)

func TestCompileGcc(t *testing.T) {

	tests := []struct {
		Name        string
		Args        map[string]interface{}
		ShouldError bool
	}{
		{
			Name: "gcc compile hello world",
			Args: map[string]interface{}{
				"source": `#include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
			},
		},
		{
			Name: "gcc compile bad hello world",
			Args: map[string]interface{}{
				"source": `include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
			},
			ShouldError: true,
		},
		{
			Name: "gcc compile hello world with args",
			Args: map[string]interface{}{
				"source": `#include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
				"cflags": "-g",
			},
		},
		{
			Name: "gcc compile hello world with bad args",
			Args: map[string]interface{}{
				"source": `include <stdio.h>

				int main() {
					printf("Hello world!\n");
				}`,
				"cflags": "-fake",
			},
			ShouldError: true,
		},
	}

	for _, test := range tests {
		checker, err := CompileGcc{}.FromArgs(test.Args)
		if err != nil {
			t.Errorf("%s: Unxpected error %s", test.Name, err)
		}

		err = checker.Check()
		if err != nil && !test.ShouldError {
			t.Errorf("%s: Got err when didn't expect one: %s", test.Name, err)
		} else if err == nil && test.ShouldError {
			t.Errorf("%s: Didn't get err when expected one: %s", test.Name, err)
		}
	}

}
