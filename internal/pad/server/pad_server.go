// Package pad_server implements the gRPC PadService server.
// It provides create, read, and delete operations for pads (ephemeral key-value blobs)
// over gRPC, with built-in recovery/logging interceptors, API key authentication via
// metadata headers, and a gRPC health check endpoint for Kubernetes probes.
package pad_server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"
	"time"

	paderrors "github.com/guionardo/gs-dev/internal/errors"
	pb "github.com/guionardo/gs-dev/internal/pad/proto"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	waitUntilReadyInterval = time.Millisecond * 100
	gracefulStopTimeout    = 15 * time.Second
	apiKeyMetadataKey      = "x-api-key"
)

// PadGrpcServer implements the PadService gRPC server.
// It wraps a PadService with gRPC transport, authentication, interceptors,
// health checking, and graceful shutdown.
type PadGrpcServer struct {
	pb.UnimplementedPadServiceServer

	config    *padservice.PadServerConfig
	service   PadService
	logger    *slog.Logger
	isRunning atomic.Bool
	health    *health.Server
}

// PadService is the business logic interface consumed by PadGrpcServer.
// It decouples gRPC transport from storage and authentication concerns.
type PadService interface {
	Post(body []byte, ttl time.Duration, metadata map[string]string) (padId string, err error)
	Get(padId string) (content []byte, headers map[string]string, err error)
	Delete(padId string) error
	IsAPIKeyValid(apiKey string) bool
}

// NewPadGrpcServer creates a new PadGrpcServer with the given config, service, and logger.
// The returned server is not running until Start is called.
func NewPadGrpcServer(config *padservice.PadServerConfig, service PadService, logger *slog.Logger) *PadGrpcServer {
	return &PadGrpcServer{
		config:  config,
		service: service,
		logger:  logger,
		health:  health.NewServer(),
	}
}

// Start starts the gRPC server on the configured port and blocks until the
// context is cancelled. On cancellation it performs a graceful shutdown with
// a 15-second timeout before hard-stopping.
func (p *PadGrpcServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.config.Port))
	if err != nil {
		return fmt.Errorf("failed to start Pad Server: %w", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			p.recoveryInterceptor,
			p.loggingInterceptor,
		),
	)
	pb.RegisterPadServiceServer(s, p)
	healthpb.RegisterHealthServer(s, p.health)

	p.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	p.health.SetServingStatus("pad.PadService", healthpb.HealthCheckResponse_SERVING)

	go func(ctx context.Context) {
		<-ctx.Done()
		p.isRunning.Store(false)
		p.health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

		p.logger.Info("gracefully shutting down Pad Server")

		stopped := make(chan struct{})

		go func() {
			s.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			p.logger.Info("Pad Server shut down gracefully")
		case <-time.After(gracefulStopTimeout):
			p.logger.Warn("graceful shutdown timed out, forcing stop")
			s.Stop()
		}
	}(ctx)

	p.logger.Info("starting Pad Server", "port", p.config.Port)
	p.isRunning.Store(true)

	if err = s.Serve(listener); err != nil {
		p.isRunning.Store(false)
		return fmt.Errorf("pad server serve error: %w", err)
	}

	return nil
}

// CreatePad handles the CreatePad RPC. It creates a new pad with the given body,
// TTL, and metadata. Returns the pad ID and active status on success.
func (p *PadGrpcServer) CreatePad(ctx context.Context, pr *pb.CreatePadRequest) (*pb.CreatedPadResponse, error) {
	if err := p.validateAuthentication(ctx); err != nil {
		return nil, err
	}

	ttl := time.Duration(pr.TtlSeconds) * time.Second

	padId, err := p.service.Post(pr.Body, ttl, pr.Metadata)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create pad: %v", err)
	}

	var validUntil time.Time
	if pr.TtlSeconds > 0 {
		validUntil = time.Now().Add(ttl)
	}

	return &pb.CreatedPadResponse{Id: padId, Status: pb.PadStatus_ACTIVE, ValidUntil: uint64(validUntil.Unix())}, nil
}

// GetPad handles the GetPad RPC. It retrieves the body and metadata of a pad by ID.
// Returns NotFound for missing or expired pads.
func (p *PadGrpcServer) GetPad(ctx context.Context, pr *pb.PadRequest) (*pb.PadResponse, error) {
	if err := p.validateAuthentication(ctx); err != nil {
		return nil, err
	}

	content, headers, err := p.service.Get(pr.GetId())
	if err != nil {
		if errors.Is(err, paderrors.ErrPostIDNotFound) {
			return nil, status.Errorf(codes.NotFound, "pad %q not found", pr.GetId())
		}

		if errors.Is(err, paderrors.ErrPostIDExpired) {
			return nil, status.Errorf(codes.NotFound, "pad %q expired", pr.GetId())
		}

		return nil, status.Errorf(codes.Internal, "get pad: %v", err)
	}

	return &pb.PadResponse{
		Body:     content,
		Metadata: headers,
		Status:   pb.PadStatus_ACTIVE,
	}, nil
}

// DeletePad handles the DeletePad RPC. It deletes a pad by ID.
// Returns NotFound if the pad does not exist.
func (p *PadGrpcServer) DeletePad(ctx context.Context, pr *pb.PadRequest) (*pb.PadResponse, error) {
	if err := p.validateAuthentication(ctx); err != nil {
		return nil, err
	}

	if err := p.service.Delete(pr.Id); err != nil {
		if errors.Is(err, paderrors.ErrPostIDNotFound) {
			return nil, status.Errorf(codes.NotFound, "pad %q not found", pr.Id)
		}

		return nil, status.Errorf(codes.Internal, "delete pad: %v", err)
	}

	return &pb.PadResponse{
		Status: pb.PadStatus_DELETED,
	}, nil
}

func (p *PadGrpcServer) validateAuthentication(ctx context.Context) error {
	if len(p.config.APIKey) == 0 {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	keys := md.Get(apiKeyMetadataKey)
	if len(keys) == 0 {
		return status.Error(codes.Unauthenticated, "missing API key")
	}

	if !p.service.IsAPIKeyValid(keys[0]) {
		return status.Error(codes.PermissionDenied, "invalid API key")
	}

	return nil
}

func (p *PadGrpcServer) recoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error("panic recovered in gRPC handler",
				"method", info.FullMethod,
				"panic", r,
			)
		}
	}()

	return handler(ctx, req)
}

func (p *PadGrpcServer) loggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	p.logger.Info("gRPC call",
		"method", info.FullMethod,
		"duration", duration,
		"code", status.Code(err).String(),
	)

	return resp, err
}

// URL returns the local address of this server.
func (p *PadGrpcServer) URL() string {
	return fmt.Sprintf("localhost:%d", p.config.Port)
}

// IsRunning returns true if the server is currently accepting connections.
func (p *PadGrpcServer) IsRunning() bool {
	return p.isRunning.Load()
}

// WaitUntilRunning blocks up to timeout waiting for the server to accept connections.
// Returns an error if the server does not become ready within the given duration.
func (p *PadGrpcServer) WaitUntilRunning(timeout time.Duration) error {
	waitUntil := time.Now().Add(timeout)
	for time.Now().Before(waitUntil) {
		if p.IsRunning() {
			return nil
		}

		time.Sleep(waitUntilReadyInterval)
	}

	return errors.New("server did not start within the expected time")
}
