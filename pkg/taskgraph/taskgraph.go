package taskgraph

import (
	"bytes"
	"fmt"
	"log"
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
type FormatFuncType func(graph *TaskGraph, format string, withTaskRef bool) string

// BuildTaskGraph creates a TaskGraph from a list of PipelineTasks
func BuildTaskGraph(tasks []v1pipeline.PipelineTask) *TaskGraph {
	graph := &TaskGraph{
		Nodes: make(map[string]*TaskNode),
	}

	// Create a node for each task and add it to the graph
	for _, task := range tasks {
		node := &TaskNode{
			Name:        task.Name,
			TaskRefName: task.TaskRef.Name,
		}
		graph.Nodes[task.Name] = node

		// Add dependencies to the node
		for _, dep := range task.RunAfter {
			depNode, ok := graph.Nodes[dep]
			if !ok {
				// Create a new node for the dependency if it doesn't already exist
				depNode = &TaskNode{
					Name: dep,
				}
				graph.Nodes[dep] = depNode
			}
			node.Dependencies = append(node.Dependencies, depNode)
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
		for _, dep := range node.Dependencies {
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\" -> \"%s\"", dep.Name, node.Name))
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
		for _, dep := range node.Dependencies {
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\n(%s)\" -> \"%s\n(%s)\"", dep.Name, dep.TaskRefName, node.Name, node.TaskRefName))
		}
	}

	return dot
}

// String converts a DOT graph to a string
func (d *DOT) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s {\n  labelloc=\"t\"\n  label=\"%s\"\n", d.Format, d.Name))
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
			plantuml += fmt.Sprintf("[*] --> %s\n", nodeName)
		}
		for _, dep := range node.Dependencies {
			// Replace dashes with underscores in node names because PlantUML doesn't like dashes
			depName := strings.ReplaceAll(dep.Name, "-", "_")
			plantuml += fmt.Sprintf("%s <-down- %s\n", nodeName, depName)
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
			plantuml += fmt.Sprintf("[*] --> %s\n", nodeName)
		}
		for _, dep := range node.Dependencies {
			// Replace dashes with underscores in node names because PlantUML doesn't like dashes
			depName := strings.ReplaceAll(dep.Name, "-", "_")
			plantuml += fmt.Sprintf("%s <-down- %s\n", nodeName, depName)
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
		for _, dep := range node.Dependencies {
			mermaid += fmt.Sprintf("   %s --> %s\n", dep.Name, node.Name)
		}
	}
	return mermaid
}

// ToMermaidWithTaskRef converts a TaskGraph to a mermaid graph with taskRefName
func (g *TaskGraph) ToMermaidWithTaskRef() string {
	mermaid := fmt.Sprintf("---\ntitle: %s\n---\nflowchart TD\n", g.PipelineName)
	for _, node := range g.Nodes {
		for _, dep := range node.Dependencies {
			mermaid += fmt.Sprintf("   %s(\"%s\n   (%s)\") --> %s(\"%s\n   (%s)\")\n", dep.Name, dep.Name, dep.TaskRefName, node.Name, node.Name, node.TaskRefName)
		}
	}
	return mermaid
}

// formatFunc generates the output format string for a TaskGraph based on the specified format
var FormatFunc FormatFuncType = func(graph *TaskGraph, format string, withTaskRef bool) string {
	switch strings.ToLower(format) {
	case "dot":
		if withTaskRef {
			return graph.ToDOTWithTaskRef().String()
		}
		return graph.ToDOT().String()
	case "puml":
		if withTaskRef {
			return graph.ToPlantUMLWithTaskRef()
		}
		return graph.ToPlantUML()
	case "mmd":
		if withTaskRef {
			return graph.ToMermaidWithTaskRef()
		}
		return graph.ToMermaid()
	default:
		log.Fatalf("Invalid output format: %s", format)
		return ""
	}
}
