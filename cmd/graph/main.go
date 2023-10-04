package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sergk/tkn-graph/pkg/client"
	"github.com/sergk/tkn-graph/pkg/taskgraph"
	"github.com/spf13/cobra"
)

type Options struct {
	Namespace    string
	TektonKind   string
	OutputFormat string
	OutputDir    string
	WithTaskRef  bool
}

func main() {
	var options Options

	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "tkn-graph",
		Short: "Generate a graph of a Tekton object",
		Long:  "tkn-graph is a command-line tool for generating graphs from Tekton kind: Pipelines and kind: PipelineRuns.",
		Example: `  graph --namespace my-namespace --kind Pipeline --output-format dot
  graph --namespace my-namespace --kind PipelineRun --output-format puml
  graph --namespace my-namespace --kind Pipeline --output-format mmd --output-dir /tmp/output`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create the Kubernetes client
			tektonClient, err := client.NewClient()
			if err != nil {
				log.Fatalf("Failed to create Tekton client: %v", err)
			}

			// Get the namespace to use
			if options.Namespace == "" {
				namespace, err := tektonClient.GetNamespace()
				if err != nil {
					log.Fatalf("Failed to get namespace: %v", err)
				}
				options.Namespace = namespace
			}

			// Build the list of task graphs
			var graphs []*taskgraph.TaskGraph

			switch options.TektonKind {
			case "Pipeline":
				pipelines, err := tektonClient.GetPipelines(options.Namespace)
				if err != nil {
					log.Fatalf("Failed to get Pipelines: %v", err)
				}
				for i := range pipelines {
					pipeline := &pipelines[i]
					graph := taskgraph.BuildTaskGraph(pipeline.Spec.Tasks)
					graph.PipelineName = pipeline.Name
					graphs = append(graphs, graph)
				}

			case "PipelineRun":
				pipelineRuns, err := tektonClient.GetPipelineRuns(options.Namespace)
				if err != nil {
					log.Fatalf("Failed to get PipelineRuns: %v", err)
				}
				for i := range pipelineRuns {
					pipelineRun := &pipelineRuns[i]
					graph := taskgraph.BuildTaskGraph(pipelineRun.Status.PipelineSpec.Tasks)
					graph.PipelineName = pipelineRun.Name
					graphs = append(graphs, graph)
				}

			default:
				log.Fatalf("Invalid kind type: %s", options.TektonKind)
			}

			// Generate graph for each object
			for _, graph := range graphs {
				// Generate the output format string
				output, err := taskgraph.FormatFunc(graph, options.OutputFormat, options.WithTaskRef)
				if err != nil {
					log.Fatalf("Failed to generate output: %v", err)
				}

				// Print or save the graph
				if options.OutputDir == "" {
					// Print the graph to the screen
					fmt.Println(output)
				} else {
					// Save the graph to a file
					err := os.MkdirAll(options.OutputDir, 0755)
					if err != nil {
						log.Fatalf("Failed to create directory %s: %v", options.OutputDir, err)
					}
					filename := filepath.Join(options.OutputDir, fmt.Sprintf("%s.%s", graph.PipelineName, options.OutputFormat))
					err = os.WriteFile(filename, []byte(output), 0600)
					if err != nil {
						log.Fatalf("Failed to write file %s: %v", filename, err)
					}
				}
			}
		},
	}

	// Define the command-line options
	rootCmd.Flags().StringVar(
		&options.Namespace, "namespace", "", "the Kubernetes namespace to use. Will try to get namespace from KUBECONFIG if not specified then fallback to 'default'")
	rootCmd.Flags().StringVar(
		&options.TektonKind, "kind", "Pipeline", "the kind of the Tekton object to parse (Pipeline or PipelineRun)")
	rootCmd.Flags().StringVar(
		&options.OutputFormat, "output-format", "dot", "the output format (dot - DOT, puml - PlantUML or mmd - Mermaid)")
	rootCmd.Flags().StringVar(
		&options.OutputDir, "output-dir", "", "the directory to save the output files. Otherwise, the output is printed to the screen")
	rootCmd.Flags().BoolVar(
		&options.WithTaskRef, "with-task-ref", false, "Include TaskRefName information in the output")

	// Parse the command-line options
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
