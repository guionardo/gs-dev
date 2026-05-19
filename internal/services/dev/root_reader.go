package devservice

import (
	"sort"

	projectdetect "github.com/guionardo/gs-dev/pkg/project_detector"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type RootReader struct {
	directory       string
	root            *Root
	projects        map[string]*projectdetect.ProjectData
	readenProjects  map[string]*projectdetect.ProjectData
	removedProjects []string
}

func NewRootReader(directory string, root *Root) *RootReader {
	removedProjects := make([]string, 0)
	projects := make(map[string]*projectdetect.ProjectData)

	for _, projectFolder := range root.Folders {
		project, err := projectdetect.DetectProject(projectFolder)
		if err != nil {
			removedProjects = append(removedProjects, projectFolder)
			continue
		}

		projects[projectFolder] = project
	}

	// Existing projects
	readenProjects := make(map[string]*projectdetect.ProjectData)

	for dir := range files.ReadDirectory(directory, root.MaxDepth) {
		if project, err := projectdetect.DetectProject(dir); err == nil {
			readenProjects[project.Folder] = project
		}
	}

	return &RootReader{directory: directory,
		root:            root,
		projects:        projects,
		readenProjects:  readenProjects,
		removedProjects: removedProjects}
}

func (r *RootReader) UpdateRoot(root *Root) {
	readenFolders := make([]string, 0, len(r.readenProjects))
	for project := range r.readenProjects {
		readenFolders = append(readenFolders, project)
	}

	sort.Strings(readenFolders)
	root.Folders = readenFolders
}

// SyncSummary returns the summary of the projects that have been removed and added
// canIncludeFunc is a function that returns true if the folder should be included in the summary
func (r *RootReader) SyncSummary(canIncludeFunc func(folder string) bool) (removed []*projectdetect.ProjectData, added []*projectdetect.ProjectData) {
	for folder := range r.readenProjects {
		if !canIncludeFunc(folder) {
			continue
		}

		if _, ok := r.projects[folder]; !ok {
			added = append(added, r.readenProjects[folder])
		}
	}

	for folder := range r.projects {
		if _, ok := r.readenProjects[folder]; !ok || !canIncludeFunc(folder) {
			removed = append(removed, r.projects[folder])
		}
	}

	return removed, added
}
