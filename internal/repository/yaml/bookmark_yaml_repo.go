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
	// ErrBookmarkNotFound is returned when an example is not found
	ErrBookmarkNotFound = errors.New("bookmark not found")
	// ErrBookmarkAlreadyExists is returned when attempting to create a duplicate example
	ErrBookmarkAlreadyExists = errors.New("example with this command already exists")
)

// YAMLBookmarkRepository implements BookmarkRepository using YAML file storage
type YAMLBookmarkRepository struct {
	filePath string
	mu       sync.RWMutex // Thread-safe operations
}

// yamlStorage represents the file structure
type yamlStorage struct {
	Bookmarks []models.Bookmark `yaml:"bookmarks"`
}

// NewYAMLBookmarkRepository creates a new YAML-based repository
func NewYAMLBookmarkRepository(filePath string) (repository.BookmarkRepository, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	repo := &YAMLBookmarkRepository{
		filePath: filePath,
	}

	// Initialize file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := repo.save(&yamlStorage{Bookmarks: []models.Bookmark{}}); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// load reads the YAML file and returns the storage structure
func (r *YAMLBookmarkRepository) load() (*yamlStorage, error) {
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
func (r *YAMLBookmarkRepository) save(storage *yamlStorage) error {
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
func (r *YAMLBookmarkRepository) Create(ctx context.Context, example *models.Bookmark) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	// Check for duplicates (command is primary key)
	for _, ex := range storage.Bookmarks {
		if ex.Command == example.Command {
			return ErrBookmarkAlreadyExists
		}
	}

	storage.Bookmarks = append(storage.Bookmarks, *example)
	return r.save(storage)
}

// GetByCommand retrieves an example by its command
func (r *YAMLBookmarkRepository) GetByCommand(ctx context.Context, command string) (*models.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	for _, ex := range storage.Bookmarks {
		if ex.Command == command {
			return &ex, nil
		}
	}

	return nil, ErrBookmarkNotFound
}

// List retrieves all examples
func (r *YAMLBookmarkRepository) List(ctx context.Context) ([]*models.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	examples := make([]*models.Bookmark, len(storage.Bookmarks))
	for i := range storage.Bookmarks {
		examples[i] = &storage.Bookmarks[i]
	}

	return examples, nil
}

// ListByToolName retrieves all examples for a specific tool name
func (r *YAMLBookmarkRepository) ListByToolName(ctx context.Context, toolName string) ([]*models.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return nil, err
	}

	var examples []*models.Bookmark
	for i := range storage.Bookmarks {
		if storage.Bookmarks[i].ToolName == toolName {
			examples = append(examples, &storage.Bookmarks[i])
		}
	}

	return examples, nil
}

// Update modifies an existing example
func (r *YAMLBookmarkRepository) Update(ctx context.Context, example *models.Bookmark) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, ex := range storage.Bookmarks {
		if ex.Command == example.Command {
			storage.Bookmarks[i] = *example
			return r.save(storage)
		}
	}

	return ErrBookmarkNotFound
}

// Delete removes an example by command
func (r *YAMLBookmarkRepository) Delete(ctx context.Context, command string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	for i, ex := range storage.Bookmarks {
		if ex.Command == command {
			storage.Bookmarks = append(storage.Bookmarks[:i], storage.Bookmarks[i+1:]...)
			return r.save(storage)
		}
	}

	return ErrBookmarkNotFound
}

// DeleteByToolName removes all examples for a tool name
func (r *YAMLBookmarkRepository) DeleteByToolName(ctx context.Context, toolName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storage, err := r.load()
	if err != nil {
		return err
	}

	// Filter out examples matching the tool name
	filtered := []models.Bookmark{}
	found := false
	for _, ex := range storage.Bookmarks {
		if ex.ToolName != toolName {
			filtered = append(filtered, ex)
		} else {
			found = true
		}
	}

	if !found {
		return ErrBookmarkNotFound
	}

	storage.Bookmarks = filtered
	return r.save(storage)
}

// Exists checks if an example with the given command exists
func (r *YAMLBookmarkRepository) Exists(ctx context.Context, command string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storage, err := r.load()
	if err != nil {
		return false, err
	}

	for _, ex := range storage.Bookmarks {
		if ex.Command == command {
			return true, nil
		}
	}

	return false, nil
}
