package graphutil_test

import (
	"errors"
	"testing"

	"github.com/sergk/tkn-graph/pkg/cmd/graphutil"
	"github.com/spf13/cobra"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

func TestNewGraphCommand(t *testing.T) {
	p := &cli.TektonParams{}
	runE := func(cmd *cobra.Command, args []string, p cli.Params, opts *graphutil.GraphOptions) error {
		return nil
	}

	cmd := graphutil.NewGraphCommand(p, runE)

	if cmd.Use != "graph" {
		t.Errorf("Expected command use to be 'graph', got '%s'", cmd.Use)
	}
}

func TestRunGraphCommand(t *testing.T) {
	tests := []struct {
		name           string
		outputFormat   string
		outputDir      string
		expectingError bool
	}{
		{
			name:           "valid output format",
			outputFormat:   "dot",
			outputDir:      "",
			expectingError: false,
		},
		{
			name:           "valid output format",
			outputFormat:   "mmd",
			outputDir:      "/tmp",
			expectingError: false,
		},
		{
			name:           "invalid output format for print to screen",
			outputFormat:   "not supported",
			outputDir:      "/tmp",
			expectingError: true,
		},
		{
			name:           "invalid output format for output dir",
			outputFormat:   "not supported",
			outputDir:      "",
			expectingError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run this test in parallel
			// Mock the GraphFetcher and cli.Params...
			fetcher := func(cs *cli.Clients, name []string, namespace string) ([]graphutil.GraphData, error) {
				// Return some mock data...
				return []graphutil.GraphData{
					{
						Name: "test",
						Spec: v1.PipelineSpec{
							Tasks: []v1.PipelineTask{
								{
									Name:    "task1",
									TaskRef: &v1.TaskRef{Name: "task1"},
								},
								{
									Name:    "task2",
									TaskRef: &v1.TaskRef{Name: "task2"},
								},
							},
						},
					},
				}, nil
			}

			p := &cli.TektonParams{}
			p.SetNamespace("default")

			opts := &graphutil.GraphOptions{
				OutputFormat: tt.outputFormat,
				OutputDir:    tt.outputDir,
				WithTaskRef:  false,
			}

			cmd := &cobra.Command{}

			err := graphutil.RunGraphCommand(cmd, []string{}, p, opts, fetcher)

			if tt.expectingError && err == nil {
				t.Errorf("Expected an error, got nil")
			}

			if !tt.expectingError && err != nil {
				t.Errorf("Expected no error, got '%s'", err)
			}
		})
	}
}

func TestRunGraphCommand_ErrorFetcher(t *testing.T) {
	// Mock the GraphFetcher to return an error
	fetcher := func(cs *cli.Clients, name []string, namespace string) ([]graphutil.GraphData, error) {
		// Return an error
		return nil, errors.New("fetcher error")
	}

	p := &cli.TektonParams{}
	p.SetNamespace("default")

	opts := &graphutil.GraphOptions{
		OutputFormat: "dot",
		OutputDir:    "",
		WithTaskRef:  false,
	}

	cmd := &cobra.Command{}

	err := graphutil.RunGraphCommand(cmd, []string{}, p, opts, fetcher)

	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}
