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

func TestNewYAMLExampleRepository(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")

	repo, err := NewYAMLExampleRepository(filePath)
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

func TestCreateExample(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	example := &models.ToolExample{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list all pods",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, example)
	if err != nil {
		t.Fatalf("Failed to create example: %v", err)
	}

	// Verify we can retrieve it
	retrieved, err := repo.GetByCommand(ctx, example.Command)
	if err != nil {
		t.Fatalf("Failed to retrieve example: %v", err)
	}

	if retrieved.ToolName != example.ToolName {
		t.Errorf("Expected tool name %s, got %s", example.ToolName, retrieved.ToolName)
	}

	if retrieved.Description != example.Description {
		t.Errorf("Expected description %s, got %s", example.Description, retrieved.Description)
	}
}

func TestCreateDuplicateCommand(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	example := &models.ToolExample{
		Command:     "lsof -i :8080",
		ToolName:    "lsof",
		Description: "check port 8080",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create first time should succeed
	if err := repo.Create(ctx, example); err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// Create duplicate command should fail
	example2 := &models.ToolExample{
		Command:     "lsof -i :8080", // Same command (primary key)
		ToolName:    "lsof",
		Description: "different description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.Create(ctx, example2)
	if err != ErrExampleAlreadyExists {
		t.Errorf("Expected ErrExampleAlreadyExists, got %v", err)
	}
}

func TestGetByCommand(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	example := &models.ToolExample{
		Command:     "docker ps -a",
		ToolName:    "docker",
		Description: "list all containers",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.Create(ctx, example)

	retrieved, err := repo.GetByCommand(ctx, "docker ps -a")
	if err != nil {
		t.Fatalf("Failed to get by command: %v", err)
	}

	if retrieved.ToolName != example.ToolName {
		t.Errorf("Expected tool name %s, got %s", example.ToolName, retrieved.ToolName)
	}
}

func TestGetByCommandNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	_, err := repo.GetByCommand(ctx, "nonexistent command")
	if err != ErrExampleNotFound {
		t.Errorf("Expected ErrExampleNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()

	// Create multiple examples
	examples := []*models.ToolExample{
		{
			Command:     "kubectl get pods",
			ToolName:    "kubectl",
			Description: "list pods",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Command:     "kubectl get nodes",
			ToolName:    "kubectl",
			Description: "list nodes",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Command:     "docker ps",
			ToolName:    "docker",
			Description: "list containers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, example := range examples {
		if err := repo.Create(ctx, example); err != nil {
			t.Fatalf("Failed to create example: %v", err)
		}
	}

	// List all examples
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 examples, got %d", len(list))
	}
}

func TestListByToolName(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()

	// Create examples for different tools
	examples := []*models.ToolExample{
		{
			Command:     "kubectl get pods",
			ToolName:    "kubectl",
			Description: "list pods",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Command:     "kubectl get nodes",
			ToolName:    "kubectl",
			Description: "list nodes",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Command:     "docker ps",
			ToolName:    "docker",
			Description: "list containers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, example := range examples {
		repo.Create(ctx, example)
	}

	// List only kubectl examples
	kubectlExamples, err := repo.ListByToolName(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to list by tool name: %v", err)
	}

	if len(kubectlExamples) != 2 {
		t.Errorf("Expected 2 kubectl examples, got %d", len(kubectlExamples))
	}

	for _, ex := range kubectlExamples {
		if ex.ToolName != "kubectl" {
			t.Errorf("Expected tool name kubectl, got %s", ex.ToolName)
		}
	}
}

func TestUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	example := &models.ToolExample{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "old description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.Create(ctx, example)

	// Update the example
	example.Description = "new description"
	example.UpdatedAt = time.Now()

	err := repo.Update(ctx, example)
	if err != nil {
		t.Fatalf("Failed to update example: %v", err)
	}

	// Verify update
	retrieved, _ := repo.GetByCommand(ctx, example.Command)
	if retrieved.Description != "new description" {
		t.Errorf("Expected updated description, got %s", retrieved.Description)
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	example := &models.ToolExample{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.Create(ctx, example)

	// Delete the example
	err := repo.Delete(ctx, example.Command)
	if err != nil {
		t.Fatalf("Failed to delete example: %v", err)
	}

	// Verify it's gone
	_, err = repo.GetByCommand(ctx, example.Command)
	if err != ErrExampleNotFound {
		t.Errorf("Expected ErrExampleNotFound after delete, got %v", err)
	}
}

func TestDeleteByToolName(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()

	// Create multiple examples for same tool
	examples := []*models.ToolExample{
		{
			Command:     "kubectl get pods",
			ToolName:    "kubectl",
			Description: "list pods",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Command:     "kubectl get nodes",
			ToolName:    "kubectl",
			Description: "list nodes",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Command:     "docker ps",
			ToolName:    "docker",
			Description: "list containers",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, example := range examples {
		repo.Create(ctx, example)
	}

	// Delete by tool name
	err := repo.DeleteByToolName(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to delete by tool name: %v", err)
	}

	// Verify kubectl examples are gone
	list, _ := repo.List(ctx)
	if len(list) != 1 {
		t.Errorf("Expected 1 example remaining, got %d", len(list))
	}

	if list[0].ToolName != "docker" {
		t.Errorf("Expected docker example to remain, got %s", list[0].ToolName)
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLExampleRepository(filePath)

	ctx := context.Background()
	example := &models.ToolExample{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Should not exist initially
	exists, err := repo.Exists(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if exists {
		t.Error("Example should not exist before creation")
	}

	// Create example
	repo.Create(ctx, example)

	// Should exist now
	exists, err = repo.Exists(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Exists check failed: %v", err)
	}
	if !exists {
		t.Error("Example should exist after creation")
	}
}
