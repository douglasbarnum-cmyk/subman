package images

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"subman/internal/models"
)

// GetImagesDir returns the path to the images directory
func GetImagesDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	imagesDir := filepath.Join(configDir, "subman", "images")
	if err := os.MkdirAll(imagesDir, 0700); err != nil {
		return "", err
	}

	return imagesDir, nil
}

// SaveImage copies an image file to the images directory with the given filename
func SaveImage(sourcePath string, filename string) error {
	imagesDir, err := GetImagesDir()
	if err != nil {
		return err
	}

	// Open source file
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	destPath := filepath.Join(imagesDir, filename)
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy file
	_, err = io.Copy(dst, src)
	return err
}

// GetImagePath returns the full path to an image file
func GetImagePath(filename string) (string, error) {
	if filename == "" {
		return "", nil
	}

	imagesDir, err := GetImagesDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(imagesDir, filename), nil
}

// DeleteImage removes an image file from the images directory
func DeleteImage(filename string) error {
	if filename == "" {
		return nil
	}

	imagesDir, err := GetImagesDir()
	if err != nil {
		return err
	}

	imagePath := filepath.Join(imagesDir, filename)
	return os.Remove(imagePath)
}

// GenerateDefaultCategoryIcon creates a simple colored square icon for a category
func GenerateDefaultCategoryIcon(category models.Category) (image.Image, error) {
	// 128x128 square
	width, height := 128, 128
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Category colors
	categoryColors := map[models.Category]color.RGBA{
		models.Streaming:  {R: 229, G: 57, B: 53, A: 255},   // Red
		models.Software:   {R: 52, G: 152, B: 219, A: 255},  // Blue
		models.Utilities:  {R: 46, G: 204, B: 113, A: 255},  // Green
		models.Gaming:     {R: 155, G: 89, B: 182, A: 255},  // Purple
		models.News:       {R: 241, G: 196, B: 15, A: 255},  // Yellow
		models.Education:  {R: 26, G: 188, B: 156, A: 255},  // Teal
		models.Creator:    {R: 230, G: 126, B: 34, A: 255},  // Orange
		models.Other:      {R: 149, G: 165, B: 166, A: 255}, // Gray
	}

	bgColor, exists := categoryColors[category]
	if !exists {
		bgColor = categoryColors[models.Other]
	}

	// Fill with solid color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	return img, nil
}

// EnsureDefaultCategoryIcons creates default icon files for all categories if they don't exist
func EnsureDefaultCategoryIcons() error {
	imagesDir, err := GetImagesDir()
	if err != nil {
		return err
	}

	categories := []models.Category{
		models.Streaming,
		models.Software,
		models.Utilities,
		models.Gaming,
		models.News,
		models.Education,
		models.Creator,
		models.Other,
	}

	for _, cat := range categories {
		filename := "default_" + string(cat) + ".png"
		imagePath := filepath.Join(imagesDir, filename)

		// Skip if already exists
		if _, err := os.Stat(imagePath); err == nil {
			continue
		}

		// Generate and save
		img, err := GenerateDefaultCategoryIcon(cat)
		if err != nil {
			return err
		}

		file, err := os.Create(imagePath)
		if err != nil {
			return err
		}

		err = png.Encode(file, img)
		file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDefaultImageForCategory returns the filename of the default image for a category
func GetDefaultImageForCategory(category models.Category) string {
	return "default_" + string(category) + ".png"
}
