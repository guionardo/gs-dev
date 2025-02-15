package install

import (
	"os"
	"testing"
)

func TestNewProfileFile(t *testing.T) {
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
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("Error reading content of file %s - %v", filename, err)
		return
	}
	backups := make([]string, 0)
	defer func() {
		_ = os.WriteFile(filename, content, 0644)
		for _, backup := range backups {
			os.Remove(backup)
		}
	}()

	p.SetFeature("echo 'testing'", true)
	err = p.Save()
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}
	backups = append(backups, p.LastBackup)
	p.SetFeature("echo 'testing'", false)
	err = p.Save()
	if err != nil {
		t.Errorf("Save() error = %v", err)
	}
	backups = append(backups, p.LastBackup)
}
