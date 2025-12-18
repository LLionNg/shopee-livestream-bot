package purchase

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/LLionNg/shopee-livestream-bot/internal/browser"
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

// ExecutePurchase executes the complete purchase flow
func (e *Executor) ExecutePurchase(productSelector string) error {
	fmt.Println("🛒 Starting purchase execution...")
	
	// Step 1: Add to cart
	if err := e.AddToCart(productSelector); err != nil {
		return fmt.Errorf("failed to add to cart: %w", err)
	}
	
	fmt.Println("Added to cart")
	
	// Small delay to mimic human behavior
	time.Sleep(500 * time.Millisecond)
	
	// Step 2: Navigate to cart (if not auto-checkout)
	if !e.cfg.Purchase.AutoCheckout {
		return nil // Stop here if auto-checkout is disabled
	}
	
	// Step 3: Proceed to checkout
	if err := e.ProceedToCheckout(); err != nil {
		return fmt.Errorf("failed to proceed to checkout: %w", err)
	}
	
	fmt.Println("Proceeded to checkout")
	
	// Step 4: Place order
	if err := e.PlaceOrder(); err != nil {
		return fmt.Errorf("failed to place order: %w", err)
	}
	
	fmt.Println("Order placed successfully!")
	
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

// ProceedToCheckout navigates to checkout
func (e *Executor) ProceedToCheckout() error {
	// Look for checkout button
	checkoutSelectors := []string{
		"button[class*='checkout']",
		"a[href*='checkout']",
		".shopee-button-solid--primary",
		"button:contains('Checkout')",
	}
	
	ctx, cancel := context.WithTimeout(e.ctx, e.cfg.Purchase.GetCheckoutTimeout())
	defer cancel()
	
	for _, selector := range checkoutSelectors {
		err := chromedp.Run(ctx,
			chromedp.WaitVisible(selector, chromedp.ByQuery),
			chromedp.Click(selector, chromedp.ByQuery),
		)
		
		if err == nil {
			// Wait for checkout page to load
			time.Sleep(2 * time.Second)
			return nil
		}
	}
	
	return fmt.Errorf("checkout button not found")
}

// PlaceOrder completes the order placement
func (e *Executor) PlaceOrder() error {
	ctx, cancel := context.WithTimeout(e.ctx, e.cfg.Purchase.GetCheckoutTimeout())
	defer cancel()
	
	// Look for place order button
	placeOrderSelectors := []string{
		"button[class*='place-order']",
		"button[class*='submit-order']",
		"button:contains('Place Order')",
		"button:contains('ยืนยันคำสั่งซื้อ')", // Thai: Confirm Order
	}
	
	for _, selector := range placeOrderSelectors {
		err := chromedp.Run(ctx,
			chromedp.WaitVisible(selector, chromedp.ByQuery),
			chromedp.Click(selector, chromedp.ByQuery),
		)
		
		if err == nil {
			// Wait for order confirmation
			time.Sleep(3 * time.Second)
			
			// Verify order success
			if e.VerifyOrderSuccess() {
				return nil
			}
			
			return fmt.Errorf("order placement failed - no confirmation")
		}
	}
	
	return fmt.Errorf("place order button not found")
}

// VerifyOrderSuccess checks if the order was successfully placed
func (e *Executor) VerifyOrderSuccess() bool {
	// Look for success indicators
	successSelectors := []string{
		"[class*='order-success']",
		"[class*='payment-success']",
		".success-icon",
		"h1:contains('Order Placed')",
	}
	
	for _, selector := range successSelectors {
		var exists bool
		err := chromedp.Run(e.ctx,
			chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector('%s')`, selector), &exists),
		)
		
		if err == nil && exists {
			return true
		}
	}
	
	// Check URL for success page
	var url string
	err := chromedp.Run(e.ctx, chromedp.Location(&url))
	if err == nil && (contains(url, "success") || contains(url, "complete")) {
		return true
	}
	
	return false
}

// QuickPurchase executes the fastest possible purchase flow
// This assumes everything is pre-configured (payment method, address, etc.)
func (e *Executor) QuickPurchase(productSelector string) error {
	fmt.Println("Executing QUICK purchase...")
	
	startTime := time.Now()
	
	// Create a tight timeout context
	ctx, cancel := context.WithTimeout(e.ctx, 3*time.Second)
	defer cancel()
	
	// Execute all steps in one go
	err := chromedp.Run(ctx,
		// Click add to cart
		chromedp.Click(productSelector, chromedp.ByQuery),
		chromedp.Sleep(300*time.Millisecond),
		
		// Click checkout (assuming it appears immediately)
		chromedp.Click("button[class*='checkout']", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		
		// Click place order
		chromedp.Click("button[class*='place-order']", chromedp.ByQuery),
	)
	
	elapsed := time.Since(startTime)
	
	if err != nil {
		fmt.Printf("Quick purchase failed in %.2fs: %v\n", elapsed.Seconds(), err)
		return err
	}
	
	fmt.Printf("Quick purchase completed in %.2fs!\n", elapsed.Seconds())
	
	return nil
}

// RetryPurchase retries purchase execution with exponential backoff
func (e *Executor) RetryPurchase(productSelector string) error {
	maxRetries := e.cfg.Purchase.MaxRetries
	retryDelay := e.cfg.Purchase.GetRetryDelay()
	
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			waitTime := retryDelay * time.Duration(i)
			fmt.Printf("Retry %d/%d after %v...\n", i+1, maxRetries, waitTime)
			time.Sleep(waitTime)
		}
		
		err := e.ExecutePurchase(productSelector)
		if err == nil {
			return nil
		}
		
		lastErr = err
		fmt.Printf("Attempt %d failed: %v\n", i+1, err)
	}
	
	return fmt.Errorf("all %d purchase attempts failed: %w", maxRetries, lastErr)
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
	if err := browser.NavigateWithRetry(e.ctx, cartURL, 3); err != nil {
		return err
	}
	
	// Select all items and delete
	err := chromedp.Run(e.ctx,
		chromedp.Click("input[type='checkbox'][class*='select-all']", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Click("button[class*='delete']", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Click("button[class*='confirm']", chromedp.ByQuery),
	)
	
	return err
}

// helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}