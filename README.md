# Tools

A command-line bookmark manager. Store and retrieve CLI commands with descriptions grouped by their root command.

## Features

- **Interactive TUI** for browsing and selecting commands
- **Add, edit, list, and remove** command bookmarks
- **Multiple bookmarks per tool** - group related commands by tool name
- **Auto-copy to clipboard** when selecting in TUI
- **YAML-based storage** following XDG Base Directory specification

## Demo

![Demo](demos/demo.gif)

## Installation

### From Source

Requires Go 1.21 or later.

```bash
git clone https://github.com/fgeck/tools.git
cd tools
make build
sudo mv tools /usr/local/bin/
```

### Via Homebrew (Coming Soon)

```bash
brew tap fgeck/tools
brew install tools
```

### Docker

Pull the latest image:

```bash
docker pull ghcr.io/fgeck/tools:latest
```

Run with a volume to persist data:

```bash
docker run -v ~/.config/tools:/config ghcr.io/fgeck/tools:latest list --cli
```

Or build locally:

```bash
docker build -t tools .
docker run -v ~/.config/tools:/config tools list --cli
```

**Note**: The Docker image uses `scratch` for minimal size (~10MB). Config must be mounted to `/config`.

## Usage

### Interactive TUI Mode (Default)

Simply run `tools` to launch the interactive terminal UI:

```bash
tools
```

**Keyboard shortcuts:**
- `↑/↓` - Navigate bookmarks
- `Enter` - Select command (copies to clipboard and prints to stdout)
- `a` - Add new bookmark
- `e` - Edit selected bookmark
- `d` - Delete selected bookmark
- `q/Esc` - Quit

When you select a bookmark with Enter, the command is:
1. Copied to clipboard using OSC 52 (supported by most modern terminals)
2. Printed to stdout

### CLI Commands

#### Add Bookmark

```bash
tools add -n <tool-name> -c <command> -d <description>
```

Example:
```bash
tools add -n lsof -c "lsof -i :8080" -d "check port 8080"
```

#### List Bookmarks

```bash
tools list --cli
# or
tools --cli
```

#### Edit Bookmark

Edit by specifying the command (primary key) and the fields to update:

```bash
tools edit -c <current-command> [--new-tool <name>] [--new-description <desc>] [--new-command <cmd>]
```

Examples:
```bash
# Change description
tools edit -c "lsof -i :8080" -d "check if port 8080 is in use"

# Change command itself
tools edit -c "lsof -i :8080" -n "lsof -t -i :8080"

# Change multiple fields
tools edit -c "lsof -i :8080" -t "lsof" -d "new description" -n "new command"
```

#### Remove Bookmark(s)

Remove specific bookmark by command:
```bash
tools rm -c <command>
```

Remove all bookmarks for a tool:
```bash
tools rm -n <tool-name>
```

Examples:
```bash
# Remove specific bookmark
tools rm -c "lsof -i :8080"

# Remove all lsof bookmarks
tools rm -n lsof
```

#### Get Help

```bash
tools --help
tools <command> --help
```

## Command Aliases

- `add` → `a`
- `list` → `l`
- `remove` → `rm`, `delete`
- `edit` → `e`, `update`

## Storage

Bookmarks are stored in `~/.config/tools/tools.yaml` by default.

The location follows XDG Base Directory specification and can be overridden with `XDG_CONFIG_HOME`:

```bash
export XDG_CONFIG_HOME=/custom/path
```

## Example Workflow

```bash
# Add some bookmarks
tools add -n kubectl -c "kubectl get pods" -d "list all pods"
tools add -n kubectl -c "kubectl get nodes" -d "list all nodes"
tools add -n docker -c "docker ps -a" -d "list all containers"

# Browse in TUI
tools
# Navigate with arrows, press Enter to copy command

# List in CLI
tools list --cli

# Edit a bookmark
tools edit -c "docker ps -a" -d "list all containers including stopped"

# Remove specific bookmark
tools rm -c "kubectl get nodes"

# Remove all kubectl bookmarks
tools rm -n kubectl
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test              # Run all tests
make unit-test         # Unit tests only
make integration-test  # Integration tests only
make coverage          # Generate coverage report
```

### Code Quality

```bash
make pre-commit        # Run all checks (tidy, fmt, vet, lint)
make install-hooks     # Install pre-commit Git hook
```

The pre-commit hook automatically runs before each commit:
- `go mod tidy`
- `go fmt`
- `go vet`
- `golangci-lint`

### Other Commands

```bash
make help              # Show all available targets
make clean             # Remove build artifacts
make install           # Install to GOPATH/bin
make deps              # Download dependencies
make docker-build      # Build Docker image
make all               # Run all checks and build
```

## Architecture

The project follows Clean Architecture principles:

```
internal/
├── cli/           # CLI commands (Cobra)
├── config/        # Configuration management
├── domain/models/ # Domain entities (Bookmark)
├── dto/           # Data transfer objects
├── repository/    # Data access layer (interface + YAML impl)
├── service/       # Business logic
└── tui/           # Terminal UI (Bubble Tea)
```

**Key Design:**
- **Repository pattern** - Storage abstraction (easy to swap YAML → PostgreSQL)
- **Service layer** - Business logic (reusable for REST API)
- **Command as primary key** - Each command string is unique
- **Tool name for grouping** - Multiple bookmarks per tool

## License

See LICENSE file for details.
