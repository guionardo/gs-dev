package environment

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/configs"
	"github.com/kardianos/osext"
)

//go:embed env_shell.sh
var env_shell string
var homeDir string

func init() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		homeDir = "~/"
	}
}

func SetupTerminal(cfg *configs.RootConfig) (err error) {
	var shellFile string
	if shellFile, err = setupShellFile(cfg); err != nil {
		return
	}
	profileFile := getProfileFile()
	if err = setupProfileFile(shellFile, profileFile); err != nil {
		return
	}
	return
}

func setupShellFile(cfg *configs.RootConfig) (shellFile string, err error) {
	shellFileName := path.Join(cfg.DataFolder, "env_shell.sh")
	executable, err := osext.Executable()
	if err != nil {
		return
	}

	content := env_shell
	replaces := map[string]string{
		"GSDEV_BIN":      executable,
		"ENV_SHELL_FILE": shellFileName,
		"DESCRIPTION":    fmt.Sprintf("script created by %s v%s", app.ToolName, app.Version),
	}
	for key, value := range replaces {
		content = strings.ReplaceAll(content, fmt.Sprintf("_%s_", key), value)
	}

	err = os.WriteFile(shellFile, []byte(content), 0744)
	return
}

func setupProfileFile(shellFileName string, profileFile string) (err error) {
	if _, err = os.Stat(profileFile); err != nil {
		profileHeader := fmt.Sprintf("# File created by %s v%s @ %v\n\n",
			app.ToolName, app.Version, time.Now())
		err = os.WriteFile(profileFile, []byte(profileHeader), 0744)
		if err != nil {
			return
		}
	}

	setupLine, err := getSetupLine(profileFile, shellFileName)
	newSetupLine := fmt.Sprintf("source %s # Setup by %s v%s @ %v\n", shellFileName, app.ToolName, app.Version, time.Now())
	var content []byte
	if err == nil {
		content, err = os.ReadFile(profileFile)
		if err != nil {
			return
		}
		stringContent := string(content)
		stringContent = strings.Replace(stringContent, setupLine, newSetupLine, -1)
		err = os.WriteFile(profileFile, []byte(stringContent), 0744)
	} else {
		var file *os.File
		if file, err = os.OpenFile(profileFile, os.O_APPEND, 0744); err == nil {
			_, err = file.WriteString("\n" + newSetupLine)
		}
	}
	return

}

func getProfileFile() (profileFile string) {
	shell := os.Getenv("SHELL")
	switch {
	case strings.HasSuffix(shell, "bash"):
		return path.Join(homeDir, ".bashrc")
	case strings.HasSuffix(shell, "zsh"):
		return path.Join(homeDir, ".zshrc")
	case strings.HasSuffix(shell, "ksh"):
		return path.Join(homeDir, ".kshrc")
	default:
		return path.Join(homeDir, ".profile")
	}
}

func getSetupLine(profileFile string, shellFileName string) (setupLine string, err error) {
	var content []byte
	if content, err = os.ReadFile(profileFile); err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, shellFileName) {
			return line, nil
		}
	}
	return "", fmt.Errorf("setup line not found - %s", shellFileName)

}
