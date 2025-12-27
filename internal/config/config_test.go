//go:build unit
// +build unit

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig should not return nil")
	}

	if cfg.StorageFilePath == "" {
		t.Error("StorageFilePath should not be empty")
	}
}

func TestGetDefaultStoragePath(t *testing.T) {
	// Save original env var
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	t.Run("with XDG_CONFIG_HOME set", func(t *testing.T) {
		testDir := "/tmp/test-config"
		os.Setenv("XDG_CONFIG_HOME", testDir)

		path := GetDefaultStoragePath()
		expected := filepath.Join(testDir, "tools", "tools.yaml")

		if path != expected {
			t.Errorf("Expected path %s, got %s", expected, path)
		}
	})

	t.Run("without XDG_CONFIG_HOME", func(t *testing.T) {
		os.Unsetenv("XDG_CONFIG_HOME")

		path := GetDefaultStoragePath()

		// Should contain .config/tools/tools.yaml
		if !filepath.IsAbs(path) {
			t.Error("Path should be absolute")
		}

		if filepath.Base(path) != "tools.yaml" {
			t.Errorf("Expected filename tools.yaml, got %s", filepath.Base(path))
		}

		if filepath.Base(filepath.Dir(path)) != "tools" {
			t.Errorf("Expected parent dir 'tools', got %s", filepath.Base(filepath.Dir(path)))
		}
	})
}

func TestStoragePathStructure(t *testing.T) {
	path := GetDefaultStoragePath()

	// Verify path structure
	dir := filepath.Dir(path)
	toolsDir := filepath.Base(dir)

	if toolsDir != "tools" {
		t.Errorf("Expected tools directory, got %s", toolsDir)
	}

	configDir := filepath.Base(filepath.Dir(dir))
	if configDir != ".config" && configDir != "test-config" {
		// Allow test-config for XDG override tests
		t.Errorf("Expected .config or test-config directory, got %s", configDir)
	}
}
