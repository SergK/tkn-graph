package taskgraph

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

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
	IsRoot       bool // Flag to indicate the the node is the root of the graph
}

type DOT struct {
	Name   string
	Edges  []string
	Format string
}

// FormatFunc is a function that generates the output format string for a TaskGraph
type formatFuncType func(graph *TaskGraph, format string, withTaskRef bool) (string, error)

func createTaskNode(task *v1pipeline.PipelineTask) *TaskNode {
	return &TaskNode{
		Name:        task.Name,
		TaskRefName: task.TaskRef.Name,
		IsRoot:      true, // we assume that the node is root until we find a parent
	}
}

// In the case where the order of tasks is arbitrary, it is necessary to create all the nodes first
// and then add the dependencies in a separate loop (since dependencies doesn't exist in TaskRef).
// BuildTaskGraph creates a TaskGraph from a list of PipelineTasks
func BuildTaskGraph(tasks []v1pipeline.PipelineTask) *TaskGraph {
	graph := &TaskGraph{
		Nodes: make(map[string]*TaskNode, len(tasks)),
	}

	// Create a node for each task and add it to the graph
	for i := range tasks {
		task := &tasks[i]
		node := createTaskNode(task)
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
			node.IsRoot = false
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
		if node.IsRoot {
			// "start" is the special node that represents the start of the pipeline
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"start\" -> \"%s\"", node.Name))
		}
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
		if node.IsRoot {
			// "start" is the special node that represents the start of the pipeline
			dot.Edges = append(dot.Edges, fmt.Sprintf("  \"start\" -> \"%s\n(%s)\"", node.Name, node.TaskRefName))
		}
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
	buf.WriteString(fmt.Sprintf("%s {\n  labelloc=\"t\"\n  label=\"%s\"\n  end [shape=\"point\" width=0.2]\n  start [shape=\"point\" width=0.2]\n", d.Format, d.Name))
	for _, edge := range d.Edges {
		buf.WriteString(fmt.Sprintf("%s\n", edge))
	}
	buf.WriteString("}\n")
	return buf.String()
}

func (g *TaskGraph) ToPlantUML(withTaskRef bool) (string, error) {
	var builder strings.Builder
	funcMap := template.FuncMap{
		"replace": strings.ReplaceAll,
	}
	var tmpl *template.Template
	var err error
	if withTaskRef {
		tmpl, err = template.New("plantuml").Funcs(funcMap).Parse(plantumlTemplateWithTaskRef)
	} else {
		tmpl, err = template.New("plantuml").Funcs(funcMap).Parse(plantumlTemplate)
	}
	if err != nil {
		return "", fmt.Errorf("failed to parse plantuml template: %w", err)
	}
	if err := tmpl.Execute(&builder, g); err != nil {
		return "", fmt.Errorf("failed to execute plantuml template: %w", err)
	}
	return builder.String(), nil
}

func (g *TaskGraph) ToMermaid() (string, error) {
	return generateMermaid(g, mermaidTemplate)
}

func (g *TaskGraph) ToMermaidWithTaskRef() (string, error) {
	return generateMermaid(g, mermaidTemplateWithTaskRef)
}

func generateMermaid(g *TaskGraph, tmpl string) (string, error) {
	var builder strings.Builder
	t, err := template.New("mermaid").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse mermaid template: %w", err)
	}
	if err := t.Execute(&builder, g); err != nil {
		return "", fmt.Errorf("failed to execute mermaid template: %w", err)
	}
	return builder.String(), nil
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
		return graph.ToPlantUML(withTaskRef)
	case "mmd":
		if withTaskRef {
			return graph.ToMermaidWithTaskRef()
		}
		return graph.ToMermaid()
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
