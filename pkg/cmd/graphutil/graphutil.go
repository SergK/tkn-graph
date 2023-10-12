package graphutil

import (
	"fmt"

	"github.com/sergk/tkn-graph/pkg/cli/prerun"
	"github.com/sergk/tkn-graph/pkg/taskgraph"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	"github.com/tektoncd/cli/pkg/flags"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type GraphOptions struct {
	OutputFormat string
	OutputDir    string
	WithTaskRef  bool
}

type GraphData struct {
	Name string
	Spec v1.PipelineSpec
}

type GraphFetcher func(cs *cli.Clients, name []string, namespace string) ([]GraphData, error)

func NewGraphCommand(p cli.Params, runE func(cmd *cobra.Command, args []string, p cli.Params, opts *GraphOptions) error) *cobra.Command {
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
			return runE(cmd, args, p, opts)
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

func RunGraphCommand(cmd *cobra.Command, args []string, p cli.Params, opts *GraphOptions, fetcher GraphFetcher) error {
	cs, err := p.Clients()
	if err != nil {
		return err
	}

	var graphs []*taskgraph.TaskGraph
	var data []GraphData

	data, err = fetcher(cs, args, p.Namespace())
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	for i := range data {
		graph := taskgraph.BuildTaskGraph(data[i].Spec.Tasks)
		graph.PipelineName = data[i].Name
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
