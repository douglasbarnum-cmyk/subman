package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"subman/internal/ui/components"
)

type DashboardView struct {
	app        *App
	monthlyLabel *widget.Label
	yearlyLabel  *widget.Label
	ytdLabel     *widget.Label
	countLabel   *widget.Label
}

func NewDashboardView(app *App) *DashboardView {
	return &DashboardView{
		app:          app,
		monthlyLabel: widget.NewLabel("$0.00"),
		yearlyLabel:  widget.NewLabel("$0.00"),
		ytdLabel:     widget.NewLabel("$0.00"),
		countLabel:   widget.NewLabel("0"),
	}
}

func (d *DashboardView) Render() fyne.CanvasObject {
	d.Refresh()

	monthlyCard := components.NewStatsCard("Monthly Total", d.monthlyLabel)
	yearlyCard := components.NewStatsCard("Yearly Total", d.yearlyLabel)
	ytdCard := components.NewStatsCard("Year to Date", d.ytdLabel)
	countCard := components.NewStatsCard("Active Subscriptions", d.countLabel)

	return container.NewHBox(
		monthlyCard,
		yearlyCard,
		ytdCard,
		countCard,
	)
}

func (d *DashboardView) Refresh() {
	summary, err := d.app.service.GetSummary()
	if err != nil {
		d.monthlyLabel.SetText("Error")
		d.yearlyLabel.SetText("Error")
		d.ytdLabel.SetText("Error")
		d.countLabel.SetText("Error")
		return
	}

	d.monthlyLabel.SetText(fmt.Sprintf("$%.2f", summary.TotalMonthly))
	d.yearlyLabel.SetText(fmt.Sprintf("$%.2f", summary.TotalYearly))
	d.ytdLabel.SetText(fmt.Sprintf("$%.2f", summary.YearToDate))
	d.countLabel.SetText(fmt.Sprintf("%d", summary.Count))
}
