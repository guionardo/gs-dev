package internal

import (
	"hash/crc64"
	"os"
	"strings"

	"github.com/guionardo/gs-dev/internal/logger"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
)

type (
	PathStatus uint8
	PathList   struct {
		Root         string           `yaml:"root"`
		Paths        map[uint64]*Path `yaml:"paths"`
		MaxSubLevels int              `yaml:"max_sub_levels"`
	}
	Path struct {
		Name   string
		Status PathStatus
	}
)

const (
	Enabled PathStatus = iota
	Disabled
	ChildrenDisabled
)

func NewPathList(pathRoot string, maxSubLevels int) *PathList {
	return &PathList{
		Root:         pathRoot,
		Paths:        make(map[uint64]*Path),
		MaxSubLevels: maxSubLevels,
	}
}

func (pl *PathList) isIgnoredFolder(path string) bool {
	for _, folder := range pl.Paths {
		if folder.Status == Enabled {
			continue
		}
		if len(path) > len(folder.Name) &&
			strings.HasPrefix(path, folder.Name) {
			return true
		}
	}
	return false
}

func (pl *PathList) Sync() error {
	existingFolders, err := pathtools.ReadFolders(pl.Root, pl.MaxSubLevels)
	if err != nil {
		return err
	}

	// Remove folders with parents disabled or children disabled
	i := 0
	last := len(existingFolders) - 1
	for i <= last {
		if pl.isIgnoredFolder(existingFolders[i]) {
			existingFolders[i] = existingFolders[last]
			last--
		} else {
			i++
		}
	}

	// Add existing folders
	for _, folder := range existingFolders[:last] {
		hashCode := GetHashCode(folder)
		if _, ok := pl.Paths[hashCode]; ok {
			continue
		}
		pl.Paths[hashCode] = &Path{
			Name:   folder,
			Status: Enabled,
		}
		logger.Info("+ %s", folder)
	}

	// Remove folders not found
	for _, folder := range pl.Paths {
		if stat, err := os.Stat(folder.Name); err != nil || !stat.IsDir() {
			delete(pl.Paths, GetHashCode(folder.Name))
			logger.Info("- %s", folder.Name)
		}
	}

	return nil

}

func (pl *PathList) Find(folder string) *Path {
	if p, ok := pl.Paths[GetHashCode(folder)]; ok {
		return p
	}
	return nil
}

func (pl *PathList) FindByPattern(args []string) []*Path {
	paths := make([]*Path, 0, len(pl.Paths))

	for _, folder := range pl.Paths {
		if isValidPattern(pl.Root, folder.Name, args) {
			paths = append(paths, folder)

		}
	}
	return paths
}

func isValidPattern(rootFolder string, folder string, args []string) bool {
	searchData, _ := strings.CutPrefix(folder, rootFolder)
	for _, arg := range args {
		if _, after, found := strings.Cut(searchData, arg); !found {
			return false
		} else {
			searchData = after
		}
	}
	return true
}

func GetHashCode(text string) uint64 {
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	h.Write([]byte(text))
	return h.Sum64()
}

func (pl *PathList) GetChildren(folder *Path) []*Path {
	children := make([]*Path, 0, len(pl.Paths))
	for _, path := range pl.Paths {
		if folder != path && strings.HasPrefix(path.Name, folder.Name) {
			children = append(children, path)
		}
	}
	return children
}
