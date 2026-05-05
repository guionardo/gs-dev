package files

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	pathtools "github.com/guionardo/go/path_tools"
)

// FindFirst finds the first file in the root directory that matches the names
// and returns the path to the file if found, otherwise returns empty string
func FindFirst(root string, names ...string) string {
	for i := range names {
		matches, err := filepath.Glob(path.Join(root, names[i]))
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// TODO: add a FindFirstInSubdirs that searches for the first file in the root directory and its subdirectories up to a certain depth
func FindFirstInSubdirs(root string, maxDepth int, names ...string) string {
	for i := range names {
		matches, err := filepath.Glob(path.Join(root, strings.Repeat("/*", maxDepth), names[i]))
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// LocateBinary locates a binary in the PATH environment variable
// and returns the path to the binary if found, otherwise returns an error
func LocateBinary(binaryName string) (string, error) {
	for path := range strings.SplitSeq(os.Getenv("PATH"), string(filepath.ListSeparator)) {
		if pathtools.FileExists(filepath.Join(path, binaryName)) {
			return filepath.Join(path, binaryName), nil
		}
	}

	return "", fmt.Errorf("binary %s not found", binaryName)
}

// FindBestBinCandidate finds the best bin directory candidate in the PATH environment variable
func FindBestBinCandidate() (string, error) {
	candidates := []string{
		filepath.Join(os.Getenv("HOME"), "bin"),
		filepath.Join(os.Getenv("HOME"), ".local", "bin"),
	}

	paths := strings.Split(os.Getenv("PATH"), string(filepath.ListSeparator))
	for _, candidate := range candidates {
		if pathtools.DirExists(candidate) && slices.Contains(paths, candidate) {
			return candidate, nil
		}
	}

	return "", errors.New("no bin candidate found")
}
