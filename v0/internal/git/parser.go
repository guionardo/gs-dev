package git

import (
	"fmt"
	"regexp"
)

type GitURL struct {
	Success bool
	Domain  string
	Repo    string
	Format  string
}

func (g GitURL) GetURL() string {
	return fmt.Sprintf(g.Format, g.Domain, g.Repo)
}

type (
	GitParser struct {
		sshRegex          *regexp.Regexp
		httpRegex         *regexp.Regexp
		gitUrlFromMatches matchesParser
	}
	matchesParser func([]string) *GitURL
)

var parsers = make(map[string]*GitParser, 0)

func Register(name string, sshRegex string, httpRegex string, getSshGitUrlFromMatches matchesParser) {
	parsers[name] = &GitParser{
		sshRegex:          regexp.MustCompile(sshRegex),
		httpRegex:         regexp.MustCompile(httpRegex),
		gitUrlFromMatches: getSshGitUrlFromMatches,
	}
}

func Parse(url string) (gitUrl *GitURL, err error) {
	var matches []string
	for _, parser := range parsers {
		if matches = parser.sshRegex.FindStringSubmatch(url); len(matches) == 0 {
			if matches = parser.httpRegex.FindStringSubmatch(url); len(matches) == 0 {
				continue
			}
		}
		gitUrl = parser.gitUrlFromMatches(matches)
		break
	}
	if gitUrl == nil {
		err = fmt.Errorf("cannot parse URL %s to a git repository", url)
	}
	return
}
