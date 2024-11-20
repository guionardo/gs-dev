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
)

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	exitVal := m.Run()

	os.Exit(exitVal)
}

type fileMock struct {
	path    string
	content string
}

func createMockFiles(root string) error {
	fileMocks := []fileMock{
		{"dev1/go_1_1/internal/some_file.go", ""},
		{"dev1/go_1_1/go.mod", ""},
		{"dev1/py_1_1/requirements.txt", ""},
		{"dev1/project_with_local_config/.gs-dev", "ignore: false\nignore_subfolders: true\ndescription: \"Ignored project with local config\""},
		{"dev2/project_with_local_config/subfolder/README.md", "# README"},
		{"dev1/py_1_1/src/__init__.py", ""},
		{"dev2/py_2_1/package/__init__.py", ""},
		{"dev2/py_2_1/pyproject.toml", ""},
		{"dev2/project_with_local_config/.gs-dev", "ignore: true\nignore_subfolders: true\ndescription: \"Ignored project with local config\""},
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
func TestNewDirReader(t *testing.T) {
	tmpDir, _ := filepath.Abs("./TestNewDirReader")
	t.Logf("Temporary dir %s", tmpDir)

	os.RemoveAll(tmpDir)

	defer func() {
		os.RemoveAll(tmpDir)
		t.Logf("Removed temporary dir %s", tmpDir)
	}()

	if err := createMockFiles(tmpDir); err != nil {
		t.Errorf("%v", err)
		return
	}
	readerFilter := NewReaderFilter()
	dirReader := NewDirReader(readerFilter, tmpDir)

	expected := map[string]int{
		"":                               0,
		"dev1":                           0,
		"dev1/go_1_1":                    0,
		"dev1/py_1_1":                    0,
		"dev1/project_with_local_config": 0,
		"dev2":                           0,
		"dev2/subdev":                    0,
		"dev2/subdev/repository_git":     0,
		"dev2/py_2_1":                    0,
	}
	for folder := range dirReader.Folders() {
		folder = strings.TrimPrefix(strings.TrimPrefix(folder, tmpDir), "/")
		if _, ok := expected[folder]; !ok {
			t.Errorf("unexpected folder -> %s", folder)
		} else {
			delete(expected, folder)
		}
	}
	for folder := range expected {
		t.Errorf("missing folder   -> %s", folder)
	}

}
