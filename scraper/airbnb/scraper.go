package airbnb

import (
	"context"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/farhanasfar/airbnb-market-scraping-system/config"
	"github.com/farhanasfar/airbnb-market-scraping-system/utils"
)

// Scraper handles Airbnb scraping operations
type Scraper struct {
	cfg    *config.ScraperConfig
	logger *utils.Logger
}

// NewScraper creates a new Airbnb scraper instance
func NewScraper(cfg *config.ScraperConfig, logger *utils.Logger) *Scraper {
	return &Scraper{
		cfg:    cfg,
		logger: logger,
	}
}

// createStealthContext creates a browser context with anti-detection settings
func (scrape *Scraper) createStealthContext(parentCtx context.Context) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// avoiding bot detection
		chromedp.Flag("headless", scrape.cfg.Headless),
		chromedp.WindowSize(1440, 900),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //not loading images to scrape fast.
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(parentCtx, opts...)
	ctx, cancelCtx := chromedp.NewContext(allocCtx)

	return ctx, func() {
		cancelCtx()
		cancelAlloc()
	}
}

// removeWebdriverProperty removes the webdriver property that sites check
func removeWebdriverProperty() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		err := chromedp.Evaluate(`
			Object.defineProperty(navigator, 'webdriver', {
				get: () => undefined
			})
		`, nil).Do(ctx)
		return err
	})
}

// randomDelay adding human-like random delay
func (scrape *Scraper) randomDelay() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		minMs := scrape.cfg.DelayMinMs
		maxMs := scrape.cfg.DelayMaxMs
		delay := time.Duration(minMs+rand.Intn(maxMs-minMs)) * time.Millisecond

		scrape.logger.Info("Waiting %v before next action...", delay)
		time.Sleep(delay)
		return nil
	})
}
