package console

import (
	"testing"
)

func TestNewTree(t *testing.T) {
	t.Parallel()

	tree := NewTree("ROOT")
	n1 := tree.Root.AddChild("Subnode 1")
	n2 := tree.Root.AddChild("Subnode 2")
	n1.AddChild("Subnode 1.1")
	n1.AddChild("Subnode 1.2")
	n1.AddChild("Subnode 1.3")
	n2.AddChild("Subnode 2.1")

	tree.Write(t.Output())
}
