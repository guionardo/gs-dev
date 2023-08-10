package pathtools

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/guionardo/go-dev/pkg/logger"
	"github.com/mattn/go-isatty"
	"github.com/schollz/progressbar/v3"
)

func ReadFolders(root string, maxSubLevel int) (subFolders []string, err error) {
	intTerm := IsRunningFromInteractiveTerminal() && !logger.IsDebugMode()
	var bar *progressbar.ProgressBar
	logger.Info("Reading folders from %s (depth=%d)", root, maxSubLevel)
	noIntLog := "Reading"
	startTime := time.Now()
	if intTerm {
		bar = progressbar.Default(-1, noIntLog)
	}
	defer func() {
		if bar != nil {
			bar.Finish()
		}
		logger.Info("%s took %v to get %d folders", noIntLog, time.Since(startTime).String(), len(subFolders))
	}()

	subFolders, err = FolderReaderReadDir(root, maxSubLevel,
		func(name string) bool {
			if intTerm {
				bar.Add(1)
			}
			return !strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "_")
		},
		func(name string) {
			logger.Debug("%s", name)
		})

	return

}

func FolderReaderReadDir(root string, maxDepth int, acceptFolder func(string) bool, notify func(string)) ([]string, error) {
	return readDirs_(root, 1, maxDepth, acceptFolder, notify)
}

func readDirs_(root string, level int, maxDepth int, acceptFolder func(string) bool, notify func(string)) ([]string, error) {
	dirs := make([]string, 0, 1000)
	entries, err := os.ReadDir(root)
	for _, entry := range entries {
		if entry.IsDir() && acceptFolder(entry.Name()) {
			dir := path.Join(root, entry.Name())
			notify(dir)
			dirs = append(dirs, dir)
			if level < maxDepth {
				subDirs, err := readDirs_(dir, level+1, maxDepth, acceptFolder, notify)
				if err == nil && len(subDirs) > 0 {
					dirs = append(dirs, subDirs...)
				}
			}
		}
	}
	return dirs, err
}

func IsRunningFromInteractiveTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) ||
		isatty.IsCygwinTerminal(os.Stdout.Fd())
}
