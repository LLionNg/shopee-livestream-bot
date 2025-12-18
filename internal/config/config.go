package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Shopee     ShopeeConfig     `mapstructure:"shopee"`
	Browser    BrowserConfig    `mapstructure:"browser"`
	Purchase   PurchaseConfig   `mapstructure:"purchase"`
	Proxy      ProxyConfig      `mapstructure:"proxy"`
	Stealth    StealthConfig    `mapstructure:"stealth"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type ShopeeConfig struct {
	BaseURL        string              `mapstructure:"base_url"`
	APIURL         string              `mapstructure:"api_url"`
	LivestreamURLs []string            `mapstructure:"livestream_urls"`
	Credentials    ShopeeCredentials   `mapstructure:"credentials"`
}

type ShopeeCredentials struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Phone    string `mapstructure:"phone"`
}

type BrowserConfig struct {
	Headless    bool           `mapstructure:"headless"`
	Timeout     int            `mapstructure:"timeout"`
	UserDataDir string         `mapstructure:"user_data_dir"`
	Viewport    ViewportConfig `mapstructure:"viewport"`
}

type ViewportConfig struct {
	Width  int `mapstructure:"width"`
	Height int `mapstructure:"height"`
}

type PurchaseConfig struct {
	MaxRetries int `mapstructure:"max_retries"`
	RetryDelay int `mapstructure:"retry_delay"`
}

type ProxyConfig struct {
	Enabled             bool   `mapstructure:"enabled"`
	Rotate              bool   `mapstructure:"rotate"`
	RotationInterval    int    `mapstructure:"rotation_interval"`
	Type                string `mapstructure:"type"`
	ListFile            string `mapstructure:"list_file"`
	TestOnStartup       bool   `mapstructure:"test_on_startup"`
	HealthCheckInterval int    `mapstructure:"health_check_interval"`
}

type StealthConfig struct {
	RandomizeFingerprint bool        `mapstructure:"randomize_fingerprint"`
	RandomDelays         bool        `mapstructure:"random_delays"`
	DelayRange           DelayRange  `mapstructure:"delay_range"`
	UserAgentsFile       string      `mapstructure:"user_agents_file"`
}

type DelayRange struct {
	Min int `mapstructure:"min"`
	Max int `mapstructure:"max"`
}

type MonitoringConfig struct {
	CheckInterval       int              `mapstructure:"check_interval"`
	MaxConcurrentStreams int             `mapstructure:"max_concurrent_streams"`
	Notifications       NotificationConfig `mapstructure:"notifications"`
}

type NotificationConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	WebhookURL string `mapstructure:"webhook_url"`
}

type LoggingConfig struct {
	Level         string `mapstructure:"level"`
	Format        string `mapstructure:"format"`
	Output        string `mapstructure:"output"`
	ConsoleOutput bool   `mapstructure:"console_output"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
}

// Load reads configuration from file and environment
func Load(configPath string) (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	// Set up viper
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with environment variables
	cfg.Shopee.Credentials.Username = getEnv("SHOPEE_USERNAME", cfg.Shopee.Credentials.Username)
	cfg.Shopee.Credentials.Password = getEnv("SHOPEE_PASSWORD", cfg.Shopee.Credentials.Password)
	cfg.Shopee.Credentials.Phone = getEnv("SHOPEE_PHONE", cfg.Shopee.Credentials.Phone)

	// Clear placeholder values if they weren't replaced
	if cfg.Shopee.Credentials.Username == "${SHOPEE_USERNAME}" {
		cfg.Shopee.Credentials.Username = ""
	}
	if cfg.Shopee.Credentials.Password == "${SHOPEE_PASSWORD}" {
		cfg.Shopee.Credentials.Password = ""
	}
	if cfg.Shopee.Credentials.Phone == "${SHOPEE_PHONE}" {
		cfg.Shopee.Credentials.Phone = ""
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.Shopee.BaseURL == "" {
		return fmt.Errorf("shopee.base_url is required")
	}
	if len(c.Shopee.LivestreamURLs) == 0 {
		return fmt.Errorf("at least one livestream URL is required")
	}
	// Credentials are optional - manual login will be used if not provided
	if c.Browser.Timeout <= 0 {
		c.Browser.Timeout = 30
	}
	if c.Purchase.MaxRetries <= 0 {
		c.Purchase.MaxRetries = 3
	}
	return nil
}

// GetTimeout returns browser timeout as duration
func (c *BrowserConfig) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// GetRetryDelay returns retry delay as duration
func (c *PurchaseConfig) GetRetryDelay() time.Duration {
	return time.Duration(c.RetryDelay) * time.Second
}

// GetCheckInterval returns monitoring check interval as duration
func (c *MonitoringConfig) GetCheckInterval() time.Duration {
	return time.Duration(c.CheckInterval) * time.Second
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}