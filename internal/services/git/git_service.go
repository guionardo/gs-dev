package gitservice

import (
	"os"

	"github.com/guionardo/gs-dev/internal/fs_tools"
	"github.com/guionardo/gs-dev/internal/git"
)

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

// GetGitStats reads the git log and produces statistics from every commit:
// * Number of commits since
func (g *GitService) GetGitStats(repositoryRoot string, filter *git.GitStatsFilter) (err error) {
	repositoryRoot, err = fs_tools.AssertDirectory(repositoryRoot)
	if err != nil {
		return err
	}

	summary, err := git.NewGitStatsSummary(repositoryRoot, filter)
	if err == nil {
		err = summary.RepositorySummary(os.Stdout)
	}

	return err
}

func (g *GitService) GetRemoteHttpURL(repositoryRoot string) (string, error) {
	return git.GetRemoteHttpURL(repositoryRoot)
}
