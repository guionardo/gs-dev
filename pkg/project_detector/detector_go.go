package projectdetector

import (
	"bufio"
	"os"
	"path"
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

// Go detects if the given folder is a Go project by looking for go.mod or main.go files.
func Go(folder string) (project *ProjectData, err error) {
	projectFile := files.FindFirst(folder, "go.mod")
	if projectFile == "" {
		return nil, os.ErrNotExist
	}

	projectName := getGoProjectName(projectFile)

	return &ProjectData{
		Folder: folder,
		Type:   GoProjectType,
		Name:   projectName,
	}, nil
}

func getGoProjectName(projectFile string) string {
	if path.Base(projectFile) == "go.mod" {
		if name := readGoProjectNameFromGoMod(projectFile); name != "" {
			return name
		}
	}

	return ""
}

func readGoProjectNameFromGoMod(projectFile string) string {
	file, err := os.Open(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}
	defer file.Close() // nolint errcheck -- file is closed in the defer function.

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if after, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(after)
		}
	}

	return ""
}
