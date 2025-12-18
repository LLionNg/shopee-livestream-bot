package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LLionNg/shopee-livestream-bot/internal/auth"
	"github.com/LLionNg/shopee-livestream-bot/internal/browser"
	"github.com/LLionNg/shopee-livestream-bot/internal/config"
	"github.com/LLionNg/shopee-livestream-bot/internal/livestream"
	"github.com/LLionNg/shopee-livestream-bot/internal/purchase"
	"github.com/LLionNg/shopee-livestream-bot/pkg/logger"
)

const (
	appName    = "Shopee Livestream Bot"
	appVersion = "1.0.0"
)

func main() {
	// Print banner
	printBanner()

	// Initialize logger
	log := logger.New("info", true)
	log.Info("Starting Shopee Livestream Bot...")

	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}
	log.Info("Configuration loaded successfully")

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Initialize browser
	log.Info("Initializing browser...")
	browserCtx, browserCancel := browser.Initialize(ctx, cfg)
	if browserCtx == nil {
		log.Fatal("Failed to initialize browser - please check Chrome installation")
	}
	defer browserCancel()

	log.Info("Browser initialized successfully")

	// Initialize authentication
	log.Info("Authenticating with Shopee...")
	authManager := auth.NewManager(browserCtx, cfg)
	if err := authManager.Login(); err != nil {
		log.Fatal("Authentication failed", "error", err)
	}
	log.Info("Authentication successful!")

	// Initialize purchase executor
	purchaseExec := purchase.NewExecutor(browserCtx, cfg)

	// Initialize livestream monitor
	log.Info("Starting livestream monitor...")
	monitor := livestream.NewMonitor(browserCtx, cfg, purchaseExec)

	// Start monitoring in a goroutine
	go func() {
		if err := monitor.Start(ctx); err != nil {
			log.Error("Monitor stopped with error", "error", err)
		}
	}()

	log.Info("Bot is now running! Monitoring livestreams...")
	log.Info("Press Ctrl+C to stop")

	// Wait for shutdown signal
	<-sigChan
	log.Info("Shutdown signal received, cleaning up...")

	// Cancel context and wait for cleanup
	cancel()
	time.Sleep(2 * time.Second)

	log.Info("Bot stopped. Goodbye!")
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║     🛒  SHOPEE LIVESTREAM AUTO-PURCHASE BOT  🛒          ║
║                                                           ║
║              Version: %s                              ║
║              Made with ❤️  in Go                          ║
║                                                           ║
║     ⚠️  Use responsibly & at your own risk ⚠️            ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Printf(banner, appVersion)
	fmt.Println()
}
