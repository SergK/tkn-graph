package taskgraph

import (
	"os"
	"path/filepath"
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
		// we can have task without any dependencies
		{
			Name: "task-with-dash",
			TaskRef: &v1pipeline.TaskRef{
				Name: "taskRef4",
			},
		},
	}
}

func TestBuildTaskGraph(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())

	// Assert that the graph has the correct number of nodes
	assert.Equal(t, 4, len(graph.Nodes))

	// Assert that the nodes have the correct names and task references
	assert.Equal(t, "taskRef1", graph.Nodes["task1"].TaskRefName)
	assert.Equal(t, "taskRef2", graph.Nodes["task2"].TaskRefName)
	assert.Equal(t, "taskRef3", graph.Nodes["task3"].TaskRefName)
	assert.Equal(t, "taskRef4", graph.Nodes["task-with-dash"].TaskRefName)

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
	assert.Contains(t, dot.Edges, "  \"task2\" -> \"task1\"")
	assert.Contains(t, dot.Edges, "  \"task3\" -> \"task1\"")
	assert.Contains(t, dot.Edges, "  \"task3\" -> \"task2\"")
	assert.Contains(t, dot.Edges, "  \"task1\" -> \"end\"")
	assert.Contains(t, dot.Edges, "  \"task-with-dash\" -> \"end\"")
	assert.Contains(t, dot.Edges, "  \"start\" -> \"task-with-dash\"")
	assert.Contains(t, dot.Edges, "  \"start\" -> \"task3\"")
}

func TestTaskGraphToDOTWithTaskRef(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToDOTWithTaskRef method
	dot := graph.ToDOTWithTaskRef()
	assert.Equal(t, "digraph", dot.Format)
	assert.Equal(t, testPipelineName, dot.Name)
	assert.Contains(t, dot.Edges, "  \"task2\n(taskRef2)\" -> \"task1\n(taskRef1)\"")
	assert.Contains(t, dot.Edges, "  \"task3\n(taskRef3)\" -> \"task1\n(taskRef1)\"")
	assert.Contains(t, dot.Edges, "  \"task3\n(taskRef3)\" -> \"task2\n(taskRef2)\"")
	assert.Contains(t, dot.Edges, "  \"task1\n(taskRef1)\" -> \"end\"")
	assert.Contains(t, dot.Edges, "  \"task-with-dash\n(taskRef4)\" -> \"end\"")
	assert.Contains(t, dot.Edges, "  \"start\" -> \"task3\n(taskRef3)\"")
	assert.Contains(t, dot.Edges, "  \"start\" -> \"task-with-dash\n(taskRef4)\"")
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
	expected := `digraph {
  labelloc="t"
  label="test-pipeline"
  end [shape="point" width=0.2]
  start [shape="point" width=0.2]
  "task1" -> "task2"
  "task1" -> "task3"
  "task2" -> "task3"
}
`
	assert.Equal(t, expected, dot.String())
}

func TestTaskGraphToPlantUML(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToPlantUML method
	plantuml, err := graph.ToPlantUML(false)
	assert.NoError(t, err)
	assert.Contains(t, plantuml, "@startuml\nhide empty description\ntitle test-pipeline\n\n")
	assert.Contains(t, plantuml, "[*] --> task3\n")
	assert.Contains(t, plantuml, "[*] --> task_with_dash\n")
	assert.Contains(t, plantuml, "task1 --> [*]\n")
	assert.Contains(t, plantuml, "task_with_dash --> [*]\n")
	assert.Contains(t, plantuml, "task2 -down-> task1\n")
	assert.Contains(t, plantuml, "task3 -down-> task1\n")
	assert.Contains(t, plantuml, "task3 -down-> task2\n")
	assert.Contains(t, plantuml, "\n@enduml\n")
}

func TestTaskGraphToPlantUMLWithTaskRef(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToPlantUMLWithTaskRef method
	plantuml, err := graph.ToPlantUML(true)
	assert.NoError(t, err)
	assert.Contains(t, plantuml, "@startuml\nhide empty description\ntitle test-pipeline\n\n")
	assert.Contains(t, plantuml, "task1 --> [*]")
	assert.Contains(t, plantuml, "task_with_dash --> [*]")
	assert.Contains(t, plantuml, "[*] --> task3\n")
	assert.Contains(t, plantuml, "[*] --> task_with_dash\n")
	assert.Contains(t, plantuml, "task2 -down-> task1")
	assert.Contains(t, plantuml, "task3 -down-> task1")
	assert.Contains(t, plantuml, "task3 -down-> task2")
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
	mermaid, err := graph.ToMermaid()
	assert.NoError(t, err)
	assert.Contains(t, mermaid, "---\ntitle: test-pipeline\n---\nflowchart TD\n")
	assert.Contains(t, mermaid, "   task2 --> task1\n")
	assert.Contains(t, mermaid, "   task3 --> task1\n")
	assert.Contains(t, mermaid, "   task3 --> task2\n")
	assert.Contains(t, mermaid, "   task1 --> stop([fa:fa-circle])\n")
	assert.Contains(t, mermaid, "   task-with-dash --> stop([fa:fa-circle])\n")
	assert.Contains(t, mermaid, "   start([fa:fa-circle]) --> task3\n")
	assert.Contains(t, mermaid, "   start([fa:fa-circle]) --> task-with-dash\n")
}

