package purchase

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/LLionNg/shopee-livestream-bot/internal/config"
)

// Executor handles the purchase execution flow
type Executor struct {
	ctx context.Context
	cfg *config.Config
}

// NewExecutor creates a new purchase executor
func NewExecutor(ctx context.Context, cfg *config.Config) *Executor {
	return &Executor{
		ctx: ctx,
		cfg: cfg,
	}
}

// ExecutePurchase adds the product to cart (items are auto-reserved once in cart)
func (e *Executor) ExecutePurchase(productSelector string) error {
	fmt.Println("🛒 Adding item to cart...")

	// Add to cart - items are automatically reserved during livestream
	if err := e.AddToCart(productSelector); err != nil {
		return fmt.Errorf("failed to add to cart: %w", err)
	}

	fmt.Println("✅ Item successfully added to cart and reserved!")

	return nil
}

// AddToCart adds the product to the cart
func (e *Executor) AddToCart(selector string) error {
	// Wait for the button to be clickable
	ctx, cancel := context.WithTimeout(e.ctx, 5*time.Second)
	defer cancel()
	
	// Click the add to cart button
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Click(selector, chromedp.ByQuery),
	)
	
	if err != nil {
		return fmt.Errorf("failed to click add to cart: %w", err)
	}
	
	// Wait for cart update animation
	time.Sleep(1 * time.Second)
	
	return nil
}


// RetryPurchase retries adding to cart with exponential backoff
func (e *Executor) RetryPurchase(productSelector string) error {
	maxRetries := e.cfg.Purchase.MaxRetries
	retryDelay := e.cfg.Purchase.GetRetryDelay()

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			waitTime := retryDelay * time.Duration(i)
			fmt.Printf("🔄 Retry %d/%d after %v...\n", i+1, maxRetries, waitTime)
			time.Sleep(waitTime)
		}

		err := e.ExecutePurchase(productSelector)
		if err == nil {
			return nil
		}

		lastErr = err
		fmt.Printf("⚠️  Attempt %d failed: %v\n", i+1, err)
	}

	return fmt.Errorf("all %d attempts to add to cart failed: %w", maxRetries, lastErr)
}

// GetCartItemCount returns the number of items in cart
func (e *Executor) GetCartItemCount() (int, error) {
	var count int
	err := chromedp.Run(e.ctx,
		chromedp.Evaluate(`parseInt(document.querySelector('[class*="cart-count"]')?.innerText || '0')`, &count),
	)
	return count, err
}

// ClearCart removes all items from the cart
func (e *Executor) ClearCart() error {
	// Navigate to cart
	cartURL := e.cfg.Shopee.BaseURL + "/cart"
	if err := chromedp.Run(e.ctx, chromedp.Navigate(cartURL)); err != nil {
		return fmt.Errorf("failed to navigate to cart: %w", err)
	}

	time.Sleep(2 * time.Second) // Wait for cart page to load

	// Select all items and delete
	err := chromedp.Run(e.ctx,
		chromedp.Click("input[type='checkbox'][class*='select-all']", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Click("button[class*='delete']", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Click("button[class*='confirm']", chromedp.ByQuery),
	)

	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	fmt.Println("🗑️  Cart cleared successfully")
	return nil
}

