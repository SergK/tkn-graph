package pipeline

import (
	common "github.com/sergk/tkn-graph/pkg/cmd/common"
	"github.com/sergk/tkn-graph/pkg/pipeline"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
)

func graphCommand(p cli.Params) *cobra.Command {
	return common.CreateGraphCommand(p, &PipelineFetcher{
		GetPipelineByNameFunc: pipeline.GetPipelineByName,
		GetAllPipelinesFunc:   pipeline.GetAllPipelines,
	})
}
