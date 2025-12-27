//go:build unit
// +build unit

package yaml

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fgeck/tools/internal/domain/models"
)

func TestNewYAMLToolRepository(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")

	repo, err := NewYAMLToolRepository(filePath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	if repo == nil {
		t.Fatal("Repository should not be nil")
	}

	// Verify file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("YAML file should have been created")
	}
}

func TestCreate(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:          "test-id-1",
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Kubernetes CLI",
		Examples:    []string{"kubectl get pods"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, tool)
	if err != nil {
		t.Fatalf("Failed to create tool: %v", err)
	}

	// Verify we can retrieve it
	retrieved, err := repo.GetByID(ctx, tool.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve tool: %v", err)
	}

	if retrieved.Name != tool.Name {
		t.Errorf("Expected name %s, got %s", tool.Name, retrieved.Name)
	}
}

func TestCreateDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:          "test-id-1",
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Kubernetes CLI",
		Examples:    []string{"kubectl get pods"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create first time should succeed
	if err := repo.Create(ctx, tool); err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// Create duplicate should fail
	tool.ID = "test-id-2" // Different ID but same name
	err := repo.Create(ctx, tool)
	if err != ErrToolAlreadyExists {
		t.Errorf("Expected ErrToolAlreadyExists, got %v", err)
	}
}

func TestGetByName(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:          "test-id-1",
		Name:        "docker",
		Command:     "/usr/bin/docker",
		Description: "Container tool",
		Examples:    []string{"docker ps"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.Create(ctx, tool)

	retrieved, err := repo.GetByName(ctx, "docker")
	if err != nil {
		t.Fatalf("Failed to get by name: %v", err)
	}

	if retrieved.Command != tool.Command {
		t.Errorf("Expected command %s, got %s", tool.Command, retrieved.Command)
	}
}

func TestGetByNameNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	_, err := repo.GetByName(ctx, "nonexistent")
	if err != ErrToolNotFound {
		t.Errorf("Expected ErrToolNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()

	// Create multiple tools
	tools := []*models.Tool{
		{
			ID:        "1",
			Name:      "kubectl",
			Command:   "/usr/bin/kubectl",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			Name:      "docker",
			Command:   "/usr/bin/docker",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "3",
			Name:      "helm",
			Command:   "/usr/bin/helm",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, tool := range tools {
		if err := repo.Create(ctx, tool); err != nil {
			t.Fatalf("Failed to create tool: %v", err)
		}
	}

	// List all tools
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 tools, got %d", len(list))
	}
}

func TestUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:          "test-id-1",
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Old description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.Create(ctx, tool)

	// Update the tool
	tool.Description = "New description"
	tool.UpdatedAt = time.Now()

	err := repo.Update(ctx, tool)
	if err != nil {
		t.Fatalf("Failed to update tool: %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetByID(ctx, tool.ID)
	if retrieved.Description != "New description" {
		t.Errorf("Expected updated description, got %s", retrieved.Description)
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:        "test-id-1",
		Name:      "kubectl",
		Command:   "/usr/bin/kubectl",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(ctx, tool)

	// Delete the tool
	err := repo.Delete(ctx, tool.ID)
	if err != nil {
		t.Fatalf("Failed to delete tool: %v", err)
	}

	// Verify it's gone
	_, err = repo.GetByID(ctx, tool.ID)
	if err != ErrToolNotFound {
		t.Errorf("Expected ErrToolNotFound after delete, got %v", err)
	}
}

func TestDeleteByName(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:        "test-id-1",
		Name:      "kubectl",
		Command:   "/usr/bin/kubectl",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.Create(ctx, tool)

	// Delete by name
	err := repo.DeleteByName(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to delete by name: %v", err)
	}

	// Verify it's gone
	_, err = repo.GetByName(ctx, "kubectl")
	if err != ErrToolNotFound {
		t.Errorf("Expected ErrToolNotFound after delete, got %v", err)
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLToolRepository(filePath)

	ctx := context.Background()
	tool := &models.Tool{
		ID:        "test-id-1",
		Name:      "kubectl",
		Command:   "/usr/bin/kubectl",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Should not exist initially
	exists, err := repo.Exists(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if exists {
		t.Error("Tool should not exist before creation")
	}

	// Create tool
	repo.Create(ctx, tool)

	// Should exist now
	exists, err = repo.Exists(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Tool should exist after creation")
	}
}
