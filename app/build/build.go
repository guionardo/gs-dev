package build

const (
	AppName        = "gs-dev"
	AppDescription = `Guiosoft Development Assistant is a tool for helping
the developer in your tasks.`
	ShortDescription = "Guiosoft Development Assistant"

	DevVersion = "develop"
	unknown    = "unknown"
)

// Vars below are defined in build time by the .github/scripts/build.sh script

var (
	Version   = DevVersion
	BuildInfo = unknown
)
