package common

import (
	"fmt"

	"github.com/sergk/tkn-graph/pkg/cli/prerun"
	"github.com/sergk/tkn-graph/pkg/taskgraph"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/flags"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// GraphOptions holds the options for the graph command
// OutputFormat: dot, puml, mmd
// OutputDir: the directory to save the output files. Otherwise, the output is printed to the screen
// WithTaskRef: Include TaskRefName information in the output
type GraphOptions struct {
	OutputFormat string
	OutputDir    string
	WithTaskRef  bool
}

// Holds the Pipeline name and the Pipeline itself, in case of PipelineRun it holds the PipelineRun name and the Pipeline
type Pipeline struct {
	Name           string
	TektonPipeline v1.Pipeline
}

// GraphFetcher is an interface that defines the methods to fetch the Pipeline
type GraphFetcher interface {
	GetByName(cs *cli.Clients, name, namespace string) (*Pipeline, error)
	GetAll(cs *cli.Clients, namespace string) ([]Pipeline, error)
}

func CreateGraphCommand(p cli.Params, fetcher GraphFetcher) *cobra.Command {
	opts := &GraphOptions{}
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
			return RunGraphCommand(p, opts, fetcher, args)
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

func RunGraphCommand(p cli.Params, opts *GraphOptions, fetcher GraphFetcher, args []string) error {
	cs, err := p.Clients()
	if err != nil {
		return err
	}

	var pipelines []Pipeline

	switch len(args) {
	case 1:
		var pipeline *Pipeline

		pipeline, err = fetcher.GetByName(cs, args[0], p.Namespace())
		if err != nil {
			return fmt.Errorf("failed to run GetByName: %w", err)
		}

		pipelines = append(pipelines, *pipeline)
	case 0:
		pipelines, err = fetcher.GetAll(cs, p.Namespace())
		if err != nil {
			return fmt.Errorf("failed to run GetAll: %w", err)
		}
	default:
		return fmt.Errorf("too many arguments. Provide either no arguments to get all Pipelines or a single Pipeline name")
	}

	// Pre-allocate the graphs slice based on the number of pipelines
	graphs := make([]*taskgraph.TaskGraph, 0, len(pipelines))

	for i := range pipelines {
		graph := taskgraph.BuildTaskGraph(pipelines[i].TektonPipeline.Spec.Tasks)
		graph.PipelineName = pipelines[i].Name
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
}
