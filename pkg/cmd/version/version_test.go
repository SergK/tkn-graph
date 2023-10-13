package version

import (
	"bytes"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := Command()

	// Create a Buffer to capture the output.
	out := new(bytes.Buffer)
	cmd.SetOut(out)

	// Execute the command.
	if err := cmd.Execute(); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	// Assert that the command is valid.
	if cmd == nil || cmd.Name() != "version" {
		t.Errorf("Command is not valid: %v", cmd)
	}
}
