package pipelinerun

import (
	"context"
	"testing"

	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	fakeclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "my-namespace"
)

func TestGetAllPipelineRuns(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	// Define the expected pipeline runs
	expectedPipelineRuns := []v1.PipelineRun{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pipeline-1",
				Namespace: namespace,
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pipeline-2",
				Namespace: namespace,
			},
		},
	}

	// Create the fake pipeline runs
	for _, pr := range expectedPipelineRuns {
		_, err := fakeClient.TektonV1().PipelineRuns(namespace).Create(context.TODO(), &pr, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error creating fake pipelineRun run: %v", err)
		}
	}

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline runs
	pipelineRuns, err := GetAllPipelineRuns(c, namespace)
	if err != nil {
		t.Fatalf("Error getting pipeline runs: %v", err)
	}

	// Check that the pipeline runs are as expected
	if len(pipelineRuns) != len(expectedPipelineRuns) {
		t.Fatalf("Expected %d pipelineRuns, got %d", len(expectedPipelineRuns), len(pipelineRuns))
	}
	for i, pr := range pipelineRuns {
		if pr.Name != expectedPipelineRuns[i].Name {
			t.Fatalf("Expected pipelineRuns %d to have name %s, got %s", i, expectedPipelineRuns[i].Name, pr.Name)
		}
	}

}

func TestGetAllPipelineRunsWithError(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline runs
	_, err := GetAllPipelineRuns(c, namespace)
	if err == nil {
		t.Fatal("GetAllPipelineRuns did not return an error, expected an error")
	}
	if err.Error() != "no PipelineRuns found in namespace my-namespace" {
		t.Fatalf("Expected error message to be 'no PipelineRuns found in namespace my-namespace', got %s", err.Error())
	}
}

func TestGetPipelineRunsByName(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	// Define the expected pipeline run
	expectedPipelineRun := &v1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pipeline-1",
			Namespace: namespace,
		},
	}

	// Create the fake pipeline run
	_, err := fakeClient.TektonV1().PipelineRuns(namespace).Create(context.TODO(), expectedPipelineRun, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating fake pipelineRun run: %v", err)
	}

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline run
	pipelineRun, err := GetPipelineRunsByName(c, expectedPipelineRun.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting pipeline run: %v", err)
	}

	// Check that the pipeline run is as expected
	if pipelineRun.Name != expectedPipelineRun.Name {
		t.Fatalf("Expected pipeline run to have name %s, got %s", expectedPipelineRun.Name, pipelineRun.Name)
	}
}

func TestPipelineByNameWithError(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline runs
	_, err := GetPipelineRunsByName(c, "fake-pipeline", namespace)
	if err == nil {
		t.Fatal("GetPipelineRunsByName did not return an error, expected an error")
	}
	if err.Error() != "failed to get PipelineRun with name fake-pipeline: pipelineruns.tekton.dev \"fake-pipeline\" not found" {
		t.Fatalf("Expected error message to be 'failed to get PipelineRun with name fake-pipeline: pipelineruns.tekton.dev \"fake-pipeline\" not found', got %s", err.Error())
	}
}
