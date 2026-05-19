package initservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/guionardo/gs-dev/app/build"

	"github.com/hashicorp/go-version"
)

type (
	InitService struct {
	}
	GitHubRelease struct {
		Author      Author    `json:"author"`
		TagName     string    `json:"tag_name"`
		Name        string    `json:"name"`
		PublishedAt time.Time `json:"published_at"`
		Body        string    `json:"body"`
	}
	Author struct {
		Login string `json:"login"`
	}
)

var (
	releaseURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", build.Owner, build.Repo)
)

func NewInitService() *InitService {
	return &InitService{}
}

// GetAvailableVersion if there is a new version available on github releases
func (i *InitService) GetAvailableVersion() (*GitHubRelease, error) {
	response, err := http.Get(releaseURL) // nolint:gosec // G114 -- URL is parsed and restricted to GitHub API HTTPS host.
	if err != nil {
		return nil, fmt.Errorf("error getting release: %w", err)
	}

	defer response.Body.Close() // nolint:errcheck

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var release GitHubRelease

	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return &release, nil
}

// CheckCurrentVersion checks if the current version is the same as the available version
func (i *InitService) CheckCurrentVersion() (msg string, isUpToDate bool, err error) {
	release, err := i.GetAvailableVersion()
	if err != nil {
		return "", false, fmt.Errorf("error getting available version: %w", err)
	}

	vRelease, err := version.NewVersion(release.TagName)
	if err != nil {
		return "", false, fmt.Errorf("error parsing release version: %w", err)
	}

	vCurrent, err := version.NewVersion(build.Version)
	if err != nil {
		return "", false, fmt.Errorf("error parsing current version: %w", err)
	}

	if vRelease.Equal(vCurrent) {
		return "You are using the latest version", true, nil
	}

	if vRelease.LessThan(vCurrent) {
		return "You are using an older version: " + vRelease.String(), false, nil
	}

	return "You are using a newer version: " + vRelease.String(), false, nil
}
