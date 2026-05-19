package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/guionardo/gs-dev/app/build"
)

type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
}

var releaseURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", build.Owner, build.Repo)

func GetLatestRelease() (*GitHubRelease, error) {
	resp, err := http.Get(releaseURL) // nolint:gosec // G114: URL is restricted to GitHub API HTTPS
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close() // nolint:errcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var release GitHubRelease
	if err = json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if release.TagName == "" {
		return nil, errors.New("unexpected response: missing tag_name")
	}

	return &release, nil
}
