//go:build unit
// +build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/fgeck/tools/internal/domain/models"
	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/repository"
)

// Error constants for mock repository
var (
	ErrExampleNotFound      = errors.New("example not found")
	ErrExampleAlreadyExists = errors.New("example already exists")
)

// Mock repository for testing
type mockExampleRepository struct {
	examples map[string]*models.ToolExample
}

func newMockExampleRepository() repository.ExampleRepository {
	return &mockExampleRepository{
		examples: make(map[string]*models.ToolExample),
	}
}

func (m *mockExampleRepository) Create(ctx context.Context, example *models.ToolExample) error {
	if _, exists := m.examples[example.Command]; exists {
		return ErrExampleAlreadyExists
	}
	m.examples[example.Command] = example
	return nil
}

func (m *mockExampleRepository) GetByCommand(ctx context.Context, command string) (*models.ToolExample, error) {
	example, ok := m.examples[command]
	if !ok {
		return nil, ErrExampleNotFound
	}
	return example, nil
}

func (m *mockExampleRepository) List(ctx context.Context) ([]*models.ToolExample, error) {
	list := make([]*models.ToolExample, 0, len(m.examples))
	for _, example := range m.examples {
		list = append(list, example)
	}
	return list, nil
}

func (m *mockExampleRepository) ListByToolName(ctx context.Context, toolName string) ([]*models.ToolExample, error) {
	list := make([]*models.ToolExample, 0)
	for _, example := range m.examples {
		if example.ToolName == toolName {
			list = append(list, example)
		}
	}
	return list, nil
}

func (m *mockExampleRepository) Update(ctx context.Context, example *models.ToolExample) error {
	if _, ok := m.examples[example.Command]; !ok {
		return ErrExampleNotFound
	}
	m.examples[example.Command] = example
	return nil
}

func (m *mockExampleRepository) Delete(ctx context.Context, command string) error {
	if _, ok := m.examples[command]; !ok {
		return ErrExampleNotFound
	}
	delete(m.examples, command)
	return nil
}

func (m *mockExampleRepository) DeleteByToolName(ctx context.Context, toolName string) error {
	found := false
	for cmd, example := range m.examples {
		if example.ToolName == toolName {
			delete(m.examples, cmd)
			found = true
		}
	}
	if !found {
		return ErrExampleNotFound
	}
	return nil
}

func (m *mockExampleRepository) Exists(ctx context.Context, command string) (bool, error) {
	_, ok := m.examples[command]
	return ok, nil
}

func TestCreateExample(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	req := dto.CreateExampleRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list all pods",
	}

	resp, err := svc.CreateExample(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create example: %v", err)
	}

	if resp.Command != req.Command {
		t.Errorf("Expected command %s, got %s", req.Command, resp.Command)
	}

	if resp.ToolName != req.ToolName {
		t.Errorf("Expected tool name %s, got %s", req.ToolName, resp.ToolName)
	}

	if resp.Description != req.Description {
		t.Errorf("Expected description %s, got %s", req.Description, resp.Description)
	}
}

