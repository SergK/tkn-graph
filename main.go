package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	v1pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	namespace = "edp-delivery-tekton-dev"
)

type TaskGraph struct {
	PipelineName string
	Nodes        map[string]*TaskNode
}

type TaskNode struct {
	Name         string
	Dependencies []*TaskNode
}

func main() {
	// Get the Kubernetes configuration
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

	// Create a Kubernetes clientset for the Tekton API group
	tektonClient, err := versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Tekton clientset: %v", err)
	}

	// Get all Tekton Pipelines
	pipelines, err := tektonClient.TektonV1().Pipelines(namespace).List(context.TODO(), v1.ListOptions{})
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

	// Get all Tekton PipelineRuns
	pipelineRuns, err := tektonClient.TektonV1().PipelineRuns(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get PipelineRuns: %v", err)
	}
	for _, pipelineRun := range pipelineRuns.Items {
		// Build the task graph
		graph := BuildTaskGraph(pipelineRun.Status.PipelineSpec.Tasks)
		graph.PipelineName = pipelineRun.Name
		// Generate the DOT format string
		dot := graph.ToDOT()

		// Print the DOT format string
		fmt.Println(dot)
	}

}

func BuildTaskGraph(tasks []v1pipeline.PipelineTask) *TaskGraph {
	graph := &TaskGraph{
		Nodes: make(map[string]*TaskNode),
	}

	// Create a node for each task
	for _, task := range tasks {
		node := &TaskNode{
			Name: task.Name,
		}
		graph.Nodes[task.Name] = node
	}

	// Add dependencies to each node
	for _, task := range tasks {
		node := graph.Nodes[task.Name]
		for _, dep := range task.RunAfter {
			depNode := graph.Nodes[dep]
			node.Dependencies = append(node.Dependencies, depNode)
		}
	}

	return graph
}

func (g *TaskGraph) ToDOT() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("digraph \"%s\" {\n", g.PipelineName))
	for _, node := range g.Nodes {
		for _, dep := range node.Dependencies {
			buf.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\"\n", dep.Name, node.Name))
		}
	}
	buf.WriteString("}\n")
	return buf.String()
}
