package ui

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"subman/internal/images"
	"subman/internal/models"
)

type SubscriptionForm struct {
	app          *App
	subscription *models.Subscription
	dialog       dialog.Dialog

	nameEntry        *widget.Entry
	costEntry        *widget.Entry
	cycleSelect      *widget.Select
	categorySelect   *widget.Select
	nextPaymentEntry *widget.Entry
	startDateEntry   *widget.Entry
	notesEntry       *widget.Entry
	pausedCheck      *widget.Check
	imageLabel       *widget.Label
	selectedImage    string // Path to selected image file
}

func NewSubscriptionForm(app *App, sub *models.Subscription) *SubscriptionForm {
	form := &SubscriptionForm{
		app:          app,
		subscription: sub,
	}

	form.buildForm()
	return form
}

func (f *SubscriptionForm) buildForm() {
	f.nameEntry = widget.NewEntry()
	f.nameEntry.SetPlaceHolder("Subscription name")

	f.costEntry = widget.NewEntry()
	f.costEntry.SetPlaceHolder("0.00")

	f.cycleSelect = widget.NewSelect([]string{"monthly", "yearly"}, nil)
	f.cycleSelect.Selected = "monthly"

	categories := []string{"Streaming", "Software", "Utilities", "Gaming", "News", "Education", "Creator", "Other"}
	f.categorySelect = widget.NewSelect(categories, nil)
	f.categorySelect.Selected = "Other"

	f.nextPaymentEntry = widget.NewEntry()
	f.nextPaymentEntry.SetPlaceHolder("YYYY-MM-DD")

	f.startDateEntry = widget.NewEntry()
	f.startDateEntry.SetPlaceHolder("YYYY-MM-DD")

	f.notesEntry = widget.NewMultiLineEntry()
	f.notesEntry.SetPlaceHolder("Additional notes...")
	f.notesEntry.Wrapping = fyne.TextWrapWord // Wrap at word boundaries, no horizontal scroll

	f.pausedCheck = widget.NewCheck("Subscription is paused", nil)

	// Image picker
	f.imageLabel = widget.NewLabel("No image selected (will use default)")
	imageButton := widget.NewButton("Choose Image", f.chooseImage)
	clearImageButton := widget.NewButton("Clear", f.clearImage)
	imageSelector := container.NewBorder(nil, nil, imageButton, clearImageButton, f.imageLabel)

	// Populate if editing
	if f.subscription != nil {
		f.nameEntry.SetText(f.subscription.Name)
		f.costEntry.SetText(strconv.FormatFloat(f.subscription.Cost, 'f', 2, 64))
		f.cycleSelect.Selected = string(f.subscription.BillingCycle)
		// Capitalize the category for display
		catStr := string(f.subscription.Category)
		f.categorySelect.Selected = strings.Title(catStr)
		f.nextPaymentEntry.SetText(f.subscription.NextPayment.Format("2006-01-02"))
		f.startDateEntry.SetText(f.subscription.StartDate.Format("2006-01-02"))
		f.notesEntry.SetText(f.subscription.Notes)
		f.pausedCheck.Checked = f.subscription.Paused

		if f.subscription.Image != "" {
			f.selectedImage = f.subscription.Image
			f.imageLabel.SetText("Current: " + f.subscription.Image)
		}
	}

	formItems := []*widget.FormItem{
		widget.NewFormItem("Name", f.nameEntry),
		widget.NewFormItem("Cost", f.costEntry),
		widget.NewFormItem("Billing Cycle", f.cycleSelect),
		widget.NewFormItem("Category", f.categorySelect),
		widget.NewFormItem("Image", imageSelector),
		widget.NewFormItem("Next Payment", f.nextPaymentEntry),
		widget.NewFormItem("Start Date", f.startDateEntry),
		widget.NewFormItem("Notes", f.notesEntry),
		widget.NewFormItem("Status", f.pausedCheck),
	}

	formWidget := widget.NewForm(formItems...)
	formWidget.OnSubmit = f.onSubmit
	formWidget.OnCancel = f.onCancel

	title := "Add Subscription"
	if f.subscription != nil {
		title = "Edit Subscription"
	}

	f.dialog = dialog.NewCustom(title, "Close", formWidget, f.app.window)
	f.dialog.Resize(fyne.NewSize(500, 600))
}

func (f *SubscriptionForm) Show() {
	f.dialog.Show()
}

func (f *SubscriptionForm) onSubmit() {
	// Validate and parse
	cost, err := strconv.ParseFloat(f.costEntry.Text, 64)
	if err != nil {
		dialog.ShowError(err, f.app.window)
		return
	}

	nextPayment, err := time.Parse("2006-01-02", f.nextPaymentEntry.Text)
	if err != nil {
		dialog.ShowError(err, f.app.window)
		return
	}

	startDate, err := time.Parse("2006-01-02", f.startDateEntry.Text)
	if err != nil {
		dialog.ShowError(err, f.app.window)
		return
	}

	// Handle image - use selected image, or default for category if none selected
	imageFilename := f.selectedImage
	if imageFilename == "" {
		category := models.Category(strings.ToLower(f.categorySelect.Selected))
		imageFilename = images.GetDefaultImageForCategory(category)
	}

	sub := &models.Subscription{
		Name:         f.nameEntry.Text,
		Cost:         cost,
		BillingCycle: models.BillingCycle(f.cycleSelect.Selected),
		Category:     models.Category(strings.ToLower(f.categorySelect.Selected)),
		NextPayment:  nextPayment,
		StartDate:    startDate,
		Notes:        f.notesEntry.Text,
		Image:        imageFilename,
		Paused:       f.pausedCheck.Checked,
	}

	if f.subscription != nil {
		// Update existing
		sub.ID = f.subscription.ID
		f.app.service.Update(sub)
	} else {
		// Create new
		f.app.service.Create(sub)
	}

	f.app.Refresh()
	f.dialog.Hide()
}

func (f *SubscriptionForm) onCancel() {
	f.dialog.Hide()
}

func (f *SubscriptionForm) chooseImage() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		sourcePath := reader.URI().Path()

		// Generate unique filename using subscription ID (or temp name for new subscriptions)
		ext := filepath.Ext(sourcePath)
		var filename string
		if f.subscription != nil {
			filename = f.subscription.ID + ext
		} else {
			filename = "temp_" + filepath.Base(sourcePath)
		}

		// Copy image to images directory
		err = images.SaveImage(sourcePath, filename)
		if err != nil {
			dialog.ShowError(err, f.app.window)
			return
		}

		f.selectedImage = filename
		f.imageLabel.SetText("Selected: " + filepath.Base(sourcePath))
	}, f.app.window)

	// Filter for image files
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}))
	fd.Show()
}

func (f *SubscriptionForm) clearImage() {
	f.selectedImage = ""
	f.imageLabel.SetText("No image selected (will use default)")
}
