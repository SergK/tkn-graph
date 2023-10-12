package cmd

import (
	"bytes"
	"testing"

	"github.com/tektoncd/cli/pkg/cli"
)

func TestRoot(t *testing.T) {
	// Create a new cobra command.
	cmd := Root(&cli.TektonParams{})

	// Create a Buffer to capture the output.
	out := new(bytes.Buffer)
	cmd.SetOut(out)

	// Execute the command.
	if err := cmd.Execute(); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	// Assert that the command is valid.
	if cmd == nil || cmd.Name() != "tkn-graph" {
		t.Errorf("Command is not valid: %v", cmd)
	}

	// Assert that the command has the expected subcommands.
	if len(cmd.Commands()) != 4 {
		t.Errorf("Command does not have the expected subcommands: %v", cmd.Commands())
	}
}