func TestCreateExampleValidation(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     dto.CreateExampleRequest
		wantErr bool
	}{
		{
			name: "empty command",
			req: dto.CreateExampleRequest{
				Command:     "",
				ToolName:    "kubectl",
				Description: "test",
			},
			wantErr: true,
		},
		{
			name: "empty tool name",
			req: dto.CreateExampleRequest{
				Command:     "kubectl get pods",
				ToolName:    "",
				Description: "test",
			},
			wantErr: true,
		},
		{
			name: "empty description",
			req: dto.CreateExampleRequest{
				Command:     "kubectl get pods",
				ToolName:    "kubectl",
				Description: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace command",
			req: dto.CreateExampleRequest{
				Command:     "   ",
				ToolName:    "kubectl",
				Description: "test",
			},
			wantErr: true,
		},
		{
			name: "valid example",
			req: dto.CreateExampleRequest{
				Command:     "kubectl get pods",
				ToolName:    "kubectl",
				Description: "list pods",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateExample(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateExample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateExampleDuplicate(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	req := dto.CreateExampleRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	// Create first example
	_, err := svc.CreateExample(ctx, req)
	if err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// Try to create duplicate command
	_, err = svc.CreateExample(ctx, req)
	if err == nil {
		t.Error("Expected error for duplicate command")
	}
}

func TestGetExample(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	// Create an example first
	req := dto.CreateExampleRequest{
		Command:     "docker ps -a",
		ToolName:    "docker",
		Description: "list all containers",
	}

	created, _ := svc.CreateExample(ctx, req)

	// Get the example
	resp, err := svc.GetExample(ctx, "docker ps -a")
	if err != nil {
		t.Fatalf("Failed to get example: %v", err)
	}

	if resp.Command != created.Command {
		t.Errorf("Expected command %s, got %s", created.Command, resp.Command)
	}

	if resp.ToolName != req.ToolName {
		t.Errorf("Expected tool name %s, got %s", req.ToolName, resp.ToolName)
	}
}

func TestGetExampleNotFound(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	_, err := svc.GetExample(ctx, "nonexistent command")
	if err == nil {
		t.Error("Expected error for nonexistent example")
	}
}

func TestListExamples(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	// Create multiple examples
	examples := []dto.CreateExampleRequest{
		{Command: "kubectl get pods", ToolName: "kubectl", Description: "list pods"},
		{Command: "kubectl get nodes", ToolName: "kubectl", Description: "list nodes"},
		{Command: "docker ps", ToolName: "docker", Description: "list containers"},
	}

	for _, req := range examples {
		svc.CreateExample(ctx, req)
	}

	// List all examples
	resp, err := svc.ListExamples(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 3 {
		t.Errorf("Expected 3 examples, got %d", resp.Count)
	}

	if len(resp.Examples) != 3 {
		t.Errorf("Expected 3 examples in response, got %d", len(resp.Examples))
	}
}

func TestListExamplesEmpty(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	resp, err := svc.ListExamples(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("Expected 0 examples, got %d", resp.Count)
	}
}

func TestUpdateExample(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	// Create an example
	req := dto.CreateExampleRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "old description",
	}

	created, _ := svc.CreateExample(ctx, req)

	// Update the example
	updateReq := dto.UpdateExampleRequest{
		Command:        "kubectl get pods",
		NewDescription: "new description",
	}

	resp, err := svc.UpdateExample(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update example: %v", err)
	}

	if resp.Description != "new description" {
		t.Errorf("Expected description 'new description', got %s", resp.Description)
	}

	// Verify CreatedAt didn't change
	if resp.CreatedAt != created.CreatedAt {
		t.Error("CreatedAt should not change on update")
	}
}

func TestUpdateExampleChangeCommand(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	// Create an example
	req := dto.CreateExampleRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	svc.CreateExample(ctx, req)

	// Update with new command (primary key change)
	updateReq := dto.UpdateExampleRequest{
		Command:    "kubectl get pods",
		NewCommand: "kubectl get pods -A",
	}

	resp, err := svc.UpdateExample(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update example: %v", err)
	}

	if resp.Command != "kubectl get pods -A" {
		t.Errorf("Expected command 'kubectl get pods -A', got %s", resp.Command)
	}

	// Verify old command doesn't exist
	_, err = svc.GetExample(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Old command should not exist after update")
	}

	// Verify new command exists
	_, err = svc.GetExample(ctx, "kubectl get pods -A")
	if err != nil {
		t.Error("New command should exist after update")
	}
}

func TestUpdateExampleNotFound(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	updateReq := dto.UpdateExampleRequest{
		Command:        "nonexistent",
		NewDescription: "test",
	}

	_, err := svc.UpdateExample(ctx, updateReq)
	if err == nil {
		t.Error("Expected error when updating nonexistent example")
	}
}

func TestDeleteExample(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	// Create an example
	req := dto.CreateExampleRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	svc.CreateExample(ctx, req)

	// Delete the example
	err := svc.DeleteExample(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Failed to delete example: %v", err)
	}

	// Verify it's gone
	_, err = svc.GetExample(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Example should not exist after deletion")
	}
}

func TestDeleteExampleNotFound(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	err := svc.DeleteExample(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error when deleting nonexistent example")
	}
}

func TestDeleteToolExamples(t *testing.T) {
	repo := newMockExampleRepository()
	svc := NewExampleService(repo)
	ctx := context.Background()

	// Create multiple examples for same tool
	examples := []dto.CreateExampleRequest{
		{Command: "kubectl get pods", ToolName: "kubectl", Description: "list pods"},
		{Command: "kubectl get nodes", ToolName: "kubectl", Description: "list nodes"},
		{Command: "docker ps", ToolName: "docker", Description: "list containers"},
	}

	for _, req := range examples {
		svc.CreateExample(ctx, req)
	}

	// Delete all kubectl examples
	err := svc.DeleteToolExamples(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to delete tool examples: %v", err)
	}

	// Verify kubectl examples are gone
	resp, _ := svc.ListExamples(ctx)
	if resp.Count != 1 {
		t.Errorf("Expected 1 example remaining, got %d", resp.Count)
	}

	if resp.Examples[0].ToolName != "docker" {
		t.Errorf("Expected docker example to remain, got %s", resp.Examples[0].ToolName)
	}
}
