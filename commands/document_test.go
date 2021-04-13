package commands

import "testing"

func TestDocumentCommand(t *testing.T) {

	Root.SetArgs([]string{"document"})
	err := Root.Execute()
	if err != nil {
		t.Errorf("Error running command: %s", err)
	}
}
