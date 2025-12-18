# Shopee Livestream Auto-Purchase Bot - Project Structure

## Project Overview

**Project Name:** shopee-livestream-bot  
**Language:** Go (Golang)  
**Target:** Shopee Thailand Livestream Flash Sales  
**Architecture:** Modular, concurrent, production-ready

---

## Directory Structure

```
shopee-livestream-bot/
├── cmd/
│   └── bot/
│       └── main.go                 # Application entry point
├── internal/
│   ├── auth/
│   │   ├── auth.go                 # Authentication logic
│   │   ├── session.go              # Session management
│   │   └── cookies.go              # Cookie handling
│   ├── browser/
│   │   ├── cdp.go                  # Chrome DevTools Protocol
│   │   ├── stealth.go              # Anti-detection measures
│   │   └── fingerprint.go          # Browser fingerprinting
│   ├── livestream/
│   │   ├── monitor.go              # Livestream monitoring
│   │   ├── detector.go             # Product availability detection
│   │   └── parser.go               # HTML/JSON parsing
│   ├── purchase/
│   │   ├── executor.go             # Purchase execution
│   │   ├── cart.go                 # Cart management
│   │   └── checkout.go             # Checkout flow
│   ├── proxy/
│   │   ├── manager.go              # Proxy rotation
│   │   └── pool.go                 # Proxy pool management
│   └── config/
│       ├── config.go               # Configuration management
│       └── env.go                  # Environment variables
├── pkg/
│   ├── logger/
│   │   └── logger.go               # Structured logging
│   └── utils/
│       ├── retry.go                # Retry logic
│       ├── delay.go                # Human-like delays
│       └── validators.go           # Input validation
├── configs/
│   ├── config.yaml                 # Main configuration
│   ├── proxies.txt                 # Proxy list
│   └── user_agents.txt             # User agent list
├── data/
│   ├── cookies/
│   │   └── session.json            # Saved cookies
│   └── logs/
│       └── app.log                 # Application logs
├── scripts/
│   ├── setup.sh                    # Setup script
│   └── build.sh                    # Build script
├── tests/
│   ├── auth_test.go
│   ├── browser_test.go
│   └── purchase_test.go
├── docs/
│   ├── SETUP.md                    # Setup instructions
│   ├── USAGE.md                    # Usage guide
│   └── API.md                      # API documentation
├── .env.example                    # Environment variables template
├── .gitignore
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
├── Makefile                        # Build automation
└── README.md                       # Project documentation
```

---

## Module Breakdown

### 1. **cmd/bot/main.go**
- Application entry point
- CLI argument parsing
- Initialization and orchestration
- Graceful shutdown handling

### 2. **internal/auth/**
- User authentication
- Session management
- Cookie persistence
- Token refresh

### 3. **internal/browser/**
- Chrome DevTools Protocol integration
- Browser automation
- Anti-detection stealth measures
- Fingerprint randomization

### 4. **internal/livestream/**
- Monitor livestream URLs
- Detect product availability
- Parse product information
- Event-driven architecture

### 5. **internal/purchase/**
- Fast checkout execution
- Cart management
- Payment processing
- Order confirmation

### 6. **internal/proxy/**
- Proxy rotation
- Health checking
- Connection pooling
- Geographic targeting

### 7. **internal/config/**
- Configuration loading
- Environment management
- Settings validation

### 8. **pkg/logger/**
- Structured logging
- Log levels
- Output formatting

### 9. **pkg/utils/**
- Retry mechanisms
- Random delays
- Validation helpers

---

## Go Dependencies (go.mod)

```go
module github.com/LLionNg/shopee-livestream-bot

go 1.21

require (
    // Browser automation
    github.com/chromedp/chromedp v0.9.3
    github.com/chromedp/cdproto v0.0.0-20231205062650-00455a960d61
    
    // HTTP & Networking
    github.com/go-resty/resty/v2 v2.11.0
    golang.org/x/net v0.19.0
    
    // Configuration
    github.com/spf13/viper v1.18.2
    github.com/joho/godotenv v1.5.1
    
    // Logging
    github.com/sirupsen/logrus v1.9.3
    go.uber.org/zap v1.26.0
    
    // CLI
    github.com/spf13/cobra v1.8.0
    
    // Utilities
    github.com/google/uuid v1.5.0
    github.com/tidwall/gjson v1.17.0
    github.com/PuerkitoBio/goquery v1.8.1
    
    // Concurrency
    golang.org/x/sync v0.5.0
    
    // Anti-detection
    github.com/EDDYCJY/fake-useragent v0.2.0
)
```

---

## Configuration Files

