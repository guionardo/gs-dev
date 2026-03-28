package projectdetect

import (
	"encoding/json"
	"os"
	"path"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

func JSDetector(folder string) (projectType string, projectName string, err error) {
	projectFile := files.FindFirst(folder, "package.json", "*.js")
	if projectFile == "" {
		return "", "", os.ErrNotExist
	}

	projectName = getJSProjectName(projectFile)

	return "js", projectName, nil
}

func getJSProjectName(projectFile string) string {
	switch path.Base(projectFile) {
	case "package.json":
		if name := readJSProjectNameFromPackageJSON(projectFile); name != "" {
			return name
		}
	default:
		return ""
	}

	return path.Base(path.Dir(projectFile))
}

func readJSProjectNameFromPackageJSON(projectFile string) string {
	content, err := os.ReadFile(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	var packageJSON map[string]any
	if err := json.Unmarshal(content, &packageJSON); err == nil {
		return packageJSON["name"].(string)
	}

	return path.Base(path.Dir(projectFile))
}
