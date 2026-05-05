package string_tools_test

import (
	"testing"

	stringtools "github.com/guionardo/gs-dev/internal/tools/strings_tools"
	"github.com/stretchr/testify/assert"
)

func TestSplitString(t *testing.T) {
	t.Parallel()

	var w1, w2, w3 string

	stringtools.SplitString("W1 W2 W3 W4", " ", &w1, &w2, &w3)

	assert.Equal(t, "W1", w1)
	assert.Equal(t, "W2", w2)
	assert.Equal(t, "W3 W4", w3)
}
