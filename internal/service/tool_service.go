package service

import (
	"context"

	"github.com/fgeck/tools/internal/dto"
)

// ToolService defines business logic operations (CLI and REST API agnostic)
type ToolService interface {
	// CreateTool adds a new tool bookmark
	CreateTool(ctx context.Context, req dto.CreateToolRequest) (*dto.ToolResponse, error)

	// GetTool retrieves a tool by name
	GetTool(ctx context.Context, name string) (*dto.ToolResponse, error)

	// ListTools retrieves all tools
	ListTools(ctx context.Context) (*dto.ListToolsResponse, error)

	// UpdateTool modifies an existing tool
	UpdateTool(ctx context.Context, name string, req dto.UpdateToolRequest) (*dto.ToolResponse, error)

	// DeleteTool removes a tool by name
	DeleteTool(ctx context.Context, name string) error
}
