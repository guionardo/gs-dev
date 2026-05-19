// Package storage provides filesystem-backed and dummy storage implementations
// for pad blobs. Posts are stored as directories containing metadata.yaml and
// content.raw files.
package storage

import (
	"context"
	"errors"

	"fmt"
	"log/slog"
	"maps"
	"os"
	"path"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/guionardo/gs-dev/internal/consts"
	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/pkg/tools/files"

	"github.com/guionardo/gs-dev/internal/interfaces"
	postid "github.com/guionardo/gs-dev/internal/pad/post_id"
	"go.yaml.in/yaml/v3"
	"golang.org/x/sync/singleflight"
)

// FileSystemStorage stores pads as directories on disk.
// Each pad gets a directory named by its post ID, containing metadata.yaml
// and content.raw. Expired pads are purged periodically by a background monitor.
type (
	FileSystemStorage struct {
		group              *singleflight.Group
		storeDirectory     string
		lock               sync.RWMutex
		removeExpiredPosts chan string
		logger             *slog.Logger
	}
	postMetadata struct {
		DataFolder   string            `yaml:"-"`
		CreatedAt    time.Time         `yaml:"created_at"`
		ValidUntil   time.Time         `yaml:"valid_until"`
		DataFilename string            `yaml:"-"`
		Headers      map[string]string `yaml:"headers"`
	}
)

const (
	metadataFilename        = "metadata.yaml"
	contentFilename         = "content.raw"
	StoreDirectoryConfigKey = "store_directory"
)

var (
	postIDOffset = uint64(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli())
	lastPostID   atomic.Uint64
)
var _ interfaces.PadStorage = &FileSystemStorage{}

// NewFileSystemStorage creates a FileSystemStorage bound to the directory
// configured in StorageConfig. It starts a background goroutine that purges
// expired posts every minute.
func NewFileSystemStorage(config FileSystemStorageConfig, ctx context.Context, logger *slog.Logger) (*FileSystemStorage, error) {
	storeDirectory, err := files.AssertDirectory(config.Directory)
	if err != nil {
		return nil, err
	}

	fs := &FileSystemStorage{group: new(singleflight.Group), storeDirectory: storeDirectory, removeExpiredPosts: make(chan string, 1000), logger: logger}
	logger.Info("FileSystemStorage created", "storeDirectory", storeDirectory)
	fs.startMonitor(ctx)

	return fs, nil
}

func (s *FileSystemStorage) startMonitor(ctx context.Context) {
	s.logger.Info("Starting monitor")

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Stopping monitor")
				return
			case postID := <-s.removeExpiredPosts:
				_ = s.Delete(postID)
			case <-ticker.C:
				_ = s.PurgeExpiredPosts()
			}
		}
	}()
}

func (s *FileSystemStorage) getNextPostID() (pID postid.PostID) {
	now := uint64(time.Now().UnixMilli() - int64(postIDOffset)) //nolint
	if lastPostID.Load() > now {
		now = lastPostID.Load() + 1
	}

	skips := 0

	for {
		pID = postid.NewPostID(now)
		if !s.postExists(pID) {
			break
		}

		now++
		skips++
	}

	s.logger.Info("Next post ID", "pID", pID.String(), "skips", skips)
	lastPostID.Store(now)

	return pID
}

func (s *FileSystemStorage) postExists(pID postid.PostID) bool {
	stat, err := os.Stat(filepath.Join(s.storeDirectory, pID.String()))
	return err == nil && stat.IsDir()
}

// Post stores content as a new pad with the given TTL and metadata headers.
// Returns the unique pad ID on success.
func (s *FileSystemStorage) Post(content []byte, ttl time.Duration, headers map[string]string) (postID string, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	pID := s.getNextPostID()

	metadata := &postMetadata{
		DataFolder:   filepath.Join(s.storeDirectory, pID.String()),
		CreatedAt:    time.Now(),
		DataFilename: filepath.Join(s.storeDirectory, pID.String(), contentFilename),
		Headers:      headers,
	}
	if ttl > 0 {
		metadata.ValidUntil = time.Now().Add(ttl)
	}

	if err = os.MkdirAll(metadata.DataFolder, consts.DirPermissions); err != nil {
		return "", fmt.Errorf("error creating data folder: %w", err)
	}

	defer func() {
		if err != nil {
			_ = os.RemoveAll(metadata.DataFolder)
		}
	}()

	var metadataContent []byte
	if metadataContent, err = yaml.Marshal(metadata); err == nil {
		if err = os.WriteFile(filepath.Join(metadata.DataFolder, metadataFilename), metadataContent, consts.FilesPermissions); err != nil {
			return "", fmt.Errorf("error writing metadata to file: %w", err)
		}
	} else {
		return "", fmt.Errorf("error marshalling metadata: %w", err)
	}

	if err = os.WriteFile(metadata.DataFilename, content, consts.FilesPermissions); err != nil {
		return "", fmt.Errorf("error writing content to file: %w", err)
	}

	return pID.String(), nil
}

