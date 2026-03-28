package projectdetect

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type (
	CargoTomlFile struct {
		Package struct {
			Name string
		}
	}
)

func RustDetector(folder string) (projectType string, projectName string, err error) {
	projectFile := files.FindFirst(folder, "Cargo.toml", "*.rs")
	if projectFile == "" {
		return "", "", os.ErrNotExist
	}

	projectName = getRustProjectName(projectFile)

	return "rust", projectName, nil
}

func getRustProjectName(projectFile string) string {
	switch path.Base(projectFile) {
	case "Cargo.toml":
		if name := readRustProjectNameFromCargoToml(projectFile); name != "" {
			return name
		}
	default:
		return ""
	}

	return ""
}

func readRustProjectNameFromCargoToml(projectFile string) string {
	content, err := os.ReadFile(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	var cargoToml CargoTomlFile
	if err := toml.Unmarshal(content, &cargoToml); err == nil {
		return cargoToml.Package.Name
	}

	return path.Base(path.Dir(projectFile))
}
