package repository

import (
	"context"

	"github.com/fgeck/tools/internal/domain/models"
)

// ToolRepository defines the interface for tool persistence
// This abstraction allows easy swapping between YAML, SQLite, PostgreSQL, etc.
type ToolRepository interface {
	// Create adds a new tool to storage
	Create(ctx context.Context, tool *models.Tool) error

	// GetByID retrieves a tool by its ID
	GetByID(ctx context.Context, id string) (*models.Tool, error)

	// GetByName retrieves a tool by its name (for user-friendly lookup)
	GetByName(ctx context.Context, name string) (*models.Tool, error)

	// List retrieves all tools
	List(ctx context.Context) ([]*models.Tool, error)

	// Update modifies an existing tool
	Update(ctx context.Context, tool *models.Tool) error

	// Delete removes a tool by ID
	Delete(ctx context.Context, id string) error

	// DeleteByName removes a tool by name (convenience for CLI)
	DeleteByName(ctx context.Context, name string) error

	// Exists checks if a tool with the given name exists
	Exists(ctx context.Context, name string) (bool, error)
}
