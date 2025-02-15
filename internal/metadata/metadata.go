package metadata

import (
	"runtime/debug"
)

var Version string

const AppName = "gs-dev"
const AppDescription = "Go development tools"
const ShortDescription = "Guiosoft Development Assistant"
const Description = `Guiosoft Development Assistant is a tool for helping
the developer in your tasks.`

// https://jerrynsh.com/3-easy-ways-to-add-version-flag-in-go

// printVersion prints the application version
func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		Version = "develop"
		return
	}

	if buildInfo.Main.Version != "" {
		Version = buildInfo.Main.Version
		return
	}
	Version = "develop"

}
