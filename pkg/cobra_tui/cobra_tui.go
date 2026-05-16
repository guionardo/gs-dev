package cobratui

import (
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

type CobraTUI struct {
	rootCommand *cobra.Command
}

func NewCobraTUI(rootCommand *cobra.Command) *CobraTUI {
	return &CobraTUI{
		rootCommand: rootCommand,
	}
}
func add(node *tview.TreeNode, cmd *cobra.Command) {
	childNode := tview.NewTreeNode(cmd.Name()).SetReference(cmd).SetSelectable(true)
	node.AddChild(childNode)

	for _, subCmd := range cmd.Commands() {
		add(childNode, subCmd)
	}
}
func (t *CobraTUI) Run() error {
	tvRoot := tview.NewTreeNode(".").SetSelectable(true)
	tvCmd := tview.NewTreeView().SetRoot(tvRoot)

	tvCmd.SetBorder(true).SetTitle("Commands")

	for _, cmd := range t.rootCommand.Commands() {
		add(tvRoot, cmd)
	}

	formBox := tview.NewBox().SetBorder(true).SetTitle("Form")

	flex := tview.NewFlex().
		AddItem(tvCmd, 0, 1, true).
		AddItem(formBox, 0, 2, false)
	if err := tview.NewApplication().SetRoot(flex, true).Run(); err != nil {
		return err
	}

	return nil
}
