package dev

import (
	"testing"
)

func Test_pathContainsPattern(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		prefix string
		words  []string
		want   bool
	}{
		{"Default", "/Users/gufurlan/dev", "/Users/gufurlan/", []string{"dev"}, true},
		{"Default2", "/Users/gufurlan/dev/guionardo/gs-dev", "/Users/gufurlan/dev", []string{"guionardo", "dev"}, true},
		{"Default3", "/Users/gufurlan/dev/guionardo/gs-dev", "/Users/gufurlan/dev", []string{"guionardo", "code"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathContainsPattern(tt.path, tt.prefix, tt.words); got != tt.want {
				t.Errorf("pathContainsPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
