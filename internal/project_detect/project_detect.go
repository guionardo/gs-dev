package projectdetect

import (
	"fmt"
	"os"
	"path"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/internal/fs_tools"
)

type (
	ProjectData struct {
		Folder          string
		IsGitRepository bool

		ProjectType       string
		ProjectName       string
		GitStatus         string // TODO: add git status, clean, dirty, etc.
		HasAtLeastOneFile bool
	}
	DetectorFunc func(folder string) (projectType string, projectName string, err error)
)

const (
	UNKNOWN = "unknown"
)

var detectors = []DetectorFunc{
	GoDetector,
	PythonDetector,
	JSDetector,
	RustDetector,
	JavaDetector,
}

func (pd ProjectData) String() string {
	return fmt.Sprintf("%s (%s) - %s", pd.ProjectName, pd.ProjectType, pd.Folder)
}

func DetectProject(folder string) (*ProjectData, error) {
	folder, err := fs_tools.AssertDirectory(folder)
	if err != nil {
		return nil, err
	}

	isGitRepository := isGitRepository(folder)
	projectType, projectName := getProjectType(folder)

	hasAtLeastOneFile := false

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			hasAtLeastOneFile = true
			break
		}
	}

	projectData := ProjectData{
		Folder:            folder,
		IsGitRepository:   isGitRepository,
		ProjectType:       projectType,
		ProjectName:       projectName,
		HasAtLeastOneFile: hasAtLeastOneFile,
	}

	return &projectData, nil
}

func isGitRepository(folder string) bool {
	return pathtools.DirExists(path.Join(folder, ".git"))
}

func getProjectType(folder string) (projectType string, projectName string) {
	var err error
	for _, detector := range detectors {
		if projectType, projectName, err = detector(folder); err == nil {
			break
		}
	}

	if projectName == "" {
		projectName = path.Base(folder)
	}

	if projectType == "" {
		projectType = UNKNOWN
	}

	return projectType, projectName
}
