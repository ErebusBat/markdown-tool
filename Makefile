# Makefile for markdown-tool

# Variables
BINARY_NAME=markdown-tool
MAIN_PACKAGE=.
BUILD_DIR=bin
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: test build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	@echo "Binaries built in $(BUILD_DIR)/"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test  ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in watch mode (requires entr: brew install entr)
.PHONY: test-watch
test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" | entr -c go test -v ./...

# Run the application in debug mode
.PHONY: debug
debug: build
	@echo "Running $(BINARY_NAME) in debug mode..."
	./$(BINARY_NAME) --verbose

# Run with sample inputs for testing
.PHONY: run-samples
run-samples: build
	@echo "Testing GitHub URL transformation:"
	@echo "https://github.com/CompanyCam/Company-Cam-API/pull/15217" | ./$(BINARY_NAME)
	@echo ""
	@echo "Testing JIRA key transformation:"
	@echo "PLAT-12345" | ./$(BINARY_NAME)
	@echo ""
	@echo "Testing JIRA URL transformation:"
	@echo "https://companycam.atlassian.net/browse/PLAT-192" | ./$(BINARY_NAME)
	@echo ""
	@echo "Testing Notion URL transformation:"
	@echo "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0" | ./$(BINARY_NAME)
	@echo ""
	@echo "Testing generic URL transformation:"
	@echo "http://ww3.domain.tld/path/to/document?query=value#anchor" | ./$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...

# Run all quality checks
.PHONY: check
check: fmt vet lint test

# Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PACKAGE)

# Create a new release (requires VERSION)
.PHONY: release
release: check build-all
	@echo "Creating release $(VERSION)..."
	@if [ -z "$(VERSION)" ]; then echo "VERSION is required"; exit 1; fi
	@echo "Release $(VERSION) ready in $(BUILD_DIR)/"

# Development setup
.PHONY: dev-setup
dev-setup:
	@echo "Setting up development environment..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
	@command -v entr >/dev/null 2>&1 || { \
		echo "entr not found. Install with: brew install entr (for test-watch)"; \
	}
	@echo "Development setup complete!"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          Build the application"
	@echo "  build-all      Build for multiple platforms"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  test-watch     Run tests in watch mode (requires entr)"
	@echo "  debug          Run application in debug mode"
	@echo "  run-samples    Test application with sample inputs"
	@echo "  clean          Clean build artifacts"
	@echo "  deps           Install dependencies"
	@echo "  fmt            Format code"
	@echo "  lint           Run linter (requires golangci-lint)"
	@echo "  vet            Vet code"
	@echo "  check          Run all quality checks"
	@echo "  install        Install binary to GOPATH/bin"
	@echo "  release        Create release build"
	@echo "  dev-setup      Setup development environment"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    # Build the application"
	@echo "  make test                     # Run tests"
	@echo "  make run-samples              # Test with sample inputs"
	@echo "  make VERSION=v1.0.0 release   # Create release"
