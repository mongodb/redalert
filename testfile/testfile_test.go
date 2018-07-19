// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved. Use of this source code is
// governed by the Apache-2.0 license that can be found in the LICENSE file.


package testfile

import "testing"

func TestTestFileSelection(t *testing.T) {
	tests := []struct {
		Name          string
		TestFile      TestFile
		Suite         string
		ExpectedNames []string
	}{
		{
			Name:          "simple selection",
			Suite:         "ubuntu1404",
			ExpectedNames: []string{"should select"},
			TestFile: TestFile{
				Tests: []Test{
					{
						Name:   "should select",
						Type:   "file-exists",
						Suites: []string{"ubuntu1404"},
					},
					{
						Name:   "should not select",
						Type:   "file-exists",
						Suites: []string{"ubuntu1604"},
					},
				},
			},
		},
		{
			Name:          "alias selection",
			Suite:         "ubuntu",
			ExpectedNames: []string{"should select", "should also select"},
			TestFile: TestFile{
				Aliases: Aliases{
					"ubuntu": []string{"ubuntu1404", "ubuntu1604"},
				},
				Tests: []Test{
					{
						Name:   "should select",
						Type:   "file-exists",
						Suites: []string{"ubuntu1404"},
					},
					{
						Name:   "should also select",
						Type:   "file-exists",
						Suites: []string{"ubuntu1604"},
					},
				},
			},
		},
		{
			Name:          "alias only suite selection",
			Suite:         "ubuntu1404",
			ExpectedNames: []string{"should select"},
			TestFile: TestFile{
				Aliases: Aliases{
					"ubuntu": []string{"ubuntu1404, ubuntu1604"},
				},
				Tests: []Test{
					{
						Name:   "should select",
						Type:   "file-exists",
						Suites: []string{"ubuntu1404"},
					},
					{
						Name:   "should not select",
						Type:   "file-exists",
						Suites: []string{"ubuntu1604"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		matched := test.TestFile.TestsToRun(test.Suite)
		if len(matched) != len(test.ExpectedNames) {
			t.Errorf("%s: matched is incorrect length expected: %d got: %d", test.Name, len(matched), len(tests))
		}

		for i := range matched {
			if matched[i].Name != test.ExpectedNames[i] {
				t.Errorf("Expected: %s Got: %s", test.ExpectedNames[i], matched[i].Name)
			}
		}
	}
}
