// pkg/cmd/common/graph_test.go
package common

import (
	"testing"

	"github.com/sergk/tkn-graph/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tektoncd/cli/pkg/cli"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// MockGraphFetcher is a mock implementation of the GraphFetcher interface
type MockGraphFetcher struct {
	mock.Mock
}

func (m *MockGraphFetcher) GetByName(cs *cli.Clients, name, namespace string) (*Pipeline, error) {
	args := m.Called(cs, name, namespace)
	return args.Get(0).(*Pipeline), args.Error(1)
}

func (m *MockGraphFetcher) GetAll(cs *cli.Clients, namespace string) ([]Pipeline, error) {
	args := m.Called(cs, namespace)
	return args.Get(0).([]Pipeline), args.Error(1)
}

// TestCreateGraphCommand tests the CreateGraphCommand function
func TestCreateGraphCommand(t *testing.T) {
	p := &test.Params{}
	fetcher := new(MockGraphFetcher)

	cmd := CreateGraphCommand(p, fetcher)

	assert.Equal(t, "graph", cmd.Use)
	assert.Equal(t, []string{"g"}, cmd.Aliases)
	assert.Equal(t, "Generates Graph", cmd.Short)
	assert.Equal(t, map[string]string{"commandType": "main"}, cmd.Annotations)
	assert.True(t, cmd.SilenceUsage)
}

// TestRunGraphCommand tests the RunGraphCommand function
func TestRunGraphCommand(t *testing.T) {
	p := &test.Params{}
	p.SetNamespace("default")

	fetcher := new(MockGraphFetcher)
	fetcher.On("GetByName", mock.Anything, "pipeline1", "default").Return(&Pipeline{
		Name: "pipeline1",
		TektonPipeline: v1.Pipeline{
			Spec: v1.PipelineSpec{
				Tasks: []v1.PipelineTask{
					{
						Name: "task1",
						TaskRef: &v1.TaskRef{
							Name: "task1",
						},
					},
				},
			},
		},
	}, nil)

	testCases := []struct {
		name         string
		outputFormat string
		expectError  bool
		errorMessage string
	}{
		{"valid output format", "dot", false, ""},
		{"invalid output format", "wrong", true, "failed to print graph: Failed to generate output: Invalid output format: wrong"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &GraphOptions{
				OutputFormat: tc.outputFormat,
			}
			args := []string{"pipeline1"}

			err := RunGraphCommand(p, opts, fetcher, args)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}

			fetcher.AssertExpectations(t)
		})
	}
}

func TestRunGraphCommandWithGetAll(t *testing.T) {
	p := &test.Params{}
	p.SetNamespace("default")

	fetcher := new(MockGraphFetcher)
	fetcher.On("GetAll", mock.Anything, "default").Return([]Pipeline{
		{
			Name: "pipeline1",
			TektonPipeline: v1.Pipeline{
				Spec: v1.PipelineSpec{
					Tasks: []v1.PipelineTask{
						{
							Name: "task1",
							TaskRef: &v1.TaskRef{
								Name: "task1",
							},
						},
					},
				},
			},
		},
		{
			Name: "pipeline2",
			TektonPipeline: v1.Pipeline{
				Spec: v1.PipelineSpec{
					Tasks: []v1.PipelineTask{
						{
							Name: "task2",
							TaskRef: &v1.TaskRef{
								Name: "task2",
							},
						},
					},
				},
			},
		},
	}, nil)

	testCases := []struct {
		name         string
		outputFormat string
		expectError  bool
		errorMessage string
	}{
		{"valid output format", "dot", false, ""},
		{"invalid output format", "wrong", true, "failed to save graph: Failed to generate output: Invalid output format: wrong"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &GraphOptions{
				OutputFormat: tc.outputFormat,
				OutputDir:    "/tmp",
			}
			args := []string{} // Empty args to trigger GetAll

			err := RunGraphCommand(p, opts, fetcher, args)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}

			fetcher.AssertExpectations(t)
		})
	}
}

func TestRunGraphCommandWithTooManyArgs(t *testing.T) {
	p := &test.Params{}
	p.SetNamespace("default")

	fetcher := new(MockGraphFetcher)
	opts := &GraphOptions{
		OutputFormat: "dot",
	}
	args := []string{"pipeline1", "pipeline2"} // Two arguments to trigger an error

	err := RunGraphCommand(p, opts, fetcher, args)

	assert.Error(t, err)
	assert.Equal(t, "too many arguments. Provide either no arguments to get all Pipelines or a single Pipeline name", err.Error())
}
