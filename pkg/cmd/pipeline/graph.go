package pipeline

import (
	"fmt"

	"github.com/sergk/tkn-graph/pkg/cli/prerun"
	pipelinepkg "github.com/sergk/tkn-graph/pkg/pipeline"
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
			var pipelines []v1.Pipeline

			switch len(args) {
			case 1:
				var pipeline *v1.Pipeline
				pipeline, err = pipelinepkg.GetPipelineByName(cs, args[0], p.Namespace())
				if err != nil {
					return fmt.Errorf("failed to run GetPipelineRunByName: %w", err)
				}
				pipelines = append(pipelines, *pipeline)
			case 0:
				pipelines, err = pipelinepkg.GetAllPipelines(cs, p.Namespace())
				if err != nil {
					return fmt.Errorf("failed to run GetAllPipelineRuns: %w", err)
				}
			default:
				return fmt.Errorf("too many arguments. Provide either no arguments to get all Pipelines or a single Pipeline name")
			}

			for i := range pipelines {
				pipeline := &pipelines[i]
				graph := taskgraph.BuildTaskGraph(pipeline.Spec.Tasks)
				graph.PipelineName = pipeline.Name
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
