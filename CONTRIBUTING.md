# Contributing to Tools CLI

Thank you for your interest in contributing to the Tools CLI project!

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- golangci-lint (optional but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/fgeck/tools.git
cd tools

# Install dependencies
make deps

# Install pre-commit hooks
make install-hooks

# Build the project
make build
```

## Development Workflow

### Pre-commit Hooks

This project uses Git pre-commit hooks to ensure code quality. The hooks automatically run:

1. **go mod tidy** - Ensures dependencies are properly managed
2. **go fmt** - Formats code according to Go standards
3. **go vet** - Checks for common Go errors
4. **golangci-lint** - Runs comprehensive linting (if installed)

#### Installing the Hook

```bash
make install-hooks
```

#### Running Checks Manually

```bash
# Run all pre-commit checks
make pre-commit

# Or run individual checks
make tidy
make fmt
make vet
make lint
```

#### Bypassing the Hook

If you need to commit without running the checks (not recommended):

```bash
git commit --no-verify -m "Your commit message"
```

## Makefile Commands

```bash
make help          # Display all available commands
make deps          # Download and tidy dependencies
make tidy          # Run go mod tidy
make fmt           # Format code with go fmt
make vet           # Run go vet
make lint          # Run golangci-lint
make pre-commit    # Run all pre-commit checks
make build         # Build the binary
make test          # Run all tests
make unit-test     # Run unit tests only
make integration-test  # Run integration tests only
make coverage      # Generate test coverage report
make clean         # Remove build artifacts
make install       # Install binary to GOPATH/bin
make install-hooks # Install Git pre-commit hook
make all           # Run all checks and build
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run only unit tests
make unit-test

# Run only integration tests
make integration-test

# Generate coverage report
make coverage
```

### Writing Tests

- Unit tests: Use `//go:build unit` build tag
- Integration tests: Use `//go:build integration` build tag
- Place tests next to the code they test
- Follow table-driven test patterns where appropriate

## Code Style

### Go Formatting

All code must be formatted with `go fmt`. The pre-commit hook will automatically check this.

### Linting

Run `golangci-lint` before submitting:

```bash
make lint
```

Install golangci-lint:
```bash
brew install golangci-lint
# Or: https://golangci-lint.run/usage/install/
```

### Code Organization

```
internal/
├── cli/           # CLI commands (Cobra)
├── config/        # Configuration management
├── domain/models/ # Domain entities
├── dto/           # Data transfer objects
├── repository/    # Data access layer
├── service/       # Business logic
└── tui/          # Terminal UI (Bubble Tea)
```

## Commit Messages

Follow conventional commits:

```
<type>: <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```bash
git commit -m "feat: add export command for examples"
git commit -m "fix: resolve TUI input validation issue"
git commit -m "docs: update README with new features"
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run all checks: `make pre-commit test`
5. Commit your changes (pre-commit hook will run automatically)
6. Push to your fork: `git push origin feature/my-feature`
7. Create a Pull Request

### PR Requirements

- [ ] All tests pass
- [ ] Pre-commit checks pass
- [ ] Code is properly formatted
- [ ] Tests added for new functionality
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventional commits

## Architecture

The project follows Clean Architecture principles:

- **Domain Layer**: Core business entities
- **Repository Layer**: Data persistence abstraction
- **Service Layer**: Business logic and use cases
- **CLI Layer**: Command-line interface
- **TUI Layer**: Terminal user interface

### Key Principles

- Dependencies point inward (CLI → Service → Repository → Domain)
- Repository pattern for storage abstraction
- DTOs for API-agnostic communication
- Interface-based design for testability

## Questions?

- Open an issue for bugs or feature requests
- Check existing issues before opening a new one
- Use discussions for questions

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.
