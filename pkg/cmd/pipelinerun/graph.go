package pipelinerun

import (
	"fmt"

	"github.com/sergk/tkn-graph/pkg/cmd/graphutil"
	pipelinerunpkg "github.com/sergk/tkn-graph/pkg/pipelinerun"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

func graphCommand(p cli.Params) *cobra.Command {
	return graphutil.NewGraphCommand(p, func(cmd *cobra.Command, args []string, p cli.Params, opts *graphutil.GraphOptions) error {
		return graphutil.RunGraphCommand(cmd, args, p, opts, func(cs *cli.Clients, args []string, namespace string) ([]graphutil.GraphData, error) {
			var pipelineruns []v1.PipelineRun
			var err error

			switch len(args) {
			case 1:
				var pipelinerun *v1.PipelineRun
				pipelinerun, err = pipelinerunpkg.GetPipelineRunsByName(cs, args[0], namespace)
				if err != nil {
					return nil, fmt.Errorf("failed to run GetPipelineRunByName: %w", err)
				}
				pipelineruns = append(pipelineruns, *pipelinerun)
			case 0:
				pipelineruns, err = pipelinerunpkg.GetAllPipelineRuns(cs, namespace)
				if err != nil {
					return nil, fmt.Errorf("failed to run GetAllPipelineRuns: %w", err)
				}
			default:
				return nil, fmt.Errorf("too many arguments. Provide either no arguments to get all PipelineRuns or a single PipelineRun name")
			}

			var data []graphutil.GraphData
			for i := range pipelineruns {
				if pipelineruns[i].Status.PipelineSpec != nil {
					data = append(data, graphutil.GraphData{Name: pipelineruns[i].Name, Spec: *pipelineruns[i].Status.PipelineSpec})
				}
			}

			return data, nil
		})
	})
}
