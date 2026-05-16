# gs-dev

Guiosoft Development Assistant

[![Go Version](https://img.shields.io/github/go-mod/go-version/guionardo/gs-dev)](https://go.dev/)
[![License](https://img.shields.io/github/license/guionardo/gs-dev)](./LICENSE)
[![CodeQL](https://github.com/guionardo/gs-dev/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/guionardo/gs-dev/actions/workflows/github-code-scanning/codeql)
[![Go Report Card](https://goreportcard.com/badge/github.com/guionardo/gs-dev)](https://goreportcard.com/report/github.com/guionardo/gs-dev)
[![Go Release](https://github.com/guionardo/gs-dev/actions/workflows/release.yml/badge.svg)](https://github.com/guionardo/gs-dev/actions/workflows/release.yml)

A CLI development assistant that provides rapid access to development folders, git statistics, URL management, and an ephemeral pad service for sharing text snippets via gRPC.

## Demo

```bash
# Start the pad server
gs-dev pad-server

# Create a pad
echo "hello world" | grpcurl -d '{"body": "aGVsbG8gd29ybGQ="}' localhost:8080 pad.PadService/CreatePad

# List available commands
gs-dev --help
```

## Getting Started

### Install

```bash
go install github.com/guionardo/gs-dev/cmd/gs-dev@latest
```

Or download a pre-built binary from [GitHub Releases](https://github.com/guionardo/gs-dev/releases/latest).

### Development Setup

```bash
make setup
```

## Features

- **Pad Service** — gRPC server and client for ephemeral key-value blobs with automatic expiry, gRPC health checks, and API key authentication via metadata headers
- **Dev Folder Access** — Rapid `cd` navigation into configured development directories
- **Favorites** — Quick access to recently used folders
- **Git Statistics** — View git commit statistics for your repositories
- **URL Management** — Open and manage URLs from the CLI
- **Shell Installation** — Install shell bindings and aliases

For full CLI documentation, see [docs/gs-dev.md](docs/gs-dev.md).

## Architecture

The pad service is organized into layered packages:

| Layer | Package | Purpose |
|-------|---------|---------|
| Proto | `internal/pad/proto` | Protobuf service definition and generated gRPC stubs |
| Server | `internal/pad/server` | gRPC server with interceptors, auth, health checks |
| Client | `internal/pad/client` | gRPC client with deadlines and metadata auth |
| Storage | `internal/pad/storage` | Filesystem-backed and dummy storage backends |
| Service | `internal/services/pad` | Business logic layer wrapping storage with compression |
| IDs | `internal/pad/post_id` | Base-62 encoded post identifiers |

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting a PR.

## License

MIT License — see [LICENSE](LICENSE) for details.
