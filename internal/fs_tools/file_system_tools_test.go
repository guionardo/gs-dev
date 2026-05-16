package fs_tools_test

import (
	"fmt"
	"os"
	"path"
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

	t.Run("AssertDirectory_with_$HOME_should_get_current_home", func(t *testing.T) {
		t.Parallel()

		firstDir, err := getFirstDirectoryInHome(t)
		if err != nil {
			return
		}

		got, err := fs_tools.AssertDirectory("$HOME/" + path.Base(firstDir))
		require.NoError(t, err)
		require.Equal(t, firstDir, got)
	})

	t.Run("AssertDirectory_with_~_should_get_current_home", func(t *testing.T) {
		t.Parallel()

		firstDir, err := getFirstDirectoryInHome(t)
		if err != nil {
			return
		}

		got, err := fs_tools.AssertDirectory("~/" + path.Base(firstDir))
		require.NoError(t, err)
		require.Equal(t, firstDir, got)
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
		require.NotEmpty(t, got)
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

func getFirstDirectoryInHome(t *testing.T) (string, error) {
	t.Helper()

	home, err := os.UserHomeDir()
	if err == nil {
		var entries []os.DirEntry

		entries, err = os.ReadDir(home)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					return path.Join(home, entry.Name()), nil
				}
			}

			err = fmt.Errorf("no directories found in %s", home)
		}
	}

	t.Skipf("Failed to get HOME directory - %v", err)

	return "", err
}
