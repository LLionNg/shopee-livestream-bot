package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/LLionNg/shopee-livestream-bot/internal/browser"
	"github.com/LLionNg/shopee-livestream-bot/internal/config"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Manager handles authentication and session management
type Manager struct {
	ctx         context.Context
	cfg         *config.Config
	sessionFile string
	cookies     []*network.Cookie
	isLoggedIn  bool
}

// NewManager creates a new authentication manager
func NewManager(ctx context.Context, cfg *config.Config) *Manager {
	return &Manager{
		ctx:         ctx,
		cfg:         cfg,
		sessionFile: "data/cookies/session.json",
		isLoggedIn:  false,
	}
}

// Login performs login to Shopee
func (m *Manager) Login() error {
	// Try to load existing session first
	if m.LoadSession() {
		if m.ValidateSession() {
			m.isLoggedIn = true
			return nil
		}
	}

	// Perform fresh login
	return m.PerformLogin()
}

// PerformLogin executes the login flow
func (m *Manager) PerformLogin() error {
	// Navigate to Shopee login page
	loginURL := m.cfg.Shopee.BaseURL + "/buyer/login"

	if err := browser.NavigateWithRetry(m.ctx, loginURL, 3); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Wait for page to load
	time.Sleep(2 * time.Second)

	// Check if already logged in (redirect to homepage)
	var currentURL string
	if err := chromedp.Run(m.ctx, chromedp.Location(&currentURL)); err != nil {
		return err
	}

	if currentURL != loginURL && !contains(currentURL, "/buyer/login") {
		// Already logged in
		return m.SaveSession()
	}

	// Fill in login form
	// Note: Shopee may require phone/email + password or phone + OTP
	// This is a simplified version - real implementation may need SMS OTP handling

	// Wait for login form
	if err := browser.WaitForElement(m.ctx, "input[type='text']", 10*time.Second); err != nil {
		return fmt.Errorf("login form not found: %w", err)
	}

	// Method 1: Username/Email + Password
	if m.cfg.Shopee.Credentials.Username != "" && m.cfg.Shopee.Credentials.Password != "" {
		// Enter username/email
		if err := browser.Type(m.ctx, "input[type='text']", m.cfg.Shopee.Credentials.Username); err != nil {
			return fmt.Errorf("failed to enter username: %w", err)
		}

		time.Sleep(500 * time.Millisecond)

		// Enter password
		if err := browser.Type(m.ctx, "input[type='password']", m.cfg.Shopee.Credentials.Password); err != nil {
			return fmt.Errorf("failed to enter password: %w", err)
		}

		time.Sleep(500 * time.Millisecond)

		// Click login button
		if err := browser.Click(m.ctx, "button[type='submit']"); err != nil {
			return fmt.Errorf("failed to click login button: %w", err)
		}

		// Wait for login to complete (check for redirect or success indicator)
		time.Sleep(5 * time.Second)

		// Check if login was successful
		if err := chromedp.Run(m.ctx, chromedp.Location(&currentURL)); err != nil {
			return err
		}

		if contains(currentURL, "/buyer/login") {
			return fmt.Errorf("login failed - still on login page")
		}

		// Save session after successful login
		return m.SaveSession()
	}

	return fmt.Errorf("no valid login credentials provided")
}

// SaveSession saves current session cookies to file
func (m *Manager) SaveSession() error {
	// Get all cookies
	var cookies []*network.Cookie
	if err := chromedp.Run(m.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		c, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}
		cookies = c
		return nil
	})); err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	m.cookies = cookies

	// Create directory if not exists
	os.MkdirAll("data/cookies", 0755)

	// Save to file
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	if err := os.WriteFile(m.sessionFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	m.isLoggedIn = true
	return nil
}

// LoadSession loads session cookies from file
func (m *Manager) LoadSession() bool {
	// Check if session file exists
	if _, err := os.Stat(m.sessionFile); os.IsNotExist(err) {
		return false
	}

	// Read session file
	data, err := os.ReadFile(m.sessionFile)
	if err != nil {
		return false
	}

	// Unmarshal cookies
	var cookies []*network.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return false
	}

	// Set cookies in browser
	if err := chromedp.Run(m.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		for _, cookie := range cookies {
			if err := network.SetCookie(cookie.Name, cookie.Value).
				WithDomain(cookie.Domain).
				WithPath(cookie.Path).
				WithHTTPOnly(cookie.HTTPOnly).
				WithSecure(cookie.Secure).
				Do(ctx); err != nil {
				return err
			}
		}
		return nil
	})); err != nil {
		return false
	}

	m.cookies = cookies
	return true
}

// ValidateSession checks if the current session is still valid
func (m *Manager) ValidateSession() bool {
	// Navigate to a page that requires authentication
	if err := browser.NavigateWithRetry(m.ctx, m.cfg.Shopee.BaseURL, 3); err != nil {
		return false
	}

	time.Sleep(2 * time.Second)

	// Check current URL
	var currentURL string
	if err := chromedp.Run(m.ctx, chromedp.Location(&currentURL)); err != nil {
		return false
	}

	// If redirected to login page, session is invalid
	if contains(currentURL, "/buyer/login") {
		return false
	}

	// Try to find user-specific elements (e.g., profile icon)
	// This is a simplified check
	var userExists bool
	err := chromedp.Run(m.ctx,
		chromedp.Evaluate(`!!document.querySelector('[data-testid="account-menu"]')`, &userExists),
	)

	return err == nil && userExists
}

// IsLoggedIn returns whether user is currently logged in
func (m *Manager) IsLoggedIn() bool {
	return m.isLoggedIn
}

// Logout performs logout
func (m *Manager) Logout() error {
	// Clear cookies
	if err := chromedp.Run(m.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return network.ClearBrowserCookies().Do(ctx)
	})); err != nil {
		return err
	}

	// Delete session file
	os.Remove(m.sessionFile)

	m.isLoggedIn = false
	m.cookies = nil

	return nil
}

// RefreshSession refreshes the current session
func (m *Manager) RefreshSession() error {
	if !m.ValidateSession() {
		return m.PerformLogin()
	}
	return nil
}

// helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) > len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
