package pipelinerun

import (
	"fmt"

	common "github.com/sergk/tkn-graph/pkg/cmd/common"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type PipelineRunFetcher struct {
	GetPipelineRunByNameFunc func(cs *cli.Clients, name, namespace string) (*v1.PipelineRun, error)
	GetAllPipelineRunsFunc   func(cs *cli.Clients, namespace string) ([]v1.PipelineRun, error)
	GetPipelineByNameFunc    func(cs *cli.Clients, name, namespace string) (*v1.Pipeline, error)
}

func (f *PipelineRunFetcher) GetByName(cs *cli.Clients, name, namespace string) (*common.Pipeline, error) {
	pr, err := f.GetPipelineRunByNameFunc(cs, name, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get PipelineRun by name: %w", err)
	}

	// Fetch the Pipeline that the PipelineRun is based on
	p, err := f.GetPipelineByNameFunc(cs, pr.Spec.PipelineRef.Name, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pipeline by name: %w", err)
	}

	return &common.Pipeline{
		Name:           name,
		TektonPipeline: *p,
	}, nil
}

func (f *PipelineRunFetcher) GetAll(cs *cli.Clients, namespace string) ([]common.Pipeline, error) {
	prs, err := f.GetAllPipelineRunsFunc(cs, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get all PipelineRuns: %w", err)
	}

	var cp []common.Pipeline
	for i := range prs {
		pipeline, err := f.GetPipelineByNameFunc(cs, prs[i].Spec.PipelineRef.Name, namespace)
		if err != nil {
			return nil, fmt.Errorf("failed to get Pipeline by name: %w", err)
		}
		cp = append(cp, common.Pipeline{
			Name:           prs[i].Name,
			TektonPipeline: *pipeline,
		})
	}

	return cp, nil
}
