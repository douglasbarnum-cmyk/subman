package components

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"subman/internal/images"
	"subman/internal/models"
)

func NewSubscriptionCard(sub models.Subscription, onEdit func(models.Subscription), onDelete func(models.Subscription), bgColor color.Color) fyne.CanvasObject {
	// Load image
	var imageWidget *canvas.Image
	imagePath, err := images.GetImagePath(sub.Image)
	if err == nil && imagePath != "" {
		imageWidget = canvas.NewImageFromFile(imagePath)
		imageWidget.FillMode = canvas.ImageFillContain
		imageWidget.SetMinSize(fyne.NewSize(64, 64))
	} else {
		// Fallback to default if image not found
		defaultImage := images.GetDefaultImageForCategory(sub.Category)
		imagePath, _ = images.GetImagePath(defaultImage)
		imageWidget = canvas.NewImageFromFile(imagePath)
		imageWidget.FillMode = canvas.ImageFillContain
		imageWidget.SetMinSize(fyne.NewSize(64, 64))
	}

	nameLabel := widget.NewLabel(sub.Name)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Show paused status
	var costText string
	if sub.Paused {
		costText = fmt.Sprintf("$%.2f / %s (PAUSED)", sub.Cost, sub.BillingCycle)
	} else {
		costText = fmt.Sprintf("$%.2f / %s", sub.Cost, sub.BillingCycle)
	}
	costLabel := widget.NewLabel(costText)

	categoryLabel := widget.NewLabel(string(sub.Category))

	var nextPaymentText string
	if sub.Paused {
		nextPaymentText = "Subscription is paused"
	} else {
		nextPaymentText = fmt.Sprintf("Next: %s", sub.NextPayment.Format("Jan 2, 2006"))
	}
	nextPaymentLabel := widget.NewLabel(nextPaymentText)

	editBtn := widget.NewButton("Edit", func() {
		onEdit(sub)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		onDelete(sub)
	})

	info := container.NewVBox(
		nameLabel,
		costLabel,
		categoryLabel,
		nextPaymentLabel,
	)

	actions := container.NewHBox(editBtn, deleteBtn)

	// Create card with image on the left
	card := container.NewBorder(
		nil,
		nil,
		container.NewHBox(imageWidget, info),
		actions,
	)

	// Create background rectangle
	background := canvas.NewRectangle(bgColor)

	// Layer content on top of background
	cardWithBg := container.NewMax(background, container.NewPadded(card))

	return cardWithBg
}
