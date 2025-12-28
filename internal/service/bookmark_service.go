package service

import (
	"context"

	"github.com/fgeck/tools/internal/dto"
)

// BookmarkService defines business logic operations (CLI and REST API agnostic)
type BookmarkService interface {
	// CreateBookmark adds a new example bookmark
	CreateBookmark(ctx context.Context, req dto.CreateBookmarkRequest) (*dto.BookmarkResponse, error)

	// GetBookmark retrieves an example by command
	GetBookmark(ctx context.Context, command string) (*dto.BookmarkResponse, error)

	// ListBookmarks retrieves all examples
	ListBookmarks(ctx context.Context) (*dto.ListBookmarksResponse, error)

	// UpdateBookmark modifies an existing example
	UpdateBookmark(ctx context.Context, req dto.UpdateBookmarkRequest) (*dto.BookmarkResponse, error)

	// DeleteBookmark removes an example by command
	DeleteBookmark(ctx context.Context, command string) error

	// DeleteToolBookmarks removes all examples for a tool name
	DeleteToolBookmarks(ctx context.Context, toolName string) error
}
