package interfaces

import "time"

type PadStorage interface {
	Post(content []byte, ttl time.Duration, headers map[string]string) (string, error)
	Get(postID string) (content []byte, headers map[string]string, err error)
	Delete(postID string) error
}
