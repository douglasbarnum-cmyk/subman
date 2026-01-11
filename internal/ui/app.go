package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"subman/internal/images"
	"subman/internal/service"
)

type App struct {
	fyneApp        fyne.App
	window         fyne.Window
	service        *service.SubscriptionService
	paymentService *service.PaymentService

	// Views
	dashboard  *DashboardView
	listView   *ListView
	filterView *FilterView
}

func NewApp(service *service.SubscriptionService, paymentService *service.PaymentService) *App {
	fyneApp := app.NewWithID("com.subman.app")
	window := fyneApp.NewWindow("Subman - Subscription Manager")

	a := &App{
		fyneApp:        fyneApp,
		window:         window,
		service:        service,
		paymentService: paymentService,
	}

	a.dashboard = NewDashboardView(a)
	a.listView = NewListView(a)
	a.filterView = NewFilterView(a)

	// Initialize default category icons
	if err := images.EnsureDefaultCategoryIcons(); err != nil {
		log.Printf("Warning: Failed to create default category icons: %v", err)
	}

	// Generate payments for all active subscriptions
	if err := paymentService.GenerateAllPayments(); err != nil {
		log.Printf("Warning: Failed to generate payment history: %v", err)
	}

	// Load saved theme preference
	a.loadThemePreference()

	// Setup menu
	a.setupMenu()

	return a
}

func (a *App) Run() {
	// Setup main layout
	content := container.NewBorder(
		a.dashboard.Render(), // Top - dashboard with stats
		nil,                   // Bottom
		nil,                   // Left
		nil,                   // Right
		container.NewVSplit(
			a.filterView.Render(), // Top section - filters
			a.listView.Render(),   // Bottom section - list
		),
	)

	a.window.SetContent(content)
	a.window.Resize(fyne.NewSize(1000, 700))
	a.window.ShowAndRun()
}

func (a *App) Refresh() {
	a.dashboard.Refresh()
	a.listView.Refresh()
}

func (a *App) setupMenu() {
	// Create Settings menu
	settingsItem := fyne.NewMenuItem("Settings", func() {
		settings := NewSettingsView(a)
		settings.Show()
	})

	settingsMenu := fyne.NewMenu("Settings", settingsItem)

	// Set the main menu
	mainMenu := fyne.NewMainMenu(settingsMenu)
	a.window.SetMainMenu(mainMenu)
}

func (a *App) loadThemePreference() {
	// Load saved theme preference
	themeName := a.fyneApp.Preferences().StringWithFallback("theme", "System Default")

	switch themeName {
	case "Dark":
		a.fyneApp.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		a.fyneApp.Settings().SetTheme(theme.LightTheme())
	default:
		// System Default - don't set anything, use default
	}
}
