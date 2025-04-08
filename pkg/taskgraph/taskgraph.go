package taskgraph

import (
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

// FormatFunc is a function that generates the output format string for a TaskGraph
type formatFuncMap func(graph *TaskGraph, format string, withTaskRef bool) (string, error)

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

func (g *TaskGraph) ToDOT(withTaskRef bool) (string, error) {
	var builder strings.Builder

	var tmpl *template.Template
	if withTaskRef {
		tmpl = template.Must(template.New("dot").Parse(dotTemplateWithTaskRef))
	} else {
		tmpl = template.Must(template.New("dot").Parse(dotTemplate))
	}

	if err := tmpl.Execute(&builder, struct {
		PipelineName string
		Nodes        map[string]*TaskNode
		Name         string
	}{
		PipelineName: g.PipelineName,
		Nodes:        g.Nodes,
		Name:         "G",
	}); err != nil {
		return "", fmt.Errorf("failed to execute dot template: %w", err)
	}

	return builder.String(), nil
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

func (g *TaskGraph) ToMermaid(withTaskRef bool) (string, error) {
	var builder strings.Builder

	tmpl := mermaidTemplate
	if withTaskRef {
		tmpl = mermaidTemplateWithTaskRef
	}

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
var formatFunc formatFuncMap = func(graph *TaskGraph, format string, withTaskRef bool) (string, error) {
	switch strings.ToLower(format) {
	case "dot":
		return graph.ToDOT(withTaskRef)
	case "puml":
		return graph.ToPlantUML(withTaskRef)
	case "mmd":
		return graph.ToMermaid(withTaskRef)
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
