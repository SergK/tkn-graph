package main

import (
	"os"

	"github.com/sergk/tkn-graph/pkg/cmd"
	"github.com/tektoncd/cli/pkg/cli"
)

func main() {
	tp := &cli.TektonParams{}
	tkn := cmd.Root(tp)

	if err := tkn.Execute(); err != nil {
		os.Exit(1)
	}
}
