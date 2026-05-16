package fs_tools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/internal/consts"
	errs "github.com/guionardo/gs-dev/internal/errors"
)

var homeDir string

func init() {
	homeDir, _ = os.UserHomeDir()
}

// AssertDirectory asserts that the path is a directory and returns the absolute path
func AssertDirectory(path string) (string, error) {
	for _, homeTmp := range []string{"~/", "$HOME/"} {
		if after, ok := strings.CutPrefix(path, homeTmp); ok {
			path = filepath.Join(homeDir, after)
		}
	}

	if !pathtools.DirExists(path) {
		return "", fmt.Errorf("path %s is not a directory", path)
	}

	return filepath.Abs(path)
}

// AssertFilename asserts that the filename is a file and returns the absolute path
func AssertFilename(filename string) (string, error) {
	if filename == "" {
		return "", errs.NewError(errors.New("filename is required"), "filename is required", false)
	}

	filename, err := filepath.Abs(filename)
	if err != nil {
		return "", errs.NewError(err, "error getting absolute filename", false)
	}

	if stat, err := os.Stat(filepath.Clean(filename)); err == nil {
		if !stat.IsDir() {
			// file exists and is not a directory
			return filename, nil
		}

		return filename, errs.NewError(fmt.Errorf("file %s is a directory", filename), "file is a directory", false)
	}
	// file does not exist, try to create it directory
	if err := os.MkdirAll(filepath.Dir(filename), consts.DirPermissions); err != nil {
		return filename, errs.NewError(err, "error creating directory", false)
	}

	// try to create the file
	if err := os.WriteFile(filename, []byte{}, consts.FilesPermissions); err != nil {
		return filename, errs.NewError(err, "error creating file", false)
	}

	_ = os.Remove(filename)

	return filename, nil
}

// PathIsRoot checks if the path is the root of the filesystem
func PathIsRoot(path string) bool {
	// Clean the path to handle trailing slashes and relative components
	cleaned := filepath.Clean(path)
	// A path is the root if its parent is itself
	return cleaned == filepath.Dir(cleaned)
}
