package models

import "time"

// ToolExample represents a single bookmarked command example
// The command itself is the unique identifier (primary key)
type ToolExample struct {
	Command     string    // PRIMARY KEY - The actual command to execute (e.g., "lsof -i :54321")
	ToolName    string    // Tool name for grouping (e.g., "lsof")
	Description string    // What this example does
	CreatedAt   time.Time // Creation timestamp
	UpdatedAt   time.Time // Last update timestamp
}
