package repository

import (
	"context"

	"github.com/fgeck/tools/internal/domain/models"
)

// BookmarkRepository defines the interface for example persistence
// Command is the primary key for all operations
type BookmarkRepository interface {
	// Create adds a new example to storage
	// Returns error if command already exists
	Create(ctx context.Context, example *models.Bookmark) error

	// GetByCommand retrieves an example by its command (primary key)
	GetByCommand(ctx context.Context, command string) (*models.Bookmark, error)

	// List retrieves all examples
	List(ctx context.Context) ([]*models.Bookmark, error)

	// ListByToolName retrieves all examples for a specific tool name
	ListByToolName(ctx context.Context, toolName string) ([]*models.Bookmark, error)

	// Update modifies an existing example (identified by command)
	Update(ctx context.Context, example *models.Bookmark) error

	// Delete removes an example by command (primary key)
	Delete(ctx context.Context, command string) error

	// DeleteByToolName removes all examples for a tool name
	DeleteByToolName(ctx context.Context, toolName string) error

	// Exists checks if an example with the given command exists
	Exists(ctx context.Context, command string) (bool, error)
}
