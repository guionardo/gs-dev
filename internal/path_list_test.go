package internal

import (
	"os"
	"path"
	"testing"

	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
)

func TestPathList_Sync(t *testing.T) {

	t.Run("Default", func(t *testing.T) {
		tmp := t.TempDir()
		dirs := []string{
			path.Join(tmp, "A", "B1"),
			path.Join(tmp, "A", "B2"),
			path.Join(tmp, "B", "00"),
			path.Join(tmp, "B", "01"),
			path.Join(tmp, "GAMA", "XY"),
			path.Join(tmp, "DELTA", "XX"),
		}

		for _, dir := range dirs {
			if err := pathtools.CreatePath(dir); err != nil {
				t.Errorf("Failed to create path: %v", err)
				return
			}
		}
		pl := NewPathList(tmp, 3)
		if err := pl.Sync(); err != nil {
			t.Errorf("PathList.Sync() error = %v", err)
			return
		}
		os.RemoveAll(path.Join(tmp, "B"))
		if err := pl.Sync(); err != nil {
			t.Errorf("PathList.Sync() error = %v", err)
			return
		}
		if len(pl.Paths) != 6 {
			t.Errorf("Expected 6 paths, found %d", len(pl.Paths))
		}
	})

}

func Test_isValidPattern(t *testing.T) {
	tests := []struct {
		name   string
		folder string
		args   []string
		want   bool
	}{

		{"dev OK", "/home/guionardo/dev/github.com/guionardo/gs-dev/app", []string{"dev"}, true},
		{"dev test FAIL", "/home/guionardo/dev/github.com/guionardo/gs-dev/app/dev", []string{"dev", "test"}, false},
		{"gs-dev setup OK", "/home/guionardo/dev/github.com/guionardo/gs-dev/app/setup", []string{"gs-dev", "setup"}, true},
		{"dev command FAIL", "/home/guionardo/dev/github.com/guionardo/gs-dev/cmd", []string{"dev", "command"}, false},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/commands", []string{"gs-dev", "commands"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/commands/dev", []string{"dev"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/config", []string{"dev"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/configs", []string{"dev"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/docs", []string{"dev"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/gen", []string{"dev"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/internal", []string{"dev"}, true},
		{"test", "/home/guionardo/dev/github.com/guionardo/gs-dev/internal/path_tools", []string{"dev"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidPattern("/home/guionardo/dev", tt.folder, tt.args); got != tt.want {
				t.Errorf("isValidPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
