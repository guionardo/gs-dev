package urlservice

import (
	"fmt"
	"net/url"

	gitservice "github.com/guionardo/gs-dev/internal/services/git"
	openurl "github.com/guionardo/gs-dev/pkg/open_url"
)

type URLService struct {
	gitService *gitservice.GitService
}

func NewURLService() *URLService {
	return &URLService{
		gitService: gitservice.NewGitService(),
	}
}

func (u *URLService) OpenURL(url string) error {
	return openurl.OpenURLInBrowser(url)
}

func (u *URLService) GetRemoteHttpURL(repositoryRoot string) (string, error) {
	return u.gitService.GetRemoteHttpURL(repositoryRoot)
}

func (u *URLService) OpenRemoteURL(repositoryRoot string) error {
	remoteURL, err := u.GetRemoteHttpURL(repositoryRoot)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return fmt.Errorf("error parsing remote URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s - %s", parsedURL.Scheme, remoteURL)
	}

	return u.OpenURL(remoteURL)
}
