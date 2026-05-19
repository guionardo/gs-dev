package projectdetector

import (
	"errors"
	"os"
	"path"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type DetectorFunc func(folder string) (projectData *ProjectData, err error)

var (
	detectors  = []DetectorFunc{Go, Python, Dotnet, JavaScript, Rust, Java, PHP}
	ignoreDirs = map[string]struct{}{
		".git": {}, "node_modules": {}, "vendor": {}, ".venv": {}, "__pycache__": {},
		"dist": {}, "build": {}, "target": {}, "bin": {}, "obj": {}, ".tox": {},
		".idea": {}, ".vscode": {}, "coverage": {},
	}
)

func DetectProject(folder string) (*ProjectData, error) {
	assertedFolder, err := files.AssertDirectory(folder)
	if err != nil {
		return nil, err
	}

	// Check if folder has files
	entries, _ := os.ReadDir(assertedFolder)

	atLeastOneFile := false

	for _, entry := range entries {
		if _, ok := ignoreDirs[entry.Name()]; ok {
			continue
		}

		if !entry.IsDir() {
			atLeastOneFile = true
			break
		}
	}

	if !atLeastOneFile {
		return nil, errors.New("folder is empty: " + assertedFolder)
	}

	for _, detector := range detectors {
		if projectData, err := detector(assertedFolder); err == nil {
			return projectData, nil
		}
	}

	projectData := &ProjectData{
		Folder: assertedFolder,
		Type:   UnknownProjectType,
		Name:   path.Base(assertedFolder),
	}

	return projectData, errors.New("could not detect project type for folder: " + assertedFolder)
}
