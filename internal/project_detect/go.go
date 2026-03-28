package projectdetect

import (
	"bufio"
	"os"
	"path"
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

func GoDetector(folder string) (projectType string, projectName string, err error) {
	projectFile := files.FindFirst(folder, "go.mod", "*.go")
	if projectFile == "" {
		return "", "", os.ErrNotExist
	}

	projectName = getGoProjectName(projectFile)

	return "go", projectName, nil
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
