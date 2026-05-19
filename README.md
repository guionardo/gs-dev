# gs-dev

Guiosoft Development Assistant

[![Go Version](https://img.shields.io/github/go-mod/go-version/guionardo/gs-dev)](https://go.dev/)
[![License](https://img.shields.io/github/license/guionardo/gs-dev)](./LICENSE)
[![CodeQL](https://github.com/guionardo/gs-dev/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/guionardo/gs-dev/actions/workflows/github-code-scanning/codeql)
[![Go Report Card](https://goreportcard.com/badge/github.com/guionardo/gs-dev)](https://goreportcard.com/report/github.com/guionardo/gs-dev)
[![Go Release](https://github.com/guionardo/gs-dev/actions/workflows/release.yml/badge.svg)](https://github.com/guionardo/gs-dev/actions/workflows/release.yml)

A CLI development assistant providing rapid folder access, git statistics, URL management, interactive TUI, shell integration, and an ephemeral gRPC pad service for sharing text snippets.

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

### Commands

| Command | Init Alias | Output Redirect | TUI Support | Description |
|---------|-----------|-----------------|-------------|-------------|
| `dev` | ✓ | ✓ | Interactive subcommand picker | Rapid dev folder access — scan roots, detect projects, find folders |
| `fav` | ✓ | ✓ | Interactive favorite picker | Quick access to most frequently chosen folders |
| `pad` | ✓ | ✓ | — | gRPC pad client — post, get, and delete text snippets |
| `pad-server` | — | — | — | Start the ephemeral gRPC pad server |
| `git-stats` | — | — | ✓ (current dir) | Show git commit statistics (files, insertions, deletions) |
| `url` | — | — | ✓ (confirm/open) | Open or show git remote URL in browser |
| `install` | — | — | ✓ (confirm/install) | Install/uninstall shell bindings into `.bashrc`/`.zshrc` |
| `init` | — | — | Info message | Print shell initialization script (`source <(gs-dev init)`) |
| `setup` | — | — | ✓ (config editor) | Edit application configuration interactively |
| `version` | — | — | — | Show application version (`--full` for build info) |

Commands with **init alias** get shell function wrappers so they can be invoked directly. Commands with **output redirect** write `cd` commands to a temp file that the shell sources, enabling directory changes from the tool.

### Interactive TUI

Run `gs-dev` with no arguments — an interactive command chooser appears. For richer interaction, use `tui` (built from `cmd/tui`):

```bash
go build ./cmd/tui
./tui
```

The TUI wraps all commands with `rivo/tview` tree views and form panels.

### Shell Integration

```bash
# Install shell bindings (adds source <(gs-dev init) to ~/.bashrc/~/.zshrc)
gs-dev install

# Uninstall
gs-dev install --uninstall
```

Once installed, commands with init aliases (`dev`, `fav`, `pad`) can be called directly as shell functions, enabling `cd` into found directories.

### Dev Folder Access

Scan development folders, detect project types, and navigate instantly:

```bash
gs-dev add /path/to/projects   # Add a root folder
gs-dev sync                    # Rescan all roots for projects
gs-dev list                    # List roots and detected projects
gs-dev find query              # Find folders matching words (with interactive picker)
```

Detects **Go**, **Python**, **JavaScript**, **Rust**, **Java**, **PHP**, and **.NET** projects by manifest files.

### gRPC Pad Service

Ephemeral text snippet sharing over gRPC:

```bash
# Start the server
gs-dev pad-server serve

# Create a pad (CLI client with flate compression)
echo "hello world" | gs-dev pad post --stdin

# Retrieve a pad
gs-dev pad get <post-id>

# Delete a pad
gs-dev pad del <post-id>
```

Features: auto-expiry via configurable TTL, API key authentication via metadata headers, gRPC health checks, filesystem-backed storage with background expiration monitoring, base-62 encoded post IDs.

### Git Statistics

```bash
gs-dev git-stats --root /path/to/repo --since "2024-01-01" --author "name"
```

Parses `git log --shortstat` into structured commit summaries.

### URL Management

```bash
gs-dev url              # Open git remote URL in browser
gs-dev url --just-show  # Just display the URL without opening
```

Supports SSH and HTTPS remote URLs.

## Architecture

gs-dev uses a **layered** architecture: commands → services → storage/transport.

### CLI Generation via Reflection

Commands are defined as tagged structs — flags, args, descriptions, and validation are auto-mapped to Cobra via `internal/cli`:

```go
type MyCommand struct {
    Name string `flag:"name,n" description:"..." validate:"required"`
}
```

See [internal/cli/README.md](internal/cli/README.md).

### Directory Structure

```
cmd/
  gs-dev/          # Main CLI entry point
  tui/             # Interactive TUI entry point
  docs/            # Documentation generator
internal/
  cli/             # Reflection-based Cobra command builder
  commands/        # Command implementations (dev, fav, git_stats, install, pad, setup, shell_init, url)
  config/          # Typed YAML config management
  services/        # Business logic (dev, git, install, config, url, init, pad)
  pad/             # gRPC pad subsystem (proto, server, client, storage, post_id)
  dialog/          # Interactive terminal dialogs
  logging/         # Structured slog logger
  errors/          # Typed error wrapper (recoverable/unrecoverable)
  interfaces/      # Command and storage interfaces
  compression/     # Flate compression for pad content
  context/         # Command context data injection
  shell/           # Shell profile file management
  git/             # Git config parsing, remote URL resolution
  find_pattern/    # Path matching utilities
pkg/
  project_detector/ # Filesystem project detection (7 languages + fallback)
  cobra_tui/       # Tview-based TUI wrapper for Cobra commands
  console/         # Terminal colors and styled output
  open_url/        # Cross-platform URL opening
  post_command/    # Post-command output redirection to shell
  plugins/         # Plugin system interfaces
  tools/           # File utilities and helpers
```

### Configuration

Per-type YAML files stored in `~/.config/gs-dev/` (overridable via `GS_DEV_CONFIG_DIR`):

- `dev.yml` / `roots.yml` — Folder scanning settings, sync interval, favorites
- `pad_server.yml` — Server port, API key, storage directory, TTL
- `pad_client.yml` — Backend URL, API key, enabled status

Types implement the `ConfigType` interface with `Key()`, `Defaults()`, `Validate()`.

## Build & Development

```bash
make setup     # Install dev tooling (golangci-lint, gocritic, gocyclo, goimports)
make lint      # Run golangci-lint --fix
make test      # Run all tests
make build     # Build via .github/scripts/build.sh
make install   # Build and install
make proto     # Regenerate gRPC stubs from internal/pad/proto/pad.proto
make docs      # Regenerate CLI docs to docs/
make version   # Show current git tag version
make release   # Generate new GitHub release
```

For full CLI documentation, see [docs/gs-dev.md](docs/gs-dev.md).

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting a PR.

## License

MIT License — see [LICENSE](LICENSE) for details.
