package checks

import "testing"

func TestRealVersionNumber(t *testing.T) {
	version := getRealVersionNumber(2010.0)
	if version != 10.0 {
		t.Errorf("Expected %f Got %f", 10.0, version)
	}

	version = getRealVersionNumber(10.0)
	if version != 10.0 {
		t.Errorf("Expected %f Got %f", 10.0, version)
	}

	version = getRealVersionNumber(2013.0)
	if version != 12.0 {
		t.Errorf("Expected %f Got %f", 12.0, version)
	}

	version = getRealVersionNumber(2015.0)
	if version != 14.0 {
		t.Errorf("Expected %f Got %f", 14.0, version)
	}

	version = getRealVersionNumber(2017.0)
	if version != 15.0 {
		t.Errorf("Expected %f Got %f", 15.0, version)
	}

	// Hopefully they'll only increment by one again for the next one.
	version = getRealVersionNumber(2019.0)
	if version != 16.0 {
		t.Errorf("Expected %f Got %f", 16.0, version)
	}
}

func TestCompileVisualStudio(t *testing.T) {
	tests := checkerTests{
		{
			Name: "visual-studio compile hello world",
			Args: Args{
				"source": `#include <iostream>

				int main() {
					std::cout << "Hello world!" << std::endl;
				}`,
			},
		},
		{
			Name: "visual-studio compile bad hello world",
			Args: Args{
				"source": `include <iostream>

				int main() {
					std::cout << "Hello world!" << std::endl;
				}`,
			},
			ShouldError: true,
		},
		{
			Name: "visual-studio compile hello world with args",
			Args: Args{
				"source": `#include <iostream>

				int main() {
					std::cout << "Hello world!" << std::endl;
				}`,
				"cflags": "/O1",
			},
		},
		{
			Name: "visual-studio compile hello world with bad args",
			Args: Args{
				"source": `include <iostream>

				int main() {
					std::cout << "Hello world!" << std::endl;
				}`,
				"cflags": "/NOPE",
			},
			ShouldError: true,
		},
	}

	runCheckerTests(t, tests, availableChecks["compile-visual-studio"])
}
