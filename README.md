# tools

A command-line bookmark manager for terminal tools. Store and retrieve CLI tools with their usage examples - each example has a description and command.

## Features

- Interactive TUI for browsing and selecting commands
- Add, list, and remove tool bookmarks
- Store examples with descriptions for easy discovery
- Examples are automatically copied to clipboard when selected
- YAML-based persistent storage following XDG Base Directory specification
- Clean architecture with separated layers (domain, repository, service, CLI)
- Easy to extend with database backends or REST API

## Installation

### From Source

Requires Go 1.25 or later.

```bash
git clone https://github.com/fgeck/tools.git
cd tools
go build -o tools ./cmd/tools
```

Move the binary to a directory in your PATH:

```bash
sudo mv tools /usr/local/bin/
```


### Docker

Build the image:

```bash
docker build -t tools .
```

Run with a volume to persist data:

```bash
docker run -v ~/.config/tools:/root/.config/tools tools list
```

Or use an alias for convenience:

```bash
alias tools='docker run -v ~/.config/tools:/root/.config/tools tools'
```

## Usage

### Interactive TUI Mode (Default)

By default, running `tools` launches an interactive terminal UI:

```bash
tools
```

The TUI displays all your tools' examples in a table format. Each row shows:
- `[Tool Name]` - The tool this example belongs to
- `Description` - What the example does
- The actual command is shown in the item description

**Keyboard shortcuts:**
- `↑/↓` or `j/k` - Navigate examples
- `Enter` - Select command (copies to clipboard)
- `/` - Filter/search examples
- `a` - Add new tool with example
- `d` - Delete tool
- `Esc` - Quit

**Example display format:**
```
[lsof] list all ports at xxx
  → lsof -i :54321
```

**How it works:**

When you select an example and press Enter, the command is automatically copied to your clipboard using OSC 52 (supported by most modern terminals). Simply paste it with Ctrl+V (or Cmd+V) and execute.

### Classic CLI Mode

Use the `--cli` flag for traditional command-line mode:

```bash
tools --cli          # List all tools in table format
tools list --cli     # Same as above
```

### Add a new tool (CLI mode)

```bash
tools add -n <name> -c <command> [-d <description>] [-e "description|command"]
```

Examples must be in the format `"description|command"`.

Example:

```bash
tools add -n lsof \
  -c /usr/bin/lsof \
  -d "List open files and network connections" \
  -e "list all ports at xxx|lsof -i :54321" \
  -e "show all network connections|lsof -i"
```

Or use the TUI mode and press `a` to add interactively (easier for multiple fields).

### Remove a tool (CLI mode)

```bash
tools remove -n <name>
```

Or use the TUI mode, select a tool, and press `d` to delete.

### Get help

```bash
tools --help
tools <command> --help
```

## Command Aliases

- `add` can be shortened to `a`
- `list` can be shortened to `l` or `ls`
- `remove` can be shortened to `rm`

## Storage

Tool bookmarks are stored in `~/.config/tools/tools.yaml` by default. The location follows the XDG Base Directory specification and can be overridden by setting the `XDG_CONFIG_HOME` environment variable.

## Development

### Building

```bash
make build
```

### Testing

Run all tests:

```bash
make test
```

Run only unit tests:

```bash
make unit-test
```

Run only integration tests:

```bash
make integration-test
```

Generate coverage report:

```bash
make coverage
```

### Code Quality

Format code:

```bash
make fmt
```

Run linter:

```bash
make lint
```

### Other Commands

```bash
make help          # Show all available targets
make clean         # Remove build artifacts
make install       # Install to GOPATH/bin
make deps          # Download dependencies
make docker-build  # Build Docker image
make all           # Run all checks and build
```

## License

See LICENSE file for details.
