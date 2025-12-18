# 🚀 Go Commands Quick Reference

## Essential Go Commands for Testing main.go

### 1. Run Without Building
```bash
# Basic run
go run cmd/bot/main.go

# Run with specific config
go run cmd/bot/main.go --config=custom-config.yaml

# Run with environment variables
SHOPEE_USERNAME=test@example.com go run cmd/bot/main.go

# Run with verbose output
go run -v cmd/bot/main.go
```

### 2. Build the Binary
```bash
# Build for current platform
go build -o bin/bot cmd/bot/main.go

# Build with optimizations
go build -ldflags="-s -w" -o bin/bot cmd/bot/main.go

# Build for Linux (from any OS)
GOOS=linux GOARCH=amd64 go build -o bin/bot-linux cmd/bot/main.go

# Build for Windows (from any OS)
GOOS=windows GOARCH=amd64 go build -o bin/bot.exe cmd/bot/main.go

# Build for macOS (from any OS)
GOOS=darwin GOARCH=amd64 go build -o bin/bot-mac cmd/bot/main.go
```

### 3. Run the Built Binary
```bash
# After building, run it directly
./bin/bot

# Run with specific config
./bin/bot --config=configs/config.yaml

# Run in background
./bin/bot &

# Run and save output to log
./bin/bot 2>&1 | tee output.log
```

### 4. Dependency Management
```bash
# Download dependencies
go mod download

# Add missing dependencies
go get

# Remove unused dependencies
go mod tidy

# Verify dependencies
go mod verify

# View dependency graph
go mod graph

# Update all dependencies
go get -u ./...

# Update specific package
go get -u github.com/chromedp/chromedp
```

### 5. Testing
```bash
# Run all tests
go test ./...

# Run tests in specific package
go test ./internal/auth/

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestLogin ./internal/auth/

# Run tests with timeout
go test -timeout 30s ./...
```

### 6. Code Quality
```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Find race conditions
go test -race ./...

# Static analysis (requires golangci-lint)
golangci-lint run

# Check for unused code
go run golang.org/x/tools/cmd/deadcode@latest ./...
```

### 7. Debugging
```bash
# Print build info
go version -m bin/bot

# Show module info
go list -m all

# Show package info
go list ./...

# Trace execution
go run -x cmd/bot/main.go

# Enable garbage collector debug
GODEBUG=gctrace=1 go run cmd/bot/main.go

# Memory profiling
go run cmd/bot/main.go -memprofile=mem.prof
go tool pprof mem.prof
```

### 8. Cleaning
```bash
# Clean build cache
go clean -cache

# Clean test cache
go clean -testcache

# Clean module cache
go clean -modcache

# Remove build artifacts
rm -rf bin/
```

## Step-by-Step: Running main.go for Testing

### Method 1: Quick Test (No Build)
```bash
# 1. Navigate to project directory
cd shopee-livestream-bot

# 2. Ensure dependencies are installed
go mod download

# 3. Run directly (fastest for testing)
go run cmd/bot/main.go
```

**When to use:** Development, quick testing, debugging

### Method 2: Build and Run (Production)
```bash
# 1. Build the binary
go build -o bin/bot cmd/bot/main.go

# 2. Verify it was created
ls -lh bin/bot

# 3. Make it executable (Linux/Mac)
chmod +x bin/bot

# 4. Run the binary
./bin/bot
```

**When to use:** Production, deployment, performance testing

### Method 3: Using Makefile (Recommended)
```bash
# 1. Run directly (development)
make run

# 2. Build binary
make build

# 3. Build for all platforms
make build-all

# 4. Run tests
make test

# 5. Clean and rebuild
make clean build
```

**When to use:** Consistent builds, team development

## Common Testing Workflows

### Workflow 1: First-Time Setup and Test
```bash
# 1. Clone and navigate
git clone <repo>
cd shopee-livestream-bot

# 2. Setup project
make setup

# 3. Configure credentials
nano .env

# 4. Install dependencies
make install

# 5. Test run
make run
```

### Workflow 2: Quick Test After Code Changes
```bash
# 1. Save your changes
# 2. Format code
go fmt ./...

# 3. Run immediately
go run cmd/bot/main.go
```

### Workflow 3: Full Test Before Deployment
```bash
# 1. Format code
go fmt ./...

# 2. Check for issues
go vet ./...

# 3. Run tests
go test ./...

# 4. Build
go build -o bin/bot cmd/bot/main.go

# 5. Test the binary
./bin/bot
```

