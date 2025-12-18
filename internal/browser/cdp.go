package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/LLionNg/shopee-livestream-bot/internal/config"
)

// Initialize creates and configures a browser context
func Initialize(ctx context.Context, cfg *config.Config) (context.Context, context.CancelFunc) {
	// Prepare Chrome options
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	}

	// Add headless mode if configured
	if cfg.Browser.Headless {
		opts = append(opts, chromedp.Headless)
	}

	// Add user data directory for session persistence
	if cfg.Browser.UserDataDir != "" {
		opts = append(opts, chromedp.UserDataDir(cfg.Browser.UserDataDir))
	}

	// Add window size
	opts = append(opts,
		chromedp.WindowSize(cfg.Browser.Viewport.Width, cfg.Browser.Viewport.Height),
	)

	// Add stealth options to avoid detection
	opts = append(opts, getStealthOptions()...)

	// Create allocator context
	allocCtx, _ := chromedp.NewExecAllocator(ctx, opts...)

	// Create browser context
	browserCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))

	// Set timeout
	browserCtx, _ = context.WithTimeout(browserCtx, cfg.Browser.GetTimeout())

	return browserCtx, cancel
}

// getStealthOptions returns options to avoid bot detection
func getStealthOptions() []chromedp.ExecAllocatorOption {
	return []chromedp.ExecAllocatorOption{
		// Disable automation flags
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),

		// Set realistic user agent
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),

		// Additional anti-detection flags
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-component-extensions-with-background-pages", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
	}
}

// NavigateWithRetry navigates to a URL with retry logic
func NavigateWithRetry(ctx context.Context, url string, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitReady("body", chromedp.ByQuery),
		)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return fmt.Errorf("failed to navigate after %d retries: %w", maxRetries, err)
}

// WaitForElement waits for an element to be visible
func WaitForElement(ctx context.Context, selector string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	)
}

// Click clicks on an element
func Click(ctx context.Context, selector string) error {
	return chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Click(selector, chromedp.ByQuery),
	)
}

// Type types text into an input field
func Type(ctx context.Context, selector, text string) error {
	return chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Clear(selector, chromedp.ByQuery),
		chromedp.SendKeys(selector, text, chromedp.ByQuery),
	)
}

// GetText retrieves text content from an element
func GetText(ctx context.Context, selector string) (string, error) {
	var text string
	err := chromedp.Run(ctx,
		chromedp.Text(selector, &text, chromedp.ByQuery),
	)
	return text, err
}

// Screenshot takes a screenshot of the current page
func Screenshot(ctx context.Context, filepath string) error {
	var buf []byte
	err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf))
	if err != nil {
		return err
	}

	// Save to file (implement file writing here if needed)
	return nil
}

// ExecuteJS executes JavaScript code
func ExecuteJS(ctx context.Context, script string, res interface{}) error {
	return chromedp.Run(ctx,
		chromedp.Evaluate(script, res),
	)
}

// RemoveWebDriverFlag removes the webdriver property to avoid detection
func RemoveWebDriverFlag(ctx context.Context) error {
	script := `
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});
	`
	return ExecuteJS(ctx, script, nil)
}

// RandomizeFingerprint attempts to randomize browser fingerprint
func RandomizeFingerprint(ctx context.Context) error {
	// Randomize various browser properties
	script := `
		// Override navigator properties
		Object.defineProperty(navigator, 'platform', {
			get: () => 'Win32'
		});
		
		Object.defineProperty(navigator, 'vendor', {
			get: () => 'Google Inc.'
		});

		// Add Chrome property
		window.chrome = {
			runtime: {}
		};

		// Randomize canvas fingerprint (simplified)
		const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
		HTMLCanvasElement.prototype.toDataURL = function() {
			const context = this.getContext('2d');
			if (context) {
				context.fillStyle = 'rgba(' + Math.random() + ',' + Math.random() + ',' + Math.random() + ',0.01)';
				context.fillRect(0, 0, 1, 1);
			}
			return originalToDataURL.apply(this, arguments);
		};
	`
	return ExecuteJS(ctx, script, nil)
}

// WaitForNavigation waits for page navigation to complete
func WaitForNavigation(ctx context.Context) error {
	return chromedp.Run(ctx,
		chromedp.WaitReady("body", chromedp.ByQuery),
	)
}