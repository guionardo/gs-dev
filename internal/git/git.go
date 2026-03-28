package git

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
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
	gu, err := Parse(url)
	if err != nil || !gu.Success {
		return "", fmt.Errorf("invalid git url: %s", url)
	}

	return gu.GetURL(), nil
}

func getRepositoryRoot(folderName string) (root string, err error) {
	root = folderName
	maxDeep := 2
	level := 0

	var stat os.FileInfo
	for level < maxDeep {
		stat, err = os.Stat(path.Join(root, ".git"))
		if err == nil && stat.IsDir() {
			if stat, err = os.Stat(path.Join(root, ".git", "config")); err != nil || stat.IsDir() {
				slog.Debug("repository doesn´t have a .git/config file", slog.String("folder", root))
				err = fmt.Errorf("repository root doesn´t have a .git/config file - %s", root)

				return "", err
			}

			if level > 0 {
				slog.Debug("repository root found on parent folder", slog.String("folder", root))
			}

			return
		}

		level++
		root = path.Dir(root)
	}

	err = fmt.Errorf("repository root not found for %s", folderName)

	return
}