### Workflow 4: Debug Mode
```bash
# 1. Enable debug logging
export LOG_LEVEL=debug

# 2. Run with output
go run cmd/bot/main.go 2>&1 | tee debug.log

# 3. Review logs
tail -f debug.log
```

### Workflow 5: Performance Testing
```bash
# 1. Build optimized binary
go build -ldflags="-s -w" -o bin/bot cmd/bot/main.go

# 2. Time execution
time ./bin/bot

# 3. Monitor resources
top -p $(pgrep bot)
```

## Troubleshooting Go Commands

### Error: "go: cannot find main module"
```bash
# Solution: Make sure you're in the project directory
cd shopee-livestream-bot
pwd  # Should show project path
```

### Error: "cannot find package"
```bash
# Solution: Download dependencies
go mod download
go mod tidy
```

### Error: "build failed"
```bash
# Solution: Check for syntax errors
go build -v cmd/bot/main.go

# If that doesn't help, clean and rebuild
go clean -cache
go build cmd/bot/main.go
```

### Error: "permission denied"
```bash
# Solution: Make binary executable
chmod +x bin/bot
```

### Error: "port already in use"
```bash
# Solution: Kill existing process
pkill -f bot
# Or find and kill manually
ps aux | grep bot
kill <PID>
```

## Environment Variables

### Set for Single Run
```bash
# Linux/Mac
SHOPEE_USERNAME=test@example.com go run cmd/bot/main.go

# Windows CMD
set SHOPEE_USERNAME=test@example.com
go run cmd/bot/main.go

# Windows PowerShell
$env:SHOPEE_USERNAME="test@example.com"
go run cmd/bot/main.go
```

### Set Permanently
```bash
# Linux/Mac
export SHOPEE_USERNAME=test@example.com
echo 'export SHOPEE_USERNAME=test@example.com' >> ~/.bashrc

# Or use .env file (recommended)
nano .env
```

## Quick Reference Card

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `go run cmd/bot/main.go` | Run without building | Quick testing, development |
| `go build cmd/bot/main.go` | Build binary | Production, deployment |
| `./bin/bot` | Run built binary | After building |
| `go test ./...` | Run all tests | Before commit |
| `go mod download` | Get dependencies | Initial setup |
| `go mod tidy` | Clean dependencies | After adding/removing code |
| `go fmt ./...` | Format code | Before commit |
| `go vet ./...` | Check for issues | Before commit |
| `make run` | Run via Makefile | Consistent development |
| `make build` | Build via Makefile | Consistent builds |

## Pro Tips

### 1. Faster Builds
```bash
# Use cache
go build -i cmd/bot/main.go

# Parallel compilation
go build -p 8 cmd/bot/main.go
```

### 2. Watch for Changes (requires external tool)
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air
```

### 3. Cross-Compilation Matrix
```bash
# Build for multiple platforms at once
GOOS=linux GOARCH=amd64 go build -o bin/bot-linux-amd64 cmd/bot/main.go &
GOOS=linux GOARCH=arm64 go build -o bin/bot-linux-arm64 cmd/bot/main.go &
GOOS=windows GOARCH=amd64 go build -o bin/bot-windows.exe cmd/bot/main.go &
GOOS=darwin GOARCH=amd64 go build -o bin/bot-mac cmd/bot/main.go &
wait
```

### 4. Build Info
```bash
# Embed version info
go build -ldflags="-X main.version=1.0.0 -X main.commit=$(git rev-parse HEAD)" cmd/bot/main.go
```

## Common Go Environment Variables

```bash
# Go root directory
GOROOT=/usr/local/go

# Go workspace
GOPATH=$HOME/go

# Go binary path
PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# Enable Go modules (default in Go 1.16+)
GO111MODULE=on

# Module proxy (for faster downloads)
GOPROXY=https://proxy.golang.org,direct

# Private modules
GOPRIVATE=github.com/yourcompany/*
```

---

## Your First Test Run

Complete these steps in order:

```bash
# 1. Navigate to project
cd shopee-livestream-bot

# 2. Setup (if not done)
make setup

# 3. Edit credentials
nano .env

# 4. Install dependencies
go mod download

# 5. Run the bot!
go run cmd/bot/main.go
```

Watch the output for:
- ✅ Configuration loaded
- ✅ Browser initialized
- ✅ Authentication successful
- ✅ Monitoring started

Press `Ctrl+C` to stop.

---

**You're ready to test! 🚀**

Start with `go run cmd/bot/main.go` and go from there!