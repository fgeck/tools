//go:build integration
// +build integration

package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fgeck/tools/internal/config"
	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/repository/yaml"
	"github.com/fgeck/tools/internal/service"
)

func setupTestCLI(t *testing.T) (string, func()) {
	// Create temp directory for test storage
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")

	// Initialize repository and service
	repo, err := yaml.NewYAMLBookmarkRepository(filePath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	testSvc := service.NewBookmarkService(repo)
	Initialize(testSvc)

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return filePath, cleanup
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestCLIAddCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	// Simulate add command
	ctx := context.Background()
	_, err := svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list all pods",
	})

	if err != nil {
		t.Errorf("Add command failed: %v", err)
	}

	// Verify example was created
	resp, err := svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 1 {
		t.Errorf("Expected 1 example, got %d", resp.Count)
	}
}

func TestCLIListCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add some examples first
	examples := []struct {
		command     string
		toolName    string
		description string
	}{
		{"kubectl get pods", "kubectl", "list pods"},
		{"docker ps", "docker", "list containers"},
		{"helm list", "helm", "list releases"},
	}

	for _, ex := range examples {
		svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
			Command:     ex.command,
			ToolName:    ex.toolName,
			Description: ex.description,
		})
	}

	// List examples
	output := captureOutput(func() {
		listExamples()
	})

	// Verify output contains tool names
	for _, ex := range examples {
		if !strings.Contains(output, ex.toolName) {
			t.Errorf("Output should contain tool name %s", ex.toolName)
		}
	}

	if !strings.Contains(output, "Total: 3 examples") {
		t.Error("Output should show total count")
	}
}

func TestCLIEditCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add an example
	svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "old description",
	})

	// Edit the example
	_, err := svc.UpdateBookmark(ctx, dto.UpdateBookmarkRequest{
		Command:        "kubectl get pods",
		NewDescription: "new description",
	})
	if err != nil {
		t.Errorf("Edit command failed: %v", err)
	}

	// Verify it was updated
	example, err := svc.GetBookmark(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Failed to get example: %v", err)
	}

	if example.Description != "new description" {
		t.Errorf("Expected 'new description', got %s", example.Description)
	}
}

func TestCLIEditCommandChangeCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add an example
	svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	})

	// Edit the command (primary key)
	_, err := svc.UpdateBookmark(ctx, dto.UpdateBookmarkRequest{
		Command:    "kubectl get pods",
		NewCommand: "kubectl get pods -A",
	})
	if err != nil {
		t.Errorf("Edit command failed: %v", err)
	}

	// Verify old command is gone
	_, err = svc.GetBookmark(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Old command should not exist")
	}

	// Verify new command exists
	example, err := svc.GetBookmark(ctx, "kubectl get pods -A")
	if err != nil {
		t.Fatalf("Failed to get example with new command: %v", err)
	}

	if example.Command != "kubectl get pods -A" {
		t.Errorf("Expected command 'kubectl get pods -A', got %s", example.Command)
	}
}

func TestCLIRemoveCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add an example
	svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	})

	// Remove the example
	err := svc.DeleteBookmark(ctx, "kubectl get pods")
	if err != nil {
		t.Errorf("Remove command failed: %v", err)
	}

	// Verify it's gone
	resp, err := svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("Expected 0 examples after removal, got %d", resp.Count)
	}
}

func TestCLIEndToEndWorkflow(t *testing.T) {
	filePath, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add multiple examples
	examples := []struct {
		command     string
		toolName    string
		description string
	}{
		{
			command:     "kubectl get pods",
			toolName:    "kubectl",
			description: "list all pods",
		},
		{
			command:     "kubectl get nodes",
			toolName:    "kubectl",
			description: "list all nodes",
		},
		{
			command:     "docker ps",
			toolName:    "docker",
			description: "list containers",
		},
	}

	for _, ex := range examples {
		_, err := svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
			Command:     ex.command,
			ToolName:    ex.toolName,
			Description: ex.description,
		})
		if err != nil {
			t.Fatalf("Failed to add example %s: %v", ex.command, err)
		}
	}

	// List and verify
	resp, err := svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 3 {
		t.Errorf("Expected 3 examples, got %d", resp.Count)
	}

	// Update one example
	_, err = svc.UpdateBookmark(ctx, dto.UpdateBookmarkRequest{
		Command:        "kubectl get pods",
		NewDescription: "updated description",
	})
	if err != nil {
		t.Fatalf("Failed to update example: %v", err)
	}

	// Verify update
	updated, err := svc.GetBookmark(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Failed to get updated example: %v", err)
	}
	if updated.Description != "updated description" {
		t.Errorf("Expected updated description, got %s", updated.Description)
	}

	// Remove one example
	err = svc.DeleteBookmark(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Failed to remove example: %v", err)
	}

	// Verify only two remain
	resp, err = svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 2 {
		t.Errorf("Expected 2 examples after removal, got %d", resp.Count)
	}

	// Verify YAML file was updated
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("YAML file should exist")
	}
}

