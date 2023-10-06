package pipelinerun

import (
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/flags"
)

func Command(p cli.Params) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipelinerun",
		Aliases: []string{"pr", "pipelineruns"},
		Short:   "Graph PipelineRuns",
		Annotations: map[string]string{
			"commandType": "main",
		},
	}

	flags.AddTektonOptions(cmd)
	cmd.AddCommand(
		graphCommand(p),
	)
	return cmd
}
