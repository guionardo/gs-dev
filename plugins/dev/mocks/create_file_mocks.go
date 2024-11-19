package mocks_test

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type fileMock struct {
	path    string
	content string
}

var fileMocks []fileMock

func createMockFiles(root string) error {
	for _, mock := range fileMocks {
		words := strings.Split(mock.path, "/")
		baseDir := path.Join(words[0 : len(words)-2]...)
		dir := path.Join(root, baseDir)
		if stat, err := os.Stat(dir); os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0644); err != nil {
				return err
			}
		} else if !stat.IsDir() {
			return fmt.Errorf("failed to create folder %s. Just exists as file", dir)
		}
		fileName := path.Join(dir, words[len(words)-1])
		if err := os.WriteFile(fileName, []byte(mock.content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s -> %v", fileName, err)
		}
	}
	return nil
}

func init() {
	fileMocks = []fileMock{
		{"dev1/go_1_1/internal/some_file.go", ""},
		{"dev1/go_1_1/go.mod", ""},
		{"dev1/py_1_1/requirements.txt", ""},
		{"dev1/py_1_1/src/__init__.py", ""},
		{"dev2/py_2_1/package/__init__.py", ""},
		{"dev2/py_2_1/pyproject.yaml", ""},
		{"dev2/subdev/project_with_local_config/.gs-dev", "ignore: true\nignore_subfolders: true\ndescription: \"Ignored project with local config\""},
		{"dev2/subdev/repository_git/.git/.gitkeep", ""},
	}
}
