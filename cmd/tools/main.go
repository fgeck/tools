package main

import (
	"fmt"
	"os"

	"github.com/fgeck/tools/internal/cli"
	"github.com/fgeck/tools/internal/config"
	"github.com/fgeck/tools/internal/repository/yaml"
	"github.com/fgeck/tools/internal/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration
	cfg := config.DefaultConfig()

	// Initialize repository
	repo, err := yaml.NewYAMLToolRepository(cfg.StorageFilePath)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Initialize service
	svc := service.NewToolService(repo)

	// Initialize and execute CLI
	cli.Initialize(svc)
	cli.Execute()

	return nil
}