func TestTaskGraphToMermaidWithTaskRef(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the ToMermaidWithTaskRef method
	mermaid, err := graph.ToMermaidWithTaskRef()
	assert.NoError(t, err)
	assert.Contains(t, mermaid, "---\ntitle: test-pipeline\n---\nflowchart TD\n")
	assert.Contains(t, mermaid, "   task2(\"task2\n   (taskRef2)\") --> task1(\"task1\n   (taskRef1)\")\n")
	assert.Contains(t, mermaid, "   task3(\"task3\n   (taskRef3)\") --> task1(\"task1\n   (taskRef1)\")\n")
	assert.Contains(t, mermaid, "   task3(\"task3\n   (taskRef3)\") --> task2(\"task2\n   (taskRef2)\")\n")
	assert.Contains(t, mermaid, "   task1(\"task1\n   (taskRef1)\") --> stop([fa:fa-circle])\n")
	assert.Contains(t, mermaid, "   task-with-dash(\"task-with-dash\n   (taskRef4)\") --> stop([fa:fa-circle])\n")
	assert.Contains(t, mermaid, "   start([fa:fa-circle]) --> task3(\"task3\n   (taskRef3)\")\n")
	assert.Contains(t, mermaid, "   start([fa:fa-circle]) --> task-with-dash(\"task-with-dash\n   (taskRef4)\")\n")
}

func TestFormatFunc(t *testing.T) {
	// Build the task graph
	graph := BuildTaskGraph(getTestTasks())
	graph.PipelineName = testPipelineName

	// Test the FormatFunc method
	dot, err := formatFunc(graph, "dot", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, dot)

	puml, err := formatFunc(graph, "puml", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, puml)

	mmd, err := formatFunc(graph, "mmd", false)
	assert.NoError(t, err)
	assert.NotEmpty(t, mmd)

	invalid, err := formatFunc(graph, "invalid", false)
	assert.Error(t, err)
	assert.Empty(t, invalid)
	assert.Equal(t, "Invalid output format: invalid", err.Error())
}

func TestPrintAllGraphs(t *testing.T) {
	// Create a test graph
	testGraph := &TaskGraph{
		Nodes: map[string]*TaskNode{
			"task1": {
				Name:        "task1",
				TaskRefName: "taskRef1",
				Dependencies: []*TaskNode{
					{
						Name:        "task2",
						TaskRefName: "taskRef2",
					},
				},
			},
			"task2": {
				Name:        "task2",
				TaskRefName: "taskRef2",
				Dependencies: []*TaskNode{
					{
						Name:        "task3",
						TaskRefName: "taskRef3",
					},
				},
			},
		},
	}

	// Create a test output format and withTaskRef value
	testOutputFormat := "dot"
	testWithTaskRef := true

	// Test the PrintAllGraphs method
	err := PrintAllGraphs([]*TaskGraph{testGraph}, testOutputFormat, testWithTaskRef)
	assert.NoError(t, err)
}

func TestPrintAllGraphsWithUnsupportedFormat(t *testing.T) {
	// Create a test graph
	testGraph := &TaskGraph{}

	// Create a test output format and withTaskRef value
	testOutputFormat := "FAIL"
	testWithTaskRef := true

	// Test the PrintAllGraphs method
	err := PrintAllGraphs([]*TaskGraph{testGraph}, testOutputFormat, testWithTaskRef)
	assert.Error(t, err)
	// contains error message
	assert.Contains(t, err.Error(), "Invalid output format: FAIL")
}

func TestWriteAllGraphs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-output")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer func() {
		// Remove the temporary directory and check for errors
		if err = os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temporary directory: %v", err)
		}
	}()

	// Create a test graph
	testGraph := &TaskGraph{
		PipelineName: "test-pipeline",
		Nodes: map[string]*TaskNode{
			"task1": {
				Name:        "task1",
				TaskRefName: "taskRef1",
				Dependencies: []*TaskNode{
					{
						Name:        "task2",
						TaskRefName: "taskRef2",
					},
				},
			},
			"task2": {
				Name:        "task2",
				TaskRefName: "taskRef2",
				Dependencies: []*TaskNode{
					{
						Name:        "task3",
						TaskRefName: "taskRef3",
					},
				},
			},
			"task3": {
				Name:        "task3",
				TaskRefName: "taskRef3",
				Dependencies: []*TaskNode{
					{
						Name:        "task4",
						TaskRefName: "taskRef4",
					},
				},
			},
			"task4": {
				Name:        "task4",
				TaskRefName: "taskRef4",
			},
		},
	}

	// Write the test graph to all supported formats
	err = WriteAllGraphs([]*TaskGraph{testGraph}, "dot", tempDir, true)
	assert.NoError(t, err)
	err = WriteAllGraphs([]*TaskGraph{testGraph}, "puml", tempDir, true)
	assert.NoError(t, err)
	err = WriteAllGraphs([]*TaskGraph{testGraph}, "mmd", tempDir, true)
	assert.NoError(t, err)

	// Check that the files were created
	_, err = os.Stat(filepath.Join(tempDir, "test-pipeline.dot"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "test-pipeline.puml"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "test-pipeline.mmd"))
	assert.NoError(t, err)
}
