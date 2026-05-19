package projectdetector

import (
	"encoding/xml"
	"os"
	"path"

	"github.com/guionardo/gs-dev/pkg/tools/files"
)

func Java(folder string) (project *ProjectData, err error) {
	projectFile := files.FindFirst(folder, "pom.xml", "build.gradle", "*.java")
	if projectFile == "" {
		return nil, os.ErrNotExist
	}

	projectName := getJavaProjectName(projectFile)

	return &ProjectData{
		Folder: folder,
		Type:   JavaProjectType,
		Name:   projectName,
	}, nil
}

func getJavaProjectName(projectFile string) string {
	switch path.Base(projectFile) {
	case "pom.xml":
		if name := readJavaProjectNameFromPomXML(projectFile); name != "" {
			return name
		}
	case "build.gradle":
		if name := readJavaProjectNameFromBuildGradle(projectFile); name != "" {
			return name
		}
	}

	return ""
}

func readJavaProjectNameFromPomXML(projectFile string) string {
	content, err := os.ReadFile(projectFile) // #nosec G304 -- path comes from local discovery in the inspected folder.
	if err != nil {
		return ""
	}

	type PomXMLFile struct {
		Project xml.Name `xml:"project"`
		Name    string   `xml:"name"`
	}

	var pomXML PomXMLFile
	if err := xml.Unmarshal(content, &pomXML); err == nil {
		return pomXML.Name
	}

	return ""
}

func readJavaProjectNameFromBuildGradle(projectFile string) string {
	// TODO: implement
	return ""
}
