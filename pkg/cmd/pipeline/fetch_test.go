package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetByName(t *testing.T) {
	fetcher := &PipelineFetcher{
		GetPipelineByNameFunc: func(cs *cli.Clients, name, namespace string) (*v1.Pipeline, error) {
			// Return a dummy pipeline
			return &v1.Pipeline{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
			}, nil
		},
	}

	p, err := fetcher.GetByName(nil, "pipeline1", "default")

	assert.NoError(t, err)
	assert.Equal(t, "pipeline1", p.Name)
}

func TestGetAll(t *testing.T) {
	fetcher := &PipelineFetcher{
		GetAllPipelinesFunc: func(cs *cli.Clients, namespace string) ([]v1.Pipeline, error) {
			// Return a slice of dummy pipelines
			return []v1.Pipeline{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pipeline1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pipeline2",
					},
				},
			}, nil
		},
	}

	ps, err := fetcher.GetAll(nil, "default")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(ps))
	assert.Equal(t, "pipeline1", ps[0].Name)
	assert.Equal(t, "pipeline2", ps[1].Name)
}
