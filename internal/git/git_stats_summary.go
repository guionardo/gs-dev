package git

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"path"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/guionardo/gs-dev/pkg/console"
	"github.com/guionardo/gs-dev/pkg/console/symbols"
)

type (
	GitStatsSummary struct {
		repositoryName       string
		remotes              string
		currentBranch        string
		commits              int
		firstCommitTimestamp time.Time
		lastCommitTimestamp  time.Time
		totalFiles           int
		totalLinesInsertions int
		totalLinesDeletions  int

		usersStats map[string]*gitUserStats
		filter     *GitStatsFilter
	}

	gitUserStats struct {
		author               string
		commits              int
		firstCommitTimestamp time.Time
		lastCommitTimestamp  time.Time
		files                int
		linesInsertions      int
		linesDeletions       int
	}
)

func NewGitStatsSummary(repositoryRoot string, filter *GitStatsFilter) (gitStats *GitStatsSummary, err error) {
	root, err := getRepositoryRoot(repositoryRoot)
	if err != nil {
		return nil, err
	}

	gitConfigFile := path.Join(root, ".git", "config")

	gitConfig, err := NewGitConfig(gitConfigFile)
	if err != nil {
		return nil, fmt.Errorf("fail reading git config file - %w", err)
	}

	currentBranch, err := getCurrentBranch(root)
	if err != nil {
		return nil, err
	}

	stats, err := ReadGitStats(root, filter)
	if err != nil {
		return nil, err
	}

	remotes := []string{}
	for _, remote := range gitConfig.Remotes {
		remotes = append(remotes, fmt.Sprintf("%s = %s", remote.Name, remote.URL))
	}

	gitStats = &GitStatsSummary{
		currentBranch:        currentBranch,
		repositoryName:       gitConfig.RepositoryName,
		remotes:              strings.Join(remotes, ", "),
		commits:              len(stats),
		firstCommitTimestamp: stats[0].Timestamp,
		lastCommitTimestamp:  stats[len(stats)-1].Timestamp,
		usersStats:           make(map[string]*gitUserStats),
		filter:               filter,
	}

	for _, stat := range stats {
		gitStats.totalFiles += stat.ChangedFiles
		gitStats.totalLinesInsertions += stat.Insertions
		gitStats.totalLinesDeletions += stat.Deletions

		authorEmail := stat.AuthorEmail()

		ustat, ok := gitStats.usersStats[authorEmail]
		if !ok {
			ustat = &gitUserStats{author: stat.author}
			gitStats.usersStats[authorEmail] = ustat
		}

		ustat.commits++
		ustat.files += stat.ChangedFiles
		ustat.linesInsertions += stat.Insertions
		ustat.linesDeletions += stat.Deletions

		if ustat.firstCommitTimestamp.IsZero() {
			ustat.firstCommitTimestamp = stat.Timestamp
		} else if ustat.firstCommitTimestamp.After(stat.Timestamp) {
			ustat.firstCommitTimestamp = stat.Timestamp
		}

		if ustat.lastCommitTimestamp.IsZero() {
			ustat.lastCommitTimestamp = stat.Timestamp
		} else if ustat.lastCommitTimestamp.Before(stat.Timestamp) {
			ustat.lastCommitTimestamp = stat.Timestamp
		}
	}

	return gitStats, nil
}

func (s *GitStatsSummary) RepositorySummary(w io.Writer) error {
	if s.commits == 0 {
		return errors.New("no commits")
	}

	tree := console.NewTree(console.Styled(s.repositoryName, console.Green))

	repo := tree.AddChild(console.Styled("Repository", console.Blue))
	repo.AddChild(symbols.Remote + " " + s.remotes)
	repo.AddChild(symbols.Branch + " " + s.currentBranch)

	if !s.filter.IsEmpty() {
		filter := tree.AddChild(console.Styled("Filters", console.Blue))

		if len(s.filter.authors) > 0 {
			filter.AddChild("authors: " + strings.Join(s.filter.authors, ", "))
		}

		if !s.filter.since.IsZero() || !s.filter.until.IsZero() {
			filter.AddChild(symbols.Date + " " + formatTime(s.filter.since) + " -> " + formatTime(s.filter.until))
		}

		if s.filter.branch != "" {
			filter.AddChild(symbols.Branch + " " + s.filter.branch)
		}
	}

	stats := tree.AddChild(console.Styled("Stats", console.Blue))
	stats.AddChild(fmt.Sprintf("%s %d commits between %s and %s",
		symbols.Commit,
		s.commits,
		formatTime(s.firstCommitTimestamp),
		formatTime(s.lastCommitTimestamp)))

	stats.AddChild(fmt.Sprintf("%d lines", s.totalLinesInsertions-s.totalLinesDeletions))

	commitsPerDay := float64(s.commits) / max(1, s.lastCommitTimestamp.Sub(s.firstCommitTimestamp).Hours()/24) //nolint:mnd
	if commitsPerDay < 1 {
		stats.AddChild(fmt.Sprintf("%.1f days between commits (average)", 1/commitsPerDay))
	} else {
		stats.AddChild(fmt.Sprintf("%.1f commits per day (average)", commitsPerDay))
	}

	contributors := tree.AddChild(console.Styled("contributors (name, e-mail, commits, added lines, removed lines)", console.Blue))

	contributions := slices.Collect(maps.Values(s.usersStats))
	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].linesInsertions+contributions[i].linesDeletions > contributions[j].linesInsertions+contributions[j].linesDeletions
	})

	index := 0
	for _, contribution := range contributions {
		index++
		contributors.AddChild(fmt.Sprintf("%d. %s (%d, +%d, -%d)",
			index, contribution.author, contribution.commits, contribution.linesInsertions, contribution.linesDeletions))
	}

	return tree.Write(w)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "♾️"
	}

	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		return t.Format(time.DateOnly)
	}

	return t.Format(time.DateTime)
}
