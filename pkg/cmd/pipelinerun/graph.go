package pipelinerun

import (
	"fmt"

	"github.com/sergk/tkn-graph/pkg/cli/prerun"
	pipelinerunpkg "github.com/sergk/tkn-graph/pkg/pipelinerun"
	"github.com/sergk/tkn-graph/pkg/taskgraph"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/flags"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type graphOptions struct {
	OutputFormat string
	OutputDir    string
	WithTaskRef  bool
}

func graphCommand(p cli.Params) *cobra.Command {

	opts := &graphOptions{}
	// Define the root command
	c := &cobra.Command{
		Use:     "graph",
		Aliases: []string{"g"},
		Short:   "Generates Graph",
		Annotations: map[string]string{
			"commandType": "main",
		},
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			// Add global args to the args list
			if err := flags.InitParams(p, cmd); err != nil {
				return err
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return prerun.ValidateGraphPreRunE(opts.OutputFormat)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cs, err := p.Clients()
			if err != nil {
				return err
			}

			var graphs []*taskgraph.TaskGraph
			var pipelineruns []v1.PipelineRun

			switch len(args) {
			case 1:
				var pipelinerun *v1.PipelineRun
				pipelinerun, err = pipelinerunpkg.GetPipelineRunsByName(cs, args[0], p.Namespace())
				if err != nil {
					return fmt.Errorf("failed to run GetPipelineRunByName: %w", err)
				}
				pipelineruns = append(pipelineruns, *pipelinerun)
			case 0:
				pipelineruns, err = pipelinerunpkg.GetAllPipelineRuns(cs, p.Namespace())
				if err != nil {
					return fmt.Errorf("failed to run GetAllPipelineRuns: %w", err)
				}
			default:
				return fmt.Errorf("too many arguments. Provide either no arguments to get all PipelineRuns or a single PipelineRuns name")
			}

			for i := range pipelineruns {
				pipelineRun := &pipelineruns[i]
				graph := taskgraph.BuildTaskGraph(pipelineRun.Status.PipelineSpec.Tasks)
				graph.PipelineName = pipelineRun.Name
				graphs = append(graphs, graph)
			}

			if opts.OutputDir != "" {
				if err = taskgraph.WriteAllGraphs(graphs, opts.OutputFormat, opts.OutputDir, opts.WithTaskRef); err != nil {
					return fmt.Errorf("failed to save graph: %w", err)
				}
			} else {
				if err = taskgraph.PrintAllGraphs(graphs, opts.OutputFormat, opts.WithTaskRef); err != nil {
					return fmt.Errorf("failed to print graph: %w", err)
				}
			}

			return nil
		},
	}

	// Define the command-line opts
	c.Flags().StringVar(
		&opts.OutputFormat, "output-format", "dot", "the output format (dot - DOT, puml - PlantUML or mmd - Mermaid)")
	c.Flags().StringVar(
		&opts.OutputDir, "output-dir", "", "the directory to save the output files. Otherwise, the output is printed to the screen")
	c.Flags().BoolVar(
		&opts.WithTaskRef, "with-task-ref", false, "Include TaskRefName information in the output")

	return c
}
