package taskgraph

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	v1pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type TaskGraph struct {
	PipelineName string
	Nodes        map[string]*TaskNode
}

type TaskNode struct {
	Name         string
	TaskRefName  string // Name of the kind: Task referenced by this task in the pipeline
	Dependencies []*TaskNode
}

type DOT struct {
	Name   string
	Edges  []string
	Format string
}

// FormatFunc is a function that generates the output format string for a TaskGraph
type formatFuncType func(graph *TaskGraph, format string, withTaskRef bool) (string, error)

// In the case where the order of tasks is arbitrary, it is necessary to create all the nodes first
// and then add the dependencies in a separate loop (since dependencies doesn't exist in TaskRef).
// BuildTaskGraph creates a TaskGraph from a list of PipelineTasks
func BuildTaskGraph(tasks []v1pipeline.PipelineTask) *TaskGraph {
	graph := &TaskGraph{
		Nodes: make(map[string]*TaskNode),
	}

	// Create a node for each task and add it to the graph
	for i := range tasks {
		task := &tasks[i]
		node := &TaskNode{
			Name:        task.Name,
			TaskRefName: task.TaskRef.Name,
		}
		graph.Nodes[task.Name] = node
	}

	// Add dependencies to the nodes
	for i := range tasks {
		task := &tasks[i]
		node := graph.Nodes[task.Name]

		// Add dependencies to the node
		for _, depName := range task.RunAfter {
			depNode := graph.Nodes[depName]
			depNode.Dependencies = append(depNode.Dependencies, node)
		}
	}

	return graph
}

// ToDOT converts a TaskGraph to a DOT graph
func (g *TaskGraph) ToDOT() *DOT {
	dot := &DOT{
		Name:   g.PipelineName,
		Format: "digraph",
	}

	for _, node := range g.Nodes {
		if len(node.Dependencies) == 0 {
			// "end" is the special node that represents the end of the pipeline
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\" -> \"end\"", node.Name))
		}
		for _, dep := range node.Dependencies {
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\" -> \"%s\"", node.Name, dep.Name))
		}
	}

	return dot
}

// ToDOTWithTaskRef converts a TaskGraph to a DOT graph
func (g *TaskGraph) ToDOTWithTaskRef() *DOT {
	dot := &DOT{
		Name:   g.PipelineName,
		Format: "digraph",
	}

	for _, node := range g.Nodes {
		if len(node.Dependencies) == 0 {
			// "end" is the special node that represents the end of the pipeline
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\n(%s)\" -> \"end\"", node.Name, node.TaskRefName))
		}
		for _, dep := range node.Dependencies {
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\n(%s)\" -> \"%s\n(%s)\"", node.Name, node.TaskRefName, dep.Name, dep.TaskRefName))
		}
	}

	return dot
}

// String converts a DOT graph to a string
func (d *DOT) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s {\n  labelloc=\"t\"\n  label=\"%s\"\n  end [shape=\"point\" width=0.2]\n", d.Format, d.Name))
	for _, edge := range d.Edges {
		buf.WriteString(fmt.Sprintf("%s\n", edge))
	}
	buf.WriteString("}\n")
	return buf.String()
}

// ToPlantUML converts a TaskGraph to a PlantUML graph
func (g *TaskGraph) ToPlantUML() string {
	plantuml := fmt.Sprintf("@startuml\nhide empty description\ntitle %s\n\n", g.PipelineName)
	for _, node := range g.Nodes {
		// Replace dashes with underscores in node names because PlantUML doesn't like dashes
		nodeName := strings.ReplaceAll(node.Name, "-", "_")
		// the root node is the one with no dependencies and that task starts the execution immediately
		if len(node.Dependencies) == 0 {
			plantuml += fmt.Sprintf("%s --> [*]\n", nodeName)
		}
		for _, dep := range node.Dependencies {
			// Replace dashes with underscores in node names because PlantUML doesn't like dashes
			depName := strings.ReplaceAll(dep.Name, "-", "_")
			plantuml += fmt.Sprintf("%s -down-> %s\n", nodeName, depName)
		}
	}
	plantuml += "\n@enduml\n"
	return plantuml
}

