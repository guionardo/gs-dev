package urlparsers

import (
	"fmt"
	"net/url"
	"regexp"
)

var gitSshUrlRe = regexp.MustCompile(`(?m)git@(.*):(.*)/(.*)\.git`)

// gitHttpUrlParser parses the url as a HTTPS URL
// https://github.com/guionardo/gs-dev -> https://github.com/guionardo/gs-dev
func gitHttpUrlParser(giturl string) (string, error) {
	// Try to parse the url as a HTTPS URL
	parsedURL, err := url.Parse(giturl)
	if err != nil {
		return "", err
	}

	switch parsedURL.Scheme {
	case "ssh":
		parsedURL, err := url.Parse(fmt.Sprintf("https://%s%s", parsedURL.Hostname(), parsedURL.Path))
		if err != nil {
			return "", err
		}

		return parsedURL.String(), nil
	case "https":
		return parsedURL.String(), nil
	default:
		return "", fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}
}

// gitSshUrlParser parses the url as a SSH URL
// git@github.com:guionardo/gs-dev.git -> https://github.com/guionardo/gs-dev
func gitSshUrlParser(giturl string) (string, error) {
	matches := gitSshUrlRe.FindStringSubmatch(giturl)
	if len(matches) == 0 {
		return "", fmt.Errorf("url is not a SSH URL: %s", giturl)
	}

	domain := matches[1]
	owner := matches[2]
	repo := matches[3]

	return fmt.Sprintf("https://%s/%s/%s", domain, owner, repo), nil
}

// ssh://
