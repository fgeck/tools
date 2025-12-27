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
	repo, err := yaml.NewYAMLToolRepository(filePath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	testSvc := service.NewToolService(repo)
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
	_, err := svc.CreateTool(ctx, dto.CreateToolRequest{
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Kubernetes CLI",
		Examples:    []string{"kubectl get pods"},
	})

	if err != nil {
		t.Errorf("Add command failed: %v", err)
	}

	// Verify tool was created
	resp, err := svc.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 1 {
		t.Errorf("Expected 1 tool, got %d", resp.Count)
	}
}

func TestCLIListCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add some tools first
	tools := []struct {
		name    string
		command string
	}{
		{"kubectl", "/usr/bin/kubectl"},
		{"docker", "/usr/bin/docker"},
		{"helm", "/usr/bin/helm"},
	}

	for _, tool := range tools {
		svc.CreateTool(ctx, dto.CreateToolRequest{
			Name:    tool.name,
			Command: tool.command,
		})
	}

	// List tools
	output := captureOutput(func() {
		listTools()
	})

	// Verify output contains tool names
	for _, tool := range tools {
		if !strings.Contains(output, tool.name) {
			t.Errorf("Output should contain tool name %s", tool.name)
		}
	}

	if !strings.Contains(output, "Total: 3 tools") {
		t.Error("Output should show total count")
	}
}

func TestCLIRemoveCommand(t *testing.T) {
	_, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add a tool
	svc.CreateTool(ctx, dto.CreateToolRequest{
		Name:    "kubectl",
		Command: "/usr/bin/kubectl",
	})

	// Remove the tool
	err := svc.DeleteTool(ctx, "kubectl")
	if err != nil {
		t.Errorf("Remove command failed: %v", err)
	}

	// Verify it's gone
	resp, err := svc.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("Expected 0 tools after removal, got %d", resp.Count)
	}
}

func TestCLIEndToEndWorkflow(t *testing.T) {
	filePath, cleanup := setupTestCLI(t)
	defer cleanup()

	ctx := context.Background()

	// Add multiple tools
	tools := []struct {
		name        string
		command     string
		description string
		examples    []string
	}{
		{
			name:        "kubectl",
			command:     "/usr/bin/kubectl",
			description: "Kubernetes CLI",
			examples:    []string{"kubectl get pods", "kubectl describe node"},
		},
		{
			name:        "docker",
			command:     "/usr/bin/docker",
			description: "Container tool",
			examples:    []string{"docker ps", "docker images"},
		},
	}

	for _, tool := range tools {
		_, err := svc.CreateTool(ctx, dto.CreateToolRequest{
			Name:        tool.name,
			Command:     tool.command,
			Description: tool.description,
			Examples:    tool.examples,
		})
		if err != nil {
			t.Fatalf("Failed to add tool %s: %v", tool.name, err)
		}
	}

	// List and verify
	resp, err := svc.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 2 {
		t.Errorf("Expected 2 tools, got %d", resp.Count)
	}

	// Remove one tool
	err = svc.DeleteTool(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to remove tool: %v", err)
	}

	// Verify only one remains
	resp, err = svc.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 1 {
		t.Errorf("Expected 1 tool after removal, got %d", resp.Count)
	}

	if resp.Tools[0].Name != "docker" {
		t.Errorf("Expected docker to remain, got %s", resp.Tools[0].Name)
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
		listTools()
	})

	if !strings.Contains(output, "No tools found") {
		t.Error("Should show 'No tools found' message when empty")
	}

	// Add a tool
	svc.CreateTool(ctx, dto.CreateToolRequest{
		Name:    "kubectl",
		Command: "/usr/bin/kubectl",
	})

	// Test non-empty list
	output = captureOutput(func() {
		listTools()
	})

	if !strings.Contains(output, "kubectl") {
		t.Error("Should show tool in list")
	}

	if !strings.Contains(output, "Total: 1 tools") {
		t.Error("Should show total count")
	}
}

func TestCLIPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tools.yaml")

	// Create first instance
	repo1, _ := yaml.NewYAMLToolRepository(filePath)
	svc1 := service.NewToolService(repo1)

	// Add a tool
	ctx := context.Background()
	svc1.CreateTool(ctx, dto.CreateToolRequest{
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Kubernetes CLI",
		Examples:    []string{"kubectl get pods"},
	})

	// Create second instance (simulating restart)
	repo2, _ := yaml.NewYAMLToolRepository(filePath)
	svc2 := service.NewToolService(repo2)

	// Verify tool persisted
	resp, err := svc2.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 1 {
		t.Errorf("Expected 1 persisted tool, got %d", resp.Count)
	}

	if resp.Tools[0].Name != "kubectl" {
		t.Errorf("Expected kubectl, got %s", resp.Tools[0].Name)
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
	repo, err := yaml.NewYAMLToolRepository(cfg.StorageFilePath)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	svc := service.NewToolService(repo)

	// Add a tool
	ctx := context.Background()
	_, err = svc.CreateTool(ctx, dto.CreateToolRequest{
		Name:    "test-tool",
		Command: "/usr/bin/test",
	})

	if err != nil {
		t.Fatalf("Failed to create tool: %v", err)
	}

	// Verify file was created in correct location
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("File should exist at %s", expectedPath)
	}
}
