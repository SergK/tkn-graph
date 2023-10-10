package pipeline

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

func TestGetAllPipelines(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	// Define the expected pipeline runs
	expectedPipelines := []v1.Pipeline{
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
	for _, pr := range expectedPipelines {
		_, err := fakeClient.TektonV1().Pipelines(namespace).Create(context.TODO(), &pr, metav1.CreateOptions{})
		if err != nil {
			t.Fatalf("Error creating fake pipeline run: %v", err)
		}
	}

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline runs
	pipelines, err := GetAllPipelines(c, namespace)
	if err != nil {
		t.Fatalf("Error getting pipeline runs: %v", err)
	}

	// Check that the pipeline runs are as expected
	if len(pipelines) != len(expectedPipelines) {
		t.Fatalf("Expected %d pipeline runs, got %d", len(expectedPipelines), len(pipelines))
	}
	for i, pr := range pipelines {
		if pr.Name != expectedPipelines[i].Name {
			t.Fatalf("Expected pipeline run %d to have name %s, got %s", i, expectedPipelines[i].Name, pr.Name)
		}
	}

}

func TestGetAllPipelineWithError(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline runs
	_, err := GetAllPipelines(c, namespace)
	if err == nil {
		t.Fatal("GetAllPipelines did not return an error, expected an error")
	}
	if err.Error() != "no Pipelines found in namespace my-namespace" {
		t.Fatalf("Expected error message to be 'no Pipelines found in namespace my-namespace', got %s", err.Error())
	}
}

func TestGetPipelineByName(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	// Define the expected pipeline run
	expectedPipeline := &v1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pipeline-1",
			Namespace: namespace,
		},
	}

	// Create the fake pipeline run
	_, err := fakeClient.TektonV1().Pipelines(namespace).Create(context.TODO(), expectedPipeline, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("Error creating fake pipelineRun run: %v", err)
	}

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline run
	pipelineRun, err := GetPipelineByName(c, expectedPipeline.Name, namespace)
	if err != nil {
		t.Fatalf("Error getting pipeline run: %v", err)
	}

	// Check that the pipeline run is as expected
	if pipelineRun.Name != expectedPipeline.Name {
		t.Fatalf("Expected pipeline run to have name %s, got %s", expectedPipeline.Name, pipelineRun.Name)
	}
}

func TestPipelineByNameWithError(t *testing.T) {
	fakeClient := fakeclient.NewSimpleClientset()

	c := &cli.Clients{
		Tekton: fakeClient,
	}

	// Get the pipeline runs
	_, err := GetPipelineByName(c, "fake-pipeline", namespace)
	if err == nil {
		t.Fatal("GetPipelineByName did not return an error, expected an error")
	}
	if err.Error() != "failed to get Pipeline with name fake-pipeline: pipelines.tekton.dev \"fake-pipeline\" not found" {
		t.Fatalf("Expected error message to be 'failed to get Pipeline with name fake-pipeline: pipelines.tekton.dev \"fake-pipeline\" not found', got %s", err.Error())
	}
}
