package padservice

import (
	"maps"
	"net/url"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/guionardo/gs-dev/internal/errors"
)

type (
	PadClientConfig struct {
		Enabled    bool      `yaml:"enabled" default:"true"`
		BackendURL string    `yaml:"backend_url" default:"http://localhost:8080"`
		APIKey     string    `yaml:"api_key" default:""`
		LastPosts  LastPosts `yaml:"last_posts"`
	}
	Post struct {
		ID         string    `yaml:"id"`
		CreatedAt  time.Time `yaml:"created_at"`
		ValidUntil time.Time `yaml:"valid_until"`
	}
	LastPosts map[string]Post
)

func (c *PadClientConfig) Defaults() {
	if _, err := url.Parse(c.BackendURL); err != nil {
		c.BackendURL = ""
		c.Enabled = false
	}

	if akUUID, err := uuid.Parse(c.APIKey); err != nil || akUUID == uuid.Nil {
		c.APIKey = ""
		c.Enabled = false
	}

	if len(c.LastPosts) == 0 {
		c.LastPosts = make(LastPosts)
	}
}

func (c PadClientConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	backendURL, err := url.Parse(c.BackendURL)
	if err != nil || backendURL.Scheme != "http" && backendURL.Scheme != "https" {
		return errors.NewError(nil, "invalid backend URL: %s", false, c.BackendURL)
	}

	// API Key can be empty, but if it's not, it should be a valid UUID
	if c.APIKey != "" {
		if akUUID, err := uuid.Parse(c.APIKey); err != nil || akUUID == uuid.Nil {
			return errors.NewError(nil, "invalid API key: %s", false, c.APIKey)
		}
	}

	return nil
}

func (c PadClientConfig) Key() string {
	return "pad_client"
}

func (lp LastPosts) Add(postId string, validUntil time.Time) {
	lp[postId] = Post{
		ID:         postId,
		CreatedAt:  time.Now(),
		ValidUntil: validUntil,
	}
}

func (lp LastPosts) Remove(postID string) {
	delete(lp, postID)
}

func (lp LastPosts) Purge() (changed bool) {
	for id, post := range lp {
		if !post.ValidUntil.IsZero() && post.ValidUntil.Before(time.Now()) {
			delete(lp, id)

			changed = true
		}
	}

	return changed
}

// IDs returns the
func (lp LastPosts) IDs() (ids []string) {
	_ = lp.Purge()
	ids = slices.Collect(maps.Keys(lp))
	slices.Sort(ids)

	return ids
}

func (lp LastPosts) Posts() (posts []Post) {
	_ = lp.Purge()
	posts = slices.Collect(maps.Values(lp))
	slices.SortFunc(posts, func(i, j Post) int {
		return int(i.CreatedAt.Sub(j.CreatedAt))
	})

	return posts
}
