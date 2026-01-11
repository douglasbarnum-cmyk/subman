package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"subman/internal/models"
	"subman/internal/ui/components"
)

type ListView struct {
	app           *App
	subscriptions []models.Subscription
	listContainer *fyne.Container
	sortBy        models.SortField
	sortOrder     models.SortOrder
	filter        *models.FilterCriteria
}

func NewListView(app *App) *ListView {
	return &ListView{
		app:       app,
		sortBy:    models.SortByName,
		sortOrder: models.Ascending,
	}
}

func (l *ListView) Render() fyne.CanvasObject {
	// Toolbar with Add, Sort, Export buttons
	addBtn := widget.NewButton("Add Subscription", l.showAddDialog)
	sortBtn := widget.NewButton("Sort", l.showSortMenu)
	exportBtn := widget.NewButton("Export", l.showExportDialog)

	toolbar := container.NewHBox(
		addBtn,
		sortBtn,
		exportBtn,
	)

	l.listContainer = container.NewVBox()
	l.Refresh()

	scrollContainer := container.NewScroll(l.listContainer)
	scrollContainer.SetMinSize(fyne.NewSize(0, 400))

	return container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		scrollContainer,
	)
}

func (l *ListView) Refresh() {
	subs, err := l.app.service.List(l.filter, l.sortBy, l.sortOrder)
	if err != nil {
		return
	}

	l.subscriptions = subs
	l.listContainer.Objects = nil

	if len(subs) == 0 {
		emptyLabel := widget.NewLabel("No subscriptions yet. Click 'Add Subscription' to get started.")
		l.listContainer.Add(emptyLabel)
	} else {
		// Get theme-aware colors for zebra striping
		baseColor := theme.BackgroundColor()
		altColor := l.getAlternateColor(baseColor)

		for i, sub := range subs {
			// Alternate background color
			bgColor := baseColor
			if i%2 == 1 {
				bgColor = altColor
			}
			card := components.NewSubscriptionCard(sub, l.onEdit, l.onDelete, bgColor)
			l.listContainer.Add(card)
		}
	}

	l.listContainer.Refresh()
}

func (l *ListView) showAddDialog() {
	form := NewSubscriptionForm(l.app, nil)
	form.Show()
}

func (l *ListView) onEdit(sub models.Subscription) {
	form := NewSubscriptionForm(l.app, &sub)
	form.Show()
}

func (l *ListView) onDelete(sub models.Subscription) {
	confirm := dialog.NewConfirm(
		"Delete Subscription",
		fmt.Sprintf("Are you sure you want to delete %s?", sub.Name),
		func(confirmed bool) {
			if confirmed {
				l.app.service.Delete(sub.ID)
				l.app.Refresh()
			}
		},
		l.app.window,
	)
	confirm.Show()
}

func (l *ListView) showSortMenu() {
	sortByName := widget.NewButton("Sort by Name", func() {
		l.sortBy = models.SortByName
		l.Refresh()
	})

	sortByCost := widget.NewButton("Sort by Cost", func() {
		l.sortBy = models.SortByCost
		l.Refresh()
	})

	sortByNextPayment := widget.NewButton("Sort by Next Payment", func() {
		l.sortBy = models.SortByNextPayment
		l.Refresh()
	})

	toggleOrder := widget.NewButton("Toggle Order", func() {
		if l.sortOrder == models.Ascending {
			l.sortOrder = models.Descending
		} else {
			l.sortOrder = models.Ascending
		}
		l.Refresh()
	})

	content := container.NewVBox(
		sortByName,
		sortByCost,
		sortByNextPayment,
		widget.NewSeparator(),
		toggleOrder,
	)

	d := dialog.NewCustom("Sort Options", "Close", content, l.app.window)
	d.Show()
}

func (l *ListView) showExportDialog() {
	exportView := NewExportView(l.app, l.subscriptions)
	exportView.Show()
}

func (l *ListView) SetFilter(filter *models.FilterCriteria) {
	l.filter = filter
	l.Refresh()
}

// getAlternateColor creates a slightly different shade for alternating rows
// For light themes: makes it slightly darker
// For dark themes: makes it slightly lighter
func (l *ListView) getAlternateColor(base color.Color) color.Color {
	r, g, b, a := base.RGBA()

	// Convert from 16-bit to 8-bit color values
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)

	// Determine if we're in a dark theme (low RGB values) or light theme (high RGB values)
	brightness := (int(r8) + int(g8) + int(b8)) / 3

	var newR, newG, newB uint8

	if brightness > 128 {
		// Light theme - make slightly darker
		newR = darken(r8, 10)
		newG = darken(g8, 10)
		newB = darken(b8, 10)
	} else {
		// Dark theme - make slightly lighter
		newR = lighten(r8, 15)
		newG = lighten(g8, 15)
		newB = lighten(b8, 15)
	}

	return color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)}
}

func darken(c uint8, amount uint8) uint8 {
	if c < amount {
		return 0
	}
	return c - amount
}

func lighten(c uint8, amount uint8) uint8 {
	result := int(c) + int(amount)
	if result > 255 {
		return 255
	}
	return uint8(result)
}
