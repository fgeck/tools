package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/fgeck/tools/internal/domain/models"
)

// MockToolRepository is a mock implementation for testing
type MockToolRepository struct {
	mu    sync.RWMutex
	tools map[string]*models.Tool // keyed by ID
}

// NewMockToolRepository creates a new mock repository
func NewMockToolRepository() ToolRepository {
	return &MockToolRepository{
		tools: make(map[string]*models.Tool),
	}
}

func (m *MockToolRepository) Create(ctx context.Context, tool *models.Tool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate name
	for _, t := range m.tools {
		if t.Name == tool.Name {
			return fmt.Errorf("tool with name '%s' already exists", tool.Name)
		}
	}

	m.tools[tool.ID] = tool
	return nil
}

func (m *MockToolRepository) GetByID(ctx context.Context, id string) (*models.Tool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tool, exists := m.tools[id]
	if !exists {
		return nil, fmt.Errorf("tool not found")
	}

	return tool, nil
}

func (m *MockToolRepository) GetByName(ctx context.Context, name string) (*models.Tool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tool := range m.tools {
		if tool.Name == name {
			return tool, nil
		}
	}

	return nil, fmt.Errorf("tool not found")
}

func (m *MockToolRepository) List(ctx context.Context) ([]*models.Tool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tools := make([]*models.Tool, 0, len(m.tools))
	for _, tool := range m.tools {
		tools = append(tools, tool)
	}

	return tools, nil
}

func (m *MockToolRepository) Update(ctx context.Context, tool *models.Tool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tools[tool.ID]; !exists {
		return fmt.Errorf("tool not found")
	}

	m.tools[tool.ID] = tool
	return nil
}

func (m *MockToolRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.tools[id]; !exists {
		return fmt.Errorf("tool not found")
	}

	delete(m.tools, id)
	return nil
}

func (m *MockToolRepository) DeleteByName(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, tool := range m.tools {
		if tool.Name == name {
			delete(m.tools, id)
			return nil
		}
	}

	return fmt.Errorf("tool not found")
}

func (m *MockToolRepository) Exists(ctx context.Context, name string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tool := range m.tools {
		if tool.Name == name {
			return true, nil
		}
	}

	return false, nil
}
