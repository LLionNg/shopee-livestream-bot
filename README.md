# Shopee Livestream Bot

A Go-based automation tool for monitoring and interacting with Shopee Thailand livestreams. This bot handles authentication, session management, livestream monitoring, and automated purchase execution.

## Features

- Browser automation using Chrome DevTools Protocol
- Manual and automated login with session persistence
- Livestream URL monitoring
- Automated purchase execution
- Configurable purchase parameters
- Structured logging

## Requirements

- Go 1.21 or higher
- Chrome/Chromium browser
- Windows/Linux/macOS

## Installation

```bash
# Clone the repository
git clone https://github.com/LLionNg/shopee-livestream-bot.git
cd shopee-livestream-bot

# Install dependencies
go mod download
```

## Configuration

### 1. Environment Variables

Copy the example environment file and configure your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your settings. Leave credentials empty to use manual login mode (supports Facebook, Google, or username/password):

```bash
SHOPEE_USERNAME=
SHOPEE_PASSWORD=
LOG_LEVEL=info
```

### 2. Bot Configuration

Edit `configs/config.yaml` to configure:
- Livestream URLs to monitor
- Browser settings (headless mode, viewport size)
- Purchase retry settings
- Logging preferences

See `configs/config.yaml` for detailed configuration options.

## Usage

### Build

```bash
# Build the binary
go build -o shopee-bot cmd/bot/main.go

# Or use make
make build
```

### Run

```bash
# Run directly with Go
go run cmd/bot/main.go

# Or run the built binary
./shopee-bot
```

### Manual Login

When credentials are not provided in `.env`, the bot will:
1. Open Chrome browser window
2. Navigate to Shopee login page
3. Wait for you to login manually using any method
4. Detect successful login automatically
5. Save session for future runs

## Project Structure

```
shopee-livestream-bot/
├── cmd/
│   └── bot/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   └── auth.go              # Authentication and session management
│   ├── browser/
│   │   └── cdp.go               # Chrome DevTools Protocol integration
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── livestream/
│   │   └── monitor.go           # Livestream monitoring
│   └── purchase/
│       └── executor.go          # Purchase execution logic
├── pkg/
│   └── logger/
│       └── logger.go            # Structured logging
├── configs/
│   └── config.yaml              # Main configuration file
├── data/
│   ├── browser/                 # Browser session data
│   └── logs/                    # Application logs
├── .env.example                 # Environment variables template
├── Makefile                     # Build automation
└── README.md
```

## Development

### Build Commands

```bash
# Build for current platform
make build

# Build for Linux
make build-linux

# Build for Windows
make build-windows

# Build for all platforms
make build-all
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Code Formatting

```bash
# Format code
make fmt

# Run go vet
make vet
```

## Dependencies

This project uses the following main dependencies:

- **chromedp/chromedp** - Chrome DevTools Protocol for browser automation
- **chromedp/cdproto** - Chrome DevTools Protocol definitions
- **spf13/viper** - Configuration management
- **joho/godotenv** - Environment variable loading
- **golang.org/x/sync** - Concurrency utilities

For a complete list, see `go.mod`.

## Troubleshooting

### Browser Not Opening

Ensure Chrome/Chromium is installed and accessible in your system PATH.

### Login Fails

- Check your credentials in `.env`
- Try manual login mode (leave credentials empty)
- Clear browser data in `data/browser/` directory

### Session Expired

The bot saves sessions in `data/browser/`. If login fails, delete this directory to force a fresh login.

## License

This project is for educational purposes only. Use responsibly and in accordance with Shopee's Terms of Service.
