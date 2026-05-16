package console

import (
	"fmt"
	"io"

	"github.com/guionardo/go/flow"
)

type (
	Tree struct {
		Root *Node
	}

	Node struct {
		Text     string
		Children []*Node
		level    int
	}

	nodeOrder uint8
)

const (
	horizontal   = "─"
	vertical     = "│"
	top_left     = "╭"
	top_right    = "╮"
	bottom_left  = "╰"
	bottom_right = "╯"
	middle_left  = "├"
	middle_right = "┤"
	top_middle   = "┬"

	firstNode nodeOrder = iota
	middleNode
	lastNode
	rootNode
)

func NewTree(rootNodeText string) *Tree {
	return &Tree{&Node{Text: rootNodeText}}
}

func (t *Tree) AddChild(childText string) *Node {
	return t.Root.AddChild(childText)
}

func (n *Node) AddChild(childText string) *Node {
	node := &Node{Text: childText, level: n.level + 1}
	n.Children = append(n.Children, node)

	return node
}

func (n *Node) Write(w io.Writer, order nodeOrder, parentPrefix string) error {
	var (
		nodeSymbol string
	)

	switch order {
	case firstNode:
		nodeSymbol = middle_left
	case middleNode:
		nodeSymbol = middle_left
	case lastNode:
		nodeSymbol = bottom_left
	case rootNode:
		nodeSymbol = " "
	}

	subnodesEntry := flow.If(len(n.Children) > 0, top_right, "")
	fmt.Fprintf(w, "%s%s %s\n", parentPrefix, nodeSymbol+subnodesEntry, n.Text)

	var childOrder nodeOrder

	if order != lastNode && order != rootNode {
		parentPrefix += vertical
	} else {
		parentPrefix += " "
	}

	for i, subNode := range n.Children {
		if len(n.Children) == 1 {
			childOrder = lastNode
		} else {
			switch i {
			case 0:
				childOrder = firstNode
			case len(n.Children) - 1:
				childOrder = lastNode
			default:
				childOrder = middleNode
			}
		}

		subNode.Write(w, childOrder, parentPrefix)
	}

	return nil
}

func (t *Tree) Write(w io.Writer) error {
	return t.Root.Write(w, rootNode, "")
}
