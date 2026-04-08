package devservice

import (
	"iter"
	"os"
	"path"
	"sort"

	projectdetect "github.com/guionardo/gs-dev/internal/project_detect"
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
	for project := range readFolder(directory, 0, root.MaxDepth) {
		readenProjects[project.Folder] = project
	}

	// Remove folders that are not projects
	for folder := range readenProjects {
		// Get the parent folder
		parentFolder := path.Dir(folder)
		if parentProject, ok := readenProjects[parentFolder]; ok {
			if parentProject.ProjectType == projectdetect.UNKNOWN {
				delete(readenProjects, parentFolder)
			}
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

func readFolder(folder string, level, maxDepth int) iter.Seq[*projectdetect.ProjectData] {
	return func(yield func(*projectdetect.ProjectData) bool) {
		if level >= maxDepth {
			return
		}

		if level > 0 {
			project, err := projectdetect.DetectProject(folder)
			if err != nil {
				return
			}

			if project.HasAtLeastOneFile {
				yield(project)
				return
			}
		}

		if level == maxDepth {
			return
		}

		entries, err := os.ReadDir(folder)
		if err != nil {
			return
		}

		for _, file := range entries {
			if file.IsDir() {
				for project := range readFolder(path.Join(folder, file.Name()), level+1, maxDepth) {
					if !yield(project) {
						return
					}
				}
			}
		}
	}
}
