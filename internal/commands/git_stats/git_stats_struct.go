package gitstats

import (
	"context"
	"io"
	"time"

	"github.com/guionardo/gs-dev/internal/git"
	gitstats "github.com/guionardo/gs-dev/internal/services/git"
)

type GitStatsStruct struct {
	service *gitstats.GitService

	RootDirectory string    `flag:"root" default:"." description:"directory of repository"`
	Since         time.Time `flag:"since" description:"run stats since date time"`
	Until         time.Time `flag:"until" description:"run stats until date time"`
	Authors       []string  `flag:"author"`
}

func (gs *GitStatsStruct) Run(ctx context.Context, output io.Writer) error {
	filter := git.GitStatsFilter{}
	if !gs.Since.IsZero() {
		filter.Since(gs.Since)
	}

	if !gs.Until.IsZero() {
		filter.To(gs.Until)
	}

	if len(gs.Authors) > 0 {
		filter.Authors(gs.Authors...)
	}

	return gs.service.GetGitStats(gs.RootDirectory, &filter)
}

func (gs *GitStatsStruct) Setup(ctx context.Context) error {
	return nil
}
