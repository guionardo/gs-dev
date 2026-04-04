package urlparsers

import (
	"fmt"
	"net/url"
	"strings"
)

type URLParserFn func(url string) (string, error)

var parsers = map[string]URLParserFn{
	"gitHTTP":   gitHttpUrlParser,
	"gitSSH":    gitSshUrlParser,
	"azureHTTP": gitAzureHttpUrlParser,
	"azureSSH":  gitAzureSshUrlParser,
}

// Parse the url from origin from git config, returning the HTTP URL
func Parse(url string) (string, error) {
	for _, parser := range parsers {
		if parsedURL, err := parser(url); err == nil {
			return removeGitSuffix(removeUsernameFromUrl(parsedURL)), nil
		}
	}

	return "", fmt.Errorf("no parser found for url: %s", url)
}

func removeGitSuffix(url string) string {
	return strings.TrimSuffix(url, ".git")
}

func removeUsernameFromUrl(originalURL string) string {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return originalURL
	}

	parsedURL.User = nil

	return parsedURL.String()
}
