package shell

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/consts"
)

type ProfileFile struct {
	Path        string
	Lines       []string
	MarkerLines []int
	Marker      string
	LastBackup  string
}

func NewProfileFile(marker string) (pf ProfileFile, err error) {
	shellInfo, err := NewShellInfo()
	if err != nil {
		return
	}

	content, err := os.ReadFile(shellInfo.RCFile)
	if err != nil {
		return
	}

	pf.Path = shellInfo.RCFile
	pf.Lines = strings.Split(string(content), "\n")
	pf.Marker = marker

	for index := range pf.Lines {
		pf.Lines[index] = strings.TrimRight(pf.Lines[index], " \n\r")
	}

	return
}

func (pf *ProfileFile) DoBackup() (err error) {
	lastBackupContent := pf.getLastBackup()

	currentProfile, err := os.ReadFile(pf.Path)
	if err != nil {
		return fmt.Errorf("error reading current profile file %s - %w", pf.Path, err)
	}

	if string(lastBackupContent) == string(currentProfile) {
		// Last backup has same content
		return nil
	}

	var stat os.FileInfo
	if stat, err = os.Stat(pf.Path); err == nil {
		fileName := fmt.Sprintf("%s.%s.bak", pf.Path, stat.ModTime().Format("20060102150405"))
		pf.LastBackup = fileName
		err = os.WriteFile(fileName, currentProfile, consts.FilesPermissions)
	}

	return
}

func (pf *ProfileFile) Save() error {
	if err := pf.DoBackup(); err != nil {
		return err
	} else {
		slog.Debug("Profile backup", slog.String("backup", pf.LastBackup))
	}

	return os.WriteFile(pf.Path, []byte(strings.Join(pf.Lines, "\n")), consts.FilesPermissions)
}

func (pf *ProfileFile) UpdateMarkLines() {
	pf.MarkerLines = []int{}
	for index := range pf.Lines {
		if strings.Contains(pf.Lines[index], "# "+pf.Marker) {
			pf.MarkerLines = append(pf.MarkerLines, index)
		}
	}
}

func (pf *ProfileFile) SetFeature(command string, enable bool) {
	pf.UpdateMarkLines()
	// Disable all marked lines
	for _, index := range pf.MarkerLines {
		if !strings.HasPrefix(strings.TrimLeft(pf.Lines[index], " "), "#") {
			pf.Lines[index] = "# " + pf.Lines[index]
		}
	}

	if !enable {
		return
	}

	commandLine := fmt.Sprintf("%s # %s", command, pf.Marker)
	if len(pf.MarkerLines) > 0 {
		// Enable just the last one
		pf.Lines[pf.MarkerLines[len(pf.MarkerLines)-1]] = commandLine
	} else { // Add a new line
		pf.Lines = append(pf.Lines, fmt.Sprintf("# %s set on %v",
			strings.ReplaceAll(build.AppDescription, "\n", " "),
			time.Now().Format(time.DateTime)), commandLine)
	}
}

func (pf *ProfileFile) HasEnabledCommandLine() (int, bool) {
	pf.UpdateMarkLines()

	for _, index := range pf.MarkerLines {
		if !strings.HasPrefix(strings.TrimLeft(pf.Lines[index], " "), "#") {
			return index, true
		}
	}

	return 0, false
}

// getLastBackup returns the last backup file and its content if it exists
func (pf *ProfileFile) getLastBackup() (backupFileContent []byte) {
	previousBackups, _ := filepath.Glob(pf.Path + ".*.bak") //nolint:errcheck
	if len(previousBackups) == 0 {
		return nil
	}

	slices.Sort(previousBackups)
	backupFile := previousBackups[len(previousBackups)-1]

	backupFileContent, err := os.ReadFile(backupFile) //nolint:errcheck,gosec
	if err == nil {
		pf.LastBackup = backupFile
		return backupFileContent
	}

	slog.Error("Error reading backup file", slog.String("backup", backupFile), slog.Any("error", err))

	return nil
}
