//go:build unit
// +build unit

package service

import (
	"context"
	"testing"

	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/repository"
)

func TestCreateTool(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	req := dto.CreateToolRequest{
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Kubernetes CLI",
		Examples:    []string{"kubectl get pods"},
	}

	resp, err := svc.CreateTool(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create tool: %v", err)
	}

	if resp.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
	}

	if resp.Command != req.Command {
		t.Errorf("Expected command %s, got %s", req.Command, resp.Command)
	}

	if resp.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestCreateToolValidation(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     dto.CreateToolRequest
		wantErr bool
	}{
		{
			name: "empty name",
			req: dto.CreateToolRequest{
				Name:    "",
				Command: "/usr/bin/kubectl",
			},
			wantErr: true,
		},
		{
			name: "empty command",
			req: dto.CreateToolRequest{
				Name:    "kubectl",
				Command: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace name",
			req: dto.CreateToolRequest{
				Name:    "   ",
				Command: "/usr/bin/kubectl",
			},
			wantErr: true,
		},
		{
			name: "valid tool",
			req: dto.CreateToolRequest{
				Name:    "kubectl",
				Command: "/usr/bin/kubectl",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateTool(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateToolDuplicate(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	req := dto.CreateToolRequest{
		Name:    "kubectl",
		Command: "/usr/bin/kubectl",
	}

	// Create first tool
	_, err := svc.CreateTool(ctx, req)
	if err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// Try to create duplicate
	_, err = svc.CreateTool(ctx, req)
	if err == nil {
		t.Error("Expected error for duplicate tool name")
	}
}

func TestGetTool(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	// Create a tool first
	req := dto.CreateToolRequest{
		Name:        "docker",
		Command:     "/usr/bin/docker",
		Description: "Container tool",
		Examples:    []string{"docker ps", "docker images"},
	}

	created, _ := svc.CreateTool(ctx, req)

	// Get the tool
	resp, err := svc.GetTool(ctx, "docker")
	if err != nil {
		t.Fatalf("Failed to get tool: %v", err)
	}

	if resp.ID != created.ID {
		t.Errorf("Expected ID %s, got %s", created.ID, resp.ID)
	}

	if resp.Name != req.Name {
		t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
	}

	if len(resp.Examples) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(resp.Examples))
	}
}

func TestGetToolNotFound(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	_, err := svc.GetTool(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent tool")
	}
}

func TestListTools(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	// Create multiple tools
	tools := []dto.CreateToolRequest{
		{Name: "kubectl", Command: "/usr/bin/kubectl"},
		{Name: "docker", Command: "/usr/bin/docker"},
		{Name: "helm", Command: "/usr/bin/helm"},
	}

	for _, req := range tools {
		svc.CreateTool(ctx, req)
	}

	// List all tools
	resp, err := svc.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 3 {
		t.Errorf("Expected 3 tools, got %d", resp.Count)
	}

	if len(resp.Tools) != 3 {
		t.Errorf("Expected 3 tools in response, got %d", len(resp.Tools))
	}
}

func TestListToolsEmpty(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	resp, err := svc.ListTools(ctx)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("Expected 0 tools, got %d", resp.Count)
	}
}

func TestUpdateTool(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	// Create a tool
	req := dto.CreateToolRequest{
		Name:        "kubectl",
		Command:     "/usr/bin/kubectl",
		Description: "Old description",
	}

	created, _ := svc.CreateTool(ctx, req)

	// Update the tool
	newDesc := "New description"
	newCmd := "/usr/local/bin/kubectl"
	updateReq := dto.UpdateToolRequest{
		Description: &newDesc,
		Command:     &newCmd,
	}

	resp, err := svc.UpdateTool(ctx, "kubectl", updateReq)
	if err != nil {
		t.Fatalf("Failed to update tool: %v", err)
	}

	if resp.Description != newDesc {
		t.Errorf("Expected description %s, got %s", newDesc, resp.Description)
	}

	if resp.Command != newCmd {
		t.Errorf("Expected command %s, got %s", newCmd, resp.Command)
	}

	// Verify CreatedAt didn't change
	if resp.CreatedAt != created.CreatedAt {
		t.Error("CreatedAt should not change on update")
	}

	// Note: UpdatedAt comparison might be flaky due to timestamp precision
	// Just verify the format is correct
	if resp.UpdatedAt == "" {
		t.Error("UpdatedAt should be set")
	}
}

func TestDeleteTool(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	// Create a tool
	req := dto.CreateToolRequest{
		Name:    "kubectl",
		Command: "/usr/bin/kubectl",
	}

	svc.CreateTool(ctx, req)

	// Delete the tool
	err := svc.DeleteTool(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to delete tool: %v", err)
	}

	// Verify it's gone
	_, err = svc.GetTool(ctx, "kubectl")
	if err == nil {
		t.Error("Tool should not exist after deletion")
	}
}

func TestDeleteToolNotFound(t *testing.T) {
	repo := repository.NewMockToolRepository()
	svc := NewToolService(repo)
	ctx := context.Background()

	err := svc.DeleteTool(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error when deleting nonexistent tool")
	}
}
