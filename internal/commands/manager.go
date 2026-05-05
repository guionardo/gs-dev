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
	"github.com/guionardo/gs-dev/internal/context"
	"github.com/guionardo/gs-dev/internal/dialog"
	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/interfaces"
	"github.com/guionardo/gs-dev/internal/logging"
	installservice "github.com/guionardo/gs-dev/internal/services/install"
	postcommand "github.com/guionardo/gs-dev/pkg/post_command"
	"github.com/spf13/cobra"
)

type CommandManager struct {
	commands   map[string]interfaces.Command
	configRoot *config.ConfigRoot
}

func NewCommandManager() (*CommandManager, error) {
	configFile, err := config.NewConfigFile(path.Join(config.GetConfigDir(), "config.json"))

	if recErr, ok := errors.AsType[errs.Error](err); ok {
		if !recErr.IsRecoverable() {
			return nil, recErr
		}
	}

	return &CommandManager{
		commands:   make(map[string]interfaces.Command),
		configRoot: configFile,
	}, nil
}

func (c *CommandManager) Register(commands ...interfaces.Command) *CommandManager {
	registeredCommands := make([]interfaces.Command, 0, len(commands))
	for _, command := range commands {
		command.Init()
		c.commands[command.GetName()] = command

		err := command.Setup(c.configRoot)
		if er, ok := errors.AsType[errs.Error](err); ok && er.IsRecoverable() {
			slog.Debug("Error setting up", slog.String("command", command.GetName()), slog.Any("error", err))
			continue
		} else if err != nil {
			slog.Error("Error setting up", slog.String("command", command.GetName()), slog.Any("error", err))
			continue
		}

		registeredCommands = append(registeredCommands, command)
	}

	// Add the registered commands to the init setup
	installservice.SetCommands(registeredCommands...)

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
			if debug {
				logging.PreSetupLog("Debug mode enabled", slog.LevelDebug)
			}

			postCommandOutput, _ := cmd.Flags().GetString(consts.PostCommandOutputFlag)
			if postCommandOutput != "" {
				postcommand.SetOutputFile(postCommandOutput)
			}

			ctxData := context.CommandContextData{
				RootConfig: c.configRoot,
				Debug:      debug,
			}
			ctx := context.GetCommandContext(cmd.Context(), ctxData)
			cmd.SetContext(ctx)
			logging.Setup(debug, os.Stdout)
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return postcommand.WriteOutput()
		},
	}

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().String(consts.PostCommandOutputFlag, "", "Output file for post command")

	for _, command := range c.commands {
		err := command.Setup(c.configRoot)
		if err == nil {
			cmd := command.GetCobraCommand()
			rootCmd.AddCommand(cmd)
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

	return rootCmd
}

func (c *CommandManager) Run(cmd *cobra.Command, args []string) error {
	var commandNames []string

	for name := range c.commands {
		if c.commands[name].GetCobraCommand() != nil {
			commandNames = append(commandNames, name)
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
