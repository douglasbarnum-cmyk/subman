package importer

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"subman/internal/models"
)

type ImportMode string

const (
	ImportModeReplace ImportMode = "replace" // Replace all existing data
	ImportModeMerge   ImportMode = "merge"   // Merge with existing data
)

type BundleImporter struct {
	imagesDir string
}

func NewBundleImporter(imagesDir string) *BundleImporter {
	return &BundleImporter{
		imagesDir: imagesDir,
	}
}

// ImportBundle extracts and imports a bundle ZIP file
func (i *BundleImporter) ImportBundle(zipPath string, mode ImportMode) (*models.SubscriptionList, error) {
	// Open the ZIP file
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open bundle: %w", err)
	}
	defer zipReader.Close()

	var subscriptionList *models.SubscriptionList

	// Extract files from ZIP
	for _, file := range zipReader.File {
		if file.Name == "subscriptions.json" {
			// Read and parse subscriptions.json
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to read subscriptions.json: %w", err)
			}

			decoder := json.NewDecoder(rc)
			var list models.SubscriptionList
			if err := decoder.Decode(&list); err != nil {
				rc.Close()
				return nil, fmt.Errorf("failed to parse subscriptions.json: %w", err)
			}
			rc.Close()

			subscriptionList = &list

		} else if strings.HasPrefix(file.Name, "images/") {
			// Extract image files
			imageName := filepath.Base(file.Name)
			if imageName == "" || imageName == "." {
				continue
			}

			destPath := filepath.Join(i.imagesDir, imageName)

			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to read image %s: %w", imageName, err)
			}

			destFile, err := os.Create(destPath)
			if err != nil {
				rc.Close()
				return nil, fmt.Errorf("failed to create image %s: %w", imageName, err)
			}

			_, err = io.Copy(destFile, rc)
			destFile.Close()
			rc.Close()

			if err != nil {
				return nil, fmt.Errorf("failed to extract image %s: %w", imageName, err)
			}
		}
	}

	if subscriptionList == nil {
		return nil, fmt.Errorf("bundle does not contain subscriptions.json")
	}

	return subscriptionList, nil
}

// ValidateBundle checks if a ZIP file is a valid subscription bundle
func (i *BundleImporter) ValidateBundle(zipPath string) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("not a valid ZIP file: %w", err)
	}
	defer zipReader.Close()

	hasSubscriptions := false
	for _, file := range zipReader.File {
		if file.Name == "subscriptions.json" {
			hasSubscriptions = true
			break
		}
	}

	if !hasSubscriptions {
		return fmt.Errorf("bundle does not contain subscriptions.json")
	}

	return nil
}
