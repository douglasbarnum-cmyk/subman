package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"subman/internal/models"
)

type FilterView struct {
	app            *App
	searchEntry    *widget.Entry
	categorySelect *widget.Select
	cycleSelect    *widget.Select
	showPausedCheck *widget.Check
}

func NewFilterView(app *App) *FilterView {
	return &FilterView{
		app: app,
	}
}

func (f *FilterView) Render() fyne.CanvasObject {
	f.searchEntry = widget.NewEntry()
	f.searchEntry.SetPlaceHolder("Search subscriptions...")
	f.searchEntry.OnChanged = func(s string) {
		f.applyFilters()
	}

	categories := []string{"All", "Streaming", "Software", "Utilities", "Gaming", "News", "Education", "Creator", "Other"}
	f.categorySelect = widget.NewSelect(categories, func(s string) {
		f.applyFilters()
	})
	f.categorySelect.Selected = "All"

	cycles := []string{"All", "Monthly", "Yearly"}
	f.cycleSelect = widget.NewSelect(cycles, func(s string) {
		f.applyFilters()
	})
	f.cycleSelect.Selected = "All"

	f.showPausedCheck = widget.NewCheck("Show Paused Subscriptions", func(checked bool) {
		f.applyFilters()
	})
	f.showPausedCheck.Checked = false // By default, hide paused subscriptions

	clearBtn := widget.NewButton("Clear Filters", f.clearFilters)

	return container.NewVBox(
		widget.NewLabel("Filter Subscriptions"),
		f.searchEntry,
		container.NewGridWithColumns(2,
			widget.NewLabel("Category:"),
			f.categorySelect,
			widget.NewLabel("Billing Cycle:"),
			f.cycleSelect,
		),
		f.showPausedCheck,
		clearBtn,
	)
}

func (f *FilterView) applyFilters() {
	// Build filter criteria
	criteria := &models.FilterCriteria{
		SearchTerm: f.searchEntry.Text,
		ShowPaused: f.showPausedCheck.Checked,
	}

	// Apply category filter
	if f.categorySelect.Selected != "All" {
		cat := models.Category(strings.ToLower(f.categorySelect.Selected))
		criteria.Category = &cat
	}

	// Apply cycle filter
	if f.cycleSelect.Selected != "All" {
		cycle := models.BillingCycle(strings.ToLower(f.cycleSelect.Selected))
		criteria.BillingCycle = &cycle
	}

	// Update list view with filters
	f.app.listView.SetFilter(criteria)
}

func (f *FilterView) clearFilters() {
	f.searchEntry.SetText("")
	f.categorySelect.Selected = "All"
	f.cycleSelect.Selected = "All"
	f.showPausedCheck.Checked = false
	f.applyFilters()
}
