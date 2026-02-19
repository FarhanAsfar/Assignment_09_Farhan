package airbnb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type LocationCard struct {
	Name string
	URL  string
}

// ScrapeHomepageLocations extracts all visible location cards from Airbnb homepage
func (s *Scraper) ScrapeHomepageLocations(ctx context.Context) ([]LocationCard, error) {
	browserCtx, cancel := s.createStealthContext(ctx)
	defer cancel()

	s.logger.Info("Visiting Airbnb homepage to extract locations...")

	var locationsJSON string

	err := chromedp.Run(browserCtx,
		removeWebdriverProperty(),
		chromedp.Navigate(s.cfg.BaseURL),

		// Wait for page to load
		chromedp.Sleep(5*time.Second),

		// Scroll down to load more location cards
		chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight / 2)`, nil),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`window.scrollTo(0, 0)`, nil),
		chromedp.Sleep(2*time.Second),

		// Extract all location cards
		chromedp.Evaluate(`
			JSON.stringify(
				Array.from(document.querySelectorAll('a[href*="/s/"]')).map(link => ({
					name: link.innerText.trim() || link.getAttribute('aria-label') || 'Unknown',
					url: link.href
				})).filter(loc => 
					loc.url.includes('/s/') && 
					loc.url.includes('/homes') &&
					loc.name !== '' &&
					loc.name !== 'Unknown'
				).slice(0, 20)
			)
		`, &locationsJSON),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scrape homepage: %w", err)
	}

	// Parse JSON
	var rawLocations []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	if err := json.Unmarshal([]byte(locationsJSON), &rawLocations); err != nil {
		return nil, fmt.Errorf("failed to parse locations JSON: %w", err)
	}

	// Convert to LocationCard structs and deduplicate
	locations := []LocationCard{}
	seen := make(map[string]bool)

	for _, raw := range rawLocations {
		// Deduplicate by URL
		if seen[raw.URL] {
			continue
		}
		seen[raw.URL] = true

		locations = append(locations, LocationCard{
			Name: raw.Name,
			URL:  raw.URL,
		})
	}

	s.logger.Success("Found %d unique locations on homepage", len(locations))

	return locations, nil
}
