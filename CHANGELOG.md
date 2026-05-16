# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- gRPC recovery and logging interceptors on PadServer
- gRPC health check service for Kubernetes probes
- API key authentication via gRPC metadata headers (moved from request body)
- 30-second deadline on all client RPC calls
- bufconn-based unit tests for PadServer (9 test cases)
- Doc comments on all exported symbols in `internal/pad/`
- Package comments for all `internal/pad/` sub-packages
- CONTRIBUTING.md, CHANGELOG.md, llms.txt
- Enhanced README with architecture table

### Changed

- `NewPadClient` now returns `(*PadClient, error)` instead of panicking
- Graceful shutdown has a 15-second timeout before hard stop
- `GetPad` and `DeletePad` no longer return response + error simultaneously
- `codes.NotFound` for missing pads, `codes.Internal` only for unexpected bugs
- `codes.PermissionDenied` for bad API keys, `codes.Unauthenticated` for missing metadata

### Fixed

- Context.Background() replaced with context.WithTimeout on all client calls
- API key moved from protobuf request body to gRPC metadata
- Pad status no longer returned alongside gRPC errors (was dead data)

## [0.1.0] - 2026-01-01

### Added

- Initial release
- gs-dev CLI with dev, fav, git-stats, init, install, setup, url, version commands
- PadService gRPC server and client
- FileSystemStorage for pad persistence with TTL expiry
- Base-62 post ID generation
- Cobra CLI framework integration
- Auto-generated CLI documentation in `docs/`
- MIT License
