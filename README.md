# tools

A command-line bookmark manager for terminal commands. Store and retrieve CLI command examples with descriptions.

## Features

- **Interactive TUI** for browsing and selecting commands
- **Add, edit, list, and remove** command examples
- **Multiple examples per tool** - group related commands by tool name
- **Auto-copy to clipboard** when selecting in TUI
- **YAML-based storage** following XDG Base Directory specification
- **Clean architecture** - easy to extend with database backends or REST API

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

Build the image:

```bash
docker build -t tools .
```

Run with a volume to persist data:

```bash
docker run -v ~/.config/tools:/root/.config/tools tools list --cli
```

## Usage

### Interactive TUI Mode (Default)

Simply run `tools` to launch the interactive terminal UI:

```bash
tools
```

**Keyboard shortcuts:**
- `↑/↓` - Navigate examples
- `Enter` - Select command (copies to clipboard and prints to stdout)
- `a` - Add new example
- `e` - Edit selected example
- `d` - Delete selected example
- `q/Esc` - Quit

When you select an example with Enter, the command is:
1. Copied to clipboard using OSC 52 (supported by most modern terminals)
2. Printed to stdout

### CLI Commands

#### Add Example

```bash
tools add -n <tool-name> -c <command> -d <description>
```

Example:
```bash
tools add -n lsof -c "lsof -i :8080" -d "check port 8080"
```

#### List Examples

```bash
tools list --cli
# or
tools --cli
```

#### Edit Example

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

#### Remove Example(s)

Remove specific example by command:
```bash
tools rm -c <command>
```

Remove all examples for a tool:
```bash
tools rm -n <tool-name>
```

Examples:
```bash
# Remove specific example
tools rm -c "lsof -i :8080"

# Remove all lsof examples
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

Examples are stored in `~/.config/tools/tools.yaml` by default.

The location follows XDG Base Directory specification and can be overridden with `XDG_CONFIG_HOME`:

```bash
export XDG_CONFIG_HOME=/custom/path
```

## Example Workflow

```bash
# Add some examples
tools add -n kubectl -c "kubectl get pods" -d "list all pods"
tools add -n kubectl -c "kubectl get nodes" -d "list all nodes"
tools add -n docker -c "docker ps -a" -d "list all containers"

# Browse in TUI
tools
# Navigate with arrows, press Enter to copy command

# List in CLI
tools list --cli

# Edit an example
tools edit -c "docker ps -a" -d "list all containers including stopped"

# Remove specific example
tools rm -c "kubectl get nodes"

# Remove all kubectl examples
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
├── domain/models/ # Domain entities (ToolExample)
├── dto/           # Data transfer objects
├── repository/    # Data access layer (interface + YAML impl)
├── service/       # Business logic
└── tui/           # Terminal UI (Bubble Tea)
```

**Key Design:**
- **Repository pattern** - Storage abstraction (easy to swap YAML → PostgreSQL)
- **Service layer** - Business logic (reusable for REST API)
- **Command as primary key** - Each command string is unique
- **Tool name for grouping** - Multiple examples per tool

## License

See LICENSE file for details.
