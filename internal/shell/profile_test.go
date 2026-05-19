package shell

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProfileFile(t *testing.T) {
	t.Parallel()

	if shell := os.Getenv("SHELL"); len(shell) == 0 {
		t.Skipf("SHELL NOT FOUND")
		return
	}

	p, err := NewProfileFile("test")
	if err != nil {
		t.Errorf("NewProfileFile() error = %v", err)
		return
	}

	filename := p.Path

	content, err := os.ReadFile(filename) // #nosec G304 -- constrained basename under config directory.
	require.NoError(t, err, "Error reading content of file %s", filename)

	backups := make([]string, 0, 2)

	defer func() {
		_ = os.WriteFile(filename, content, 0600)

		for _, backup := range backups {
			_ = os.Remove(backup)
		}
	}()

	p.SetFeature("echo 'testing'", true)

	err = p.Save()
	require.NoError(t, err, "Save() error = %v", err)

	backups = append(backups, p.LastBackup)
	p.SetFeature("echo 'testing'", false)

	err = p.Save()
	require.NoError(t, err, "Save() error = %v", err)

	backups = append(backups, p.LastBackup)
}
