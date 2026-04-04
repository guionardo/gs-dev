package gitservice

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/fs_tools"
	"github.com/guionardo/gs-dev/internal/git"
)

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

func (g *GitService) GetGitStats(repositoryRoot string) error {
	repositoryRoot, err := fs_tools.AssertDirectory(repositoryRoot)
	if err != nil {
		return err
	}

	remoteURL, err := git.GetRemoteHttpURL(repositoryRoot)
	if err != nil {
		return err
	}

	fmt.Println(remoteURL)
	// TODO: Get git stats from remote URL

	return nil
}

func (g *GitService) GetRemoteHttpURL(repositoryRoot string) (string, error) {
	return git.GetRemoteHttpURL(repositoryRoot)
}
