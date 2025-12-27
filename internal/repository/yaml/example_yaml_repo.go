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
	// ErrExampleNotFound is returned when an example is not found
	ErrExampleNotFound = errors.New("example not found")
	// ErrExampleAlreadyExists is returned when attempting to create a duplicate example
	ErrExampleAlreadyExists = errors.New("example with this command already exists")
)

// YAMLExampleRepository implements ExampleRepository using YAML file storage
type YAMLExampleRepository struct {
	filePath string
	mu       sync.RWMutex // Thread-safe operations
}

// yamlStorage represents the file structure
type yamlStorage struct {
	Examples []models.ToolExample `yaml:"examples"`
}

// NewYAMLExampleRepository creates a new YAML-based repository
func NewYAMLExampleRepository(filePath string) (repository.ExampleRepository, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	repo := &YAMLExampleRepository{
		filePath: filePath,
	}

	// Initialize file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := repo.save(&yamlStorage{Examples: []models.ToolExample{}}); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// load reads the YAML file and returns the storage structure
func (r *YAMLExampleRepository) load() (*yamlStorage, error) {
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
func (r *YAMLExampleRepository) save(storage *yamlStorage) error {
	data, err := yaml.Marshal(storage)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write storage file: %w", err)
	}

	return nil
}

// Create adds a new example to storage
func (r *YAMLExampleRepository) Create(ctx context.Context, example *models.ToolExample) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	// Check for duplicates (command is primary key)
	for _, ex := range storage.Examples {
		if ex.Command == example.Command {
			return ErrExampleAlreadyExists
		}
	}

	storage.Examples = append(storage.Examples, *example)
	return r.save(storage)
}

// GetByCommand retrieves an example by its command
func (r *YAMLExampleRepository) GetByCommand(ctx context.Context, command string) (*models.ToolExample, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, ex := range storage.Examples {
		if ex.Command == command {
			return &ex, nil
		}
	}

	return nil, ErrExampleNotFound
}

// List retrieves all examples
func (r *YAMLExampleRepository) List(ctx context.Context) ([]*models.ToolExample, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	examples := make([]*models.ToolExample, len(storage.Examples))
	for i := range storage.Examples {
		examples[i] = &storage.Examples[i]
	}

	return examples, nil
}

// ListByToolName retrieves all examples for a specific tool name
func (r *YAMLExampleRepository) ListByToolName(ctx context.Context, toolName string) ([]*models.ToolExample, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	var examples []*models.ToolExample
	for i := range storage.Examples {
		if storage.Examples[i].ToolName == toolName {
			examples = append(examples, &storage.Examples[i])
		}
	}

	return examples, nil
}

// Update modifies an existing example
func (r *YAMLExampleRepository) Update(ctx context.Context, example *models.ToolExample) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, ex := range storage.Examples {
		if ex.Command == example.Command {
			storage.Examples[i] = *example
			return r.save(storage)
		}
	}

	return ErrExampleNotFound
}

// Delete removes an example by command
func (r *YAMLExampleRepository) Delete(ctx context.Context, command string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, ex := range storage.Examples {
		if ex.Command == command {
			storage.Examples = append(storage.Examples[:i], storage.Examples[i+1:]...)
			return r.save(storage)
		}
	}

	return ErrExampleNotFound
}

// DeleteByToolName removes all examples for a tool name
func (r *YAMLExampleRepository) DeleteByToolName(ctx context.Context, toolName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	// Filter out examples matching the tool name
	filtered := []models.ToolExample{}
	found := false
	for _, ex := range storage.Examples {
		if ex.ToolName != toolName {
			filtered = append(filtered, ex)
		} else {
			found = true
		}
	}

	if !found {
		return ErrExampleNotFound
	}

	storage.Examples = filtered
	return r.save(storage)
}

// Exists checks if an example with the given command exists
func (r *YAMLExampleRepository) Exists(ctx context.Context, command string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return false, err
	}

	for _, ex := range storage.Examples {
		if ex.Command == command {
			return true, nil
		}
	}

	return false, nil
}
