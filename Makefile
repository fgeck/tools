.PHONY: help fmt lint build test unit-test integration-test clean install coverage

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=tools
BINARY_PATH=./cmd/tools

# Build variables
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

fmt: ## Run go fmt on all source files
	@echo "Running go fmt..."
	@go fmt ./...
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run ./...

build: fmt ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "Binary created: $(BINARY_NAME)"

test: unit-test integration-test ## Run all tests

unit-test: ## Run unit tests
	@echo "Running unit tests..."
	@go test -tags=unit -v -race -coverprofile=coverage-unit.out ./...

integration-test: ## Run integration tests
	@echo "Running integration tests..."
	@go test -tags=integration -v -race -coverprofile=coverage-integration.out ./...

coverage: ## Generate and display test coverage report
	@echo "Running tests with coverage..."
	@go test -tags=unit -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out

clean: ## Remove build artifacts and test files
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage*.out coverage.html
	@go clean -testcache

install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(BINARY_PATH)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "Docker image built: $(BINARY_NAME):$(VERSION)"

docker-run: docker-build ## Run the CLI in Docker
	@docker run -v ~/.config/tools:/root/.config/tools $(BINARY_NAME):latest

run: build ## Build and run the binary
	@./$(BINARY_NAME)

all: clean deps fmt lint test build ## Run all checks and build
