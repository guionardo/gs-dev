package git

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	urlparsers "github.com/guionardo/gs-dev/internal/git/url_parsers"
	"github.com/guionardo/gs-dev/internal/logging"
)

func GetRemoteHttpURL(folderName string) (string, error) {
	folderName, err := filepath.Abs(folderName)
	if err != nil {
		return "", err
	}
	// Check if folder has a .git subfolder
	gitFolder, err := getRepositoryRoot(folderName)
	if err != nil {
		return "", err
	}

	gitConfigFile := path.Join(gitFolder, ".git", "config")

	gitConfig, err := NewGitConfig(gitConfigFile)
	if err != nil {
		return "", fmt.Errorf("fail reading git config file - %w", err)
	}

	return getHttpUrl(gitConfig.Remotes[0].URL)
}

func getHttpUrl(url string) (string, error) {
	parsedURL, err := urlparsers.Parse(url)
	return parsedURL, err
}

func getRepositoryRoot(folderName string) (root string, err error) {
	root = folderName
	maxDeep := 3
	level := 0

	var stat os.FileInfo
	for level < maxDeep {
		stat, err = os.Stat(path.Join(root, ".git"))
		if err == nil && stat.IsDir() {
			if stat, err = os.Stat(path.Join(root, ".git", "config")); err != nil || stat.IsDir() {
				logging.Debug("repository doesn´t have a .git/config file", slog.String("folder", root))
				err = fmt.Errorf("repository root doesn´t have a .git/config file - %s", root)

				return "", err
			}

			if level > 0 {
				logging.Debug("repository root found on parent folder", slog.String("folder", root))
			}

			return
		}

		level++
		root = path.Dir(root)
	}

	err = fmt.Errorf("repository root not found for %s", folderName)

	return
}

func getCurrentBranch(folderName string) (string, error) {
	root, err := getRepositoryRoot(folderName)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = root

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch - %w", err)
	}

	branch := strings.TrimSuffix(string(output), "\n")

	return branch, nil
}
