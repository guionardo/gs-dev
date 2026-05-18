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

func TestFindFirstInSubdirs(t *testing.T) {
	t.Parallel()

	t.Run("find_file_at_root_with_maxdepth_0", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.WriteFile(path.Join(root, "file.txt"), []byte("content"), 0600)
		got := files.FindFirstInSubdirs(root, 0, "file.txt")
		assert.Equal(t, path.Join(root, "file.txt"), got)
	})

	t.Run("find_file_at_root_with_maxdepth_1", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.WriteFile(path.Join(root, "file.txt"), []byte("content"), 0600)
		got := files.FindFirstInSubdirs(root, 1, "file.txt")
		assert.Equal(t, path.Join(root, "file.txt"), got)
	})

	t.Run("find_file_one_level_deep", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.MkdirAll(path.Join(root, "subdir"), 0750)
		_ = os.WriteFile(path.Join(root, "subdir", "file.txt"), []byte("content"), 0600)
		got := files.FindFirstInSubdirs(root, 1, "file.txt")
		assert.Equal(t, path.Join(root, "subdir", "file.txt"), got)
	})

	t.Run("shallow_depth_wins_over_deeper", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.WriteFile(path.Join(root, "file.txt"), []byte("root"), 0600)
		_ = os.MkdirAll(path.Join(root, "subdir"), 0750)
		_ = os.WriteFile(path.Join(root, "subdir", "file.txt"), []byte("deep"), 0600)
		got := files.FindFirstInSubdirs(root, 1, "file.txt")
		assert.Equal(t, path.Join(root, "file.txt"), got)
	})

	t.Run("find_file_two_levels_deep", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.MkdirAll(path.Join(root, "a", "b"), 0750)
		_ = os.WriteFile(path.Join(root, "a", "b", "file.txt"), []byte("content"), 0600)
		got := files.FindFirstInSubdirs(root, 2, "file.txt")
		assert.Equal(t, path.Join(root, "a", "b", "file.txt"), got)
	})

	t.Run("not_found_returns_empty", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.MkdirAll(path.Join(root, "subdir"), 0750)
		got := files.FindFirstInSubdirs(root, 2, "nonexistent.txt")
		assert.Empty(t, got)
	})

	t.Run("first_matching_name_wins", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.MkdirAll(path.Join(root, "subdir"), 0750)
		_ = os.WriteFile(path.Join(root, "subdir", "alpha.txt"), []byte("content"), 0600)
		_ = os.WriteFile(path.Join(root, "subdir", "beta.txt"), []byte("content"), 0600)
		got := files.FindFirstInSubdirs(root, 1, "alpha.txt", "beta.txt")
		assert.Equal(t, path.Join(root, "subdir", "alpha.txt"), got)
	})

	t.Run("limited_maxdepth_prevents_deeper_match", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		_ = os.MkdirAll(path.Join(root, "a", "b"), 0750)
		_ = os.WriteFile(path.Join(root, "a", "b", "file.txt"), []byte("content"), 0600)
		got := files.FindFirstInSubdirs(root, 1, "file.txt")
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
