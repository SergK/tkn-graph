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
	Namespace  string
	ObjectKind string
}

func main() {
	var options Options

	// Define the root command
	rootCmd := &cobra.Command{
		Use: "graph",
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

			switch options.ObjectKind {
			case "Pipeline":
				pipelines, err := tektonClient.TektonV1().Pipelines(options.Namespace).List(context.TODO(), v1.ListOptions{})
				if err != nil {
					log.Fatalf("Failed to get Pipelines: %v", err)
				}
				for _, pipeline := range pipelines.Items {
					// Build the task graph
					graph := BuildTaskGraph(pipeline.Spec.Tasks)
					graph.PipelineName = pipeline.Name

					// Generate the DOT format string
					dot := graph.ToDOT()

					// Print the DOT format string
					fmt.Println(dot)
				}
			case "PipelineRun":
				pipelineRuns, err := tektonClient.TektonV1().PipelineRuns(options.Namespace).List(context.TODO(), v1.ListOptions{})
				if err != nil {
					log.Fatalf("Failed to get PipelineRuns: %v", err)
				}
				for _, pipelineRun := range pipelineRuns.Items {
					// Build the task graph
					graph := BuildTaskGraph(pipelineRun.Status.PipelineSpec.Tasks)
					plantuml := graph.ToPlantUML()
					fmt.Println(plantuml)
					// graph.PipelineName = pipelineRun.Name

					// Generate the DOT format string
					// dot := graph.ToDOT()

					// Print the DOT format string
					// fmt.Println(dot)
				}
			default:
				log.Fatalf("Invalid kind type: %s", options.ObjectKind)
			}
		},
	}

	// Define the command-line options
	rootCmd.Flags().StringVar(&options.Namespace, "namespace", "", "the Kubernetes namespace to use")
	rootCmd.Flags().StringVar(&options.ObjectKind, "kind", "Pipeline", "the kind of the Tekton object to parse (Pipeline or PipelineRun)")

	// Parse the command-line options
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
