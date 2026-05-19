package build

import (
	"fmt"
	"os"
	"strings"
)

const (
	AppName        = "gs-dev"
	AppDescription = `Guiosoft Development Assistant is a tool for helping
the developer in your tasks.`
	ShortDescription = "Guiosoft Development Assistant"

	DevVersion = "develop"
	unknown    = "unknown"

	Owner = "guionardo"
	Repo  = "gs-dev"
)

// Vars below are defined in build time by the .github/scripts/build.sh script

var (
	Version        = DevVersion
	BuildInfo      = unknown
	ExecutableName = ""
	IsValidBinary  = true
)

func init() {
	var err error

	ExecutableName, err = os.Executable()
	if err != nil {
		panic(fmt.Errorf("error getting executable name: %w", err))
	}

	if strings.Contains(ExecutableName, "__debug_bin") {
		// Running from vscode
		BuildInfo = "running from vscode"
		IsValidBinary = false
	} else if strings.HasPrefix(ExecutableName, os.TempDir()) {
		// Running from temp directory like go run
		BuildInfo = "running from temp directory"
		IsValidBinary = false
	}
}
