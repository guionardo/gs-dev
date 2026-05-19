// Package pad_client provides a gRPC client for the PadService.
// It handles connection lifecycle, API key authentication via metadata headers,
// and enforces a 30-second deadline on every RPC call.
package pad_client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/guionardo/gs-dev/internal/pad/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	defaultRPCTimeout = 30 * time.Second
	apiKeyMetadataKey = "x-api-key"
)

// PadClient is a gRPC client for the PadService.
// It reuses a single connection and sends the API key in gRPC metadata headers.
type PadClient struct {
	padServerURL string
	apiKey       string
	conn         *grpc.ClientConn
	grpcClient   pb.PadServiceClient
}

var defaultGrpcOptions = []grpc.CallOption{
	grpc.WaitForReady(true),
}

// NewPadClient creates a new PadClient connected to the given server URL.
// Authentication is performed by the server; callers should handle the error
// instead of panicking.
func NewPadClient(padServerURL string, apiKey string) (*PadClient, error) {
	conn, err := grpc.NewClient(padServerURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	grpcClient := pb.NewPadServiceClient(conn)

	return &PadClient{padServerURL: padServerURL, apiKey: apiKey, conn: conn, grpcClient: grpcClient}, nil
}

// Close closes the underlying gRPC connection.
func (p *PadClient) Close() {
	if p.conn != nil {
		_ = p.conn.Close()
	}
}

func (p *PadClient) newContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRPCTimeout)
	if p.apiKey != "" {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(apiKeyMetadataKey, p.apiKey))
	}

	return ctx, cancel
}

// Post creates a new pad with the given content, TTL, and metadata headers.
// Returns the assigned pad ID on success.
func (p *PadClient) Post(content []byte, ttl time.Duration, headers map[string]string) (string, error) {
	req := &pb.CreatePadRequest{
		Body:       content,
		TtlSeconds: uint32(ttl.Seconds()),
		Metadata:   headers,
	}

	ctx, cancel := p.newContext()
	defer cancel()

	resp, err := p.grpcClient.CreatePad(ctx, req, defaultGrpcOptions...)
	if err != nil {
		return "", fmt.Errorf("failed to create pad: %w", err)
	}

	return resp.GetId(), nil
}

// Get retrieves the body and metadata of a pad by ID.
// Returns NotFound if the pad does not exist or is expired.
func (p *PadClient) Get(id string) ([]byte, map[string]string, error) {
	req := &pb.PadRequest{
		Id: id,
	}

	ctx, cancel := p.newContext()
	defer cancel()

	resp, err := p.grpcClient.GetPad(ctx, req, defaultGrpcOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get pad: %w", err)
	}

	return resp.GetBody(), resp.GetMetadata(), nil
}

// Delete deletes a pad by ID. Returns nil on success or NotFound if the pad
// does not exist.
func (p *PadClient) Delete(id string) error {
	req := &pb.PadRequest{
		Id: id,
	}

	ctx, cancel := p.newContext()
	defer cancel()

	_, err := p.grpcClient.DeletePad(ctx, req, defaultGrpcOptions...)
	if err != nil {
		return fmt.Errorf("failed to delete pad: %w", err)
	}

	return nil
}
