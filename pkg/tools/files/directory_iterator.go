package files

import (
	"iter"
	"os"
	"path"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/go/set"
)

var (
	ignoreDirsSet = set.New(".git", "node_modules", "vendor", ".venv", "__pycache__", "dist", "build", "target", "bin", "obj", ".tox", ".idea", ".vscode", "coverage")
)

// ReadDirectory reads the directory at the given path and yields the paths of all subdirectories up to the specified max depth.
// It returns a generator with only the directory paths limited to maxDept
func ReadDirectory(basePath string, maxDepth int) iter.Seq[string] {
	return func(yield func(string) bool) {
		type frame struct {
			path  string
			level int
			dirs  []string
			index int
		}

		rootDirs := getSubDirs(basePath)
		stack := []frame{{basePath, 0, rootDirs, 0}}

		for len(stack) > 0 {
			top := &stack[len(stack)-1]

			if top.index >= len(top.dirs) {
				stack = stack[:len(stack)-1]
				continue
			}

			childPath := top.dirs[top.index]
			top.index++

			if !yield(childPath) {
				return
			}

			if top.level < maxDepth-1 {
				childDirs := getSubDirs(childPath)
				if len(childDirs) > 0 {
					stack = append(stack, frame{childPath, top.level + 1, childDirs, 0})
				}
			}
		}
	}
}

func directoryHasIgnoredSubFolder(directory string) bool {
	for ignored := range ignoreDirsSet.Iter() {
		if pathtools.DirExists(path.Join(directory, ignored)) {
			return true
		}
	}

	return false
}

// getSubDirs get the sub directories of the directory if it has no ignored entries
func getSubDirs(directory string) (subDirectories []string) {
	if directoryHasIgnoredSubFolder(directory) {
		return nil
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subDirectories = append(subDirectories, path.Join(directory, entry.Name()))
		}
	}

	return
}
