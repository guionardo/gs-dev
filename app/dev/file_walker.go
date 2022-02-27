package dev

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func GetPathLevel(pathname string) int {
	absolute, err := filepath.Abs(pathname)
	if err != nil {
		return 0
	}

	return len(strings.Split(absolute, string(os.PathSeparator)))
}

func filterFileinfo(fileinfo os.DirEntry, parent string, maxPathLevel int) bool {
	return fileinfo.IsDir() && !strings.HasPrefix(fileinfo.Name(), ".") && GetPathLevel(path.Join(parent, fileinfo.Name())) <= maxPathLevel
}

func readFolders(pathname string, maxPathLevel int) []string {
	f, err := os.Open(pathname)
	if err != nil {
		return nil
	}
	defer f.Close()
	files, err := f.ReadDir(0)
	if err != nil {
		return nil
	}
	filelist := make([]string, 0)
	go_sub_paths := GetPathLevel(pathname) < maxPathLevel
	for _, file := range files {
		if !filterFileinfo(file, pathname, maxPathLevel) {
			continue
		}
		itemname := path.Join(pathname, file.Name())
		filelist = append(filelist, itemname)
		if !go_sub_paths {
			continue
		}
		subfilelist := readFolders(itemname, maxPathLevel)
		if len(subfilelist) > 0 {
			filelist = append(filelist, subfilelist...)
		}
	}
	return filelist
}

func ReadFolders(pathname string, maxPathLevel int) []string {
	maxPathLevel += GetPathLevel(pathname)
	items := readFolders(pathname, maxPathLevel)
	sort.Strings(items)
	return items
}
