package main

import (
	"bytes"
	"fmt"
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

func (d *DOT) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s \"%s\" {\n", d.Format, d.Name))
	for _, edge := range d.Edges {
		buf.WriteString(fmt.Sprintf("%s\n", edge))
	}
	buf.WriteString("}\n")
	return buf.String()
}

func (g *TaskGraph) ToPlantUML() string {
	plantuml := "@startuml\n\n"
	for _, node := range g.Nodes {
		nodeName := strings.ReplaceAll(node.Name, "-", "_")
		if len(node.Dependencies) == 0 {
			plantuml += fmt.Sprintf("[*] --> %s\n", nodeName)
		}
		for _, dep := range node.Dependencies {
			depName := strings.ReplaceAll(dep.Name, "-", "_")
			plantuml += fmt.Sprintf("%s <-down- %s\n", nodeName, depName)
		}
		plantuml += "\n"
	}
	plantuml += "@enduml\n"
	return plantuml
}
