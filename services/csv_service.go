package services

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/farhanasfar/airbnb-market-scraping-system/models"
	"github.com/farhanasfar/airbnb-market-scraping-system/storage"
	"github.com/farhanasfar/airbnb-market-scraping-system/utils"
)

// CSVService handles CSV export operations
type CSVService struct {
	db     *storage.DB
	logger *utils.Logger
}

// NewCSVService creates a new CSV service
func NewCSVService(db *storage.DB, logger *utils.Logger) *CSVService {
	return &CSVService{
		db:     db,
		logger: logger,
	}
}

// ExportToCSV exports all listings from database to CSV file
func (s *CSVService) ExportToCSV(filename string) error {
	s.logger.Info("Exporting listings to CSV: %s", filename)

	// Get all listings from database
	listings, err := s.db.GetAllListings()
	if err != nil {
		return fmt.Errorf("failed to get listings: %w", err)
	}

	if len(listings) == 0 {
		s.logger.Warning("No listings to export")
		return nil
	}

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID",
		"Title",
		"Price",
		"Location",
		"Rating",
		"Bedrooms",
		"Bathrooms",
		"Guests",
		"URL",
		"Created At",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, listing := range listings {
		row := []string{
			fmt.Sprintf("%d", listing.ID),
			listing.Title,
			fmt.Sprintf("%.2f", listing.Price),
			listing.Location,
			fmt.Sprintf("%.2f", listing.Rating),
			fmt.Sprintf("%d", listing.Bedrooms),
			fmt.Sprintf("%d", listing.Bathrooms),
			fmt.Sprintf("%d", listing.Guests),
			listing.URL,
			listing.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	s.logger.Success("Exported %d listings to %s", len(listings), filename)
	return nil
}

// ExportListingsToCSV exports a specific list of listings to CSV
func (s *CSVService) ExportListingsToCSV(listings []models.Listing, filename string) error {
	s.logger.Info("Exporting %d listings to CSV: %s", len(listings), filename)

	if len(listings) == 0 {
		s.logger.Warning("No listings to export")
		return nil
	}

	// Create CSV file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Title",
		"Price",
		"Location",
		"Rating",
		"Bedrooms",
		"Bathrooms",
		"Guests",
		"URL",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, listing := range listings {
		row := []string{
			listing.Title,
			fmt.Sprintf("%.2f", listing.Price),
			listing.Location,
			fmt.Sprintf("%.2f", listing.Rating),
			fmt.Sprintf("%d", listing.Bedrooms),
			fmt.Sprintf("%d", listing.Bathrooms),
			fmt.Sprintf("%d", listing.Guests),
			listing.URL,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	s.logger.Success("Exported %d listings to %s", len(listings), filename)
	return nil
}
