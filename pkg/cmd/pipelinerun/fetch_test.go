package pipelinerun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetByName(t *testing.T) {
	fetcher := &PipelineRunFetcher{
		GetPipelineRunByNameFunc: func(cs *cli.Clients, name, namespace string) (*v1.PipelineRun, error) {
			// Return a dummy pipeline run
			return &v1.PipelineRun{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: v1.PipelineRunSpec{
					PipelineRef: &v1.PipelineRef{
						Name: "pipeline1",
					},
				},
			}, nil
		},
		GetPipelineByNameFunc: func(cs *cli.Clients, name, namespace string) (*v1.Pipeline, error) {
			// Return a dummy pipeline
			return &v1.Pipeline{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
			}, nil
		},
	}

	p, err := fetcher.GetByName(nil, "pipelinerun1", "default")

	assert.NoError(t, err)
	assert.Equal(t, "pipelinerun1", p.Name)
}

func TestGetAll(t *testing.T) {
	fetcher := &PipelineRunFetcher{
		GetAllPipelineRunsFunc: func(cs *cli.Clients, namespace string) ([]v1.PipelineRun, error) {
			// Return a slice of dummy pipeline runs
			return []v1.PipelineRun{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pipelinerun1",
					},
					Spec: v1.PipelineRunSpec{
						PipelineRef: &v1.PipelineRef{
							Name: "pipeline1",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pipelinerun2",
					},
					Spec: v1.PipelineRunSpec{
						PipelineRef: &v1.PipelineRef{
							Name: "pipeline2",
						},
					},
				},
			}, nil
		},
		GetPipelineByNameFunc: func(cs *cli.Clients, name, namespace string) (*v1.Pipeline, error) {
			// Return a dummy pipeline
			return &v1.Pipeline{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
			}, nil
		},
	}

	ps, err := fetcher.GetAll(nil, "default")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(ps))
	assert.Equal(t, "pipelinerun1", ps[0].Name)
	assert.Equal(t, "pipelinerun2", ps[1].Name)
}
