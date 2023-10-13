package cmd

import (
	"github.com/sergk/tkn-graph/pkg/cmd/pipeline"
	"github.com/sergk/tkn-graph/pkg/cmd/pipelinerun"
	"github.com/sergk/tkn-graph/pkg/cmd/version"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
)

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (eq .Annotations.commandType "main")}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Other Commands:{{range .Commands}}{{if (or (eq .Annotations.commandType "utility") (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func Root(p cli.Params) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "tkn-graph",
		Short:        "Generate a graph of a Tekton Pipelines",
		Long:         "tkn-graph is a command-line tool for generating graphs from Tekton Pipelines and PipelineRuns.",
		SilenceUsage: true,
	}
	cmd.SetUsageTemplate(usageTemplate)

	cmd.AddCommand(
		pipeline.Command(p),
		pipelinerun.Command(p),
		version.Command(),
	)

	return cmd
}