func TestCLIDefaultListBehavior(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Test empty list
	output := captureOutput(func() {
		listExamples()
	})

	if !strings.Contains(output, "No examples found") {
		t.Error("Should show 'No examples found' message when empty")
	}

	// Add an example
	svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	})

	// Test non-empty list
	output = captureOutput(func() {
		listExamples()
	})

	if !strings.Contains(output, "kubectl") {
		t.Error("Should show tool in list")
	}

	if !strings.Contains(output, "Total: 1 examples") {
		t.Error("Should show total count")
	}
}

func TestCLIPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")

	// Create first instance
	repo1, _ := yaml.NewYAMLBookmarkRepository(filePath)
	svc1 := service.NewBookmarkService(repo1)

	// Add examples
	ctx := context.Background()
	svc1.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	})
	svc1.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "kubectl get nodes",
		ToolName:    "kubectl",
		Description: "list nodes",
	})

	// Create second instance (simulating restart)
	repo2, _ := yaml.NewYAMLBookmarkRepository(filePath)
	svc2 := service.NewBookmarkService(repo2)

	// Verify examples persisted
	resp, err := svc2.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 2 {
		t.Errorf("Expected 2 persisted examples, got %d", resp.Count)
	}
}

func TestCLIWithXDGConfigHome(t *testing.T) {
	// Save original
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Set custom XDG_CONFIG_HOME
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Get storage path
	cfg := config.DefaultConfig()
	expectedPath := filepath.Join(tmpDir, "tools", "tools.yaml")

	if cfg.StorageFilePath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, cfg.StorageFilePath)
	}

	// Create repository with this path
	repo, err := yaml.NewYAMLBookmarkRepository(cfg.StorageFilePath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	svc := service.NewBookmarkService(repo)

	// Add an example
	ctx := context.Background()
	_, err = svc.CreateBookmark(ctx, dto.CreateBookmarkRequest{
		Command:     "test command",
		ToolName:    "test-tool",
		Description: "test description",
	})

	if err != nil {
		t.Fatalf("Failed to create example: %v", err)
	}

	// Verify file was created in correct location
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("File should exist at %s", expectedPath)
	}
}

func TestCLIMultipleExamplesForSameTool(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add multiple examples for same tool
	examples := []dto.CreateBookmarkRequest{
		{
			Command:     "lsof -i :8080",
			ToolName:    "lsof",
			Description: "check port 8080",
		},
		{
			Command:     "lsof -i :3000",
			ToolName:    "lsof",
			Description: "check port 3000",
		},
		{
			Command:     "lsof -t -i :8080 | xargs kill -9",
			ToolName:    "lsof",
			Description: "kill process on port 8080",
		},
	}

	for _, req := range examples {
		_, err := svc.CreateBookmark(ctx, req)
		if err != nil {
			t.Fatalf("Failed to add example: %v", err)
		}
	}

	// Verify all examples exist
	resp, err := svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 3 {
		t.Errorf("Expected 3 examples, got %d", resp.Count)
	}

	// Verify all are for lsof
	for _, ex := range resp.Examples {
		if ex.ToolName != "lsof" {
			t.Errorf("Expected tool name lsof, got %s", ex.ToolName)
		}
	}

	// Delete one example by command
	err = svc.DeleteBookmark(ctx, "lsof -i :3000")
	if err != nil {
		t.Fatalf("Failed to delete example: %v", err)
	}

	// Verify only 2 remain
	resp, err = svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 2 {
		t.Errorf("Expected 2 examples after deletion, got %d", resp.Count)
	}
}
