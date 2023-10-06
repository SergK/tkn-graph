package pipelinerun

import (
	"context"
	"fmt"

	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAllPipelineRuns(c *cli.Clients, ns string) ([]v1.PipelineRun, error) {
	pipelineruns, err := c.Tekton.TektonV1().PipelineRuns(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get PipelineRuns: %w", err)
	}
	if len(pipelineruns.Items) == 0 {
		return nil, fmt.Errorf("no PipelineRuns found in namespace %s", ns)
	}
	return pipelineruns.Items, nil
}
