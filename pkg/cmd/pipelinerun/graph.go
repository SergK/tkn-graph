package pipelinerun

import (
	common "github.com/sergk/tkn-graph/pkg/cmd/common"
	"github.com/sergk/tkn-graph/pkg/pipeline"
	"github.com/sergk/tkn-graph/pkg/pipelinerun"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
)

func graphCommand(p cli.Params) *cobra.Command {
	return common.CreateGraphCommand(p, &PipelineRunFetcher{
		GetPipelineRunByNameFunc: pipelinerun.GetPipelineRunsByName,
		GetAllPipelineRunsFunc:   pipelinerun.GetAllPipelineRuns,
		GetPipelineByNameFunc:    pipeline.GetPipelineByName,
	})
}
