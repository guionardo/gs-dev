package installservice

import (
	"fmt"
	"strings"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/shell"
)

type InstallService struct {
	sourceCommand string
}

func NewInstallService(configuration *config.ConfigFile) *InstallService {
	sourceCommand := fmt.Sprintf("source <(%s init)", build.ExecutableName)

	return &InstallService{
		sourceCommand: sourceCommand,
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
	initSetup := commands.NewInitSetup(build.AppName)

	initCommand, err := initSetup.GenerateInitSetup(build.ExecutableName)
	if err != nil {
		return ""
	}

	return string(initCommand)
}
