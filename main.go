package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/farhanasfar/airbnb-market-scraping-system/config"
	"github.com/farhanasfar/airbnb-market-scraping-system/models"
	"github.com/farhanasfar/airbnb-market-scraping-system/scraper/airbnb"
	"github.com/farhanasfar/airbnb-market-scraping-system/services"
	"github.com/farhanasfar/airbnb-market-scraping-system/storage"
	"github.com/farhanasfar/airbnb-market-scraping-system/utils"
)

func main() {
	logger := utils.NewLogger()
	logger.Info("Starting Airbnb Multi-Location Scraper...")

	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	logger.Success("Configuration loaded")

	// Connect to database
	db, err := storage.NewDB(cfg.Database.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create services
	listingService := services.NewListingService(db, logger)
	scraper := airbnb.NewScraper(&cfg.Scraper, logger)
	ctx := context.Background()

	// Step 1: Scrape homepage to get location URLs
	logger.Info("\n=== STEP 1: EXTRACTING LOCATIONS FROM HOMEPAGE ===")
	locations, err := scraper.ScrapeHomepageLocations(ctx)
	if err != nil {
		log.Fatal("Failed to scrape homepage:", err)
	}

	if len(locations) == 0 {
		logger.Warning("No locations found on homepage")
		return
	}

	logger.Info("Found %d locations:", len(locations))
	for i, loc := range locations {
		logger.Info("  %d. %s", i+1, loc.Name)
	}

	// Step 2: Scrape properties from each location
	logger.Info("\n=== STEP 2: SCRAPING PROPERTIES FROM EACH LOCATION ===")

	allRawListings := []models.RawListing{}
	totalProperties := 0

	for i, location := range locations {
		logger.Info("\n[%d/%d] Scraping: %s", i+1, len(locations), location.Name)
		logger.Info("URL: %s", location.URL)

		// Scrape this location (2 pages × 5 properties = 10 per location)
		rawListings, err := scraper.ScrapeListings(ctx, location.URL)
		if err != nil {
			logger.Error("Failed to scrape %s: %v", location.Name, err)
			continue
		}

		if len(rawListings) == 0 {
			logger.Warning("No listings found for %s", location.Name)
			continue
		}

		logger.Success("Got %d properties from %s", len(rawListings), location.Name)
		allRawListings = append(allRawListings, rawListings...)
		totalProperties += len(rawListings)
	}

	logger.Success("\n=== SCRAPED %d TOTAL PROPERTIES FROM %d LOCATIONS ===",
		totalProperties, len(locations))

	if totalProperties == 0 {
		logger.Warning("No properties scraped, exiting")
		return
	}

	// Print summary by location if JSON console enabled
	if cfg.Output.JSONConsole {
		logger.Info("\n=== RAW LISTINGS SUMMARY ===")
		jsonData, _ := json.MarshalIndent(allRawListings, "", "  ")
		fmt.Println(string(jsonData))
	}

	// Step 3: Scrape detail pages for bedrooms/bathrooms/guests
	logger.Info("\n=== STEP 3: SCRAPING DETAIL PAGES ===")
	urls := make([]string, 0, len(allRawListings))
	for _, listing := range allRawListings {
		if listing.URL != "" {
			normalizedURL := utils.NormalizeURL(listing.URL)
			urls = append(urls, normalizedURL)
		}
	}

	logger.Info("Scraping details for %d properties...", len(urls))
	detailResults := scraper.ScrapeDetailsWithWorkers(ctx, urls)

	// Merge detail data
	for i := range allRawListings {
		normalizedURL := utils.NormalizeURL(allRawListings[i].URL)
		if detail, ok := detailResults[normalizedURL]; ok && detail.Error == nil {
			allRawListings[i].Bedrooms = detail.Bedrooms
			allRawListings[i].Bathrooms = detail.Bathrooms
			allRawListings[i].Guests = detail.Guests
		}
	}

	// Step 4: Save to database
	logger.Info("\n=== STEP 4: SAVING TO DATABASE ===")
	savedCount, err := listingService.NormalizeAndSave(allRawListings)
	if err != nil {
		logger.Error("Failed to save listings: %v", err)
	}

	// Final summary
	logger.Success("\n=== SCRAPING COMPLETE ===")
	logger.Info("Locations scraped: %d", len(locations))
	logger.Info("Total properties found: %d", totalProperties)
	logger.Info("Successfully saved: %d", savedCount)
	logger.Info("Duplicates/errors: %d", totalProperties-savedCount)
	logger.Info("\nAverage properties per location: %.1f", float64(totalProperties)/float64(len(locations)))
}
