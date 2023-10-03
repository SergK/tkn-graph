package taskgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const (
	testPipelineName = "test-pipeline"
)

func getTestTasks() []v1pipeline.PipelineTask {
	return []v1pipeline.PipelineTask{
		{
			Name: "task1",
			TaskRef: &v1pipeline.TaskRef{
				Name: "taskRef1",
			},
			RunAfter: []string{"task2", "task3"},
		},
		{
			Name: "task2",
			TaskRef: &v1pipeline.TaskRef{
				Name: "taskRef2",
			},
			RunAfter: []string{"task3"},
		},
		{
			Name: "task3",
			TaskRef: &v1pipeline.TaskRef{
				Name: "taskRef3",
			},
		},
	}
}

func TestBuildTaskGraph(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())

	// Assert that the graph has the correct number of nodes
	assert.Equal(t, 3, len(graph.Nodes))

	// Assert that the nodes have the correct names and task references
	assert.Equal(t, "taskRef1", graph.Nodes["task1"].TaskRefName)
	assert.Equal(t, "taskRef2", graph.Nodes["task2"].TaskRefName)
	assert.Equal(t, "taskRef3", graph.Nodes["task3"].TaskRefName)

	// Assert that the nodes have the correct dependencies
	// Task3 has two downstream dependencies Task1 and Task2
	assert.Equal(t, []*TaskNode{graph.Nodes["task1"], graph.Nodes["task2"]}, graph.Nodes["task3"].Dependencies)
	// Task2 has one downstream dependency Task1
	assert.Equal(t, []*TaskNode{graph.Nodes["task1"]}, graph.Nodes["task2"].Dependencies)
	// Task1 has no downstream dependencies
	assert.Empty(t, graph.Nodes["task1"].Dependencies)
}

func TestTaskGraphToDOT(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToDOT method
	dot := graph.ToDOT()
	assert.Equal(t, "digraph", dot.Format)
	assert.Equal(t, testPipelineName, dot.Name)
	assert.Contains(t, dot.Edges, "  \"task1\" -> \"task2\"")
	assert.Contains(t, dot.Edges, "  \"task1\" -> \"task3\"")
	assert.Contains(t, dot.Edges, "  \"task2\" -> \"task3\"")
}

func TestTaskGraphToDOTWithTaskRef(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToDOTWithTaskRef method
	dot := graph.ToDOTWithTaskRef()
	assert.Equal(t, "digraph", dot.Format)
	assert.Equal(t, testPipelineName, dot.Name)
	assert.Contains(t, dot.Edges, "  \"task1\n(taskRef1)\" -> \"task2\n(taskRef2)\"")
	assert.Contains(t, dot.Edges, "  \"task1\n(taskRef1)\" -> \"task3\n(taskRef3)\"")
	assert.Contains(t, dot.Edges, "  \"task2\n(taskRef2)\" -> \"task3\n(taskRef3)\"")
}

func TestDOTString(t *testing.T) {
	// Define some test edges
	edges := []string{
		"  \"task1\" -> \"task2\"",
		"  \"task1\" -> \"task3\"",
		"  \"task2\" -> \"task3\"",
	}

	// Create a DOT object with the test edges
	dot := DOT{
		Format: "digraph",
		Name:   testPipelineName,
		Edges:  edges,
	}

	// Test the String method
	expected := "digraph {\n  labelloc=\"t\"\n  label=\"test-pipeline\"\n  \"task1\" -> \"task2\"\n  \"task1\" -> \"task3\"\n  \"task2\" -> \"task3\"\n}\n"
	assert.Equal(t, expected, dot.String())
}

func TestTaskGraphToPlantUML(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToPlantUML method
	plantuml := graph.ToPlantUML()
	assert.Contains(t, plantuml, "@startuml\nhide empty description\ntitle test-pipeline\n\n")
	assert.Contains(t, plantuml, "[*] --> task1\n")
	assert.Contains(t, plantuml, "task2 <-down- task1\n")
	assert.Contains(t, plantuml, "task3 <-down- task1\n")
	assert.Contains(t, plantuml, "task3 <-down- task2\n")
	assert.Contains(t, plantuml, "\n@enduml\n")
}

func TestTaskGraphToPlantUMLWithTaskRef(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToPlantUMLWithTaskRef method
	plantuml := graph.ToPlantUMLWithTaskRef()
	assert.Contains(t, plantuml, "@startuml\nhide empty description\ntitle test-pipeline\n\n")
	assert.Contains(t, plantuml, "[*] --> task1")
	assert.Contains(t, plantuml, "task2 <-down- task1")
	assert.Contains(t, plantuml, "task3 <-down- task1")
	assert.Contains(t, plantuml, "task3 <-down- task2")
	assert.Contains(t, plantuml, "task1: taskRef1\n")
	assert.Contains(t, plantuml, "task2: taskRef2\n")
	assert.Contains(t, plantuml, "task3: taskRef3\n")
	assert.Contains(t, plantuml, "\n@enduml\n")
}

func TestTaskGraphToMermaid(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToMermaid method
	mermaid := graph.ToMermaid()
	assert.Contains(t, mermaid, "---\ntitle: test-pipeline\n---\nflowchart TD\n")
	assert.Contains(t, mermaid, "   task1 --> task2\n")
	assert.Contains(t, mermaid, "   task1 --> task3\n")
	assert.Contains(t, mermaid, "   task2 --> task3\n")
}

func TestTaskGraphToMermaidWithTaskRef(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToMermaidWithTaskRef method
	mermaid := graph.ToMermaidWithTaskRef()
	assert.Contains(t, mermaid, "---\ntitle: test-pipeline\n---\nflowchart TD\n")
	assert.Contains(t, mermaid, "   task1(\"task1\n   (taskRef1)\") --> task2(\"task2\n   (taskRef2)\")\n")
	assert.Contains(t, mermaid, "   task1(\"task1\n   (taskRef1)\") --> task3(\"task3\n   (taskRef3)\")\n")
	assert.Contains(t, mermaid, "   task2(\"task2\n   (taskRef2)\") --> task3(\"task3\n   (taskRef3)\")\n")
}

func TestFormatFunc(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the FormatFunc method
	dot, err := FormatFunc(graph, "dot", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, dot)

	puml, err := FormatFunc(graph, "puml", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, puml)

	mmd, err := FormatFunc(graph, "mmd", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, mmd)

	invalid, err := FormatFunc(graph, "invalid", false)
	assert.Error(t, err)
	assert.Empty(t, invalid)
	assert.Equal(t, "Invalid output format: invalid", err.Error())
}
