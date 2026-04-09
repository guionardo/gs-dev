package installservice

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/interfaces"
	"github.com/guionardo/gs-dev/internal/shell"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type (
	InstallService struct {
		sourceCommand string
	}
)

var commandsMap map[string]interfaces.Command

func NewInstallService(configuration *config.ConfigFile) *InstallService {
	sourceCommand := fmt.Sprintf("source <(%s init)", build.ExecutableName)

	return &InstallService{
		sourceCommand: sourceCommand,
	}
}

func SetCommands(registeredCommands ...interfaces.Command) {
	commandsMap = make(map[string]interfaces.Command, len(registeredCommands))
	for _, command := range registeredCommands {
		commandsMap[command.GetName()] = command
	}
}

func (i *InstallService) IsInstalled() (msg string, isInstalled bool) {
	profile, err := shell.NewProfileFile(build.AppName)
	if err != nil {
		return fmt.Sprintf("error getting profile file: %v", err), false
	}

	if line, ok := profile.HasEnabledCommandLine(); ok {
		return fmt.Sprintf("binding was just installed into shell profile %s at line %d", profile.Path, line), true
	}

	return "binding was not installed into shell profile " + profile.Path, false
}

func (i *InstallService) Install() error {
	// source <(./gs-dev init)
	if strings.Contains(i.sourceCommand, "__debug_bin") {
		// Running from vscode
		return fmt.Errorf("bad executable name %s", i.sourceCommand)
	}

	profile, err := shell.NewProfileFile(build.AppName)
	if err != nil {
		return err
	}

	if line, ok := profile.HasEnabledCommandLine(); ok {
		return fmt.Errorf("binding was just installed into shell profile %s at line %d",
			profile.Path, line)
	}

	profile.SetFeature(i.sourceCommand, true)

	if err := profile.Save(); err != nil {
		return fmt.Errorf("error saving profile %s - %w", profile.Path, err)
	}

	fmt.Printf("binding installed at file %s\nbackup done at %s", profile.Path, profile.LastBackup)

	return nil
}

func (i *InstallService) Uninstall() error {
	profile, err := shell.NewProfileFile(build.AppName)
	if err != nil {
		return err
	}

	if _, ok := profile.HasEnabledCommandLine(); !ok {
		return fmt.Errorf("binding was not installed into shell profile %s",
			profile.Path)
	}

	profile.SetFeature(i.sourceCommand, false)

	if err := profile.Save(); err != nil {
		return fmt.Errorf("error saving profile %s - %w", profile.Path, err)
	}

	fmt.Printf("binding uninstalled at file %s\nbackup done at %s", profile.Path, profile.LastBackup)

	return nil
}

func (i *InstallService) GetSourceCommand() string {
	return i.sourceCommand
}

func (i *InstallService) GetInitCommand() string {
	initCommand, err := i.GenerateInitSetup(build.ExecutableName, commandsMap)
	if err != nil {
		return ""
	}

	return string(initCommand)
}

// GenerateInitSetup generates the init setup for the commands
func (i *InstallService) GenerateInitSetup(toolBinaryPath string, commands map[string]interfaces.Command) ([]byte, error) {
	if len(commands) == 0 {
		return nil, errors.New("no commands to generate init setup")
	}

	content := fmt.Appendf([]byte{}, "#!/usr/bin/env bash\n")

	treatOutputFunctionName := strings.ReplaceAll(fmt.Sprintf("__%s_treat_output", build.AppName), "-", "_")

	// create output file name
	mktemp, err := files.LocateBinary("mktemp")

	var tmpFileAttr string

	if err == nil {
		tmpFileAttr = fmt.Sprintf("$(%s -u)", mktemp)
	} else {
		tmpFileAttr = path.Join(os.TempDir(), build.AppName+"."+uuid.New().String())
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
	for _, command := range commands {
		if !command.HasInitAlias() {
			continue
		}

		aliasCommands = append(aliasCommands, command.GetName())

		content = fmt.Appendf(content, "%s() {\n", command.GetName())
		if command.UsesOutput() {
			content = fmt.Appendf(content, `  tmp_file="%s"  
  %s --%s $tmp_file %s %s $@ && %s $tmp_file
  `, tmpFileAttr, toolBinaryPath, consts.PostCommandOutputFlag, command.GetName(), command.GetArguments(), treatOutputFunctionName)
		} else {
			content = fmt.Appendf(content, "  %s %s $@ \n", toolBinaryPath, command.GetName())
		}

		content = fmt.Appendf(content, "}\n")
	}

	// write command aliases
	if len(aliasCommands) > 0 {
		content = fmt.Appendf(content, "echo '%s is ready to use (%s)'\n", build.AppName, strings.Join(aliasCommands, ", "))
	}

	return content, nil
}
