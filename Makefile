.PHONY: help fmt lint build test unit-test integration-test clean install coverage pre-commit install-hooks release-patch release-minor release-major delete-release list-releases

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

tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	@go mod tidy
	@git diff --exit-code go.mod go.sum || (echo "go.mod or go.sum changed, please commit the changes" && exit 1)

fmt: ## Run go fmt on all source files
	@echo "Running go fmt..."
	@go fmt ./...
	@git diff --exit-code || (echo "Code not formatted, please commit the changes" && exit 1)

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@golangci-lint run ./...

pre-commit: tidy fmt vet lint ## Run all pre-commit checks (tidy, fmt, vet, lint)
	@echo "✅ All pre-commit checks passed!"

build: ## Build the binary
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

install-hooks: ## Install Git pre-commit hook
	@echo "Installing pre-commit hook..."
	@chmod +x scripts/pre-commit.sh
	@cp scripts/pre-commit.sh .git/hooks/pre-commit
	@echo "✅ Pre-commit hook installed successfully!"
	@echo "To skip the hook, use: git commit --no-verify"

demo: build ## Generate demo GIF using VHS
	@echo "Generating demo GIF..."
	@command -v vhs >/dev/null 2>&1 || { echo "Error: VHS not installed. Run: brew install vhs"; exit 1; }
	@vhs demos/demo.tape
	@echo "✅ Demo generated: demos/demo.gif"

all: clean deps pre-commit test build ## Run all checks and build

# Release management
# Semantic versioning for releases (vMAJOR.MINOR.PATCH)
# - PATCH: Backwards-compatible bug fixes (v1.2.3 -> v1.2.4)
# - MINOR: Backwards-compatible new features (v1.2.3 -> v1.3.0)
# - MAJOR: Breaking changes (v1.2.3 -> v2.0.0)
#
# Usage:
#   make release-patch  - Create a patch release
#   make release-minor  - Create a minor release
#   make release-major  - Create a major release
#   make list-releases  - List all existing releases
#   make delete-release TAG=v1.2.3 - Delete a specific release
#
# Get the current version from git tags, default to v0.0.0 if no tags exist
CURRENT_VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

# Function to increment semantic version
# Usage: $(call increment_version,major|minor|patch)
define get_next_version
$(shell echo $(CURRENT_VERSION) | sed 's/^v//' | awk -F. -v type=$(1) '\
	BEGIN { OFS="." } \
	{ \
		if (type == "major") { print $$1+1, 0, 0 } \
		else if (type == "minor") { print $$1, $$2+1, 0 } \
		else if (type == "patch") { print $$1, $$2, $$3+1 } \
	}')
endef

list-releases: ## List all releases/tags
	@echo "Current releases:"
	@git tag -l --sort=-v:refname | head -20

release-patch: ## Create a new patch release (e.g., v1.2.3 -> v1.2.4)
	@echo "Current version: $(CURRENT_VERSION)"
	$(eval NEW_VERSION := v$(call get_next_version,patch))
	@echo "Creating new patch release: $(NEW_VERSION)"
	@read -p "Create tag $(NEW_VERSION)? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a $(NEW_VERSION) -m "Release $(NEW_VERSION)"; \
		echo "✅ Tag $(NEW_VERSION) created locally"; \
		read -p "Push tag to remote? [y/N] " -n 1 -r; \
		echo; \
		if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
			git push origin $(NEW_VERSION); \
			echo "✅ Tag pushed to remote"; \
			echo "View at: https://github.com/$$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases/tag/$(NEW_VERSION)"; \
		else \
			echo "⚠️  Tag created locally but not pushed. Use 'git push origin $(NEW_VERSION)' to push later."; \
		fi; \
	else \
		echo "❌ Release cancelled"; \
	fi

release-minor: ## Create a new minor release (e.g., v1.2.3 -> v1.3.0)
	@echo "Current version: $(CURRENT_VERSION)"
	$(eval NEW_VERSION := v$(call get_next_version,minor))
	@echo "Creating new minor release: $(NEW_VERSION)"
	@read -p "Create tag $(NEW_VERSION)? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a $(NEW_VERSION) -m "Release $(NEW_VERSION)"; \
		echo "✅ Tag $(NEW_VERSION) created locally"; \
		read -p "Push tag to remote? [y/N] " -n 1 -r; \
		echo; \
		if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
			git push origin $(NEW_VERSION); \
			echo "✅ Tag pushed to remote"; \
			echo "View at: https://github.com/$$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases/tag/$(NEW_VERSION)"; \
		else \
			echo "⚠️  Tag created locally but not pushed. Use 'git push origin $(NEW_VERSION)' to push later."; \
		fi; \
	else \
		echo "❌ Release cancelled"; \
	fi

release-major: ## Create a new major release (e.g., v1.2.3 -> v2.0.0)
	@echo "Current version: $(CURRENT_VERSION)"
	$(eval NEW_VERSION := v$(call get_next_version,major))
	@echo "Creating new major release: $(NEW_VERSION)"
	@echo "⚠️  WARNING: Major version change indicates breaking changes!"
	@read -p "Create tag $(NEW_VERSION)? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a $(NEW_VERSION) -m "Release $(NEW_VERSION)"; \
		echo "✅ Tag $(NEW_VERSION) created locally"; \
		read -p "Push tag to remote? [y/N] " -n 1 -r; \
		echo; \
		if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
			git push origin $(NEW_VERSION); \
			echo "✅ Tag pushed to remote"; \
			echo "View at: https://github.com/$$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases/tag/$(NEW_VERSION)"; \
		else \
			echo "⚠️  Tag created locally but not pushed. Use 'git push origin $(NEW_VERSION)' to push later."; \
		fi; \
	else \
		echo "❌ Release cancelled"; \
	fi

delete-release: ## Delete a release/tag (usage: make delete-release TAG=v1.2.3)
	@if [ -z "$(TAG)" ]; then \
		echo "❌ Error: TAG parameter is required"; \
		echo "Usage: make delete-release TAG=v1.2.3"; \
		exit 1; \
	fi
	@echo "⚠️  WARNING: This will delete tag $(TAG)!"
	@read -p "Delete tag $(TAG) locally? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -d $(TAG) 2>/dev/null && echo "✅ Tag deleted locally" || echo "⚠️  Tag not found locally"; \
		read -p "Delete tag from remote? [y/N] " -n 1 -r; \
		echo; \
		if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
			git push origin :refs/tags/$(TAG) 2>/dev/null && echo "✅ Tag deleted from remote" || echo "⚠️  Tag not found on remote"; \
		else \
			echo "⚠️  Tag deleted locally but not from remote. Use 'git push origin :refs/tags/$(TAG)' to delete from remote later."; \
		fi; \
	else \
		echo "❌ Deletion cancelled"; \
	fi
