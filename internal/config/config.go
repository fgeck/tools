package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	StorageFilePath string
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		StorageFilePath: GetDefaultStoragePath(),
	}
}

// GetDefaultStoragePath returns the default YAML storage path
// Following XDG Base Directory specification
func GetDefaultStoragePath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, _ := os.UserHomeDir()
		configDir = filepath.Join(home, ".config")
	}

	return filepath.Join(configDir, "tools", "tools.yaml")
}
