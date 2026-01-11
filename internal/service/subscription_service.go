package service

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"subman/internal/models"
	"subman/internal/storage"
	"subman/pkg/calculator"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidID            = errors.New("invalid subscription ID")
)

type SubscriptionService struct {
	storage storage.Storage
}

func NewSubscriptionService(storage storage.Storage) *SubscriptionService {
	return &SubscriptionService{
		storage: storage,
	}
}

// Create adds a new subscription
func (s *SubscriptionService) Create(sub *models.Subscription) error {
	sub.ID = uuid.New().String()
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()

	list, err := s.storage.Load()
	if err != nil {
		return err
	}

	list.Subscriptions = append(list.Subscriptions, *sub)
	return s.storage.Save(list)
}

// Update modifies an existing subscription
func (s *SubscriptionService) Update(sub *models.Subscription) error {
	if sub.ID == "" {
		return ErrInvalidID
	}

	list, err := s.storage.Load()
	if err != nil {
		return err
	}

	found := false
	for i, existing := range list.Subscriptions {
		if existing.ID == sub.ID {
			sub.CreatedAt = existing.CreatedAt
			sub.UpdatedAt = time.Now()
			list.Subscriptions[i] = *sub
			found = true
			break
		}
	}

	if !found {
		return ErrSubscriptionNotFound
	}

	return s.storage.Save(list)
}

// Delete marks a subscription as deleted (soft delete)
// Payments history is preserved for YTD calculations
func (s *SubscriptionService) Delete(id string) error {
	list, err := s.storage.Load()
	if err != nil {
		return err
	}

	found := false
	for i, sub := range list.Subscriptions {
		if sub.ID == id {
			list.Subscriptions[i].Deleted = true
			list.Subscriptions[i].DeletedAt = time.Now()
			list.Subscriptions[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return ErrSubscriptionNotFound
	}

	return s.storage.Save(list)
}

// Get retrieves a subscription by ID
func (s *SubscriptionService) Get(id string) (*models.Subscription, error) {
	list, err := s.storage.Load()
	if err != nil {
		return nil, err
	}

	for _, sub := range list.Subscriptions {
		if sub.ID == id {
			return &sub, nil
		}
	}

	return nil, ErrSubscriptionNotFound
}

// List returns all subscriptions with optional filtering and sorting
func (s *SubscriptionService) List(filter *models.FilterCriteria, sortBy models.SortField, order models.SortOrder) ([]models.Subscription, error) {
	list, err := s.storage.Load()
	if err != nil {
		return nil, err
	}

	// Apply filters
	filtered := s.filterSubscriptions(list.Subscriptions, filter)

	// Apply sorting
	s.sortSubscriptions(filtered, sortBy, order)

	return filtered, nil
}

// GetSummary calculates cost statistics including YTD from payment history
func (s *SubscriptionService) GetSummary() (*models.CostSummary, error) {
	list, err := s.storage.Load()
	if err != nil {
		return nil, err
	}

	return calculator.CalculateSummary(list.Subscriptions, list.Payments), nil
}

// filterSubscriptions applies filter criteria
func (s *SubscriptionService) filterSubscriptions(subs []models.Subscription, filter *models.FilterCriteria) []models.Subscription {
	if filter == nil {
		return subs
	}

	var result []models.Subscription

	for _, sub := range subs {
		// Always exclude deleted subscriptions
		if sub.Deleted {
			continue
		}

		// Paused filter - if ShowPaused is false, skip paused subscriptions
		if !filter.ShowPaused && sub.Paused {
			continue
		}

		// Search term filter (name or notes)
		if filter.SearchTerm != "" {
			term := strings.ToLower(filter.SearchTerm)
			if !strings.Contains(strings.ToLower(sub.Name), term) &&
				!strings.Contains(strings.ToLower(sub.Notes), term) {
				continue
			}
		}

		// Category filter
		if filter.Category != nil && sub.Category != *filter.Category {
			continue
		}

		// Billing cycle filter
		if filter.BillingCycle != nil && sub.BillingCycle != *filter.BillingCycle {
			continue
		}

		// Cost range filter
		if filter.MinCost != nil && sub.Cost < *filter.MinCost {
			continue
		}
		if filter.MaxCost != nil && sub.Cost > *filter.MaxCost {
			continue
		}

		result = append(result, sub)
	}

	return result
}

// sortSubscriptions sorts by field and order
func (s *SubscriptionService) sortSubscriptions(subs []models.Subscription, sortBy models.SortField, order models.SortOrder) {
	sort.Slice(subs, func(i, j int) bool {
		var less bool

		switch sortBy {
		case models.SortByName:
			less = strings.ToLower(subs[i].Name) < strings.ToLower(subs[j].Name)
		case models.SortByCost:
			less = subs[i].Cost < subs[j].Cost
		case models.SortByNextPayment:
			less = subs[i].NextPayment.Before(subs[j].NextPayment)
		default:
			less = subs[i].Name < subs[j].Name
		}

		if order == models.Descending {
			return !less
		}
		return less
	})
}
