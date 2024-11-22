package build

import (
	"runtime/debug"
)

const AppName = "gs-dev"
const AppDescription = "Go development tools"

const DevVersion = "develop"
const unknown = "unknown"

var Version = DevVersion
var BuildInfo = unknown

func init() {
	if Version == DevVersion {
		return
	}
	// https://jerrynsh.com/3-easy-ways-to-add-version-flag-in-go
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		Version = DevVersion
		return
	}

	if buildInfo.Main.Version != "" {
		Version = buildInfo.Main.Version
		return
	}
}
