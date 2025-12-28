//go:build unit
// +build unit

package yaml

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fgeck/tools/internal/domain/models"
)

func TestNewYAMLBookmarkRepository(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")

	repo, err := NewYAMLBookmarkRepository(filePath)
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

func TestCreateBookmark(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list all pods",
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
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "lsof -i :8080",
		ToolName:    "lsof",
		Description: "check port 8080",
	}

	// Create first time should succeed
	if err := repo.Create(ctx, example); err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// Create duplicate command should fail
	example2 := &models.Bookmark{
		Command:     "lsof -i :8080", // Same command (primary key)
		ToolName:    "lsof",
		Description: "different description",
	}
	err := repo.Create(ctx, example2)
	if err != ErrBookmarkAlreadyExists {
		t.Errorf("Expected ErrBookmarkAlreadyExists, got %v", err)
	}
}

func TestGetByCommand(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "docker ps -a",
		ToolName:    "docker",
		Description: "list all containers",
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
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	_, err := repo.GetByCommand(ctx, "nonexistent command")
	if err != ErrBookmarkNotFound {
		t.Errorf("Expected ErrBookmarkNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Create multiple examples
	examples := []*models.Bookmark{
		{
			Command:     "kubectl get pods",
			ToolName:    "kubectl",
			Description: "list pods",
		},
		{
			Command:     "kubectl get nodes",
			ToolName:    "kubectl",
			Description: "list nodes",
		},
		{
			Command:     "docker ps",
			ToolName:    "docker",
			Description: "list containers",
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
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Create examples for different tools
	examples := []*models.Bookmark{
		{
			Command:     "kubectl get pods",
			ToolName:    "kubectl",
			Description: "list pods",
		},
		{
			Command:     "kubectl get nodes",
			ToolName:    "kubectl",
			Description: "list nodes",
		},
		{
			Command:     "docker ps",
			ToolName:    "docker",
			Description: "list containers",
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
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "old description",
	}

	repo.Create(ctx, example)

	// Update the example
	example.Description = "new description"

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
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	repo.Create(ctx, example)

	// Delete the example
	err := repo.Delete(ctx, example.Command)
	if err != nil {
		t.Fatalf("Failed to delete example: %v", err)
	}

	// Verify it's gone
	_, err = repo.GetByCommand(ctx, example.Command)
	if err != ErrBookmarkNotFound {
		t.Errorf("Expected ErrBookmarkNotFound after delete, got %v", err)
	}
}

func TestDeleteByToolName(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Create multiple examples for same tool
	examples := []*models.Bookmark{
		{
			Command:     "kubectl get pods",
			ToolName:    "kubectl",
			Description: "list pods",
		},
		{
			Command:     "kubectl get nodes",
			ToolName:    "kubectl",
			Description: "list nodes",
		},
		{
			Command:     "docker ps",
			ToolName:    "docker",
			Description: "list containers",
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
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
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

// Additional tests to improve coverage

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Test concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			example := &models.Bookmark{
				Command:     fmt.Sprintf("command-%d", n),
				ToolName:    "test",
				Description: fmt.Sprintf("description-%d", n),
			}
			repo.Create(ctx, example)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all were created
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list after concurrent writes: %v", err)
	}

	if len(list) != 10 {
		t.Errorf("Expected 10 bookmarks, got %d", len(list))
	}
}

func TestConcurrentReads(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Create some data first
	for i := 0; i < 5; i++ {
		example := &models.Bookmark{
			Command:     fmt.Sprintf("command-%d", i),
			ToolName:    "test",
			Description: fmt.Sprintf("description-%d", i),
		}
		repo.Create(ctx, example)
	}

	// Test concurrent reads
	done := make(chan bool)
	errors := make(chan error, 20)

	for i := 0; i < 20; i++ {
		go func() {
			_, err := repo.List(ctx)
			if err != nil {
				errors <- err
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	close(errors)
	for err := range errors {
		t.Errorf("Concurrent read error: %v", err)
	}
}

func TestUpdateNonExistentBookmark(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()
	example := &models.Bookmark{
		Command:     "nonexistent",
		ToolName:    "test",
		Description: "test",
	}

	err := repo.Update(ctx, example)
	if err == nil {
		t.Error("Expected error when updating non-existent bookmark")
	}
}

func TestDeleteByToolNameNoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Create a bookmark with different tool name
	example := &models.Bookmark{
		Command:     "docker ps",
		ToolName:    "docker",
		Description: "list containers",
	}
	repo.Create(ctx, example)

	// Try to delete bookmarks for different tool
	err := repo.DeleteByToolName(ctx, "kubectl")
	if err == nil {
		t.Error("Expected error when deleting non-existent tool bookmarks")
	}

	// Verify docker bookmark still exists
	list, _ := repo.ListByToolName(ctx, "docker")
	if len(list) != 1 {
		t.Errorf("Expected docker bookmark to still exist, got %d bookmarks", len(list))
	}
}

func TestEmptyRepository(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Test List on empty repository
	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List failed on empty repository: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d bookmarks", len(list))
	}

	// Test ListByToolName on empty repository
	list, err = repo.ListByToolName(ctx, "kubectl")
	if err != nil {
		t.Fatalf("ListByToolName failed on empty repository: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %d bookmarks", len(list))
	}

	// Test Exists on empty repository
	exists, err := repo.Exists(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("Exists failed on empty repository: %v", err)
	}
	if exists {
		t.Error("Expected bookmark to not exist in empty repository")
	}
}

func TestCreateUpdateDeleteCycle(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")
	repo, _ := NewYAMLBookmarkRepository(filePath)

	ctx := context.Background()

	// Create
	example := &models.Bookmark{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "original description",
	}

	err := repo.Create(ctx, example)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update
	example.Description = "updated description"
	example.ToolName = "k8s"

	err = repo.Update(ctx, example)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	retrieved, err := repo.GetByCommand(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("GetByCommand failed: %v", err)
	}

	if retrieved.Description != "updated description" {
		t.Errorf("Expected description 'updated description', got %s", retrieved.Description)
	}
	if retrieved.ToolName != "k8s" {
		t.Errorf("Expected tool name 'k8s', got %s", retrieved.ToolName)
	}

	// Delete
	err = repo.Delete(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.GetByCommand(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Expected error after deletion")
	}
}
