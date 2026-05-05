package projectdetector

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type CargoTomlFile struct {
	Package struct {
		Name string
	}
}

func Rust(folder string) (project *ProjectData, err error) {
	projectFile := files.FindFirst(folder, "Cargo.toml")
	if projectFile == "" {
		return nil, os.ErrNotExist
	}

	projectName := getRustProjectName(projectFile)

	return &ProjectData{
		Folder: folder,
		Type:   RustProjectType,
		Name:   projectName,
	}, nil
}

func getRustProjectName(projectFile string) string {
	if name := readRustProjectNameFromCargoToml(projectFile); name != "" {
		return name
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
