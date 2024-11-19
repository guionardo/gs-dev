package install

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/guionardo/gs-dev/internal/metadata"
	"github.com/guionardo/gs-dev/internal/shell"
)

type ProfileFile struct {
	Path        string
	Lines       []string
	MarkerLines []int
	Marker      string
	LastBackup  string
}

func NewProfileFile(marker string) (pf ProfileFile, err error) {
	shellInfo, err := shell.NewShellInfo()
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
	pf.UpdateLastBackup()
	lastBackupContent, err := os.ReadFile(pf.LastBackup)
	if !os.IsNotExist(err) {
		return
	}

	currentProfile, err := os.ReadFile(pf.Path)
	if err != nil {
		return
	}
	if string(lastBackupContent) == string(currentProfile) {
		// Last backup has same content
		return nil
	}
	if stat, err := os.Stat(pf.Path); err == nil {
		fileName := fmt.Sprintf("%s.%s.bak", pf.Path, stat.ModTime().Format("20060102150405"))
		pf.LastBackup = fileName
		err = os.WriteFile(fileName, currentProfile, 0644)
	}
	return
}

func (pf *ProfileFile) UpdateLastBackup() {
	previousBackups, err := filepath.Glob(fmt.Sprintf("%s.*.bak", pf.Path))
	if err == nil && len(previousBackups) > 0 {
		slices.Sort(previousBackups)
		pf.LastBackup = previousBackups[len(previousBackups)-1]
	}
}

func (pf *ProfileFile) Save() error {
	if err := pf.DoBackup(); err != nil {
		return err
	} else {
		fmt.Printf("Backup saved as %s\n", pf.LastBackup)
	}
	return os.WriteFile(pf.Path, []byte(strings.Join(pf.Lines, "\n")), 0644)
}

func (pf *ProfileFile) UpdateMarkLines() {
	pf.MarkerLines = []int{}
	for index := range pf.Lines {
		if strings.Contains(pf.Lines[index], fmt.Sprintf("# %s", pf.Marker)) {
			pf.MarkerLines = append(pf.MarkerLines, index)
		}
	}
}

func (pf *ProfileFile) SetFeature(command string, enable bool) {
	pf.UpdateMarkLines()
	// Disable all marked lines
	for _, index := range pf.MarkerLines {
		if !strings.HasPrefix(strings.TrimLeft(pf.Lines[index], " "), "#") {
			pf.Lines[index] = fmt.Sprintf("# %s", pf.Lines[index])
		}
	}

	if !enable {
		return
	}
	commandLine := fmt.Sprintf("%s # %s", command, pf.Marker)
	if len(pf.MarkerLines) > 0 {
		// Enable just the last one
		pf.Lines[pf.MarkerLines[len(pf.MarkerLines)-1]] = commandLine
	} else {
		// Add a new line
		pf.Lines = append(pf.Lines, fmt.Sprintf("# %s set on %v", metadata.AppDescription, time.Now().Format(time.DateTime)), commandLine)
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
