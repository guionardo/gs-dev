package dev

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/guionardo/gs-dev/pkg/tools/files"
	"gopkg.in/yaml.v3"
)

type (
	DirReader struct {
		sources []string
		filter  ReaderFilter
	}
	ValidatorFunc func(folder string) (bool, string)
	ReaderFilter  struct {
		maxLevel         int
		folderValidators []ValidatorFunc
	}
)

const (
	ReasonGsDevFile     = "gs-dev file"
	ReasonGitRepository = "git repository"
	ReasonProjectFolder = "project folder"
	ReasonPrefix        = "prefix"
	DefaultMaxLevel     = 3
)

func NewDirReader(filter ReaderFilter, sources ...string) *DirReader {
	validSources := make([]string, 0, len(sources))
	for i := range sources {
		if stat, err := os.Stat(sources[i]); err == nil && stat.IsDir() {
			if source, err := filepath.Abs(sources[i]); err == nil {
				validSources = append(validSources, source)
			}
		}
	}

	return &DirReader{
		sources: validSources,
		filter:  filter,
	}
}

func NewReaderFilter() ReaderFilter {
	validators := []ValidatorFunc{
		func(folder string) (bool, string) {
			for _, prefix := range []string{".", "_"} {
				if strings.HasPrefix(path.Base(folder), prefix) {
					return false, ReasonPrefix
				}
			}

			return true, ""
		},
		func(folder string) (bool, string) {
			if stat, err := os.Stat(path.Join(folder, ".git")); err == nil && stat.IsDir() {
				return false, ReasonGitRepository
			}

			return true, ""
		},
		func(folder string) (bool, string) {
			if files.FindFirst(folder, "pyproject.toml", "requirements.txt") != "" {
				return false, fmt.Sprintf("%s %s", ReasonProjectFolder, "python")
			}

			if files.FindFirst(folder, "go.mod") != "" {
				return false, fmt.Sprintf("%s %s", ReasonProjectFolder, "golang")
			}

			if files.FindFirst(folder, "package.json") != "" {
				return false, fmt.Sprintf("%s %s", ReasonProjectFolder, "node")
			}

			return true, ""
		},
	}

	return ReaderFilter{
		maxLevel:         DefaultMaxLevel,
		folderValidators: validators,
	}
}

func (dr *DirReader) Folders() func(func(string) bool) {
	return func(yield func(string) bool) {
		for _, source := range dr.sources {
			for folder := range dr.folders(source, 0) {
				if !yield(folder) {
					return
				}
			}
		}
	}
}

func (dr *DirReader) folders(root string, level int) func(func(string) bool) {
	return func(yield func(string) bool) {
		// Verify gs-dev file
		var gsDev LocalConfig
		if err := getGsDev(root, &gsDev); err == nil && gsDev.Ignore {
			return
		}

		if !yield(root) {
			return
		}

		if gsDev.IgnoreSubfolders {
			return
		}

		// Verify level
		if level >= dr.filter.maxLevel {
			slog.Debug("ignored sub folders due level", slog.String("folder", root), slog.Int("level", level))
			return
		}
		// Verify validators
		ok, reason := false, ""
		for _, validator := range dr.filter.folderValidators {
			ok, reason = validator(root)
			if !ok {
				break
			}
		}

		if !ok {
			slog.Debug("ignored sub folders", slog.String("folder", root), slog.String("reason", reason))
			return
		}

		// Verify subfolders
		sf, err := os.ReadDir(root)
		if err != nil {
			// check if the error is a fs.PathError
			if errors.Is(err, fs.ErrPermission) {
				slog.Debug("failed to read subfolders", slog.String("folder", root), slog.Any("error", err))
			} else {
				slog.Error("failed to read subfolders", slog.String("folder", root), slog.Any("error", err))
			}

			return
		}

		for i := range sf {
			if sf[i].IsDir() {
				for subFolder := range dr.folders(path.Join(root, sf[i].Name()), level+1) {
					if !yield(subFolder) {
						return
					}
				}
			}
		}
	}
}

func getGsDev(root string, gsDev *LocalConfig) error {
	gsDevFile := files.FindFirst(root, ".gsdev.yaml", ".gsdev.yml", ".gs_dev.yaml", ".gs_dev.yml", ".gs-dev.yaml", ".gs-dev.yml", ".gsdev", ".gs-dev", ".gs_dev")
	if gsDevFile == "" {
		return os.ErrNotExist
	}

	content, err := os.ReadFile(filepath.Clean(gsDevFile))
	if err != nil {
		return fmt.Errorf("error reading gs-dev file: %s - %w", gsDevFile, err)
	}

	if err = yaml.Unmarshal(content, gsDev); err != nil {
		return fmt.Errorf("error unmarshalling gs-dev file: %s - %w", gsDevFile, err)
	}

	return nil
}
