package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
}

func main() {
	var options Options

	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "graph",
		Short: "Generate a graph of a Tekton object",
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
			var graphs []*TaskGraph

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
					graph := BuildTaskGraph(pipeline.Spec.Tasks)
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
					graph := BuildTaskGraph(pipelineRun.Status.PipelineSpec.Tasks)
					graph.PipelineName = pipelineRun.Name
					graphs = append(graphs, graph)
				}

			default:
				log.Fatalf("Invalid kind type: %s", options.ObjectKind)
			}

			// Generate graph for each object
			for _, graph := range graphs {
				// Generate the output format string
				output := formatFunc(graph, options.OutputFormat)

				// Print the graph
				fmt.Println(output)
			}
		},
	}

	// Define the command-line options
	rootCmd.Flags().StringVar(&options.Namespace, "namespace", "", "the Kubernetes namespace to use")
	rootCmd.Flags().StringVar(&options.ObjectKind, "kind", "Pipeline", "the kind of the Tekton object to parse (Pipeline or PipelineRun)")
	rootCmd.Flags().StringVar(&options.OutputFormat, "output-format", "dot", "the output format (dot or plantuml)")

	// Parse the command-line options
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
