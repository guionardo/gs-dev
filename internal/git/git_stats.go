package git

import (
	"bytes"
	"errors"
	"fmt"
	"iter"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/guionardo/go/flow"
)

type (
	GitStatsFilter struct {
		authors []string
		since   time.Time
		until   time.Time
		branch  string
	}
	GitCommit struct {
		Timestamp    time.Time
		ChangedFiles int
		Insertions   int
		Deletions    int

		author   string
		hash     string
		meta     string
		comments []string
	}
)

func NewGitStatsFilter() *GitStatsFilter {
	return &GitStatsFilter{
		since: time.Time{},
		until: time.Now().AddDate(0, 0, 1),
	}
}

func (f *GitStatsFilter) IsEmpty() bool {
	return f == nil || (len(f.authors) == 0 && f.branch == "" && f.since.IsZero() && f.until.IsZero())
}

func (f *GitStatsFilter) Since(t time.Time) *GitStatsFilter {
	f.since = t
	return f
}

func (f *GitStatsFilter) To(t time.Time) *GitStatsFilter {
	f.until = t
	return f
}

func (f *GitStatsFilter) Authors(authors ...string) *GitStatsFilter {
	f.authors = append(f.authors, authors...)
	return f
}

func (c *GitCommit) Comments() []string {
	for len(c.comments) > 0 && c.comments[0] == "" {
		c.comments = c.comments[1:]
	}

	for len(c.comments) > 0 && c.comments[len(c.comments)-1] == "" {
		c.comments = c.comments[:len(c.comments)]
	}

	return c.comments
}

func (c GitCommit) String() string {
	if c.hash == "" {
		return "⚠️ Empty"
	}

	icon := flow.If(c.isComplete(), "✅", "❌")

	return fmt.Sprintf("%s %s %s %s %s (changes: %d = +%d -%d)", icon, c.hash[:4], c.Timestamp.Format(time.DateTime), c.author, c.meta, c.ChangedFiles, c.Insertions, c.Deletions)
}

// AuthorEmail get the email from author string: Guionardo Furlan <guionardo@gmail.com>
func (c *GitCommit) AuthorEmail() string {
	w := strings.SplitN(c.author, "<", 2)
	if len(w) < 2 {
		return c.author
	}

	w = strings.SplitN(w[1], ">", 2)
	if len(w) < 1 {
		return c.author
	}

	return w[0]
}

func ReadGitStats(rootDirectory string, filter *GitStatsFilter) (stats []GitCommit, err error) {
	args := []string{"log", "--shortstat", "--date=iso-strict", "--reverse"}

	if filter == nil {
		filter = &GitStatsFilter{}
	}

	if !filter.since.IsZero() {
		args = append(args, "--since="+filter.since.Format(time.RFC3339))
	}

	if !filter.until.IsZero() {
		args = append(args, "--until="+filter.until.Format(time.RFC3339))
	}

	for _, author := range filter.authors {
		args = append(args, "--author="+author)
	}

	if filter.branch != "" {
		args = append(args, filter.branch)
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = rootDirectory

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err // TODO: Implementar tratativa do erro de execução
	}

	stats = slices.Collect(parseGitStats(output))
	if len(stats) == 0 {
		err = errors.New("git log is empty")
	}

	return
}

// parseGitStats splits the output in lines, parsing the commit data from git log --shortstat
func parseGitStats(output []byte) iter.Seq[GitCommit] {
	const LF byte = 10

	return func(yield func(GitCommit) bool) {
		var commit GitCommit

		doYield := func() bool {
			yieldResult := true
			if commit.isComplete() {
				yieldResult = yield(commit)
				commit.reset()
			}

			return yieldResult
		}

		for lineBytes := range bytes.SplitSeq(output, []byte{LF}) {
			if !doYield() {
				return
			}

			line := strings.ReplaceAll(strings.TrimSuffix(string(lineBytes), "\n"), "\t", "")

			// Try to parse summary
			if readSummary(line, &commit) {
				continue
			}

			words := splitInWords(line)
			if len(words) == 0 {
				continue
			}

			if parseCommit(words, &commit) {
				continue
			}

			if parseAuthor(words, &commit) {
				continue
			}

			if parseDate(words, &commit) {
				continue
			}

			commit.comments = append(commit.comments, line)
		}

		doYield() // Send last commit if still valid
	}
}

// readSummary parses the text:  5 files changed, 1 insertion(+), 179 deletions(-)
func readSummary(line string, commit *GitCommit) bool {
	words := splitInWords(line)
	if len(words) != 7 {
		return false
	}

	var (
		changedFiles, insertions, deletions int
		err                                 error
	)
	if changedFiles, err = strconv.Atoi(words[0]); err != nil {
		return false
	}

	if insertions, err = strconv.Atoi(words[3]); err != nil {
		return false
	}

	if deletions, err = strconv.Atoi(words[5]); err != nil {
		return false
	}

	commit.ChangedFiles = changedFiles
	commit.Insertions = insertions
	commit.Deletions = deletions

	return true
}

func parseCommit(words []string, commit *GitCommit) bool {
	if words[0] != "commit" {
		return false
	}

	commit.reset()
	commit.hash = words[1]
	commit.meta = strings.Join(words[2:], " ")

	return true
}

func parseAuthor(words []string, commit *GitCommit) bool {
	if words[0] != "Author:" {
		return false
	}

	commit.author = strings.Join(words[1:], " ")

	return true
}

func parseDate(words []string, commit *GitCommit) bool {
	if words[0] != "Date:" {
		return false
	}

	timeStamp, err := time.Parse(time.RFC3339, words[1])
	if err != nil {
		return false
	}

	commit.Timestamp = timeStamp

	return true
}
func splitInWords(s string) (out []string) {
	for word := range strings.SplitSeq(s, " ") {
		word = strings.TrimSpace(strings.ReplaceAll(word, "\t", ""))
		if len(word) > 0 {
			out = append(out, word)
		}
	}

	return out
}

func (c *GitCommit) isComplete() bool {
	return c.hash != "" && len(c.comments) > 0 && c.author != "" && !c.Timestamp.IsZero() && c.ChangedFiles > 0
}

func (c *GitCommit) reset() {
	c.hash = ""
	c.meta = ""
	c.comments = nil
	c.author = ""
	c.Timestamp = time.Time{}
	c.ChangedFiles = 0
	c.Insertions = 0
	c.Deletions = 0
}
