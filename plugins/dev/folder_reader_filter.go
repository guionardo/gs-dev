package dev

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type (
	FolderFilter struct {
		ignoreFilters []IgnoreFunc
	}
	IgnoreFunc func(root string, entry fs.DirEntry) (IgnoreKind, string)
	IgnoreKind byte
)

const (
	NoIgnore         IgnoreKind = 0
	IgnoreSubfolders IgnoreKind = 1
	IgnoreFolder     IgnoreKind = 2

	ReasonGsDevFile     = "gs-dev file"
	ReasonProjectFolder = "project folder"
	ReasonGitRepository = "git repository"
	ReasonIsFile        = "is file"
	ReasonPrefix        = "prefix"
	ReasonFileExists    = "file exists"
)

var projectFolderByContent = map[string][]string{
	"git repository": {".git"},
	"python project": {"pyproject.toml", "requirements.txt"},
	"golang project": {"go.mod"},
	"node project":   {"package.json"},
}

func NewFolderFilter() *FolderFilter {
	return (&FolderFilter{}).
		WithIgnore(ignoreByGsDev).
		WithIgnore(ignoreByProjectFolder).
		WithIgnorePrefixes("_", ".", "node_modules")
}

func findFirst(root string, names ...string) (fileName string, ok bool) {
	for i := range names {
		matches, err := filepath.Glob(path.Join(root, names[i]))
		if err == nil && len(matches) > 0 {
			ok = true
			fileName = matches[0]
			break
		}
	}
	return
}

func ignoreByGsDev(root string, name fs.DirEntry) (IgnoreKind, string) {
	gsDevFile, ok := findFirst(path.Join(root, name.Name()), ".gsdev.yaml", ".gsdev.yml", ".gs_dev.yaml", ".gs_dev.yml", ".gs-dev.yaml", ".gs-dev.yml")
	if !ok {
		return NoIgnore, ""
	}
	content, err := os.ReadFile(gsDevFile)
	if err != nil {
		slog.Error("Error reading gs-dev file", slog.String("file", gsDevFile), slog.Any("error", err))
		return NoIgnore, ""
	}
	var gsDev LocalConfig
	if err = yaml.Unmarshal(content, &gsDev); err != nil {
		slog.Error("Error unmarshalling gs-dev file", slog.String("file", gsDevFile), slog.Any("error", err))
		return NoIgnore, ""
	}
	if gsDev.Ignore {
		return IgnoreFolder, ReasonGsDevFile
	}
	if gsDev.IgnoreSubfolders {
		return IgnoreSubfolders, ReasonGsDevFile
	}
	return NoIgnore, ""
}

func ignoreByProjectFolder(root string, name fs.DirEntry) (IgnoreKind, string) {
	for projectType := range projectFolderByContent {
		if fileExists(path.Join(root, name.Name()), projectFolderByContent[projectType]...) {
			slog.Debug("Found", slog.String("type", projectType), slog.String("dir", root))
			return IgnoreSubfolders, fmt.Sprintf("%s - %s", ReasonProjectFolder, projectType)
		}
	}
	return NoIgnore, ""
}

func (f *FolderFilter) WithIgnore(ignoreFunc IgnoreFunc) *FolderFilter {
	f.ignoreFilters = append(f.ignoreFilters, ignoreFunc)
	return f
}

func (f *FolderFilter) WithIgnorePrefixes(prefixes ...string) *FolderFilter {
	return f.WithIgnore(func(root string, name fs.DirEntry) (IgnoreKind, string) {
		for i := range prefixes {
			if strings.HasPrefix(name.Name(), prefixes[i]) {
				return IgnoreFolder, fmt.Sprintf("%s %s", ReasonPrefix, strings.Join(prefixes, ", "))
			}
		}
		return NoIgnore, ""
	})
}

func (f *FolderFilter) WithIgnoreFileExists(names ...string) *FolderFilter {
	return f.WithIgnore(func(root string, name fs.DirEntry) (IgnoreKind, string) {
		for i := range names {
			matches, err := filepath.Glob(path.Join(root, names[i]))
			if err == nil && len(matches) > 0 {
				return IgnoreFolder, fmt.Sprintf("%s %s", ReasonFileExists, strings.Join(matches, ", "))
			}
		}
		return NoIgnore, ""
	})
}

func (f *FolderFilter) Accept(root string, name fs.DirEntry) (ignore IgnoreKind, reason string) {
	folder := path.Join(root, name.Name())
	defer func() {
		slog.Debug("Accept", slog.String("folder", folder), slog.String("reason", reason), slog.Int("kind", int(ignore)), slog.Bool("accepted", ignore != IgnoreFolder))
	}()

	if ignore, reason = ignoreByGsDev(root, name); ignore == IgnoreFolder {
		return
	}
	if ignore, reason = ignoreGitRepository(folder); ignore == IgnoreSubfolders {
		return
	}

	for projectType, files := range projectFolderByContent {
		if found, _ := findFile(folder, files...); found {
			ignore = IgnoreSubfolders
			reason = projectType
			return
		}
	}

	return

	// for i := range f.ignoreFilters {
	// 	if ignoring, reason := f.ignoreFilters[i](root, name); ignoring != NoIgnore {
	// 		slog.Debug("Ignoring", slog.String("folder", name.Name()), slog.String("reason", reason), slog.Int("kind", int(ignoring)))
	// 		return ignoring, reason
	// 	}
	// }
	// slog.Debug("Accepting", slog.String("folder", name.Name()))
	// return NoIgnore, ""
}

func ignoreGitRepository(folder string) (ignore IgnoreKind, reason string) {
	if fileExists(folder, ".git") {
		ignore = IgnoreSubfolders
		reason = "git repository"
	}
	return
}

func fileExists(dir string, names ...string) bool {
	found, _ := findFile(dir, names...)
	return found
}

func findFile(folder string, files ...string) (found bool, fileName string) {
	for i := range files {
		if matches, err := filepath.Glob(path.Join(folder, files[i])); err == nil && len(matches) > 0 {
			return true, matches[0]
		}
	}
	return
}