// Get retrieves the content and headers of a pad by its post ID.
// Returns ErrPostIDNotFound or ErrPostIDExpired for missing or expired pads.
func (s *FileSystemStorage) Get(postID string) (content []byte, headers map[string]string, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if _, err := postid.ParsePostID(postID); err != nil {
		return nil, nil, err
	}

	key := "get:" + postID

	value, err, _ := s.group.Do(key, func() (any, error) {
		metadata, err := s.getPostMetadata(postID)
		if err != nil {
			return nil, err
		}

		return metadata, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if metadata, ok := value.(*postMetadata); ok {
		return metadata.Content()
	}

	return nil, nil, errs.NewError(nil, "post not found", true)
}

// Delete removes a pad by its post ID. Returns ErrPostIDNotFound if the pad
// does not exist or was already removed.
func (s *FileSystemStorage) Delete(postID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, err := postid.ParsePostID(postID); err != nil {
		return err
	}

	key := "delete:" + postID
	_, err, _ := s.group.Do(key, func() (any, error) {
		metadata, err := s.getPostMetadata(postID)
		if metadata != nil {
			if err = metadata.Delete(); err == nil {
				s.logger.Info("Post deleted", "postID", postID)
				return nil, nil
			}
		}

		if os.IsNotExist(err) || errors.Is(err, errs.ErrPostIDNotFound) {
			return nil, nil
		}

		s.logger.Error("Error deleting post", "postID", postID, "error", err)

		return nil, err
	})

	return err
}

func (s *FileSystemStorage) getFilePath(postID string) string {
	return filepath.Join(s.storeDirectory, postID)
}

func (s *FileSystemStorage) getPostMetadata(postID string) (metadata *postMetadata, err error) {
	metadata, err = readMetadata(s.getFilePath(postID))
	if err == nil {
		if metadata.ValidUntil.IsZero() || metadata.ValidUntil.After(time.Now()) {
			return metadata, nil
		}

		defer func() {
			s.removeExpiredPosts <- postID // delete expired post
		}()

		return metadata, errs.ErrPostIDExpired
	}

	return nil, errs.ErrPostIDNotFound
}

// PurgeExpiredPosts scans all stored pads and removes those past their TTL.
// It is called automatically every minute by the background monitor.
func (s *FileSystemStorage) PurgeExpiredPosts() error {
	_, err, _ := s.group.Do("purgeExpiredPosts", func() (any, error) {
		return nil, s.purgeExpiredPosts()
	})

	return err
}
func (s *FileSystemStorage) purgeExpiredPosts() error {
	files, err := os.ReadDir(s.storeDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if metadata, err := s.getPostMetadata(file.Name()); err == nil && !metadata.IsValid() {
			s.removeExpiredPosts <- file.Name()
		}
	}

	return nil
}

func readMetadata(filePath string) (metadata *postMetadata, err error) {
	var content []byte

	content, err = os.ReadFile(path.Join(filePath, metadataFilename)) // nolint:gosec
	if err == nil {
		metadata = new(postMetadata)

		err = yaml.Unmarshal(content, metadata)
		if err == nil {
			metadata.DataFolder = filePath
			metadata.DataFilename = path.Join(filePath, contentFilename)

			return metadata, nil
		}
	}

	return nil, err
}

func (p *postMetadata) Content() (content []byte, headers map[string]string, err error) {
	content, err = os.ReadFile(p.DataFilename)
	if err == nil {
		headers = maps.Clone(p.Headers)
		headers["valid-until"] = p.ValidUntil.Format(time.RFC1123Z)

		return content, headers, nil
	}

	return nil, nil, err
}

func (p *postMetadata) IsValid() bool {
	return p.ValidUntil.IsZero() || p.ValidUntil.After(time.Now())
}

func (p *postMetadata) Delete() error {
	return os.RemoveAll(p.DataFolder)
}