### config.yaml
```yaml
app:
  name: "Shopee Livestream Bot"
  version: "1.0.0"
  environment: "development"

shopee:
  base_url: "https://shopee.co.th"
  api_url: "https://shopee.co.th/api/v4"
  livestream_urls:
    - "https://live.shopee.co.th/..."
  
  credentials:
    username: "${SHOPEE_USERNAME}"
    password: "${SHOPEE_PASSWORD}"
    phone: "${SHOPEE_PHONE}"

browser:
  headless: false  # Set to true for production
  timeout: 30
  user_data_dir: "./data/browser"
  
  viewport:
    width: 1920
    height: 1080

purchase:
  max_retries: 3
  retry_delay: 1  # seconds
  checkout_timeout: 5  # seconds
  
  auto_checkout: true
  pre_fill_cart: true
  
  payment_method: "ShopeePay"  # or "COD", "Credit Card"

proxy:
  enabled: true
  rotate: true
  rotation_interval: 300  # seconds
  type: "residential"  # or "datacenter"
  
  list_file: "./configs/proxies.txt"
  
  test_on_startup: true
  health_check_interval: 60

stealth:
  randomize_fingerprint: true
  random_delays: true
  delay_range:
    min: 100  # milliseconds
    max: 500
  
  user_agents_file: "./configs/user_agents.txt"

monitoring:
  check_interval: 1  # seconds
  max_concurrent_streams: 5
  
  notifications:
    enabled: true
    webhook_url: "${WEBHOOK_URL}"

logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json or text
  output: "./data/logs/app.log"
  
  console_output: true
  max_size: 100  # MB
  max_backups: 5
  max_age: 30  # days
```

### .env.example
```bash
# Shopee Credentials
SHOPEE_USERNAME=your_email@example.com
SHOPEE_PASSWORD=your_password
SHOPEE_PHONE=+66812345678

# Proxy Configuration (optional)
PROXY_URL=http://username:password@proxy.example.com:8080

# Notifications (optional)
WEBHOOK_URL=https://discord.com/api/webhooks/...
TELEGRAM_BOT_TOKEN=your_token
TELEGRAM_CHAT_ID=your_chat_id

# Application Settings
ENV=development
LOG_LEVEL=info
```

---

## Core Features

### Phase 1 - MVP (Minimum Viable Product)
- ✅ Browser automation with CDP
- ✅ Login and session management
- ✅ Single livestream monitoring
- ✅ Basic purchase execution
- ✅ Configuration management
- ✅ Logging

### Phase 2 - Enhanced Features
- ✅ Multiple livestream monitoring
- ✅ Proxy rotation
- ✅ Anti-detection measures
- ✅ Retry mechanisms
- ✅ Error recovery

### Phase 3 - Advanced Features
- ✅ Predictive monitoring (ML-based timing)
- ✅ Multi-account support
- ✅ Advanced fingerprinting
- ✅ Performance optimization
- ✅ Dashboard/UI

---

## Development Workflow

### 1. Setup Development Environment
```bash
# Install Go (1.21+)
# Install Chrome/Chromium

# Clone repository
git clone <repo>
cd shopee-livestream-bot

# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Edit configuration
nano .env
nano configs/config.yaml
```

### 2. Build Project
```bash
# Build binary
go build -o bin/bot cmd/bot/main.go

# Or use Makefile
make build
```

### 3. Run Tests
```bash
# Run all tests
go test ./...

# Run specific module tests
go test ./internal/auth/...

# Run with coverage
go test -cover ./...
```

### 4. Run Application
```bash
# Development mode
go run cmd/bot/main.go

# Production mode (built binary)
./bin/bot

# With flags
./bin/bot --config=configs/config.yaml --env=production
```

---

## Testing Strategy

### Unit Tests
- Test individual functions
- Mock external dependencies
- Focus on business logic

### Integration Tests
- Test module interactions
- Real browser automation (slower)
- Configuration loading

### End-to-End Tests
- Test complete purchase flow
- Use test livestreams
- Verify success metrics

---

## Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| Startup Time | < 5s | - |
| Login Time | < 3s | - |
| Purchase Execution | < 2s | - |
| Memory Usage | < 500MB | - |
| CPU Usage | < 50% | - |
| Success Rate | > 80% | - |

---

## Security Considerations

1. **Credentials**: Never commit credentials to git
2. **Cookies**: Encrypt session cookies
3. **Proxies**: Validate proxy sources
4. **Logs**: Sanitize sensitive data from logs
5. **Rate Limiting**: Respect Shopee's limits

---

## Deployment

### Local Development
```bash
go run cmd/bot/main.go
```

### Production (VPS)
```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/bot-linux cmd/bot/main.go

# Transfer to server
scp bin/bot-linux user@server:/opt/shopee-bot/

# Run as service (systemd)
sudo systemctl start shopee-bot
```

### Docker (Optional)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bot cmd/bot/main.go

FROM chromium:latest
COPY --from=builder /app/bot /bot
CMD ["/bot"]
```

---

## Maintenance

### Regular Tasks
- Update Go dependencies monthly
- Monitor proxy health daily
- Check Shopee's UI changes weekly
- Review logs for detection patterns
- Backup session cookies

### When Things Break
1. Check logs: `tail -f data/logs/app.log`
2. Verify proxy connectivity
3. Test authentication manually
4. Update stealth measures
5. Check Shopee's TOS updates

---

## Next Steps

1. ✅ Review project structure
2. ⏭️ Implement core modules (next phase)
3. ⏭️ Write unit tests
4. ⏭️ Test on real livestreams
5. ⏭️ Optimize performance
6. ⏭️ Deploy to production

---

**Status:** Project Structure Complete  
**Ready for:** Module Implementation