// ToPlantUMLWithTaskRef converts a TaskGraph to a PlantUML graph with taskRefName
func (g *TaskGraph) ToPlantUMLWithTaskRef() string {
	plantuml := fmt.Sprintf("@startuml\nhide empty description\ntitle %s\n\n", g.PipelineName)

	// Create a map to store the unique nodes and their TaskRefName values
	uniqueNodes := make(map[string]string)

	for _, node := range g.Nodes {
		// Replace dashes with underscores in node names because PlantUML doesn't like dashes
		nodeName := strings.ReplaceAll(node.Name, "-", "_")
		// the root node is the one with no dependencies and that task starts the execution immediately
		if len(node.Dependencies) == 0 {
			plantuml += fmt.Sprintf("%s --> [*]\n", nodeName)
		}
		for _, dep := range node.Dependencies {
			// Replace dashes with underscores in node names because PlantUML doesn't like dashes
			depName := strings.ReplaceAll(dep.Name, "-", "_")
			plantuml += fmt.Sprintf("%s -down-> %s\n", nodeName, depName)
		}
		// Add the node to the uniqueNodes map if it doesn't already exist
		if _, ok := uniqueNodes[nodeName]; !ok {
			uniqueNodes[nodeName] = node.TaskRefName
		}
	}
	plantuml += "\n"
	// Add the unique nodes to the output
	for nodeName, taskRefName := range uniqueNodes {
		plantuml += fmt.Sprintf("%s: %s\n", nodeName, taskRefName)
	}
	plantuml += "\n@enduml\n"
	return plantuml
}

// ToMermaid converts a TaskGraph to a mermaid graph
func (g *TaskGraph) ToMermaid() string {
	mermaid := fmt.Sprintf("---\ntitle: %s\n---\nflowchart TD\n", g.PipelineName)
	for _, node := range g.Nodes {
		if len(node.Dependencies) == 0 {
			mermaid += fmt.Sprintf("   %s --> id([fa:fa-circle])\n", node.Name)
		}
		for _, dep := range node.Dependencies {
			mermaid += fmt.Sprintf("   %s --> %s\n", node.Name, dep.Name)
		}
	}
	return mermaid
}

// ToMermaidWithTaskRef converts a TaskGraph to a mermaid graph with taskRefName
func (g *TaskGraph) ToMermaidWithTaskRef() string {
	mermaid := fmt.Sprintf("---\ntitle: %s\n---\nflowchart TD\n", g.PipelineName)
	for _, node := range g.Nodes {
		if len(node.Dependencies) == 0 {
			mermaid += fmt.Sprintf("   %s(\"%s\n   (%s)\") --> id([fa:fa-circle])\n", node.Name, node.Name, node.TaskRefName)
		}
		for _, dep := range node.Dependencies {
			mermaid += fmt.Sprintf("   %s(\"%s\n   (%s)\") --> %s(\"%s\n   (%s)\")\n", node.Name, node.Name, node.TaskRefName, dep.Name, dep.Name, dep.TaskRefName)
		}
	}
	return mermaid
}

// formatFunc generates the output format string for a TaskGraph based on the specified format
var formatFunc formatFuncType = func(graph *TaskGraph, format string, withTaskRef bool) (string, error) {
	switch strings.ToLower(format) {
	case "dot":
		if withTaskRef {
			return graph.ToDOTWithTaskRef().String(), nil
		}
		return graph.ToDOT().String(), nil
	case "puml":
		if withTaskRef {
			return graph.ToPlantUMLWithTaskRef(), nil
		}
		return graph.ToPlantUML(), nil
	case "mmd":
		if withTaskRef {
			return graph.ToMermaidWithTaskRef(), nil
		}
		return graph.ToMermaid(), nil
	default:
		return "", fmt.Errorf("Invalid output format: %s", format)
	}
}

// Function that prints graph to stdout
func PrintAllGraphs(graphs []*TaskGraph, outputFormat string, withTaskRef bool) error {
	for _, graph := range graphs {
		output, err := formatFunc(graph, outputFormat, withTaskRef)
		if err != nil {
			return fmt.Errorf("Failed to generate output: %w", err)
		}
		fmt.Println(output)
	}
	return nil
}

// Function that writes graph to file
func WriteAllGraphs(graphs []*TaskGraph, outputFormat string, outputDir string, withTaskRef bool) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create directory %s: %w", outputDir, err)
	}
	for _, graph := range graphs {
		output, err := formatFunc(graph, outputFormat, withTaskRef)
		if err != nil {
			return fmt.Errorf("Failed to generate output: %w", err)
		}
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.%s", graph.PipelineName, outputFormat))
		err = os.WriteFile(filename, []byte(output), 0600)
		if err != nil {
			return fmt.Errorf("Failed to write file %s: %w", filename, err)
		}
	}
	return nil
}
