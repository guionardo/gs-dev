package install

import (
	"fmt"
	"os"
	"strings"

	"github.com/guionardo/gs-dev/shell"
)

func getAssets() (sourceCommand string, profile string, profileLines []string, err error) {
	if sourceCommand, err = getSourceCommand(); err != nil {
		return
	}
	if profile, profileLines, err = getProfile(); err != nil {
		return
	}
	return
}
func getSourceCommand() (source string, err error) {
	//	source <(./go-dev init)
	var executableName string
	if executableName, err = os.Executable(); err != nil {
		return
	}

	if strings.Contains(executableName, "__debug_bin") && false {
		// Running from vscode
		err = fmt.Errorf("bad executable name %s", executableName)
	} else {
		source = fmt.Sprintf("source <(%s init)", executableName)
	}
	return
}

func getProfile() (profile string, profileLines []string, err error) {
	var shellInfo *shell.ShellInfo
	if shellInfo, err = shell.NewShellInfo(); err != nil {
		return
	}
	profile = shellInfo.RCFile
	var content []byte
	if content, err = os.ReadFile(profile); err != nil {
		return
	}
	profileLines = strings.Split(string(content), "\n")
	for index, line := range profileLines {
		profileLines[index] = strings.TrimRight(line, " \n\r")
	}
	return
}

func findContent(lines []string, content string) int {
	for index, line := range lines {
		if strings.HasPrefix(line, content) {
			return index + 1
		}
	}
	return 0
}

func RunInstall() error {
	//	source <(./gs-dev init)

	sourceCommand, profile, profileLines, err := getAssets()
	if err != nil {
		return err
	}
	if index := findContent(profileLines, sourceCommand); index > 0 {
		return fmt.Errorf("binding was just installed into shell profile %s at line %d",
			profile, index)
	}

	profileLines = append(profileLines, sourceCommand)

	if err := os.WriteFile(profile, []byte(strings.Join(profileLines, "\n")), 0644); err != nil {
		return err
	}
	fmt.Printf("binding installed at file %s\n", profile)
	return nil

}

func RunUninstall() error {
	sourceCommand, profile, profileLines, err := getAssets()
	if err != nil {
		return err
	}
	if index := findContent(profileLines, sourceCommand); index == 0 {
		return fmt.Errorf("binding was not installed into shell profile %s",
			profile)
	} else {
		index--
		if index < len(profileLines) {
			profileLines = append(profileLines[:index], profileLines[index+1:]...)
		} else {
			profileLines = profileLines[:index]
		}
	}

	if err := os.WriteFile(profile, []byte(strings.Join(profileLines, "\n")), 0644); err != nil {
		return err
	}
	fmt.Printf("binding uninstalled at file %s\n", profile)
	return nil
}
