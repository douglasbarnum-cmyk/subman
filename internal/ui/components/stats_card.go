package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewStatsCard(title string, valueLabel *widget.Label) fyne.CanvasObject {
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	valueLabel.TextStyle = fyne.TextStyle{Bold: true}
	valueLabel.Alignment = fyne.TextAlignCenter

	card := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		valueLabel,
	)

	return container.NewPadded(card)
}
