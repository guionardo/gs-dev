package projectdetector

import (
	"os"
	"path"
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

func Dotnet(folder string) (project *ProjectData, err error) {
	solutionFile := files.FindFirst(folder, "*.sln", "*.slnx")

	var projectName string
	if solutionFile != "" {
		projectName = getFirstDotnetProjectName(solutionFile)
	}

	if projectName == "" {
		projectFile := files.FindFirstInSubdirs(folder, 1, "*.csproj")
		if projectFile == "" {
			return nil, os.ErrNotExist
		}

		projectName = strings.TrimSuffix(path.Base(projectFile), ".csproj")
	}

	return &ProjectData{
		Folder: folder,
		Type:   DotnetProjectType,
		Name:   projectName,
	}, nil
}

// getFirstDotnetProjectName gets the first project file name from the solution file
func getFirstDotnetProjectName(solutionFile string) string {
	if strings.HasSuffix(solutionFile, ".sln") {
		return getFirstDotnetProjectNameFromSln(solutionFile)
	}

	if strings.HasSuffix(solutionFile, ".slnx") {
		return getFirstDotnetProjectNameFromSlnx(solutionFile)
	}

	return ""
}

// getFirstDotnetProjectNameFromSln gets the first project file name from the solution file
// sln files have a format where the project file is specified in a line like this:
// Project("{GUID}") = "ProjectName", "ProjectFile.csproj", "{GUID}"
func getFirstDotnetProjectNameFromSln(solutionFile string) string {
	content, err := os.ReadFile(solutionFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.HasPrefix(line, "Project(") {
			parts := strings.Split(line, "=")
			if len(parts) > 1 {
				projectPart := strings.TrimSpace(parts[1])
				projectParts := strings.Split(projectPart, ",")

				return strings.ReplaceAll(strings.TrimSpace(projectParts[0]), "\"", "")
			}
		}
	}

	return ""
}

// getFirstDotnetProjectNameFromSlnx gets the first project file name from the solution file
// slnx files have a format where the project file is specified in a line like this:
// <Project Path="ProjectFile.csproj" />
func getFirstDotnetProjectNameFromSlnx(solutionFile string) string {
	content, err := os.ReadFile(solutionFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	lines := strings.SplitSeq(string(content), "\n")
	for line := range lines {
		if strings.HasPrefix(line, "<Project Path=") {
			parts := strings.Split(line, "\"")
			if len(parts) > 1 {
				pn, _ := strings.CutSuffix(parts[1], ".csproj")
				return pn
			}
		}
	}

	return ""
}
