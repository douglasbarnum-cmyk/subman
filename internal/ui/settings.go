package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SettingsView struct {
	app *App
}

func NewSettingsView(app *App) *SettingsView {
	return &SettingsView{
		app: app,
	}
}

func (s *SettingsView) Show() {
	// Get current theme preference from saved preferences
	selectedTheme := s.app.fyneApp.Preferences().StringWithFallback("theme", "System Default")

	// Create theme selection radio group
	themeRadio := widget.NewRadioGroup([]string{"Light", "Dark", "System Default"}, func(value string) {
		s.applyTheme(value)
	})
	themeRadio.Selected = selectedTheme

	content := container.NewVBox(
		widget.NewLabel("Theme:"),
		themeRadio,
	)

	d := dialog.NewCustom("Settings", "Close", content, s.app.window)
	d.Resize(fyne.NewSize(300, 200))
	d.Show()
}

func (s *SettingsView) applyTheme(themeName string) {
	// Apply the theme
	switch themeName {
	case "Dark":
		s.app.fyneApp.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		s.app.fyneApp.Settings().SetTheme(theme.LightTheme())
	default:
		// System Default - use default theme
		s.app.fyneApp.Settings().SetTheme(nil)
	}

	// Save the preference
	s.app.fyneApp.Preferences().SetString("theme", themeName)

	// Refresh the UI to apply the new theme
	s.app.Refresh()
}
