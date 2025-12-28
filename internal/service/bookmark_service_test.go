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
	ErrBookmarkNotFound      = errors.New("bookmark not found")
	ErrBookmarkAlreadyExists = errors.New("bookmark already exists")
)

// Mock repository for testing
type mockBookmarkRepository struct {
	examples map[string]*models.Bookmark
}

func newMockBookmarkRepository() repository.BookmarkRepository {
	return &mockBookmarkRepository{
		examples: make(map[string]*models.Bookmark),
	}
}

func (m *mockBookmarkRepository) Create(ctx context.Context, example *models.Bookmark) error {
	if _, exists := m.examples[example.Command]; exists {
		return ErrBookmarkAlreadyExists
	}
	m.examples[example.Command] = example
	return nil
}

func (m *mockBookmarkRepository) GetByCommand(ctx context.Context, command string) (*models.Bookmark, error) {
	example, ok := m.examples[command]
	if !ok {
		return nil, ErrBookmarkNotFound
	}
	return example, nil
}

func (m *mockBookmarkRepository) List(ctx context.Context) ([]*models.Bookmark, error) {
	list := make([]*models.Bookmark, 0, len(m.examples))
	for _, example := range m.examples {
		list = append(list, example)
	}
	return list, nil
}

func (m *mockBookmarkRepository) ListByToolName(ctx context.Context, toolName string) ([]*models.Bookmark, error) {
	list := make([]*models.Bookmark, 0)
	for _, example := range m.examples {
		if example.ToolName == toolName {
			list = append(list, example)
		}
	}
	return list, nil
}

func (m *mockBookmarkRepository) Update(ctx context.Context, example *models.Bookmark) error {
	if _, ok := m.examples[example.Command]; !ok {
		return ErrBookmarkNotFound
	}
	m.examples[example.Command] = example
	return nil
}

func (m *mockBookmarkRepository) Delete(ctx context.Context, command string) error {
	if _, ok := m.examples[command]; !ok {
		return ErrBookmarkNotFound
	}
	delete(m.examples, command)
	return nil
}

func (m *mockBookmarkRepository) DeleteByToolName(ctx context.Context, toolName string) error {
	found := false
	for cmd, example := range m.examples {
		if example.ToolName == toolName {
			delete(m.examples, cmd)
			found = true
		}
	}
	if !found {
		return ErrBookmarkNotFound
	}
	return nil
}

func (m *mockBookmarkRepository) Exists(ctx context.Context, command string) (bool, error) {
	_, ok := m.examples[command]
	return ok, nil
}

func TestCreateBookmark(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list all pods",
	}

	resp, err := svc.CreateBookmark(ctx, req)
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

