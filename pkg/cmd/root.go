package cmd

import (
	"os"

	"github.com/sergk/tkn-graph/pkg/cmd/pipeline"
	"github.com/sergk/tkn-graph/pkg/cmd/pipelinerun"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tektoncd/cli/pkg/cli"
)

const usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
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
	// Reset CommandLine so we don't get the flags from the libraries, i.e:
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

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
	)

	return cmd
}
