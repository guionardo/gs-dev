package files

import (
	"iter"
	"os"
	"path"
)

var ignoreDirs = map[string]struct{}{
	".git": {}, "node_modules": {}, "vendor": {}, ".venv": {}, "__pycache__": {},
	"dist": {}, "build": {}, "target": {}, "bin": {}, "obj": {}, ".tox": {},
	".idea": {}, ".vscode": {}, "coverage": {},
}

// ReadDirectory reads the directory at the given path and yields the paths of all subdirectories up to the specified max depth.
// It returns a generator with only the directory paths limited to maxDept
func ReadDirectory(basePath string, maxDept int) iter.Seq[string] {
	return readDirectory(basePath, 0, maxDept)
}

func readDirectory(basePath string, level, maxDepth int) iter.Seq[string] {
	return func(yield func(string) bool) {
		if level >= maxDepth {
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

			if _, ok := ignoreDirs[entry.Name()]; ok {
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
