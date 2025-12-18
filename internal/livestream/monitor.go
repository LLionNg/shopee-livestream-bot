package livestream

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/LLionNg/shopee-livestream-bot/internal/browser"
	"github.com/LLionNg/shopee-livestream-bot/internal/config"
	"github.com/LLionNg/shopee-livestream-bot/internal/purchase"
	"golang.org/x/sync/errgroup"
)

// Monitor monitors livestreams for product availability
type Monitor struct {
	ctx      context.Context
	cfg      *config.Config
	executor *purchase.Executor
	streams  []string
}

// NewMonitor creates a new livestream monitor
func NewMonitor(ctx context.Context, cfg *config.Config, executor *purchase.Executor) *Monitor {
	return &Monitor{
		ctx:      ctx,
		cfg:      cfg,
		executor: executor,
		streams:  cfg.Shopee.LivestreamURLs,
	}
}

// Start begins monitoring all configured livestreams
func (m *Monitor) Start(ctx context.Context) error {
	fmt.Println("Starting livestream monitoring...")
	fmt.Printf("Monitoring %d livestream(s)\n", len(m.streams))

	// Create error group for concurrent monitoring
	g, ctx := errgroup.WithContext(ctx)

	// Monitor each livestream concurrently
	for i, streamURL := range m.streams {
		streamURL := streamURL // capture variable
		streamID := i + 1

		g.Go(func() error {
			return m.monitorStream(ctx, streamURL, streamID)
		})
	}

	// Wait for all monitors to complete or error
	if err := g.Wait(); err != nil {
		return fmt.Errorf("monitoring error: %w", err)
	}

	return nil
}

// monitorStream monitors a single livestream
func (m *Monitor) monitorStream(ctx context.Context, streamURL string, streamID int) error {
	fmt.Printf("🎥 [Stream %d] Starting monitor: %s\n", streamID, streamURL)

	// Navigate to livestream
	if err := browser.NavigateWithRetry(m.ctx, streamURL, 3); err != nil {
		return fmt.Errorf("failed to navigate to stream %d: %w", streamID, err)
	}

	fmt.Printf("✅ [Stream %d] Successfully loaded livestream\n", streamID)

	// Start monitoring loop
	ticker := time.NewTicker(m.cfg.Monitoring.GetCheckInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 [Stream %d] Stopping monitor\n", streamID)
			return ctx.Err()

		case <-ticker.C:
			// Check for product availability
			if err := m.checkProductAvailability(streamID); err != nil {
				fmt.Printf("⚠️  [Stream %d] Check error: %v\n", streamID, err)
			}
		}
	}
}

// checkProductAvailability checks if products are available for purchase
func (m *Monitor) checkProductAvailability(streamID int) error {
	// Look for "Add to Cart" or "Buy Now" buttons
	// This is a simplified check - real implementation would be more sophisticated

	var buttonExists bool
	
	// Check for various possible selectors
	selectors := []string{
		"button[class*='add-to-cart']",
		"button[class*='buy-now']",
		"button[class*='add-cart']",
		"div[class*='shop-bag'] button",
		".shopee-button-solid",
	}

	for _, selector := range selectors {
		err := chromedp.Run(m.ctx,
			chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, selector), &buttonExists),
		)
		
		if err == nil && buttonExists {
			fmt.Printf("[Stream %d] Product available! Attempting purchase...\n", streamID)
			
			// Attempt to purchase
			if err := m.executor.ExecutePurchase(selector); err != nil {
				fmt.Printf("❌ [Stream %d] Purchase failed: %v\n", streamID, err)
				return err
			}
			
			fmt.Printf("[Stream %d] Purchase successful!\n", streamID)
			return nil
		}
	}

	return nil
}

// CheckFlashSale checks for flash sale countdown
func (m *Monitor) CheckFlashSale(streamID int) (*FlashSale, error) {
	// Look for flash sale timer/countdown
	var hasTimer bool
	
	err := chromedp.Run(m.ctx,
		chromedp.Evaluate(`!!document.querySelector('[class*="countdown"]')`, &hasTimer),
	)
	
	if err != nil || !hasTimer {
		return nil, nil
	}

	// Extract countdown time
	var countdownText string
	err = chromedp.Run(m.ctx,
		chromedp.Text(`[class*="countdown"]`, &countdownText, chromedp.ByQuery),
	)
	
	if err != nil {
		return nil, err
	}

	fmt.Printf("[Stream %d] Flash sale detected: %s\n", streamID, countdownText)

	return &FlashSale{
		StreamID:  streamID,
		Countdown: countdownText,
		Detected:  time.Now(),
	}, nil
}

// GetProductInfo extracts product information from livestream
func (m *Monitor) GetProductInfo() (*ProductInfo, error) {
	var info ProductInfo

	// Extract product name
	var name string
	err := chromedp.Run(m.ctx,
		chromedp.Text(`[class*="product-name"], [class*="product-title"]`, &name, chromedp.ByQuery),
	)
	if err == nil {
		info.Name = name
	}

	// Extract price
	var price string
	err = chromedp.Run(m.ctx,
		chromedp.Text(`[class*="price"], [class*="amount"]`, &price, chromedp.ByQuery),
	)
	if err == nil {
		info.Price = price
	}

	// Extract stock info
	var stock string
	err = chromedp.Run(m.ctx,
		chromedp.Text(`[class*="stock"], [class*="quantity"]`, &stock, chromedp.ByQuery),
	)
	if err == nil {
		info.Stock = stock
	}

	return &info, nil
}

// ProductInfo holds product information
type ProductInfo struct {
	Name  string
	Price string
	Stock string
}

// FlashSale represents a flash sale event
type FlashSale struct {
	StreamID  int
	Countdown string
	Detected  time.Time
}