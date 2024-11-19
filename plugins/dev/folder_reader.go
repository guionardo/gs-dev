package dev

import (
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/guionardo/gs-dev/internal/colors"
)

func readFolders(root string, maxSubLevel int) (subFolders []string, err error) {
	startTime := time.Now()
	defer func() {
		if err != nil {
			colors.Red("Error reading folders %s: %v\n", root, err)
		} else {
			colors.Blue("Synced folders %s: %d subfolders @ %v\n", root, len(subFolders), time.Since(startTime))
		}
	}()
	subFolders, err = FolderReaderReadDir(root, maxSubLevel, func(name string) {
		slog.Debug("Accepted", slog.String("folder", name))
	})
	return
}

func FolderReaderReadDir(root string, maxDepth int, notify func(string)) ([]string, error) {
	return readDirs(root, 1, maxDepth, notify)
}

func readDirs(root string, level int, maxDepth int, notify func(string)) ([]string, error) {
	dirs := make([]string, 0, 1000)
	entries, err := os.ReadDir(root)
	filter := NewFolderFilter()
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		ignoreKind, reason := filter.Accept(root, entry)
		if ignoreKind == IgnoreFolder {
			slog.Debug("Ignored", slog.String("folder", entry.Name()), slog.String("reason", reason))
			continue
		}

		dir := path.Join(root, entry.Name())
		notify(dir)

		dirs = append(dirs, dir)

		if ignoreKind == NoIgnore && level < maxDepth {
			subDirs, err := readDirs(dir, level+1, maxDepth, notify)
			if err == nil && len(subDirs) > 0 {
				dirs = append(dirs, subDirs...)
			}
		}

	}
	return dirs, err
}
