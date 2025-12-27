package yaml

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fgeck/tools/internal/domain/models"
	"github.com/fgeck/tools/internal/repository"
	"gopkg.in/yaml.v3"
)

var (
	// ErrToolNotFound is returned when a tool is not found
	ErrToolNotFound = errors.New("tool not found")
	// ErrToolAlreadyExists is returned when attempting to create a duplicate tool
	ErrToolAlreadyExists = errors.New("tool already exists")
)

// YAMLToolRepository implements ToolRepository using YAML file storage
type YAMLToolRepository struct {
	filePath string
	mu       sync.RWMutex // Thread-safe operations
}

// yamlStorage represents the file structure
type yamlStorage struct {
	Tools []models.Tool `yaml:"tools"`
}

// NewYAMLToolRepository creates a new YAML-based repository
func NewYAMLToolRepository(filePath string) (repository.ToolRepository, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	repo := &YAMLToolRepository{
		filePath: filePath,
	}

	// Initialize file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := repo.save(&yamlStorage{Tools: []models.Tool{}}); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// load reads the YAML file and returns the storage structure
func (r *YAMLToolRepository) load() (*yamlStorage, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage file: %w", err)
	}

	var storage yamlStorage
	if err := yaml.Unmarshal(data, &storage); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &storage, nil
}

// save writes the storage structure to the YAML file
func (r *YAMLToolRepository) save(storage *yamlStorage) error {
	data, err := yaml.Marshal(storage)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}

	return nil
}

// Create adds a new tool to storage
func (r *YAMLToolRepository) Create(ctx context.Context, tool *models.Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	// Check for duplicates
	for _, t := range storage.Tools {
		if t.Name == tool.Name {
			return ErrToolAlreadyExists
		}
	}

	storage.Tools = append(storage.Tools, *tool)
	return r.save(storage)
}

// GetByID retrieves a tool by its ID
func (r *YAMLToolRepository) GetByID(ctx context.Context, id string) (*models.Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, t := range storage.Tools {
		if t.ID == id {
			return &t, nil
		}
	}

	return nil, ErrToolNotFound
}

// GetByName retrieves a tool by its name
func (r *YAMLToolRepository) GetByName(ctx context.Context, name string) (*models.Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, t := range storage.Tools {
		if t.Name == name {
			return &t, nil
		}
	}

	return nil, ErrToolNotFound
}

// List retrieves all tools
func (r *YAMLToolRepository) List(ctx context.Context) ([]*models.Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	tools := make([]*models.Tool, len(storage.Tools))
	for i := range storage.Tools {
		tools[i] = &storage.Tools[i]
	}

	return tools, nil
}

// Update modifies an existing tool
func (r *YAMLToolRepository) Update(ctx context.Context, tool *models.Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, t := range storage.Tools {
		if t.ID == tool.ID {
			storage.Tools[i] = *tool
			return r.save(storage)
		}
	}

	return ErrToolNotFound
}

// Delete removes a tool by ID
func (r *YAMLToolRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, t := range storage.Tools {
		if t.ID == id {
			storage.Tools = append(storage.Tools[:i], storage.Tools[i+1:]...)
			return r.save(storage)
		}
	}

	return ErrToolNotFound
}

// DeleteByName removes a tool by name
func (r *YAMLToolRepository) DeleteByName(ctx context.Context, name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, t := range storage.Tools {
		if t.Name == name {
			storage.Tools = append(storage.Tools[:i], storage.Tools[i+1:]...)
			return r.save(storage)
		}
	}

	return ErrToolNotFound
}

// Exists checks if a tool with the given name exists
func (r *YAMLToolRepository) Exists(ctx context.Context, name string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return false, err
	}

	for _, t := range storage.Tools {
		if t.Name == name {
			return true, nil
		}
	}

	return false, nil
}
