package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"subman/internal/models"
	"subman/pkg/export"
)

type ExportView struct {
	app           *App
	subscriptions []models.Subscription
}

func NewExportView(app *App, subs []models.Subscription) *ExportView {
	return &ExportView{
		app:           app,
		subscriptions: subs,
	}
}

func (e *ExportView) Show() {
	formatSelect := widget.NewRadioGroup([]string{"CSV", "JSON"}, nil)
	formatSelect.Selected = "CSV"

	content := widget.NewForm(
		widget.NewFormItem("Format", formatSelect),
	)

	confirm := dialog.NewCustomConfirm("Export Subscriptions", "Export", "Cancel", content, func(ok bool) {
		if ok {
			e.doExport(formatSelect.Selected)
		}
	}, e.app.window)

	confirm.Show()
}

func (e *ExportView) doExport(format string) {
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		var exporter export.Exporter
		if format == "CSV" {
			exporter = export.NewCSVExporter()
		} else {
			exporter = export.NewJSONExporter()
		}

		if err := exporter.Export(e.subscriptions, writer); err != nil {
			dialog.ShowError(err, e.app.window)
		}
	}, e.app.window)

	if format == "CSV" {
		saveDialog.SetFileName("subscriptions.csv")
	} else {
		saveDialog.SetFileName("subscriptions.json")
	}

	saveDialog.Show()
}
