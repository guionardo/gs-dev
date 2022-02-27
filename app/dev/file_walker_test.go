package dev

import (
	"fmt"
	"testing"
)

func TestReadFolders(t *testing.T) {
	type args struct {
		pathname     string
		maxPathLevel int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Test",
			args: args{
				pathname:     "/home/guionardo/dev",
				maxPathLevel: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathLevel := GetPathLevel(tt.args.pathname)
			list := ReadFolders(tt.args.pathname, tt.args.maxPathLevel)
			for _, pathname := range list {
				itemPathLevel := GetPathLevel(pathname)
				if itemPathLevel > pathLevel+2 {
					t.Errorf("Pathlevel overflow %s (%d > %d)", pathname, itemPathLevel, pathLevel+2)
				}
			}
			fmt.Println(list)
		})
	}
}