func TestCreateBookmarkValidation(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	tests := []struct {
		name    string
		req     dto.CreateBookmarkRequest
		wantErr bool
	}{
		{
			name: "empty command",
			req: dto.CreateBookmarkRequest{
				Command:     "",
				ToolName:    "kubectl",
				Description: "test",
			},
			wantErr: true,
		},
		{
			name: "empty tool name",
			req: dto.CreateBookmarkRequest{
				Command:     "kubectl get pods",
				ToolName:    "",
				Description: "test",
			},
			wantErr: true,
		},
		{
			name: "empty description",
			req: dto.CreateBookmarkRequest{
				Command:     "kubectl get pods",
				ToolName:    "kubectl",
				Description: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace command",
			req: dto.CreateBookmarkRequest{
				Command:     "   ",
				ToolName:    "kubectl",
				Description: "test",
			},
			wantErr: true,
		},
		{
			name: "valid example",
			req: dto.CreateBookmarkRequest{
				Command:     "kubectl get pods",
				ToolName:    "kubectl",
				Description: "list pods",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateBookmark(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBookmark() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateBookmarkDuplicate(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	// Create first example
	_, err := svc.CreateBookmark(ctx, req)
	if err != nil {
		t.Fatalf("First create should succeed: %v", err)
	}

	// Try to create duplicate command
	_, err = svc.CreateBookmark(ctx, req)
	if err == nil {
		t.Error("Expected error for duplicate command")
	}
}

func TestGetBookmark(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create an example first
	req := dto.CreateBookmarkRequest{
		Command:     "docker ps -a",
		ToolName:    "docker",
		Description: "list all containers",
	}

	created, _ := svc.CreateBookmark(ctx, req)

	// Get the example
	resp, err := svc.GetBookmark(ctx, "docker ps -a")
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

func TestGetBookmarkNotFound(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	_, err := svc.GetBookmark(ctx, "nonexistent command")
	if err == nil {
		t.Error("Expected error for nonexistent example")
	}
}

func TestListBookmarks(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create multiple examples
	examples := []dto.CreateBookmarkRequest{
		{Command: "kubectl get pods", ToolName: "kubectl", Description: "list pods"},
		{Command: "kubectl get nodes", ToolName: "kubectl", Description: "list nodes"},
		{Command: "docker ps", ToolName: "docker", Description: "list containers"},
	}

	for _, req := range examples {
		svc.CreateBookmark(ctx, req)
	}

	// List all examples
	resp, err := svc.ListBookmarks(ctx)
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

func TestListBookmarksEmpty(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	resp, err := svc.ListBookmarks(ctx)
	if err != nil {
		t.Fatalf("Failed to list examples: %v", err)
	}

	if resp.Count != 0 {
		t.Errorf("Expected 0 examples, got %d", resp.Count)
	}
}

func TestUpdateBookmark(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create an example
	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "old description",
	}

	_, _ = svc.CreateBookmark(ctx, req)

	// Update the example
	updateReq := dto.UpdateBookmarkRequest{
		Command:        "kubectl get pods",
		NewDescription: "new description",
	}

	resp, err := svc.UpdateBookmark(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update example: %v", err)
	}

	if resp.Description != "new description" {
		t.Errorf("Expected description 'new description', got %s", resp.Description)
	}
}

func TestUpdateBookmarkChangeCommand(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create an example
	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	svc.CreateBookmark(ctx, req)

	// Update with new command (primary key change)
	updateReq := dto.UpdateBookmarkRequest{
		Command:    "kubectl get pods",
		NewCommand: "kubectl get pods -A",
	}

	resp, err := svc.UpdateBookmark(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update example: %v", err)
	}

	if resp.Command != "kubectl get pods -A" {
		t.Errorf("Expected command 'kubectl get pods -A', got %s", resp.Command)
	}

	// Verify old command doesn't exist
	_, err = svc.GetBookmark(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Old command should not exist after update")
	}

	// Verify new command exists
	_, err = svc.GetBookmark(ctx, "kubectl get pods -A")
	if err != nil {
		t.Error("New command should exist after update")
	}
}

func TestUpdateBookmarkNotFound(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	updateReq := dto.UpdateBookmarkRequest{
		Command:        "nonexistent",
		NewDescription: "test",
	}

	_, err := svc.UpdateBookmark(ctx, updateReq)
	if err == nil {
		t.Error("Expected error when updating nonexistent example")
	}
}

func TestDeleteBookmark(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create an example
	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	svc.CreateBookmark(ctx, req)

	// Delete the example
	err := svc.DeleteBookmark(ctx, "kubectl get pods")
	if err != nil {
		t.Fatalf("Failed to delete example: %v", err)
	}

	// Verify it's gone
	_, err = svc.GetBookmark(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Example should not exist after deletion")
	}
}

func TestDeleteBookmarkNotFound(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	err := svc.DeleteBookmark(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error when deleting nonexistent example")
	}
}

func TestDeleteToolBookmarks(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create multiple examples for same tool
	examples := []dto.CreateBookmarkRequest{
		{Command: "kubectl get pods", ToolName: "kubectl", Description: "list pods"},
		{Command: "kubectl get nodes", ToolName: "kubectl", Description: "list nodes"},
		{Command: "docker ps", ToolName: "docker", Description: "list containers"},
	}

	for _, req := range examples {
		svc.CreateBookmark(ctx, req)
	}

	// Delete all kubectl examples
	err := svc.DeleteToolBookmarks(ctx, "kubectl")
	if err != nil {
		t.Fatalf("Failed to delete tool examples: %v", err)
	}

	// Verify kubectl examples are gone
	resp, _ := svc.ListBookmarks(ctx)
	if resp.Count != 1 {
		t.Errorf("Expected 1 example remaining, got %d", resp.Count)
	}

	if resp.Examples[0].ToolName != "docker" {
		t.Errorf("Expected docker example to remain, got %s", resp.Examples[0].ToolName)
	}
}

// Additional tests to improve coverage

func TestUpdateBookmarkCommandConflict(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create two bookmarks
	req1 := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}
	req2 := dto.CreateBookmarkRequest{
		Command:     "kubectl get nodes",
		ToolName:    "kubectl",
		Description: "list nodes",
	}

	svc.CreateBookmark(ctx, req1)
	svc.CreateBookmark(ctx, req2)

	// Try to update first bookmark to use the same command as second
	updateReq := dto.UpdateBookmarkRequest{
		Command:    "kubectl get pods",
		NewCommand: "kubectl get nodes", // This already exists
	}

	_, err := svc.UpdateBookmark(ctx, updateReq)
	if err == nil {
		t.Error("Expected error when updating to an existing command")
	}
	if err != nil && !contains(err.Error(), "already exists") {
		t.Errorf("Expected 'already exists' error, got: %v", err)
	}
}

func TestUpdateBookmarkOnlyDescription(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create bookmark
	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "old description",
	}

	created, _ := svc.CreateBookmark(ctx, req)

	// Update only description
	updateReq := dto.UpdateBookmarkRequest{
		Command:        "kubectl get pods",
		NewDescription: "new description",
	}

	updated, err := svc.UpdateBookmark(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update bookmark: %v", err)
	}

	if updated.Description != "new description" {
		t.Errorf("Expected description 'new description', got %s", updated.Description)
	}
	if updated.ToolName != created.ToolName {
		t.Errorf("Tool name should not change, expected %s, got %s", created.ToolName, updated.ToolName)
	}
}

func TestUpdateBookmarkOnlyToolName(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create bookmark
	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	svc.CreateBookmark(ctx, req)

	// Update only tool name
	updateReq := dto.UpdateBookmarkRequest{
		Command:     "kubectl get pods",
		NewToolName: "k8s",
	}

	updated, err := svc.UpdateBookmark(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update bookmark: %v", err)
	}

	if updated.ToolName != "k8s" {
		t.Errorf("Expected tool name 'k8s', got %s", updated.ToolName)
	}
}

func TestUpdateBookmarkAllFields(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Create bookmark
	req := dto.CreateBookmarkRequest{
		Command:     "kubectl get pods",
		ToolName:    "kubectl",
		Description: "list pods",
	}

	svc.CreateBookmark(ctx, req)

	// Update all fields
	updateReq := dto.UpdateBookmarkRequest{
		Command:        "kubectl get pods",
		NewCommand:     "k get pods -A",
		NewToolName:    "k8s",
		NewDescription: "list all pods in all namespaces",
	}

	updated, err := svc.UpdateBookmark(ctx, updateReq)
	if err != nil {
		t.Fatalf("Failed to update bookmark: %v", err)
	}

	if updated.Command != "k get pods -A" {
		t.Errorf("Expected command 'k get pods -A', got %s", updated.Command)
	}
	if updated.ToolName != "k8s" {
		t.Errorf("Expected tool name 'k8s', got %s", updated.ToolName)
	}
	if updated.Description != "list all pods in all namespaces" {
		t.Errorf("Expected new description, got %s", updated.Description)
	}

	// Verify old command doesn't exist
	_, err = svc.GetBookmark(ctx, "kubectl get pods")
	if err == nil {
		t.Error("Old command should not exist after update")
	}
}

func TestDeleteToolBookmarksNotFound(t *testing.T) {
	repo := newMockBookmarkRepository()
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	// Try to delete bookmarks for non-existent tool
	err := svc.DeleteToolBookmarks(ctx, "nonexistent")
	if err == nil {
		t.Error("Expected error when deleting bookmarks for non-existent tool")
	}
}

func TestCreateBookmarkRepositoryError(t *testing.T) {
	// Use a mock that returns errors for Exists check
	repo := &errorMockRepository{shouldErrorOnExists: true}
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	req := dto.CreateBookmarkRequest{
		Command:     "test command",
		ToolName:    "test",
		Description: "test description",
	}

	_, err := svc.CreateBookmark(ctx, req)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

func TestListBookmarksRepositoryError(t *testing.T) {
	repo := &errorMockRepository{shouldErrorOnList: true}
	svc := NewBookmarkService(repo)
	ctx := context.Background()

	_, err := svc.ListBookmarks(ctx)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsRec(s, substr))
}

func containsRec(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Error mock repository for testing error paths
type errorMockRepository struct {
	shouldErrorOnExists bool
	shouldErrorOnList   bool
}

func (m *errorMockRepository) Create(ctx context.Context, example *models.Bookmark) error {
	return errors.New("mock create error")
}

func (m *errorMockRepository) GetByCommand(ctx context.Context, command string) (*models.Bookmark, error) {
	return nil, errors.New("mock get error")
}

func (m *errorMockRepository) List(ctx context.Context) ([]*models.Bookmark, error) {
	if m.shouldErrorOnList {
		return nil, errors.New("mock list error")
	}
	return []*models.Bookmark{}, nil
}

func (m *errorMockRepository) ListByToolName(ctx context.Context, toolName string) ([]*models.Bookmark, error) {
	return nil, errors.New("mock list by tool error")
}

func (m *errorMockRepository) Update(ctx context.Context, example *models.Bookmark) error {
	return errors.New("mock update error")
}

func (m *errorMockRepository) Delete(ctx context.Context, command string) error {
	return errors.New("mock delete error")
}

func (m *errorMockRepository) DeleteByToolName(ctx context.Context, toolName string) error {
	return errors.New("mock delete by tool error")
}

func (m *errorMockRepository) Exists(ctx context.Context, command string) (bool, error) {
	if m.shouldErrorOnExists {
		return false, errors.New("mock exists error")
	}
	return false, nil
}
