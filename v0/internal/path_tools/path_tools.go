package pathtools

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func PathExists(path string) bool {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}
	return false
}

func AbsPath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return path
}

func ValidatePath(path string) error {

	if len(path) == 0 {
		return errors.New("path can't be empty")
	}
	if stat, err := os.Stat(path); err == nil {
		if stat.IsDir() {
			return nil
		}
		return fmt.Errorf("path '%s' is a file", path)
	} else if os.IsNotExist(err) {
		return CreatePath(path)
	} else {
		return err
	}

}
