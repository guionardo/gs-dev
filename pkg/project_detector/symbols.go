package projectdetector

import "github.com/guionardo/gs-dev/pkg/console"

// symbols and color styles inspired by https://starship.rs/config

const (
	UnknownProjectType = "unknown"
	GoProjectType      = "go"
	PythonProjectType  = "python"
	JSProjectType      = "nodejs"
	RustProjectType    = "rust"
	JavaProjectType    = "java"
	DotnetProjectType  = "dotnet"

	OSLinux   = "os_linux"
	OSMacOS   = "os_macos"
	OSWindows = "os_windows"
	OSUnknown = "os_unknown"
)

var symbols = map[string]string{
	UnknownProjectType: " ",
	DotnetProjectType:  " ",
	GoProjectType:      " ",
	PythonProjectType:  " ",
	JSProjectType:      " ",
	JavaProjectType:    " ",
	RustProjectType:    "󱘗 ",

	OSLinux:   " ",
	OSMacOS:   " ",
	OSWindows: "󰍲 ",
	OSUnknown: " ",

	// "dockerfile":   " ",
	// "makefile":     " ",
	// "gradle":       " ",
	// "pom.xml":      " ",
	// "package.json": "󰏗 ",
	// "aws":          " ",
	// "buf":          " ",
	// "bun":          " ",
	// "c":            " ",
	// "cpp":          " ",
	// "cmake":        " ",
	// "conda":        " ",
	// "crystal":      " ",
	// "dart":         " ",
	// "deno":         " ",
	// "docker":       " ",
	// "elixir":       " ",
	// "elm":          " ",
	// "fennel":       " ",
	// "fortran":      " ",
	// "gcloud":       " ",
	// "git_branch":   " ",
	// "git_commit":   " ",
	// "guix_shell":   " ",
	// "haskell":      " ",
	// "haxe":         " ",

	// "julia":        " ",
	// "kotlin":       " ",
	// "lua":          " ",
	// "hostname":     " ",
	// "memory_usage": "󰍛 ",
	// "meson":        "󰔷 ",
	// "nim":          "󰆥 ",
	// "nix_shell":    " ",
	// "ocaml":        " ",

	// "perl":    " ",
	// "php":     " ",
	// "package": "󰏗 ",
	// "rlang":   "󰟔 ",
	// "ruby":    " ",
	// "scala":   " ",
}

var styles = map[string]string{
	GoProjectType:     console.Cyan,
	PythonProjectType: console.Yellow,
	JSProjectType:     console.Green,
	RustProjectType:   console.Red,
	JavaProjectType:   console.Red,
	DotnetProjectType: console.Blue,
}
