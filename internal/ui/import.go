package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"subman/internal/images"
	"subman/pkg/importer"
)

type ImportView struct {
	app *App
}

func NewImportView(app *App) *ImportView {
	return &ImportView{
		app: app,
	}
}

func (i *ImportView) Show() {
	// Create file open dialog for ZIP files
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		i.importBundle(filePath)
	}, i.app.window)

	// Filter for ZIP files
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".zip"}))
	fileDialog.Show()
}

func (i *ImportView) importBundle(zipPath string) {
	// Get images directory
	imagesDir, err := images.GetImagesDir()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to get images directory: %w", err), i.app.window)
		return
	}

	// Create importer
	bundleImporter := importer.NewBundleImporter(imagesDir)

	// Validate the bundle first
	if err := bundleImporter.ValidateBundle(zipPath); err != nil {
		dialog.ShowError(fmt.Errorf("invalid bundle: %w", err), i.app.window)
		return
	}

	// Ask user how they want to import (replace or merge)
	i.showImportModeDialog(zipPath, bundleImporter)
}

func (i *ImportView) showImportModeDialog(zipPath string, bundleImporter *importer.BundleImporter) {
	modeSelect := widget.NewRadioGroup([]string{"Replace all data", "Merge with existing data"}, nil)
	modeSelect.Selected = "Merge with existing data"

	content := widget.NewForm(
		widget.NewFormItem("Import Mode", modeSelect),
	)

	infoLabel := widget.NewLabel("Replace: Deletes all existing subscriptions and replaces with imported data.\nMerge: Adds imported subscriptions to existing data (duplicates may occur).")
	infoLabel.Wrapping = fyne.TextWrapWord

	fullContent := container.NewVBox(
		infoLabel,
		content,
	)

	confirm := dialog.NewCustomConfirm("Import Bundle", "Import", "Cancel", fullContent, func(ok bool) {
		if ok {
			var mode importer.ImportMode
			if modeSelect.Selected == "Replace all data" {
				mode = importer.ImportModeReplace
			} else {
				mode = importer.ImportModeMerge
			}
			i.doImport(zipPath, bundleImporter, mode)
		}
	}, i.app.window)

	confirm.Show()
}

func (i *ImportView) doImport(zipPath string, bundleImporter *importer.BundleImporter, mode importer.ImportMode) {
	// Import the bundle
	importedList, err := bundleImporter.ImportBundle(zipPath, mode)
	if err != nil {
		dialog.ShowError(fmt.Errorf("import failed: %w", err), i.app.window)
		return
	}

	// Get current storage
	storage := i.app.service.GetStorage()

	if mode == importer.ImportModeReplace {
		// Replace mode: save the imported list directly
		if err := storage.Save(importedList); err != nil {
			dialog.ShowError(fmt.Errorf("failed to save imported data: %w", err), i.app.window)
			return
		}
	} else {
		// Merge mode: load current data and merge
		currentList, err := storage.Load()
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to load current data: %w", err), i.app.window)
			return
		}

		// Merge subscriptions
		currentList.Subscriptions = append(currentList.Subscriptions, importedList.Subscriptions...)

		// Merge payments
		currentList.Payments = append(currentList.Payments, importedList.Payments...)

		if err := storage.Save(currentList); err != nil {
			dialog.ShowError(fmt.Errorf("failed to save merged data: %w", err), i.app.window)
			return
		}
	}

	// Regenerate payments after import
	if err := i.app.paymentService.GenerateAllPayments(); err != nil {
		dialog.ShowError(fmt.Errorf("failed to regenerate payment history: %w", err), i.app.window)
		return
	}

	// Refresh the UI
	i.app.Refresh()

	// Show success message
	successMsg := fmt.Sprintf("Successfully imported %d subscriptions and %d payments",
		len(importedList.Subscriptions), len(importedList.Payments))
	dialog.ShowInformation("Import Complete", successMsg, i.app.window)
}
