package integration

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/guionardo/gs-dev/internal/config"
	pad_client "github.com/guionardo/gs-dev/internal/pad/client"
	pad_server "github.com/guionardo/gs-dev/internal/pad/server"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"

	"github.com/guionardo/gs-dev/internal/pad/storage"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

const defaultTimeout = time.Second * 5

func TestPadIntegration(t *testing.T) { // nolint:funlen
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	logger := slog.New(slog.NewTextHandler(t.Output(), nil))
	configRoot, _ := config.NewConfigRoot(t.TempDir())
	storageDirectory := t.TempDir()
	storageConfig := storage.FileSystemStorageConfig{
		Directory: storageDirectory,
	}

	port, err := GetFreePort()
	require.NoError(t, err, "there is no available port")
	config.SetValue(configRoot, &padservice.PadServerConfig{
		Port:          port,
		StorageConfig: storageConfig,
	})
	require.NoError(t, configRoot.Save())

	storage, err := storage.NewFileSystemStorage(storageConfig, ctx, logger)
	require.NoError(t, err, "unexpected error on NewFileSystemStorage")

	service, err := padservice.NewPadServerService(configRoot, storage)
	require.NoError(t, err, "unexpected error on NewPadServerService")

	padConfig := padservice.NewPadServerConfig()

	server := pad_server.NewPadGrpcServer(padConfig, service, logger)

	go server.Start(ctx) // nolint:errcheck

	require.NoError(t, server.WaitUntilRunning(defaultTimeout), "server did not start within the expected time")

	client, err := pad_client.NewPadClient(server.URL(), "")
	require.NoError(t, err, "failed to create pad client")

	wg := sync.WaitGroup{}
	wg.Add(2)

	t.Run("Create, Get, and Expire Pad", func(t *testing.T) { // nolint:paralleltest
		require.NoError(t, server.WaitUntilRunning(defaultTimeout), "server did not start within the expected time")

		postID, err := client.Post([]byte("test"), time.Second*1, map[string]string{})
		require.NoError(t, err)
		require.NotEmpty(t, postID)

		content, _, err := client.Get(postID) // TODO: Check headers
		require.NoError(t, err)
		require.Equal(t, "test", string(content))

		time.Sleep(time.Second * 3)

		content, _, err = client.Get(postID) // TODO: Check headers
		require.Error(t, err)
		require.Empty(t, content)

		wg.Done()
	})

	t.Run("Run multiple concurrent requests", func(t *testing.T) { // nolint:paralleltest
		require.NoError(t, server.WaitUntilRunning(defaultTimeout), "server did not start within the expected time")

		type result struct {
			postID string
			err    error
		}

		const numRequests = 100

		results := make(chan result, numRequests)

		postEG := errgroup.Group{}
		postEG.SetLimit(runtime.NumCPU())

		startTime := time.Now()

		for range numRequests {
			postEG.Go(func() error {
				payload := []byte(rand.Text())

				postID, err := client.Post(payload, time.Hour, map[string]string{})
				results <- result{postID: postID, err: err}

				return err
			})
		}

		require.NoError(t, postEG.Wait(), "failed to complete all Post requests")

		postTime := time.Since(startTime)

		close(results)

		getEG := errgroup.Group{}
		getEG.SetLimit(runtime.NumCPU())

		startTime = time.Now()

		for result := range results {
			getEG.Go(func() error {
				content, _, err := client.Get(result.postID) // TODO: Check headers
				if err == nil && len(content) == 0 {
					err = fmt.Errorf("expected non-empty content for postID %s", result.postID)
				}

				return err
			})
		}

		require.NoError(t, getEG.Wait(), "failed to complete all Get requests")

		getTime := time.Since(startTime)

		t.Logf("Completed %d concurrent Post requests in %s: %f posts/s", numRequests, postTime, float64(numRequests)/postTime.Seconds())
		t.Logf("Completed %d concurrent Get requests in %s: %f gets/s", numRequests, getTime, float64(numRequests)/getTime.Seconds())
		wg.Done()
	})
	wg.Wait()
	cancel()
	time.Sleep(time.Second * 3)
	require.False(t, server.IsRunning())
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}

	return
}
