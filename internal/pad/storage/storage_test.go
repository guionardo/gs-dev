package storage_test

import (
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/guionardo/gs-dev/internal/pad/storage"
	"github.com/stretchr/testify/require"
)

func TestFileSystemStorage(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	ctx := t.Context()
	logger := slog.New(slog.NewTextHandler(t.Output(), nil)).With("test", "FileSystemStorage")
	config := storage.FileSystemStorageConfig{
		Directory:  folder,
		DefaultTTL: time.Hour,
	}

	fs, err := storage.NewFileSystemStorage(config, ctx, logger)
	require.NoError(t, err)

	const postCount = 1000

	postIDs := make([]string, 0, postCount)

	t.Run("01_Post_should_return_post_id", func(t *testing.T) { //nolint:paralleltest
		for i := range postCount {
			postID, err := fs.Post([]byte("test"), time.Second*2, map[string]string{"name": "test", "id": strconv.Itoa(i)})
			require.NoError(t, err)

			postIDs = append(postIDs, postID)
		}

		require.Len(t, postIDs, postCount)
		t.Logf("Posted %d posts", postCount)
	})
	t.Run("02_Get_should_return_post_content", func(t *testing.T) { //nolint:paralleltest
		for index, postID := range postIDs {
			content, headers, err := fs.Get(postID)
			require.NoErrorf(t, err, "error getting post #%d: %s", index, postID)
			require.Equal(t, "test", string(content))
			require.Len(t, headers, 3)
		}
	})

	t.Run("03_Get_should_return_error_if_post_is_expired", func(t *testing.T) { //nolint:paralleltest
		t.Logf("Sleeping until all posts are expired")
		time.Sleep(time.Second * 3)

		for index, postID := range postIDs {
			content, headers, err := fs.Get(postID)
			require.Error(t, err, "expected error getting post #%d: %s", index, postID)
			require.Nil(t, content)
			require.Nil(t, headers)
		}
	})
}
