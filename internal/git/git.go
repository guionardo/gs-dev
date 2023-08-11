package git

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func GetRemoteHttpURL(folderName string) (string, error) {
	folderName, err := filepath.Abs(folderName)
	if err != nil {
		return "", err
	}
	// Check if folder has a .git subfolder
	gitFolder := path.Join(folderName, ".git")
	if _, err := os.Stat(gitFolder); err != nil {
		return "", fmt.Errorf("folder %s is not a git repository", folderName)
	}
	// Run git config --get remote.origin.url
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir, _ = filepath.Abs(folderName)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("current repository has no remote origin - %v", err)
	}
	output := strings.ReplaceAll(strings.SplitN(string(out), "\n", 1)[0], "\n", "")
	return getHttpUrl(output)
}

func getHttpUrl(url string) (string, error) {
	gu, err := Parse(url)
	if !gu.Success || err != nil {
		return "", fmt.Errorf("invalid git url: %s", url)
	}
	return gu.GetURL(), nil
}
