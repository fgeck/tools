package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fgeck/tools/internal/domain/models"
	"github.com/fgeck/tools/internal/dto"
	"github.com/fgeck/tools/internal/repository"
)

type exampleServiceImpl struct {
	repo repository.ExampleRepository
}

// NewExampleService creates a new example service instance
func NewExampleService(repo repository.ExampleRepository) ExampleService {
	return &exampleServiceImpl{
		repo: repo,
	}
}

// CreateExample implements business logic for creating an example
func (s *exampleServiceImpl) CreateExample(ctx context.Context, req dto.CreateExampleRequest) (*dto.ExampleResponse, error) {
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
	now := time.Now()
	example := &models.ToolExample{
		Command:     req.Command,
		ToolName:    req.ToolName,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Persist
	if err := s.repo.Create(ctx, example); err != nil {
		return nil, fmt.Errorf("failed to create example: %w", err)
	}

	// Convert to DTO
	return s.modelToDTO(example), nil
}

// GetExample retrieves an example by command
func (s *exampleServiceImpl) GetExample(ctx context.Context, command string) (*dto.ExampleResponse, error) {
	example, err := s.repo.GetByCommand(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to get example: %w", err)
	}

	return s.modelToDTO(example), nil
}

// ListExamples retrieves all examples
func (s *exampleServiceImpl) ListExamples(ctx context.Context) (*dto.ListExamplesResponse, error) {
	examples, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list examples: %w", err)
	}

	responses := make([]dto.ExampleResponse, len(examples))
	for i, example := range examples {
		responses[i] = *s.modelToDTO(example)
	}

	return &dto.ListExamplesResponse{
		Examples: responses,
		Count:    len(responses),
	}, nil
}

// UpdateExample modifies an existing example
func (s *exampleServiceImpl) UpdateExample(ctx context.Context, req dto.UpdateExampleRequest) (*dto.ExampleResponse, error) {
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
			existing.UpdatedAt = time.Now()
			if err := s.repo.Create(ctx, existing); err != nil {
				return nil, fmt.Errorf("failed to create updated example: %w", err)
			}
			return s.modelToDTO(existing), nil
		}
	}

	// Update timestamp
	existing.UpdatedAt = time.Now()

	// Persist changes
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update example: %w", err)
	}

	return s.modelToDTO(existing), nil
}

// DeleteExample removes an example by command
func (s *exampleServiceImpl) DeleteExample(ctx context.Context, command string) error {
	if err := s.repo.Delete(ctx, command); err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}

	return nil
}

// DeleteToolExamples removes all examples for a tool name
func (s *exampleServiceImpl) DeleteToolExamples(ctx context.Context, toolName string) error {
	if err := s.repo.DeleteByToolName(ctx, toolName); err != nil {
		return fmt.Errorf("failed to delete tool examples: %w", err)
	}

	return nil
}

// validateCreateRequest validates the create example request
func (s *exampleServiceImpl) validateCreateRequest(req dto.CreateExampleRequest) error {
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
func (s *exampleServiceImpl) modelToDTO(example *models.ToolExample) *dto.ExampleResponse {
	return &dto.ExampleResponse{
		Command:     example.Command,
		ToolName:    example.ToolName,
		Description: example.Description,
		CreatedAt:   example.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   example.UpdatedAt.Format(time.RFC3339),
	}
}
