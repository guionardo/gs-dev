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

	_, _ = fmt.Fprintf(w, "# %s\n\n", s.repositoryName)
	_, _ = fmt.Fprintf(w, "remote: %s\n\n", s.remotes)

	_, _ = fmt.Fprintf(w, "branch: %s\n\n", s.currentBranch)
	if s.filter != nil {
		_, _ = fmt.Fprintf(w, "filters:\n")
		if len(s.filter.authors) > 0 {
			_, _ = fmt.Fprintf(w, "\tauthors: %s\n", strings.Join(s.filter.authors, ", "))
		}

		if !s.filter.since.IsZero() {
			_, _ = fmt.Fprintf(w, "\tsince: %s\n", s.filter.since.Format(time.DateTime))
		}

		if !s.filter.until.IsZero() {
			_, _ = fmt.Fprintf(w, "\tuntil: %s\n", s.filter.until.Format(time.DateTime))
		}

		if s.filter.branch != "" {
			_, _ = fmt.Fprintf(w, "\tbranch: %s\n", s.filter.branch)
		}

		_, _ = fmt.Fprintln(w)
	}

	_, _ = fmt.Fprintf(w, "* %d commits from %s to %s\n",
		s.commits,
		s.firstCommitTimestamp.Format(time.DateTime),
		s.lastCommitTimestamp.Format(time.DateTime))
	_, _ = fmt.Fprintf(w, "* %d lines\n", s.totalLinesInsertions-s.totalLinesDeletions)

	commitsPerDay := float64(s.commits) / max(1, s.lastCommitTimestamp.Sub(s.firstCommitTimestamp).Hours()/24)
	if commitsPerDay < 1 {
		_, _ = fmt.Fprintf(w, "* %.1f days between commits (average)\n", 1/commitsPerDay)
	} else {
		_, _ = fmt.Fprintf(w, "* %.1f commits per day (average)\n", commitsPerDay)
	}

	fmt.Fprintf(w, "\n")

	_, _ = fmt.Fprintf(w, "## Contributors\n\n")

	contributions := slices.Collect(maps.Values(s.usersStats))
	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].linesInsertions+contributions[i].linesDeletions > contributions[j].linesInsertions+contributions[j].linesDeletions
	})

	index := 0
	for _, contribution := range contributions {
		index++
		_, _ = fmt.Fprintf(w, "%d. %s (%d commits, %d added lines, %d removed lines)\n",
			index, contribution.author, contribution.commits, contribution.linesInsertions, contribution.linesDeletions)
	}

	return nil
}
