package projectdetector

import (
	"os"
	"path"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

// ProjectNamerFn is a function type that takes a project file path and returns the project name parsing its content.
type ProjectNamerFn func(projectFile string) string

// NewSimpleDetector creates a simple project detector that checks for the existence of specific files to determine the project type and name.
// The project name is derived from the folder name, and the project type is provided as an argument.
func NewSimpleDetector(projectType string, existingFiles []string, projectNamerFn ProjectNamerFn) DetectorFunc {
	if projectNamerFn == nil {
		projectNamerFn = func(projectFile string) string {
			return path.Base(path.Dir(projectFile))
		}
	}

	return func(folder string) (projectData *ProjectData, err error) {
		if file := files.FindFirst(folder, existingFiles...); file != "" {
			projectName := projectNamerFn(file)
			if projectName == "" {
				projectName = path.Base(folder)
			}

			return &ProjectData{
				Folder: folder,
				Type:   projectType,
				Name:   projectName,
			}, nil
		}

		return nil, os.ErrNotExist
	}
}
