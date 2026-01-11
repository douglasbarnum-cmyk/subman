package storage

import (
	"subman/internal/models"
)

// Storage defines the interface for subscription persistence
type Storage interface {
	// Load reads all subscriptions from storage
	Load() (*models.SubscriptionList, error)

	// Save writes all subscriptions to storage
	Save(list *models.SubscriptionList) error

	// GetPath returns the storage file path
	GetPath() string
}
