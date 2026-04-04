package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sort"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/dialog"
	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/logging"
	postcommand "github.com/guionardo/gs-dev/pkg/post_command"
	"github.com/spf13/cobra"
)

type (
	CommandManager struct {
		commands   map[string]Command
		configFile *config.ConfigFile
	}

	Command interface {
		// Initialize the command
		Init()
		Setup(configuration *config.ConfigFile) error
		GetName() string
		GetDescription() string

		// if the command is a cobra command, it will be added to the root command
		GetCobraCommand() *cobra.Command

		// if the command is a TUI command, it will be added to the TUI command
		GetTUICommand() func() error

		// if the command has an init alias, it will be added to the init script
		HasInitAlias() bool

		// Use by the init alias
		GetArguments() string

		UsesOutput() bool
	}
)

func NewCommandManager() (*CommandManager, error) {
	configFile, err := config.NewConfigFile(path.Join(config.GetConfigDir(), "config.json"))

	if recErr, ok := errors.AsType[errs.Error](err); ok {
		if !recErr.IsRecoverable() {
			return nil, recErr
		}
	}

	return &CommandManager{
		commands:   make(map[string]Command),
		configFile: configFile,
	}, nil
}

func (c *CommandManager) Register(commands ...Command) *CommandManager {
	for _, command := range commands {
		command.Init()
		c.commands[command.GetName()] = command

		err := command.Setup(c.configFile)
		if err != nil {
			slog.Error("Error setting up command", slog.String("command", command.GetName()), slog.Any("error", err))
		}
	}

	// Add the registered commands to the init setup
	_ = NewInitSetup(build.AppName, commands...)

	return c
}

func (c *CommandManager) GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gs-dev",
		Short: "gs-dev",
		Long:  "gs-dev",
		RunE:  c.Run,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.Flags().GetBool("debug")

			logging.PreSetupLog("Debug mode enabled", slog.LevelDebug)
			logging.Setup(debug, os.Stdout)
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return postcommand.WriteOutput()
		},
	}

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	aliasCommands := []postcommand.Command{}

	for _, command := range c.commands {
		cmd := command.GetCobraCommand()
		rootCmd.AddCommand(cmd)

		if cmd.Annotations != nil && len(cmd.Annotations[consts.UseOutputAnnotation]) > 0 {
			aliasCommands = append(aliasCommands, postcommand.NewCommand(command.GetName(), true, command.GetName()))
		}
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			if full, _ := cmd.Flags().GetBool("full"); full {
				cmd.Printf("%s %s\n%s", build.AppName, build.Version, build.BuildInfo)
			} else {
				cmd.Printf("%s", build.Version)
			}

			if !build.IsValidBinary {
				cmd.Printf(" %s (development binary)", build.ExecutableName)
			}
		},

		Version: build.Version,
	}
	versionCmd.Flags().BoolP("full", "f", false, "Show full version information")
	rootCmd.AddCommand(versionCmd)

	initSetup := postcommand.NewInitSetup(build.AppName, aliasCommands...)
	if err := initSetup.UpdateRootCommand(rootCmd); err != nil {
		slog.Error("Error updating root command", slog.Any("error", err))
		os.Exit(1)
	}

	return rootCmd
}

func (c *CommandManager) Run(cmd *cobra.Command, args []string) error {
	var commandNames, commandDescriptions []string

	for name := range c.commands {
		if c.commands[name].GetCobraCommand() != nil {
			commandNames = append(commandNames, name)
			commandDescriptions = append(commandDescriptions, c.commands[name].GetDescription())
		}
	}

	sort.Strings(commandNames)

	commands := []dialog.ChooseItem{}
	for name, command := range c.commands {
		if command.GetCobraCommand() != nil {
			commands = append(commands, dialog.ChooseItem{
				Name:        name,
				Description: command.GetDescription(),
			})
		}
	}

	fmt.Println(colors.Parse("{primary}[%s]{normal} v %s %s\nCtrl+C to exit", build.AppName, build.Version, build.BuildInfo))

	commandName, err := dialog.Choose(">", commands...)
	if err != nil {
		return err
	}

	command := c.commands[commandName]

	tuiCommand := command.GetTUICommand()
	if tuiCommand != nil {
		return tuiCommand()
	}

	return command.GetCobraCommand().Execute()
}
