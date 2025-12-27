package service

import (
	"context"

	"github.com/fgeck/tools/internal/dto"
)

// ExampleService defines business logic operations (CLI and REST API agnostic)
type ExampleService interface {
	// CreateExample adds a new example bookmark
	CreateExample(ctx context.Context, req dto.CreateExampleRequest) (*dto.ExampleResponse, error)

	// GetExample retrieves an example by command
	GetExample(ctx context.Context, command string) (*dto.ExampleResponse, error)

	// ListExamples retrieves all examples
	ListExamples(ctx context.Context) (*dto.ListExamplesResponse, error)

	// UpdateExample modifies an existing example
	UpdateExample(ctx context.Context, req dto.UpdateExampleRequest) (*dto.ExampleResponse, error)

	// DeleteExample removes an example by command
	DeleteExample(ctx context.Context, command string) error

	// DeleteToolExamples removes all examples for a tool name
	DeleteToolExamples(ctx context.Context, toolName string) error
}
