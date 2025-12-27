# tools

A command-line bookmark manager for terminal tools. Store and retrieve information about CLI tools including their paths, descriptions, and usage examples.

## Features

- Add, list, and remove tool bookmarks
- Store tool metadata: name, command path, description, examples
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

### List all tools

```bash
tools
# or
tools list
```

### Add a new tool

```bash
tools add -n <name> -c <command> [-d <description>] [-e <example>]
```

Example:

```bash
tools add -n kubectl \
  -c /usr/local/bin/kubectl \
  -d "Kubernetes command-line tool" \
  -e "kubectl get pods" \
  -e "kubectl describe node"
```

### Remove a tool

```bash
tools remove -n <name>
```

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
