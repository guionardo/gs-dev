package commands

import (
	"path/filepath"

	"github.com/guionardo/gs-dev/internal/git"
)

type CommonCommand struct {
	hasInitAlias bool
	arguments    string
	usesOutput   bool
	name         string
	description  string
}

func (c *CommonCommand) InitCommandVariables(name string, description string, hasInitAlias bool, arguments string, usesOutput bool) {
	c.name = name
	c.description = description
	c.hasInitAlias = hasInitAlias
	c.arguments = arguments
	c.usesOutput = usesOutput
}
func (c *CommonCommand) HasInitAlias() bool {
	return c.hasInitAlias
}

func (c *CommonCommand) GetArguments() string {
	return c.arguments
}

func (c *CommonCommand) UsesOutput() bool {
	return c.usesOutput
}

func (c *CommonCommand) GetName() string {
	return c.name
}

func (c *CommonCommand) GetDescription() string {
	return c.description
}
func (c *CommonCommand) GetRepositoryName() string {
	gitConfig, err := git.NewGitConfig(".git/config")
	if err == nil {
		return gitConfig.RepositoryName
	}

	dir, err := filepath.Abs(".")
	if err != nil {
		return ""
	}

	return filepath.Base(dir)
}
