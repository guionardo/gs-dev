package gitservice

import (
	"os"

	"github.com/guionardo/gs-dev/internal/git"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type GitService struct {
}

func NewGitService() *GitService {
	return &GitService{}
}

// GetGitStats reads the git log and produces statistics from every commit:
// * Number of commits since
func (g *GitService) GetGitStats(repositoryRoot string, filter *git.GitStatsFilter) (err error) {
	repositoryRoot, err = files.AssertDirectory(repositoryRoot)
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
