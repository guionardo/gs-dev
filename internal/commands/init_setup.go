package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/pkg/tools/files"
	"github.com/spf13/cobra"
)

type InitSetup struct {
	toolName string
	commands map[string]Command // command name -> command arguments
}

const (
	flagOutput  = "post-command-output"
	commandInit = "init2"
)

var (
	postCommandOutputFile string
	initSetupInstance     *InitSetup
)

func NewInitSetup(toolName string, commands ...Command) *InitSetup {
	if initSetupInstance == nil {
		initSetupInstance = &InitSetup{
			toolName: toolName,
			commands: make(map[string]Command),
		}
		for _, command := range commands {
			initSetupInstance.commands[command.GetName()] = command
		}
	}

	return initSetupInstance
}

// GenerateInitSetup generates the init setup for the commands
// TODO: move to the service
func (i *InitSetup) GenerateInitSetup(toolBinaryPath string) ([]byte, error) {
	if len(i.commands) == 0 {
		return nil, errors.New("no commands to generate init setup")
	}

	commandNames := make([]string, 0)
	for _, command := range i.commands {
		commandNames = append(commandNames, command.GetName())
	}

	content := fmt.Appendf([]byte{}, "#!/usr/bin/env bash\n")

	treatOutputFunctionName := strings.ReplaceAll(fmt.Sprintf("__%s_treat_output", i.toolName), "-", "_")

	// create output file name
	mktemp, err := files.LocateBinary("mktemp")

	var tmpFileAttr string

	if err == nil {
		tmpFileAttr = fmt.Sprintf("$(%s -u)", mktemp)
	} else {
		tmpFileAttr = path.Join(os.TempDir(), i.toolName+"."+uuid.New().String())
	}

	// write treat output function
	content = fmt.Appendf(content, "%s() {\n", treatOutputFunctionName)
	content = fmt.Appendf(content, "  tmp_file=$1\n")
	content = fmt.Appendf(content, "  if [[ ! -f $tmp_file ]]; then\n")
	content = fmt.Appendf(content, "    echo \"Missing output file: $tmp_file\"\n")
	content = fmt.Appendf(content, "    return\n")
	content = fmt.Appendf(content, "  fi\n")
	content = fmt.Appendf(content, "  source $tmp_file\n")
	content = fmt.Appendf(content, "  rm -f $tmp_file\n")
	content = fmt.Appendf(content, "  stty sane\n")
	content = fmt.Appendf(content, "}\n")

	aliasCommands := make([]string, 0)
	// write command functions
	for _, command := range i.commands {
		if !command.HasInitAlias() {
			continue
		}

		aliasCommands = append(aliasCommands, command.GetName())

		content = fmt.Appendf(content, "%s() {\n", command.GetName())
		if command.UsesOutput() {
			content = fmt.Appendf(content, `  tmp_file="%s"  
  %s --%s $tmp_file %s %s $@ && %s $tmp_file
  `, tmpFileAttr, toolBinaryPath, flagOutput, command.GetName(), command.GetArguments(), treatOutputFunctionName)
		} else {
			content = fmt.Appendf(content, "  %s %s $@ \n", toolBinaryPath, command.GetName())
		}

		content = fmt.Appendf(content, "}\n")
	}

	// write command aliases
	content = fmt.Appendf(content, "echo '%s is ready to use (%s)'\n", i.toolName, strings.Join(aliasCommands, ", "))

	return content, nil
}

func (i *InitSetup) UpdateRootCommand(rootCmd *cobra.Command) error {
	// Check if flagOutput is already set
	if rootCmd.PersistentFlags().Lookup(flagOutput) == nil {
		rootCmd.PersistentFlags().StringVar(&postCommandOutputFile, flagOutput, "", "Output file for post command")
	}

	for _, command := range rootCmd.Commands() {
		if command.Name() == commandInit {
			return fmt.Errorf("command %s already exists", commandInit)
		}
	}

	initCmd := &cobra.Command{
		Use:   commandInit,
		Short: "Initialization for shell alias",
		Long: fmt.Sprintf(`Add to your profile script (.bashrc, etc)

		source <(%s %s)`, i.toolName, commandInit),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := i.GenerateInitSetup(build.ExecutableName)
			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(content)

			return err
		},
	}

	rootCmd.AddCommand(initCmd)

	return nil
}
