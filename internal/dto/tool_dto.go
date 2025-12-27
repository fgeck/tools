package dto

// CreateToolRequest - DTO for creating a new tool
type CreateToolRequest struct {
	Name        string   `json:"name" yaml:"name"`
	Command     string   `json:"command" yaml:"command"`
	Description string   `json:"description" yaml:"description"`
	Examples    []string `json:"examples" yaml:"examples"`
}

// UpdateToolRequest - DTO for updating a tool
type UpdateToolRequest struct {
	Name        *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Command     *string   `json:"command,omitempty" yaml:"command,omitempty"`
	Description *string   `json:"description,omitempty" yaml:"description,omitempty"`
	Examples    *[]string `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// ToolResponse - DTO for returning tool data
type ToolResponse struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Command     string   `json:"command" yaml:"command"`
	Description string   `json:"description" yaml:"description"`
	Examples    []string `json:"examples" yaml:"examples"`
	CreatedAt   string   `json:"created_at" yaml:"created_at"`
	UpdatedAt   string   `json:"updated_at" yaml:"updated_at"`
}

// ListToolsResponse - DTO for listing multiple tools
type ListToolsResponse struct {
	Tools []ToolResponse `json:"tools" yaml:"tools"`
	Count int            `json:"count" yaml:"count"`
}
