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
func ReadDirectory(basePath string, maxDept int) iter.Seq[string] {
	return readDirectory2(basePath, 0, maxDept)
}

func readDirectory(basePath string, level, maxDepth int) iter.Seq[string] {
	return func(yield func(string) bool) {
		if level >= maxDepth {
			return
		}
		// Check if the basepath has one of ignoreDirsSet to avoid sub reading
		if directoryHasIgnoredSubFolder(basePath) {
			return
		}

		entries, err := os.ReadDir(basePath)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			dirPath := path.Join(basePath, entry.Name())
			if !yield(dirPath) {
				return
			}

			if level < maxDepth-1 {
				for dir := range readDirectory(dirPath, level+1, maxDepth) {
					if !yield(dir) {
						return
					}
				}
			}
		}
	}
}

func readDirectory2(basePath string, level, maxDepth int) iter.Seq[string] {
	return func(yield func(string) bool) {
		if level >= maxDepth {
			return
		}

		type frame struct {
			path  string
			level int
			dirs  []string
			index int
		}

		rootDirs := getSubDirs(basePath)
		stack := []frame{{basePath, level, rootDirs, 0}}

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
