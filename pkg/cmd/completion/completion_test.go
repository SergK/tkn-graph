package completion

import (
	"strings"
	"testing"

	"github.com/sergk/tkn-graph/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestCompletion_Empty(t *testing.T) {
	completion := Command()
	out, err := test.ExecuteCommand(completion)
	if err == nil {
		t.Errorf("No errors was defined. Output: %s", out)
	}
	expected := "accepts 1 arg(s), received 0"
	assert.Contains(t, out, expected)
}

func TestCompletionZSH(t *testing.T) {
	cmd := Command()
	output := genZshCompletion(cmd)
	assert.True(t, strings.HasPrefix(output, "#compdef"))
}
