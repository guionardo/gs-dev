package projectdetector

import (
	"encoding/json"
	"os"
	"path"
)

var JavaScript = NewSimpleDetector(JSProjectType, []string{"package.json"}, getJSProjectName)

func getJSProjectName(projectFile string) string {
	content, err := os.ReadFile(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	var (
		projectName string
		packageJSON map[string]any
	)
	if err := json.Unmarshal(content, &packageJSON); err == nil {
		if pn, ok := packageJSON["name"]; ok {
			projectName = pn.(string)
			return projectName
		}
	}

	return path.Base(path.Dir(projectFile))
}
