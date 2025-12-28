package models

// Bookmark represents a single bookmarked command
// The command string itself is the unique identifier (primary key)
type Bookmark struct {
	Command     string // PRIMARY KEY - The actual command to execute (e.g., "lsof -i :54321")
	ToolName    string // Tool name for grouping (e.g., "lsof")
	Description string // What this bookmark does
}
