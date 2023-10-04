package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	v1pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	fakeclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/fake"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace    = "test-namespace"
	pipelineName = "test-pipeline"
)

func TestClient_GetPipelinesNoError(t *testing.T) {
	// Create a fake Tekton clientset for testing
	fakeClient := fakeclient.NewSimpleClientset()

	expectedPipelines := []v1pipeline.Pipeline{
		{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pipeline-1",
				Namespace: namespace,
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pipeline-2",
				Namespace: namespace,
			},
		},
	}

	// Add the test pipelines to the fake clientset
	for i := range expectedPipelines {
		_, err := fakeClient.TektonV1().Pipelines(namespace).Create(context.TODO(), &expectedPipelines[i], v1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed to create Pipeline: %v", err)
		}
	}

	// Create a client instance with the fake clientset
	client := &Client{
		tektonClient: fakeClient,
	}

	pipelines, err := client.GetPipelines(namespace)
	assert.NoError(t, err)
	// assert number of pipelines
	assert.Equal(t, 2, len(pipelines))
	// Assert that the returned pipelines are the same as the expected pipelines
	assert.Equal(t, expectedPipelines, pipelines)
}

func TestClient_GetPipelinesError(t *testing.T) {
	// Create a fake Tekton clientset for testing
	fakeClient := fakeclient.NewSimpleClientset()

	// Create a client instance with the fake clientset
	client := &Client{
		tektonClient: fakeClient,
	}

	// Call the GetPipelines function
	_, err := client.GetPipelines(namespace)

	// Assert that error occurred
	assert.Error(t, err)

	// Assert that the error message is as expected
	assert.Equal(t, "no Pipelines found in namespace test-namespace", err.Error())
}

func TestClient_GetPipelineRunsNoError(t *testing.T) {
	// Create a fake Tekton clientset for testing
	fakeClient := fakeclient.NewSimpleClientset()

	expectedPipelineRuns := []v1pipeline.PipelineRun{
		{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pipeline-1",
				Namespace: namespace,
			},
		},
		{
			ObjectMeta: v1.ObjectMeta{
				Name:      "pipeline-2",
				Namespace: namespace,
			},
		},
	}

	// Add the test pipelineruns to the fake clientset
	for i := range expectedPipelineRuns {
		_, err := fakeClient.TektonV1().PipelineRuns(namespace).Create(context.TODO(), &expectedPipelineRuns[i], v1.CreateOptions{})
		if err != nil {
			t.Fatalf("failed to create PipelineRun: %v", err)
		}
	}

	// Create a client instance with the fake clientset
	client := &Client{
		tektonClient: fakeClient,
	}

	pipelineruns, err := client.GetPipelineRuns(namespace)
	assert.NoError(t, err)
	// assert number of pipelineruns
	assert.Equal(t, 2, len(pipelineruns))
	// Assert that the returned pipelineruns are the same as the expected pipelineruns
	assert.Equal(t, expectedPipelineRuns, pipelineruns)
}

func TestClient_GetPipelineRunsError(t *testing.T) {
	// Create a fake Tekton clientset for testing
	fakeClient := fakeclient.NewSimpleClientset()

	// Create a client instance with the fake clientset
	client := &Client{
		tektonClient: fakeClient,
	}

	// Call the GetPipelines function
	_, err := client.GetPipelineRuns(namespace)

	// Assert that error occurred
	assert.Error(t, err)

	// Assert that the error message is as expected
	assert.Equal(t, "no PipelineRuns found in namespace test-namespace", err.Error())
}
