package dev

import "strings"

func pathContainsPattern(path string, prefix string, words []string) bool {
	searchData, _ := strings.CutPrefix(path, prefix)
	for index := range words {
		if _, after, found := strings.Cut(searchData, words[index]); !found {
			return false
		} else {
			searchData = after
		}
	}
	return true
}
