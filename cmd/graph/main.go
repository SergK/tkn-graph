package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sergk/tkn-graph/pkg/taskgraph"

	"github.com/spf13/cobra"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Options struct {
	Namespace    string
	ObjectKind   string
	OutputFormat string
	OutputDir    string
	WithTaskRef  bool
}

func main() {
	var options Options

	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "graph",
		Short: "Generate a graph of a Tekton object",
		Long:  "graph is a command-line tool for generating graphs from Tekton kind: Pipelines and kind: PipelineRuns.",
		Example: `  graph --namespace my-namespace --kind Pipeline --output-format dot
  graph --namespace my-namespace --kind PipelineRun --output-format puml
  graph --namespace my-namespace --kind Pipeline --output-format mmd --output-dir /tmp/output`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create the Kubernetes client
			config, err := rest.InClusterConfig()
			if err != nil {
				kubeconfig := os.Getenv("KUBECONFIG")
				if kubeconfig == "" {
					kubeconfig = clientcmd.RecommendedHomeFile
				}
				config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
				if err != nil {
					log.Fatalf("Failed to get Kubernetes configuration: %v", err)
				}
			}
			tektonClient, err := versioned.NewForConfig(config)
			if err != nil {
				log.Fatalf("Failed to create Tekton client: %v", err)
			}

			// Get the namespace to use
			if options.Namespace == "" {
				namespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
					clientcmd.NewDefaultClientConfigLoadingRules(),
					&clientcmd.ConfigOverrides{},
				).Namespace()
				if err != nil {
					log.Fatalf("Failed to get namespace from kubeconfig: %v", err)
				}
				if namespace == "" {
					namespace = "default"
				}
				options.Namespace = namespace
			}

			// Build the list of task graphs
			var graphs []*taskgraph.TaskGraph

			switch options.ObjectKind {
			case "Pipeline":
				pipelines, err := tektonClient.TektonV1().Pipelines(options.Namespace).List(context.TODO(), v1.ListOptions{})
				if err != nil {
					log.Fatalf("Failed to get Pipelines: %v", err)
				}
				if len(pipelines.Items) == 0 {
					log.Fatalf("No Pipelines found in namespace %s", options.Namespace)
				}
				for _, pipeline := range pipelines.Items {
					graph := taskgraph.BuildTaskGraph(pipeline.Spec.Tasks)
					graph.PipelineName = pipeline.Name
					graphs = append(graphs, graph)
				}

			case "PipelineRun":
				pipelineRuns, err := tektonClient.TektonV1().PipelineRuns(options.Namespace).List(context.TODO(), v1.ListOptions{})
				if err != nil {
					log.Fatalf("Failed to get PipelineRuns: %v", err)
				}
				if len(pipelineRuns.Items) == 0 {
					log.Fatalf("No PipelineRuns found in namespace %s", options.Namespace)
				}
				for _, pipelineRun := range pipelineRuns.Items {
					graph := taskgraph.BuildTaskGraph(pipelineRun.Status.PipelineSpec.Tasks)
					graph.PipelineName = pipelineRun.Name
					graphs = append(graphs, graph)
				}

			default:
				log.Fatalf("Invalid kind type: %s", options.ObjectKind)
			}

			// Generate graph for each object
			for _, graph := range graphs {
				// Generate the output format string
				output := taskgraph.FormatFunc(graph, options.OutputFormat, options.WithTaskRef)

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
					err = os.WriteFile(filename, []byte(output), 0644)
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
		&options.ObjectKind, "kind", "Pipeline", "the kind of the Tekton object to parse (Pipeline or PipelineRun)")
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
