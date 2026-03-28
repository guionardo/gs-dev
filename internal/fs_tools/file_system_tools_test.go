package fs_tools_test

import (
	"os"
	"testing"

	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/fs_tools"
	"github.com/stretchr/testify/require"
)

func TestAssertDirectory(t *testing.T) {
	t.Parallel()

	t.Run("AssertDirectory_with_existing_directory_should_return_absolute_path", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		require.NoError(t, os.MkdirAll(dir+"/tmp_dir", consts.DirPermissions))
		got, err := fs_tools.AssertDirectory(dir + "/tmp_dir")
		require.NoError(t, err)
		require.Equal(t, dir+"/tmp_dir", got)
	})

	t.Run("AssertDirectory_with_non_existing_directory_should_return_error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		got, err := fs_tools.AssertDirectory(dir + "/tmp_dir")
		require.Error(t, err)
		require.Empty(t, got)
	})
}

func TestAssertFilename(t *testing.T) {
	t.Parallel()

	t.Run("AssertFilename_with_existing_file_should_return_absolute_path", func(t *testing.T) {
		t.Parallel()

		filename := t.TempDir() + "/tmp_file.txt"
		_ = os.WriteFile(filename, []byte("test"), 0600)
		got, err := fs_tools.AssertFilename(filename)
		require.NoError(t, err)
		require.Equal(t, filename, got)
	})

	t.Run("AssertFilename_with_non_existing_file_should_return_absolute_path", func(t *testing.T) {
		t.Parallel()

		filename := t.TempDir() + "/tmp_file.txt"
		got, err := fs_tools.AssertFilename(filename)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("AssertFilename_with_existing_directory_should_return_error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		got, err := fs_tools.AssertFilename(dir)
		require.Error(t, err)
		require.NotEmpty(t, got)
	})

	t.Run("AssertFilename_with_empty_filename_should_return_error", func(t *testing.T) {
		t.Parallel()

		got, err := fs_tools.AssertFilename("")
		require.Error(t, err)
		require.Empty(t, got)
	})
}
