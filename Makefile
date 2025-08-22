# Auto-Spotify Makefile

.PHONY: build clean install test fmt vet lint run help

# Default target
.DEFAULT_GOAL := help

# Build the application
build: ## Build the auto-spotify binary
	@echo "Building auto-spotify..."
	go build -ldflags="-s -w" -o auto-spotify .

# Clean build artifacts
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	rm -f auto-spotify auto-spotify.exe
	go clean

# Install dependencies
install: ## Install/update dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Run a quick test build
test: fmt vet unit-test ## Run formatting, vetting, unit tests, and test build
	@echo "Running test build..."
	go build -o /tmp/auto-spotify-test .
	rm -f /tmp/auto-spotify-test
	@echo "All checks passed!"

# Run unit tests
unit-test: ## Run unit tests
	@echo "Running unit tests..."
	go test ./... -v

# Run unit tests with coverage
test-coverage: ## Run unit tests with coverage report
	@echo "Running tests with coverage..."
	go test ./... -cover

# Run benchmark tests
benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	go test ./... -bench=. -benchmem

# Run the application with a test prompt
run: build ## Build and run with a test prompt
	@echo "Running auto-spotify with test prompt..."
	./auto-spotify "chill indie rock for coding"

# Run with custom prompt (usage: make run-prompt PROMPT="your prompt here")
run-prompt: build ## Build and run with custom prompt
	@echo "Running auto-spotify with prompt: $(PROMPT)"
	./auto-spotify "$(PROMPT)"

# Development setup
setup: ## Set up development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Created .env file from env.example"; \
		echo "Please edit .env with your API keys"; \
		echo ""; \
		echo "Note: Default configuration uses HTTPS (https://localhost:8080/callback)"; \
		echo "Make sure to configure this URL in your Spotify app settings."; \
	else \
		echo ".env file already exists"; \
	fi
	$(MAKE) install

# Release build (cross-platform)
release: clean ## Build release binaries for multiple platforms
	@echo "Building release binaries..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/auto-spotify-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/auto-spotify-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/auto-spotify-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/auto-spotify-windows-amd64.exe .
	@echo "Release binaries created in dist/"

# Help target
help: ## Show this help message
	@echo "Auto-Spotify Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make setup              # Set up development environment"
	@echo "  make build              # Build the application"
	@echo "  make run                # Build and run with test prompt"
	@echo "  make run-prompt PROMPT=\"jazz for studying\"  # Run with custom prompt"
	@echo "  make test               # Run all checks"
	@echo "  make release            # Build release binaries"
