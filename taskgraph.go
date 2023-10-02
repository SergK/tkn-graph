package main

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
	Dependencies []*TaskNode
}

type DOT struct {
	Name   string
	Edges  []string
	Format string
}

// formatFunc is a function that generates the output format string for a TaskGraph
type formatFuncType func(graph *TaskGraph, format string) string

// BuildTaskGraph creates a TaskGraph from a list of PipelineTasks
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

// String converts a DOT graph to a string
func (d *DOT) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s \"%s\" {\n", d.Format, d.Name))
	for _, edge := range d.Edges {
		buf.WriteString(fmt.Sprintf("%s\n", edge))
	}
	buf.WriteString("}\n")
	return buf.String()
}

// ToPlantUML converts a TaskGraph to a PlantUML graph
func (g *TaskGraph) ToPlantUML() string {
	plantuml := fmt.Sprintf("@startuml\ntitle %s\n\n", g.PipelineName)
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

// formatFunc generates the output format string for a TaskGraph based on the specified format
var formatFunc formatFuncType = func(graph *TaskGraph, format string) string {
	switch strings.ToLower(format) {
	case "dot":
		return graph.ToDOT().String()
	case "puml":
		return graph.ToPlantUML()
	case "mmd":
		return graph.ToMermaid()
	default:
		log.Fatalf("Invalid output format: %s", format)
		return ""
	}
}
