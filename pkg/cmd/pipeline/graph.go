package pipeline

import (
	"fmt"

	"github.com/sergk/tkn-graph/pkg/cmd/graphutil"
	pipelinepkg "github.com/sergk/tkn-graph/pkg/pipeline"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

func graphCommand(p cli.Params) *cobra.Command {
	return graphutil.NewGraphCommand(p, func(cmd *cobra.Command, args []string, p cli.Params, opts *graphutil.GraphOptions) error {
		return graphutil.RunGraphCommand(args, p, opts, func(cs *cli.Clients, args []string, namespace string) ([]graphutil.GraphData, error) {
			var pipelines []v1.Pipeline
			var err error

			switch len(args) {
			case 1:
				var pipeline *v1.Pipeline
				pipeline, err = pipelinepkg.GetPipelineByName(cs, args[0], namespace)
				if err != nil {
					return nil, fmt.Errorf("failed to run GetPipelineByName: %w", err)
				}
				pipelines = append(pipelines, *pipeline)
			case 0:
				pipelines, err = pipelinepkg.GetAllPipelines(cs, namespace)
				if err != nil {
					return nil, fmt.Errorf("failed to run GetAllPipelines: %w", err)
				}
			default:
				return nil, fmt.Errorf("too many arguments. Provide either no arguments to get all Pipelines or a single Pipeline name")
			}

			var data []graphutil.GraphData
			for i := range pipelines {
				data = append(data, graphutil.GraphData{Name: pipelines[i].Name, Spec: pipelines[i].Spec})
			}

			return data, nil
		})
	})
}
