.PHONY: help build run test clean install dev lint

# Variables
BINARY_NAME=bot
MAIN_PATH=cmd/bot/main.go
BUILD_DIR=bin
GO=go

# Colors for terminal output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Show this help message
	@echo '$(YELLOW)Available commands:$(NC)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo '$(YELLOW)Installing dependencies...$(NC)'
	$(GO) mod download
	$(GO) mod tidy
	@echo '$(GREEN)Dependencies installed!$(NC)'

build: ## Build the application
	@echo '$(YELLOW)Building application...$(NC)'
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo '$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)'

build-linux: ## Build for Linux
	@echo '$(YELLOW)Building for Linux...$(NC)'
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo '$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux$(NC)'

build-windows: ## Build for Windows
	@echo '$(YELLOW)Building for Windows...$(NC)'
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo '$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME).exe$(NC)'

build-all: build build-linux build-windows ## Build for all platforms
	@echo '$(GREEN)All builds complete!$(NC)'

run: ## Run the application
	@echo '$(YELLOW)Running application...$(NC)'
	$(GO) run $(MAIN_PATH)

dev: ## Run in development mode with live reload
	@echo '$(YELLOW)Running in development mode...$(NC)'
	@echo '$(RED)Note: Install air for live reload: go install github.com/cosmtrek/air@latest$(NC)'
	air

test: ## Run tests
	@echo '$(YELLOW)Running tests...$(NC)'
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo '$(YELLOW)Running tests with coverage...$(NC)'
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo '$(GREEN)Coverage report: coverage.html$(NC)'

lint: ## Run linter
	@echo '$(YELLOW)Running linter...$(NC)'
	@echo '$(RED)Note: Install golangci-lint first: https://golangci-lint.run/usage/install/$(NC)'
	golangci-lint run

fmt: ## Format code
	@echo '$(YELLOW)Formatting code...$(NC)'
	$(GO) fmt ./...
	@echo '$(GREEN)Code formatted!$(NC)'

vet: ## Run go vet
	@echo '$(YELLOW)Running go vet...$(NC)'
	$(GO) vet ./...

clean: ## Clean build artifacts
	@echo '$(YELLOW)Cleaning build artifacts...$(NC)'
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo '$(GREEN)Clean complete!$(NC)'

setup: ## Initial setup (copy .env.example to .env)
	@echo '$(YELLOW)Setting up project...$(NC)'
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo '$(GREEN).env file created. Please edit it with your credentials.$(NC)'; \
	else \
		echo '$(YELLOW).env file already exists.$(NC)'; \
	fi
	@mkdir -p data/cookies data/logs data/browser
	@echo '$(GREEN)Setup complete!$(NC)'
	@echo '$(YELLOW)Next steps:$(NC)'
	@echo '  1. Edit .env with your Shopee credentials'
	@echo '  2. Edit configs/config.yaml with livestream URLs'
	@echo '  3. Run: make install'
	@echo '  4. Run: make run'

docker-build: ## Build Docker image
	@echo '$(YELLOW)Building Docker image...$(NC)'
	docker build -t shopee-bot:latest .
	@echo '$(GREEN)Docker image built!$(NC)'

docker-run: ## Run Docker container
	@echo '$(YELLOW)Running Docker container...$(NC)'
	docker run --rm -it -v $(PWD)/data:/app/data shopee-bot:latest

deps-update: ## Update dependencies
	@echo '$(YELLOW)Updating dependencies...$(NC)'
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo '$(GREEN)Dependencies updated!$(NC)'

check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo '$(GREEN)All checks passed!$(NC)'