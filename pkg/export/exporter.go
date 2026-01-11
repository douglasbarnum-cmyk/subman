package export

import (
	"io"
	"subman/internal/models"
)

type Exporter interface {
	Export(subscriptions []models.Subscription, writer io.Writer) error
}
