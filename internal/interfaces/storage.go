package interfaces

import "time"

type Storage interface {
	Post(content []byte, ttl time.Duration, headers map[string]string) (postID string, err error)
	Get(postID string) (content []byte, headers map[string]string, err error)
	Delete(postID string) error
}
