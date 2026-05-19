package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guionardo/gs-dev/pkg/tools/files"
	"github.com/stretchr/testify/assert"
)

func TestReadDirectory(t *testing.T) {
	t.Parallel()

	// generate test cases for ReadDirectory function
	tmp := t.TempDir()
	cases := []string{"dir1/file1.txt", "dir1/file2.txt", "dir2/file3.txt", "dir3/dir4/file4.txt", "dir5/dir6/dir7/file5.txt",
		"dir1/dir2/dir3/dir4/file6.txt"}

	for _, cas := range cases {
		createFile(t, filepath.Join(tmp, cas))
	}

	var got []string
	for dir := range files.ReadDirectory(tmp, 2) {
		got = append(got, dir)
	}

	assert.Len(t, got, 7)
}

func createFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		t.Fatalf("failed to create directory for file: %v", err)
	}

	if err := os.WriteFile(path, []byte{}, 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
