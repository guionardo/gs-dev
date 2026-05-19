package console

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTree(t *testing.T) {
	t.Parallel()

	output := bytes.NewBufferString("")
	tree := NewTree("ROOT")
	n1 := tree.Root.AddChild("Subnode 1")
	n2 := tree.Root.AddChild("Subnode 2")
	n1.AddChild("Subnode 1.1")
	n1.AddChild("Subnode 1.2")
	n1.AddChild("Subnode 1.3")
	n2.AddChild("Subnode 2.1")

	require.NoError(t, tree.Write(output))

	const expected = ` ╮ ROOT
 ├╮ Subnode 1
 │├ Subnode 1.1
 │├ Subnode 1.2
 │╰ Subnode 1.3
 ╰╮ Subnode 2
  ╰ Subnode 2.1
`
	assert.Equal(t, expected, output.String())
}
