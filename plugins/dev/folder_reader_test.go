package dev

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
	"github.com/stretchr/testify/assert"
)

type fileMock struct {
	path    string
	content string
}

func createMockFiles(root string) error {
	fileMocks := []fileMock{
		{"dev1/go_1_1/internal/some_file.go", ""},
		{"dev1/go_1_1/go.mod", ""},
		{"dev1/py_1_1/requirements.txt", ""},
		{"dev1/py_1_1/src/__init__.py", ""},
		{"dev2/py_2_1/package/__init__.py", ""},
		{"dev2/py_2_1/pyproject.toml", ""},
		{"dev2/subdev/project_with_local_config/.gs-dev", "ignore: true\nignore_subfolders: true\ndescription: \"Ignored project with local config\""},
		{"dev2/subdev/repository_git/.git/.gitkeep", ""},
	}
	for _, mock := range fileMocks {
		words := strings.Split(mock.path, "/")
		baseDir := path.Join(words[0 : len(words)-1]...)
		dir := path.Join(root, baseDir)
		var err error
		if err = pathtools.CreatePath(dir); err != nil {
			return err
		}

		fileName := path.Join(dir, words[len(words)-1])
		if err = os.WriteFile(fileName, []byte(mock.content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s -> %v", fileName, err)
		}
	}
	return nil
}

func Test_readFolders(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	tmpDir, _ := filepath.Abs("./_testingReadFolders")

	os.RemoveAll(tmpDir)

	if err := pathtools.CreatePath(tmpDir); err != nil {
		t.Errorf("Failed to create temporary folder %s - %v", tmpDir, err)
		return
	}
	defer os.RemoveAll(tmpDir)
	if err := createMockFiles(tmpDir); err != nil {
		t.Errorf("%v", err)
		return
	}
	type args struct {
		root        string
		maxSubLevel int
	}
	tests := []struct {
		name           string
		args           args
		wantSubFolders []string
		wantErr        bool
	}{
		{
			// 	name: "Default",
			// 	args: args{"dev1", 3},
			// 	wantSubFolders: []string{
			// 		path.Join(tmpDir, "dev1", "go_1_1"),
			// 		path.Join(tmpDir, "dev1", "py_1_1")},
			// }, {
			name: "dev2",
			args: args{"dev2", 3},
			wantSubFolders: []string{
				path.Join(tmpDir, "dev2", "py_2_1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubFolders, err := readFolders(path.Join(tmpDir, tt.args.root), tt.args.maxSubLevel)
			if (err != nil) != tt.wantErr {
				t.Errorf("readFolders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, gotSubFolders, tt.wantSubFolders)
			// if !reflect.DeepEqual(gotSubFolders, tt.wantSubFolders) {
			// 	t.Errorf("readFolders() = %v, want %v", gotSubFolders, tt.wantSubFolders)
			// }
		})
	}
}
