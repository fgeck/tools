package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/fgeck/tools/internal/domain/models"
	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/repository"
)

type bookmarkServiceImpl struct {
	repo repository.BookmarkRepository
}

// NewBookmarkService creates a new example service instance
func NewBookmarkService(repo repository.BookmarkRepository) BookmarkService {
	return &bookmarkServiceImpl{
		repo: repo,
	}
}

// CreateBookmark implements business logic for creating an example
func (s *bookmarkServiceImpl) CreateBookmark(ctx context.Context, req dto.CreateBookmarkRequest) (*dto.BookmarkResponse, error) {
	// Validation
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check if command already exists
	exists, err := s.repo.Exists(ctx, req.Command)
	if err != nil {
		return nil, fmt.Errorf("failed to check example existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("example with command '%s' already exists", req.Command)
	}

	// Create domain model
	example := &models.Bookmark{
		Command:     req.Command,
		ToolName:    req.ToolName,
		Description: req.Description,
	}

	// Persist
	if err := s.repo.Create(ctx, example); err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	// Convert to DTO
	return s.modelToDTO(example), nil
}

// GetBookmark retrieves an example by command
func (s *bookmarkServiceImpl) GetBookmark(ctx context.Context, command string) (*dto.BookmarkResponse, error) {
	example, err := s.repo.GetByCommand(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	return s.modelToDTO(example), nil
}

// ListBookmarks retrieves all examples
func (s *bookmarkServiceImpl) ListBookmarks(ctx context.Context) (*dto.ListBookmarksResponse, error) {
	examples, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list examples: %w", err)
	}

	responses := make([]dto.BookmarkResponse, len(examples))
	for i, example := range examples {
		responses[i] = *s.modelToDTO(example)
	}

	return &dto.ListBookmarksResponse{
		Examples: responses,
		Count:    len(responses),
	}, nil
}

// UpdateBookmark modifies an existing example
func (s *bookmarkServiceImpl) UpdateBookmark(ctx context.Context, req dto.UpdateBookmarkRequest) (*dto.BookmarkResponse, error) {
	// Get existing example
	existing, err := s.repo.GetByCommand(ctx, req.Command)
	if err != nil {
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	// Update fields if provided
	if req.NewToolName != "" {
		existing.ToolName = req.NewToolName
	}
	if req.NewDescription != "" {
		existing.Description = req.NewDescription
	}
	if req.NewCommand != "" {
		// If changing the command (primary key), check for conflicts
		if req.NewCommand != req.Command {
			exists, err := s.repo.Exists(ctx, req.NewCommand)
			if err != nil {
				return nil, fmt.Errorf("failed to check new command existence: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("example with command '%s' already exists", req.NewCommand)
			}
			// Delete old entry and create new one with new command
			if err := s.repo.Delete(ctx, req.Command); err != nil {
				return nil, fmt.Errorf("failed to delete old example: %w", err)
			}
			existing.Command = req.NewCommand
			if err := s.repo.Create(ctx, existing); err != nil {
				return nil, fmt.Errorf("failed to create updated example: %w", err)
			}
			return s.modelToDTO(existing), nil
		}
	}

	// Persist changes
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update example: %w", err)
	}

	return s.modelToDTO(existing), nil
}

// DeleteBookmark removes an example by command
func (s *bookmarkServiceImpl) DeleteBookmark(ctx context.Context, command string) error {
	if err := s.repo.Delete(ctx, command); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

	return nil
}

// DeleteToolBookmarks removes all examples for a tool name
func (s *bookmarkServiceImpl) DeleteToolBookmarks(ctx context.Context, toolName string) error {
	if err := s.repo.DeleteByToolName(ctx, toolName); err != nil {
		return fmt.Errorf("failed to delete tool examples: %w", err)
	}

	return nil
}

// validateCreateRequest validates the create example request
func (s *bookmarkServiceImpl) validateCreateRequest(req dto.CreateBookmarkRequest) error {
	if strings.TrimSpace(req.Command) == "" {
		return fmt.Errorf("command cannot be empty")
	}
	if strings.TrimSpace(req.ToolName) == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if strings.TrimSpace(req.Description) == "" {
		return fmt.Errorf("description cannot be empty")
	}
	return nil
}

// modelToDTO converts a domain model to a DTO
func (s *bookmarkServiceImpl) modelToDTO(example *models.Bookmark) *dto.BookmarkResponse {
	return &dto.BookmarkResponse{
		Command:     example.Command,
		ToolName:    example.ToolName,
		Description: example.Description,
	}
}
