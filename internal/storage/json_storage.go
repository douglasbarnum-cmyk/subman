package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"subman/internal/models"
)

const (
	defaultFileName = "subscriptions.json"
	dataVersion     = "1.0"
)

type JSONStorage struct {
	filePath string
	mu       sync.RWMutex
}

// NewJSONStorage creates a new JSON storage instance
// Stores in user's config directory (platform-aware)
func NewJSONStorage() (*JSONStorage, error) {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// Create subman directory
	submanDir := filepath.Join(configDir, "subman")
	if err := os.MkdirAll(submanDir, 0700); err != nil {
		return nil, err
	}

	filePath := filepath.Join(submanDir, defaultFileName)

	return &JSONStorage{
		filePath: filePath,
	}, nil
}

// NewJSONStorageWithPath creates storage at a specific path
func NewJSONStorageWithPath(path string) *JSONStorage {
	return &JSONStorage{
		filePath: path,
	}
}

func (s *JSONStorage) Load() (*models.SubscriptionList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// If file doesn't exist, return empty list
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return &models.SubscriptionList{
			Subscriptions: []models.Subscription{},
			Version:       dataVersion,
		}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var list models.SubscriptionList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (s *JSONStorage) Save(list *models.SubscriptionList) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	list.Version = dataVersion

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0600)
}

func (s *JSONStorage) GetPath() string {
	return s.filePath
}
