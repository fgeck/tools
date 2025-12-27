package models

import "time"

// Tool represents the core domain entity for a CLI tool bookmark
type Tool struct {
	ID          string    // Unique identifier (UUID)
	Name        string    // Display name
	Command     string    // Executable path
	Description string    // Tool description
	Examples    []string  // Usage examples
	CreatedAt   time.Time // Creation timestamp
	UpdatedAt   time.Time // Last update timestamp
}
