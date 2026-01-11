package models

import (
	"time"
)

// BillingCycle represents the frequency of billing
type BillingCycle string

const (
	Monthly BillingCycle = "monthly"
	Yearly  BillingCycle = "yearly"
)

// Category represents subscription categories
type Category string

const (
	Streaming  Category = "streaming"
	Software   Category = "software"
	Utilities  Category = "utilities"
	Gaming     Category = "gaming"
	News       Category = "news"
	Education  Category = "education"
	Creator    Category = "creator"
	Other      Category = "other"
)

// Subscription represents a single subscription
type Subscription struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Cost         float64      `json:"cost"`
	BillingCycle BillingCycle `json:"billing_cycle"`
	NextPayment  time.Time    `json:"next_payment"`
	StartDate    time.Time    `json:"start_date"`
	Category     Category     `json:"category"`
	Notes        string       `json:"notes"`
	Image        string       `json:"image"` // Filename only (e.g., "abc-123.png"), stored in images/ folder
	Paused       bool         `json:"paused"`
	Deleted      bool         `json:"deleted"`      // Soft delete flag
	DeletedAt    time.Time    `json:"deleted_at"`   // When subscription was deleted
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// Payment represents a single payment made for a subscription
type Payment struct {
	ID             string    `json:"id"`
	SubscriptionID string    `json:"subscription_id"`
	Amount         float64   `json:"amount"`
	PaymentDate    time.Time `json:"payment_date"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

// SubscriptionList is a collection of subscriptions and payments
type SubscriptionList struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Payments      []Payment      `json:"payments"`
	Version       string         `json:"version"`
}

// FilterCriteria defines search/filter parameters
type FilterCriteria struct {
	SearchTerm   string
	Category     *Category
	BillingCycle *BillingCycle
	MinCost      *float64
	MaxCost      *float64
	ShowPaused   bool // If true, show paused subscriptions; if false, hide them
}

// SortField defines sortable fields
type SortField string

const (
	SortByName        SortField = "name"
	SortByCost        SortField = "cost"
	SortByNextPayment SortField = "next_payment"
)

// SortOrder defines sort direction
type SortOrder string

const (
	Ascending  SortOrder = "asc"
	Descending SortOrder = "desc"
)

// CostSummary represents aggregated cost statistics
type CostSummary struct {
	TotalMonthly float64
	TotalYearly  float64
	YearToDate   float64 // Actual payments made from Jan 1 to today
	ByCategory   map[Category]float64
	Count        int // Total count of active (non-paused, non-deleted) subscriptions
	PausedCount  int // Count of paused subscriptions
}
