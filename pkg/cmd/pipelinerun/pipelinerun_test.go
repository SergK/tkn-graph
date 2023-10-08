package pipelinerun

import (
	"testing"

	"github.com/tektoncd/cli/pkg/cli"
)

func TestRoot(t *testing.T) {
	// Create a new cobra command.
	cmd := Command(&cli.TektonParams{})

	// Execute the command.
	if err := cmd.Execute(); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	// Assert that the command is valid.
	if cmd == nil || cmd.Name() != "pipelinerun" {
		t.Errorf("Command is not valid: %v", cmd)
	}

	// Assert that the command has the expected subcommands.
	if len(cmd.Commands()) != 3 {
		t.Errorf("Command does not have the expected subcommands: %v", cmd.Commands())
	}
}
