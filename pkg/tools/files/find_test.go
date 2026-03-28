package files_test

import (
	"os"
	"path"
	"testing"

	"github.com/guionardo/gs-dev/pkg/tools/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindFirst(t *testing.T) {
	t.Parallel()

	t.Run("find_first_file_in_root_directory", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.WriteFile(path.Join(root, "file1.txt"), []byte("file1 content"), 0600)
		_ = os.WriteFile(path.Join(root, "file2.txt"), []byte("file2 content"), 0600)
		_ = os.WriteFile(path.Join(root, "file3.txt"), []byte("file3 content"), 0600)
		got := files.FindFirst(root, "file1.txt", "file2.txt", "file3.txt")
		assert.Equal(t, path.Join(root, "file1.txt"), got)
	})

	t.Run("find_first_file_not_found_should_return_empty_string", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		got := files.FindFirst(root, "file1.txt", "file2.txt", "file3.txt")
		assert.Empty(t, got)
	})
}

func TestLocateBinary(t *testing.T) {
	t.Parallel()

	t.Run("locate_ls_should_return_path_to_ls", func(t *testing.T) {
		t.Parallel()

		got, err := files.LocateBinary("ls")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
	})

	t.Run("locate_not_found_should_return_error", func(t *testing.T) {
		t.Parallel()

		got, err := files.LocateBinary("not_found")
		require.Error(t, err)
		assert.Empty(t, got)
	})
}
