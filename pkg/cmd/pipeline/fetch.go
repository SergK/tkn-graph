package pipeline

import (
	"fmt"

	common "github.com/sergk/tkn-graph/pkg/cmd/common"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

type PipelineFetcher struct {
	GetPipelineByNameFunc func(cs *cli.Clients, name, namespace string) (*v1.Pipeline, error)
	GetAllPipelinesFunc   func(cs *cli.Clients, namespace string) ([]v1.Pipeline, error)
}

func (f *PipelineFetcher) GetByName(cs *cli.Clients, name, namespace string) (*common.Pipeline, error) {
	p, err := f.GetPipelineByNameFunc(cs, name, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pipeline by name: %w", err)
	}

	return &common.Pipeline{
		Name:           name,
		TektonPipeline: *p,
	}, nil
}

func (f *PipelineFetcher) GetAll(cs *cli.Clients, namespace string) ([]common.Pipeline, error) {
	ps, err := f.GetAllPipelinesFunc(cs, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get all Pipelines: %w", err)
	}

	cp := make([]common.Pipeline, 0, len(ps))
	for i := range ps {
		cp = append(cp, common.Pipeline{
			Name:           ps[i].Name,
			TektonPipeline: ps[i],
		})
	}

	return cp, nil
}
