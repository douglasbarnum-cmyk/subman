package service

import (
	"time"

	"github.com/google/uuid"
	"subman/internal/models"
	"subman/internal/storage"
)

type PaymentService struct {
	storage storage.Storage
}

func NewPaymentService(storage storage.Storage) *PaymentService {
	return &PaymentService{
		storage: storage,
	}
}

// GeneratePaymentsForSubscription creates payment records for a subscription
// based on its billing cycle and start date up to the current date
func (p *PaymentService) GeneratePaymentsForSubscription(sub *models.Subscription) error {
	list, err := p.storage.Load()
	if err != nil {
		return err
	}

	// Don't generate payments for paused or deleted subscriptions
	if sub.Paused || sub.Deleted {
		return nil
	}

	// Get existing payments for this subscription
	existingPayments := p.getPaymentsForSubscription(list.Payments, sub.ID)

	// Find the last payment date
	lastPaymentDate := sub.StartDate
	if len(existingPayments) > 0 {
		// Find most recent payment
		for _, payment := range existingPayments {
			if payment.PaymentDate.After(lastPaymentDate) {
				lastPaymentDate = payment.PaymentDate
			}
		}
	}

	// Generate payments from last payment date to now
	now := time.Now()
	currentDate := lastPaymentDate

	for currentDate.Before(now) {
		// Calculate next payment date
		if sub.BillingCycle == models.Monthly {
			currentDate = currentDate.AddDate(0, 1, 0)
		} else {
			currentDate = currentDate.AddDate(1, 0, 0)
		}

		// Only create payment if it's in the past
		if currentDate.Before(now) || currentDate.Equal(now.Truncate(24*time.Hour)) {
			// Check if payment already exists for this date
			if !p.paymentExistsForDate(list.Payments, sub.ID, currentDate) {
				payment := models.Payment{
					ID:             uuid.New().String(),
					SubscriptionID: sub.ID,
					Amount:         sub.Cost,
					PaymentDate:    currentDate,
					Notes:          "Auto-generated",
					CreatedAt:      time.Now(),
				}
				list.Payments = append(list.Payments, payment)
			}
		}
	}

	// Calculate and update the next payment date
	nextPaymentDate := currentDate
	if !nextPaymentDate.After(now) {
		// If currentDate is not after now, calculate the next billing cycle
		if sub.BillingCycle == models.Monthly {
			nextPaymentDate = currentDate.AddDate(0, 1, 0)
		} else {
			nextPaymentDate = currentDate.AddDate(1, 0, 0)
		}
	}

	// Update the subscription's NextPayment field
	for i := range list.Subscriptions {
		if list.Subscriptions[i].ID == sub.ID {
			list.Subscriptions[i].NextPayment = nextPaymentDate
			list.Subscriptions[i].UpdatedAt = time.Now()
			break
		}
	}

	return p.storage.Save(list)
}

// GenerateAllPayments generates payments for all active subscriptions
func (p *PaymentService) GenerateAllPayments() error {
	list, err := p.storage.Load()
	if err != nil {
		return err
	}

	for _, sub := range list.Subscriptions {
		if !sub.Paused && !sub.Deleted {
			if err := p.GeneratePaymentsForSubscription(&sub); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetPaymentsForSubscription returns all payments for a subscription
func (p *PaymentService) GetPaymentsForSubscription(subscriptionID string) ([]models.Payment, error) {
	list, err := p.storage.Load()
	if err != nil {
		return nil, err
	}

	return p.getPaymentsForSubscription(list.Payments, subscriptionID), nil
}

// GetYTDPayments returns all payments from January 1 of current year to today
func (p *PaymentService) GetYTDPayments() ([]models.Payment, error) {
	list, err := p.storage.Load()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	var ytdPayments []models.Payment
	for _, payment := range list.Payments {
		if payment.PaymentDate.After(startOfYear) || payment.PaymentDate.Equal(startOfYear) {
			if payment.PaymentDate.Before(now) || payment.PaymentDate.Equal(now.Truncate(24*time.Hour)) {
				ytdPayments = append(ytdPayments, payment)
			}
		}
	}

	return ytdPayments, nil
}

// Helper functions

func (p *PaymentService) getPaymentsForSubscription(payments []models.Payment, subscriptionID string) []models.Payment {
	var result []models.Payment
	for _, payment := range payments {
		if payment.SubscriptionID == subscriptionID {
			result = append(result, payment)
		}
	}
	return result
}

func (p *PaymentService) paymentExistsForDate(payments []models.Payment, subscriptionID string, date time.Time) bool {
	dateOnly := date.Truncate(24 * time.Hour)
	for _, payment := range payments {
		if payment.SubscriptionID == subscriptionID {
			paymentDateOnly := payment.PaymentDate.Truncate(24 * time.Hour)
			if paymentDateOnly.Equal(dateOnly) {
				return true
			}
		}
	}
	return false
}
