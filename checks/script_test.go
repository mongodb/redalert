package checks

import (
	"testing"
)

func TestRunScript(t *testing.T) {

	tests := []struct {
		Name        string
		Args        map[string]interface{}
		ShouldError bool
	}{
		{
			Name: "exit 0 check for exit code",
			Args: map[string]interface{}{
				"source": "exit 0",
			},
		},
		{
			Name: "exit 1 check for exit code",
			Args: map[string]interface{}{
				"source": "exit 1",
			},
			ShouldError: true,
		},
		{
			Name: "run bad command",
			Args: map[string]interface{}{
				"source": "ls -l /fake_dir",
			},
			ShouldError: true,
		},
		{
			Name: "echo 123 checks for 123",
			Args: map[string]interface{}{
				"output": "123",
				"source": "echo 123",
			},
		},
		{
			Name: "echo 123 checks for 111",
			Args: map[string]interface{}{
				"output": "111",
				"source": "echo 123",
			},
			ShouldError: true,
		},
		{
			Name: "echo 123 with good interpreter",
			Args: map[string]interface{}{
				"source":      "echo 123",
				"interpreter": "/bin/sh",
			},
		},
		{
			Name: "echo 123 with bad interpreter",
			Args: map[string]interface{}{
				"source":      "echo 123",
				"interpreter": "/bin/shhh",
			},
			ShouldError: true,
		},
		{
			Name: "python print 123 expecting 123",
			Args: map[string]interface{}{
				"source":      "print('123')",
				"interpreter": "python",
				"output":      "123",
			},
		},
		{
			Name: "python print 123 expecting 111",
			Args: map[string]interface{}{
				"source":      "print('123')",
				"interpreter": "python",
				"output":      "111",
			},
			ShouldError: true,
		},
		{
			Name: "python exit 0 check for exit code",
			Args: map[string]interface{}{
				"source":      "exit(0)",
				"interpreter": "python",
			},
		},
		{
			Name: "python exit 1 check for exit code",
			Args: map[string]interface{}{
				"source":      "exit(1)",
				"interpreter": "python",
			},
			ShouldError: true,
		},
		{
			Name: "python import module check exit code",
			Args: map[string]interface{}{
				"source": `import datetime
print(datetime.date.today())`,
				"interpreter": "python",
			},
		},
		{
			Name: "python import bad module check exit code",
			Args: map[string]interface{}{
				"source":      "import datetimes",
				"interpreter": "python",
			},
			ShouldError: true,
		},
	}

	for _, test := range tests {
		checker, err := RunScript{}.FromArgs(test.Args)
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
