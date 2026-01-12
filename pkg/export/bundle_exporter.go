package export

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"subman/internal/models"
)

type BundleExporter struct {
	imagesDir string
}

func NewBundleExporter(imagesDir string) *BundleExporter {
	return &BundleExporter{
		imagesDir: imagesDir,
	}
}

// ExportBundle creates a ZIP archive containing subscriptions.json and all user-supplied images
func (e *BundleExporter) ExportBundle(list *models.SubscriptionList, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// 1. Add subscriptions.json to the ZIP
	subscriptionsFile, err := zipWriter.Create("subscriptions.json")
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(subscriptionsFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(list); err != nil {
		return err
	}

	// 2. Collect all unique image filenames (excluding default category icons)
	imageFiles := make(map[string]bool)
	for _, sub := range list.Subscriptions {
		if sub.Image != "" && !strings.HasPrefix(sub.Image, "default_") {
			imageFiles[sub.Image] = true
		}
	}

	// 3. Add each image file to the ZIP under images/ folder
	for filename := range imageFiles {
		imagePath := filepath.Join(e.imagesDir, filename)

		// Check if file exists before trying to add it
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			// Skip missing images rather than failing the entire export
			continue
		}

		// Create entry in ZIP under images/ folder
		zipPath := filepath.Join("images", filename)
		imageEntry, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		// Copy image file into ZIP
		imageFile, err := os.Open(imagePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(imageEntry, imageFile)
		imageFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
