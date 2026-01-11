package calculator

import (
	"time"

	"subman/internal/models"
)

// CalculateSummary computes cost statistics from subscriptions and payments
// Paused and deleted subscriptions are counted separately and excluded from cost totals
// YTD is calculated from actual payment records
func CalculateSummary(subscriptions []models.Subscription, payments []models.Payment) *models.CostSummary {
	summary := &models.CostSummary{
		ByCategory:  make(map[models.Category]float64),
		Count:       0,
		PausedCount: 0,
		YearToDate:  0,
	}

	for _, sub := range subscriptions {
		// Skip deleted subscriptions
		if sub.Deleted {
			continue
		}

		if sub.Paused {
			summary.PausedCount++
			continue
		}

		summary.Count++
		monthlyCost := ToMonthlyCost(sub.Cost, sub.BillingCycle)
		summary.TotalMonthly += monthlyCost
		summary.ByCategory[sub.Category] += monthlyCost
	}

	summary.TotalYearly = summary.TotalMonthly * 12

	// Calculate YTD from actual payments (Jan 1 to today)
	now := time.Now()
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	for _, payment := range payments {
		if payment.PaymentDate.After(startOfYear) || payment.PaymentDate.Equal(startOfYear) {
			if payment.PaymentDate.Before(now) || payment.PaymentDate.Equal(now.Truncate(24*time.Hour)) {
				summary.YearToDate += payment.Amount
			}
		}
	}

	return summary
}

// ToMonthlyCost converts any cost to monthly equivalent
func ToMonthlyCost(cost float64, cycle models.BillingCycle) float64 {
	if cycle == models.Yearly {
		return cost / 12
	}
	return cost
}

// ToYearlyCost converts any cost to yearly equivalent
func ToYearlyCost(cost float64, cycle models.BillingCycle) float64 {
	if cycle == models.Monthly {
		return cost * 12
	}
	return cost
}

// CalculateNextPayment calculates the next payment date based on current date and cycle
func CalculateNextPayment(lastPayment time.Time, cycle models.BillingCycle) time.Time {
	now := time.Now()
	next := lastPayment

	for next.Before(now) {
		if cycle == models.Monthly {
			next = next.AddDate(0, 1, 0)
		} else {
			next = next.AddDate(1, 0, 0)
		}
	}

	return next
}
