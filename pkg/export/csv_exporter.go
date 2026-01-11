package export

import (
	"encoding/csv"
	"io"
	"strconv"

	"subman/internal/models"
)

type CSVExporter struct{}

func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) Export(subscriptions []models.Subscription, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{"Name", "Cost", "Billing Cycle", "Next Payment", "Start Date", "Category", "Notes"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write data
	for _, sub := range subscriptions {
		record := []string{
			sub.Name,
			strconv.FormatFloat(sub.Cost, 'f', 2, 64),
			string(sub.BillingCycle),
			sub.NextPayment.Format("2006-01-02"),
			sub.StartDate.Format("2006-01-02"),
			string(sub.Category),
			sub.Notes,
		}
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	return nil
}
