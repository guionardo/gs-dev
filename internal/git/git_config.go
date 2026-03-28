package git

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

/*
[core]

	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true

[remote "origin"]

	url = git@github.com:guionardo/gs-dev.git
	fetch = +refs/heads/*:refs/remotes/origin/*

[branch "develop"]

	remote = origin
	merge = refs/heads/develop
	vscode-merge-base = origin/develop

[branch "feature/new-version"]

	vscode-merge-base = origin/develop
*/
type (
	// GitConfig represents the git config file
	GitConfig struct {
		// t       *toml.Tree
		// Core    GitConfigCore
		Remotes []GitConfigRemote
	}
	// GitConfigCore struct {
	// 	RepositoryFormatVersion int
	// 	FileMode                int
	// 	Bare                    bool
	// 	LogAllRefUpdates        bool
	// }

	// GitConfigRemote represents a remote repository in the git config file
	GitConfigRemote struct {
		Name  string
		URL   string
		Fetch string
	}
)

func NewGitConfig(filename string) (*GitConfig, error) {
	remotes, err := GetGitRemotes(filename)
	if err != nil {
		return nil, err
	}

	return &GitConfig{
		Remotes: remotes,
	}, nil
}

func GetGitRemotes(filename string) ([]GitConfigRemote, error) {
	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}

	remotes := make([]GitConfigRemote, 0, 1)

	var (
		remoteName  string
		remoteURL   string
		remoteFetch string
		inRemote    bool
	)

	flushRemote := func() {
		if remoteName != "" && remoteURL != "" && remoteFetch != "" {
			remotes = append(remotes, GitConfigRemote{
				Name:  remoteName,
				URL:   remoteURL,
				Fetch: remoteFetch,
			})
		}
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			flushRemote()

			inRemote = false
			remoteName = ""
			remoteURL = ""
			remoteFetch = ""

			if strings.HasPrefix(line, "[remote \"") {
				if w := strings.Split(line, "\""); len(w) >= 2 { //nolint:mnd
					remoteName = strings.TrimSpace(w[1])
					inRemote = len(remoteName) > 0
				}
			}

			continue
		}

		if !inRemote {
			continue
		}

		if strings.HasPrefix(line, "url = ") {
			remoteURL, _ = strings.CutPrefix(line, "url = ")
			remoteURL = strings.TrimSpace(remoteURL)

			continue
		}

		if strings.HasPrefix(line, "fetch = ") {
			remoteFetch, _ = strings.CutPrefix(line, "fetch = ")
			remoteFetch = strings.TrimSpace(remoteFetch)

			continue
		}
	}

	flushRemote()

	if len(remotes) == 0 {
		return nil, errors.New("no remotes in this repository")
	}

	return remotes, nil
}
