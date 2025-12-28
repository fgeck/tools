package dto

// CreateBookmarkRequest - DTO for creating a new example
type CreateBookmarkRequest struct {
	Command     string `json:"command" yaml:"command"`         // The actual command (primary key)
	ToolName    string `json:"tool_name" yaml:"tool_name"`     // Tool name for grouping
	Description string `json:"description" yaml:"description"` // What this example does
}

// BookmarkResponse - DTO for returning example data
type BookmarkResponse struct {
	Command     string `json:"command" yaml:"command"`
	ToolName    string `json:"tool_name" yaml:"tool_name"`
	Description string `json:"description" yaml:"description"`
}

// UpdateBookmarkRequest - DTO for updating an existing example
type UpdateBookmarkRequest struct {
	Command        string `json:"command" yaml:"command"`                 // The command to update (primary key)
	NewToolName    string `json:"new_tool_name" yaml:"new_tool_name"`     // New tool name (optional)
	NewDescription string `json:"new_description" yaml:"new_description"` // New description (optional)
	NewCommand     string `json:"new_command" yaml:"new_command"`         // New command (optional)
}

// ListBookmarksResponse - DTO for listing multiple examples
type ListBookmarksResponse struct {
	Examples []BookmarkResponse `json:"examples" yaml:"examples"`
	Count    int                `json:"count" yaml:"count"`
}
