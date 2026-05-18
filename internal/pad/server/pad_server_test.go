package pad_server_test

import (
	"context"
	"log/slog"
	"net"
	"testing"
	"time"

	paderrors "github.com/guionardo/gs-dev/internal/errors"
	pb "github.com/guionardo/gs-dev/internal/pad/proto"
	padserver "github.com/guionardo/gs-dev/internal/pad/server"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/stretchr/testify/require"
)

const bufSize = 1024 * 1024

type mockPadService struct {
	posts    map[string]mockPost
	validKey string
	postFunc func(body []byte, ttl time.Duration, metadata map[string]string) (string, error)
}

type mockPost struct {
	body    []byte
	headers map[string]string
	expires time.Time
}

func newMockPadService(validKey string) *mockPadService {
	return &mockPadService{
		posts:    make(map[string]mockPost),
		validKey: validKey,
	}
}

func (m *mockPadService) Post(body []byte, ttl time.Duration, metadata map[string]string) (string, error) {
	if m.postFunc != nil {
		return m.postFunc(body, ttl, metadata)
	}

	id := "mock-" + time.Now().Format("150405.000000000")

	var expires time.Time
	if ttl > 0 {
		expires = time.Now().Add(ttl)
	}

	m.posts[id] = mockPost{
		body:    body,
		headers: metadata,
		expires: expires,
	}

	return id, nil
}

func (m *mockPadService) Get(padID string) ([]byte, map[string]string, error) {
	p, ok := m.posts[padID]
	if !ok {
		return nil, nil, paderrors.ErrPostIDNotFound
	}

	if !p.expires.IsZero() && time.Now().After(p.expires) {
		return nil, nil, paderrors.ErrPostIDExpired
	}

	return p.body, p.headers, nil
}

func (m *mockPadService) Delete(padID string) error {
	if _, ok := m.posts[padID]; !ok {
		return paderrors.ErrPostIDNotFound
	}

	delete(m.posts, padID)

	return nil
}

func (m *mockPadService) IsAPIKeyValid(apiKey string) bool {
	return m.validKey == "" || m.validKey == apiKey
}

func setupTestServer(t *testing.T, cfg *padservice.PadServerConfig, svc padserver.PadService) (pb.PadServiceClient, context.Context) {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(t.Output(), nil))

	server := padserver.NewPadGrpcServer(cfg, svc, logger)

	listener := bufconn.Listen(bufSize)

	go func() {
		s := grpc.NewServer()
		pb.RegisterPadServiceServer(s, server)

		if err := s.Serve(listener); err != nil {
			t.Logf("bufconn server stopped: %v", err)
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = conn.Close() })

	client := pb.NewPadServiceClient(conn)
	ctx := context.Background()

	return client, ctx
}

func TestCreatePad_Success(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	client, ctx := setupTestServer(t, cfg, mock)

	resp, err := client.CreatePad(ctx, &pb.CreatePadRequest{
		Body:       []byte("hello"),
		TtlSeconds: 60,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)
	require.Equal(t, pb.PadStatus_ACTIVE, resp.Status)
	require.Positive(t, resp.ValidUntil)
}

func TestCreatePad_Unauthenticated(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	cfg.APIKey = "valid-key-uuid"
	mock := newMockPadService("valid-key-uuid")
	client, ctx := setupTestServer(t, cfg, mock)

	resp, err := client.CreatePad(ctx, &pb.CreatePadRequest{
		Body: []byte("hello"),
	})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestCreatePad_WithValidAPIKey(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	cfg.APIKey = "valid-key-uuid"
	mock := newMockPadService("valid-key-uuid")
	client, baseCtx := setupTestServer(t, cfg, mock)

	ctx := metadata.NewOutgoingContext(baseCtx, metadata.Pairs("x-api-key", "valid-key-uuid"))
	resp, err := client.CreatePad(ctx, &pb.CreatePadRequest{
		Body: []byte("hello"),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)
}

func TestGetPad_Success(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	client, ctx := setupTestServer(t, cfg, mock)

	created, err := client.CreatePad(ctx, &pb.CreatePadRequest{
		Body:       []byte("world"),
		TtlSeconds: 3600,
	})
	require.NoError(t, err)

	resp, err := client.GetPad(ctx, &pb.PadRequest{Id: created.Id})
	require.NoError(t, err)
	require.Equal(t, []byte("world"), resp.Body)
	require.Equal(t, pb.PadStatus_ACTIVE, resp.Status)
}

func TestGetPad_NotFound(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	client, ctx := setupTestServer(t, cfg, mock)

	resp, err := client.GetPad(ctx, &pb.PadRequest{Id: "nonexistent"})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestGetPad_Expired(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	client, ctx := setupTestServer(t, cfg, mock)

	created, err := client.CreatePad(ctx, &pb.CreatePadRequest{
		Body:       []byte("expires"),
		TtlSeconds: 0,
	})
	require.NoError(t, err)

	// Manually expire the post
	p := mock.posts[created.Id]
	p.expires = time.Now().Add(-time.Second)
	mock.posts[created.Id] = p

	resp, err := client.GetPad(ctx, &pb.PadRequest{Id: created.Id})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeletePad_Success(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	client, ctx := setupTestServer(t, cfg, mock)

	created, err := client.CreatePad(ctx, &pb.CreatePadRequest{Body: []byte("delete me")})
	require.NoError(t, err)

	resp, err := client.DeletePad(ctx, &pb.PadRequest{Id: created.Id})
	require.NoError(t, err)
	require.Equal(t, pb.PadStatus_DELETED, resp.Status)

	// Verify it's gone
	getResp, err := client.GetPad(ctx, &pb.PadRequest{Id: created.Id})
	require.Nil(t, getResp)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestDeletePad_NotFound(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	client, ctx := setupTestServer(t, cfg, mock)

	resp, err := client.DeletePad(ctx, &pb.PadRequest{Id: "nonexistent"})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestCreatePad_InternalError(t *testing.T) {
	t.Parallel()

	cfg := padservice.NewPadServerConfig()
	mock := newMockPadService("")
	mock.postFunc = func(body []byte, ttl time.Duration, metadata map[string]string) (string, error) {
		return "", paderrors.ErrPostIDInvalid
	}
	client, ctx := setupTestServer(t, cfg, mock)

	resp, err := client.CreatePad(ctx, &pb.CreatePadRequest{Body: []byte("oops")})
	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
}
