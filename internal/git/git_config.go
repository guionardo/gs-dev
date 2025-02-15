package git

import (
	"fmt"
	"os"
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
	GitConfigRemote struct {
		Name  string
		URL   string
		Fetch string
	}
)

func NewGitConfig(filename string) (*GitConfig, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var remoteName, remoteUrl, remoteFetch string
	var remotes = make([]GitConfigRemote, 0, 1)
	for _, line := range strings.Split(string(content), "\n") {
		if len(remoteName) > 0 && len(remoteUrl) > 0 && len(remoteFetch) > 0 {
			remotes = append(remotes, GitConfigRemote{
				Name:  remoteName,
				URL:   remoteUrl,
				Fetch: remoteFetch,
			})
			remoteName = ""
			remoteUrl = ""
			remoteFetch = ""
		}
		line = strings.Trim(line, " \n\t")
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "[remote \"") {
			if w := strings.Split(line, "\""); len(w) == 3 {
				remoteName = w[1]
			}

			continue
		}
		if len(remoteName) > 0 && strings.HasPrefix(line, "url = ") {
			remoteUrl, _ = strings.CutPrefix(line, "url = ")
			continue
		}
		if len(remoteName) > 0 && strings.HasPrefix(line, "fetch = ") {
			remoteFetch, _ = strings.CutPrefix(line, "fetch = ")
			continue
		}

	}
	if len(remotes) == 0 {
		return nil, fmt.Errorf("no remotes in this repository")
	}

	return &GitConfig{
		// Core: GitConfigCore{
		// 	RepositoryFormatVersion: tree.GetPath([]string{"core", "repositoryformatversion"}).(int),
		// },
		Remotes: remotes,
	}, nil
}
