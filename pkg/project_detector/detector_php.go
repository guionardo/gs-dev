package projectdetector

import (
	"encoding/json"
	"os"
	"path"
)

var PHP = NewSimpleDetector("php", []string{"composer.json"}, getPHPProjectName)

func getPHPProjectName(projectFile string) string {
	content, err := os.ReadFile(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	var composerJSON map[string]any

	if err := json.Unmarshal(content, &composerJSON); err == nil {
		if pn, ok := composerJSON["name"]; ok {
			return pn.(string)
		}
	}

	return path.Base(path.Dir(projectFile))
}
