package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	// Version is the current version of the CLI
	devVersion = "dev"
)

var (
	cliVersion = devVersion
)

// Command returns the version command
func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Annotations: map[string]string{
			"commandType": "utility",
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", cliVersion)
		},
	}
}
