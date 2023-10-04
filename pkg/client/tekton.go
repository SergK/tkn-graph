package client

import (
	"context"
	"fmt"
	"os"

	v1pipeline "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	tektonClient versioned.Interface
}

func NewClient() (*Client, error) {
	// Create the Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get Kubernetes configuration: %w", err)
		}
	}
	tektonClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tekton client: %w", err)
	}

	return &Client{
		tektonClient: tektonClient,
	}, nil
}

func (c *Client) GetNamespace() (string, error) {
	namespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).Namespace()
	if err != nil {
		return "", fmt.Errorf("failed to get namespace from kubeconfig: %w", err)
	}
	if namespace == "" {
		namespace = "default"
	}
	return namespace, nil
}

func (c *Client) GetPipelines(namespace string) ([]v1pipeline.Pipeline, error) {
	pipelines, err := c.tektonClient.TektonV1().Pipelines(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Pipelines: %w", err)
	}
	if len(pipelines.Items) == 0 {
		return nil, fmt.Errorf("no Pipelines found in namespace %s", namespace)
	}
	return pipelines.Items, nil
}

func (c *Client) GetPipelineRuns(namespace string) ([]v1pipeline.PipelineRun, error) {
	pipelineRuns, err := c.tektonClient.TektonV1().PipelineRuns(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PipelineRuns: %w", err)
	}
	if len(pipelineRuns.Items) == 0 {
		return nil, fmt.Errorf("no PipelineRuns found in namespace %s", namespace)
	}
	return pipelineRuns.Items, nil
}
