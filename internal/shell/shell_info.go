package shell

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

type ShellInfo struct {
	Name   string
	RCFile string
}

func NewShellInfo() (si *ShellInfo, err error) {
	si = new(ShellInfo)
	// Detect shell
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		err = errors.New("no SHELL environment detected")
		return
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("error getting user home dir - %w", err)
		return
	}

	for _, sh := range []string{
		"bash",
		"zsh",
		"ksh",
	} {
		if strings.HasSuffix(shell, sh) {
			si.Name = sh
			si.RCFile = path.Join(homePath, "."+sh+"rc")

			break
		}
	}

	if len(si.RCFile) == 0 {
		err = fmt.Errorf("unexpected shell - %s", shell)
	}

	return
}
