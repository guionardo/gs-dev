package commands

import (
	"path/filepath"
	"sync"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/git"
	"github.com/spf13/cobra"
)

type CommonCommand struct {
	hasInitAlias    bool
	defaultArgument string
	usesOutput      bool
	name            string
	description     string
	buildCommand    *cobra.Command
	self            any

	setupFunc  func(*config.ConfigRoot) error
	setupDone  bool
	setupLock  sync.Mutex
	setupError error
}

func (c *CommonCommand) InitCommandVariables(
	name string,
	description string,
	self any,
	setupFunc func(*config.ConfigRoot) error,
) *CommonCommand {
	c.name = name
	c.description = description
	c.self = self
	c.setupFunc = setupFunc

	return c
}

// WithDefaultArgument adds a extra command argument to the initialization script
func (c *CommonCommand) WithDefaultArgument(defaultArgument string) *CommonCommand {
	c.defaultArgument = defaultArgument
	return c
}

// WithInitAlias is used to build an function alias in shell to run the command
func (c *CommonCommand) WithInitAlias() *CommonCommand {
	c.hasInitAlias = true
	return c
}

func (c *CommonCommand) WithOutput() *CommonCommand {
	c.usesOutput = true
	return c
}

func (c *CommonCommand) HasInitAlias() bool {
	return c.hasInitAlias
}

func (c *CommonCommand) GetArguments() string {
	return c.defaultArgument
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

func (c *CommonCommand) GetCobraCommand() *cobra.Command {
	if c.buildCommand != nil {
		return c.buildCommand
	}

	if cmd, ok := c.self.(interface{ BuildCommand() *cobra.Command }); ok {
		c.buildCommand = cmd.BuildCommand()
		if c.hasInitAlias {
			if c.buildCommand.Annotations == nil {
				c.buildCommand.Annotations = make(map[string]string)
			}

			c.buildCommand.Annotations["initAlias"] = "true"
		}

		return c.buildCommand
	}

	panic("GetCobraCommand not implemented for " + c.name)
}

func (c *CommonCommand) Setup(configuration *config.ConfigRoot) error {
	c.setupLock.Lock()
	defer c.setupLock.Unlock()

	if c.setupDone {
		return c.setupError
	}

	c.setupError = c.setupFunc(configuration)
	c.setupDone = c.setupError == nil

	return c.setupError
}
