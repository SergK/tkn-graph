package taskgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const (
	testPipelineName = "test-pipeline"
)

func TestBuildTaskGraph(t *testing.T) {
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)

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
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
	graph.PipelineName = testPipelineName

	// Test the ToDOT method
	dot := graph.ToDOT()
	assert.Equal(t, "digraph", dot.Format)
	assert.Equal(t, testPipelineName, dot.Name)
	assert.ElementsMatch(t, []string{
		"  \"task1\" -> \"task2\"",
		"  \"task1\" -> \"task3\"",
		"  \"task2\" -> \"task3\"",
	}, dot.Edges)
}

func TestTaskGraphToDOTWithTaskRef(t *testing.T) {
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
	graph.PipelineName = testPipelineName

	// Test the ToDOTWithTaskRef method
	dot := graph.ToDOTWithTaskRef()
	assert.Equal(t, "digraph", dot.Format)
	assert.Equal(t, testPipelineName, dot.Name)
	// dot.Edges = append(dot.Edges, fmt.Sprintf("  \"%s\n(%s)\" -> \"%s\n(%s)\"", dep.Name, dep.TaskRefName, node.Name, node.TaskRefName))
	assert.ElementsMatch(t, []string{
		"  \"task1\n(taskRef1)\" -> \"task2\n(taskRef2)\"",
		"  \"task1\n(taskRef1)\" -> \"task3\n(taskRef3)\"",
		"  \"task2\n(taskRef2)\" -> \"task3\n(taskRef3)\"",
	}, dot.Edges)
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
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
	graph.PipelineName = testPipelineName

	// Test the ToPlantUML method
	plantuml := graph.ToPlantUML()
	expected := "@startuml\nhide empty description\ntitle test-pipeline\n\n[*] --> task1\n" +
		"task2 <-down- task1\n" +
		"task3 <-down- task1\n" +
		"task3 <-down- task2\n\n@enduml\n"
	assert.Equal(t, expected, plantuml)
}

func TestTaskGraphToPlantUMLWithTaskRef(t *testing.T) {
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
	graph.PipelineName = testPipelineName

	// Test the ToPlantUMLWithTaskRef method
	plantuml := graph.ToPlantUMLWithTaskRef()
	expected := "@startuml\nhide empty description\ntitle test-pipeline\n\n[*] --> task1\n" +
		"task2 <-down- task1\n" +
		"task3 <-down- task1\n" +
		"task3 <-down- task2\n\n" +
		"task1: taskRef1\n" +
		"task2: taskRef2\n" +
		"task3: taskRef3\n\n@enduml\n"
	assert.Equal(t, expected, plantuml)
}

func TestTaskGraphToMermaid(t *testing.T) {
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
	graph.PipelineName = testPipelineName

	// Test the ToMermaid method
	mermaid := graph.ToMermaid()
	expected := "---\ntitle: test-pipeline\n---\nflowchart TD\n   task1 --> task2\n   task1 --> task3\n   task2 --> task3\n"
	assert.Equal(t, expected, mermaid)
}

func TestTaskGraphToMermaidWithTaskRef(t *testing.T) {
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
	graph.PipelineName = testPipelineName

	// Test the ToMermaidWithTaskRef method
	mermaid := graph.ToMermaidWithTaskRef()
	expected := "---\ntitle: test-pipeline\n---\nflowchart TD\n" +
		"   task1(\"task1\n   (taskRef1)\") --> task2(\"task2\n   (taskRef2)\")\n" +
		"   task1(\"task1\n   (taskRef1)\") --> task3(\"task3\n   (taskRef3)\")\n" +
		"   task2(\"task2\n   (taskRef2)\") --> task3(\"task3\n   (taskRef3)\")\n"
	assert.Equal(t, expected, mermaid)
}

func TestFormatFunc(t *testing.T) {
	// Define some test tasks
	tasks := []v1pipeline.PipelineTask{
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

	// Build the task graph
	graph := BuildTaskGraph(tasks)
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
