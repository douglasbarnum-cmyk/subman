package export

import (
	"encoding/json"
	"io"

	"subman/internal/models"
)

type JSONExporter struct{}

func NewJSONExporter() *JSONExporter {
	return &JSONExporter{}
}

func (e *JSONExporter) Export(subscriptions []models.Subscription, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(subscriptions)
}
