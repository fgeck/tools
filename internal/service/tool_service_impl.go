package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fgeck/tools/internal/domain/models"
	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/repository"
	"github.com/google/uuid"
)

type toolServiceImpl struct {
	repo repository.ToolRepository
}

// NewToolService creates a new tool service instance
func NewToolService(repo repository.ToolRepository) ToolService {
	return &toolServiceImpl{
		repo: repo,
	}
}

// CreateTool implements business logic for creating a tool
func (s *toolServiceImpl) CreateTool(ctx context.Context, req dto.CreateToolRequest) (*dto.ToolResponse, error) {
	// Validation
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Check for duplicates
	exists, err := s.repo.Exists(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check tool existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("tool with name '%s' already exists", req.Name)
	}

	// Create domain model
	now := time.Now()
	tool := &models.Tool{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Command:     req.Command,
		Description: req.Description,
		Examples:    req.Examples,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Persist
	if err := s.repo.Create(ctx, tool); err != nil {
		return nil, fmt.Errorf("failed to create tool: %w", err)
	}

	// Convert to DTO
	return s.modelToDTO(tool), nil
}

// GetTool retrieves a tool by name
func (s *toolServiceImpl) GetTool(ctx context.Context, name string) (*dto.ToolResponse, error) {
	tool, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}

	return s.modelToDTO(tool), nil
}

// ListTools retrieves all tools
func (s *toolServiceImpl) ListTools(ctx context.Context) (*dto.ListToolsResponse, error) {
	tools, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	responses := make([]dto.ToolResponse, len(tools))
	for i, tool := range tools {
		responses[i] = *s.modelToDTO(tool)
	}

	return &dto.ListToolsResponse{
		Tools: responses,
		Count: len(responses),
	}, nil
}

// UpdateTool modifies an existing tool
func (s *toolServiceImpl) UpdateTool(ctx context.Context, name string, req dto.UpdateToolRequest) (*dto.ToolResponse, error) {
	// Get existing tool
	tool, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		tool.Name = *req.Name
	}
	if req.Command != nil {
		tool.Command = *req.Command
	}
	if req.Description != nil {
		tool.Description = *req.Description
	}
	if req.Examples != nil {
		tool.Examples = *req.Examples
	}
	tool.UpdatedAt = time.Now()

	// Persist
	if err := s.repo.Update(ctx, tool); err != nil {
		return nil, fmt.Errorf("failed to update tool: %w", err)
	}

	return s.modelToDTO(tool), nil
}

// DeleteTool removes a tool by name
func (s *toolServiceImpl) DeleteTool(ctx context.Context, name string) error {
	if err := s.repo.DeleteByName(ctx, name); err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}

	return nil
}

// validateCreateRequest validates the create tool request
func (s *toolServiceImpl) validateCreateRequest(req dto.CreateToolRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if strings.TrimSpace(req.Command) == "" {
		return fmt.Errorf("tool command cannot be empty")
	}
	return nil
}

// modelToDTO converts a domain model to a DTO
func (s *toolServiceImpl) modelToDTO(tool *models.Tool) *dto.ToolResponse {
	return &dto.ToolResponse{
		ID:          tool.ID,
		Name:        tool.Name,
		Command:     tool.Command,
		Description: tool.Description,
		Examples:    tool.Examples,
		CreatedAt:   tool.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   tool.UpdatedAt.Format(time.RFC3339),
	}
}
