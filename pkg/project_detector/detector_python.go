package projectdetector

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type (
	PyProjectFile struct {
		Project struct {
			Name string
		}
	}
)

var Python = NewSimpleDetector(PythonProjectType, []string{"pyproject.toml", "*.py"}, getPythonProjectName)

// getPythonProjectName gets the project name from the project pyproject.toml or requirements.txt file
// the project name is the name of the project in the pyproject.toml or basename of the path file
func getPythonProjectName(projectFile string) string {
	switch path.Base(projectFile) {
	case "pyproject.toml":
		if name := readPythonProjectNameFromPyProject(projectFile); name != "" {
			return name
		}

	default:
		return ""
	}

	return path.Base(path.Dir(projectFile))
}

func readPythonProjectNameFromPyProject(projectFile string) string {
	content, err := os.ReadFile(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	var pyProjectFile PyProjectFile
	if err := toml.Unmarshal(content, &pyProjectFile); err == nil {
		return pyProjectFile.Project.Name
	}

	return path.Base(path.Dir(projectFile))
}
