package pipeline

import (
	"context"
	"fmt"

	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAllPipelines(c *cli.Clients, ns string) ([]v1.Pipeline, error) {
	pipelines, err := c.Tekton.TektonV1().Pipelines(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Pipelines: %w", err)
	}
	if len(pipelines.Items) == 0 {
		return nil, fmt.Errorf("no Pipelines found in namespace %s", ns)
	}
	return pipelines.Items, nil
}
