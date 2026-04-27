# justfile for markdown-tool

binary_name := "markdown-tool"
main_package := "."
build_dir    := "bin"
version      := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
ldflags      := '-ldflags "-X main.Version=' + version + '"'

# Build and test (default)
[private]
default: test build

# Build the application
[group('build')]
build:
    @echo "Building {{binary_name}}..."
    go build {{ldflags}} -o {{binary_name}} {{main_package}}

# Build for multiple platforms
[group('build')]
build-all: clean
    @echo "Building for multiple platforms..."
    mkdir -p {{build_dir}}
    GOOS=linux  GOARCH=amd64 go build {{ldflags}} -o {{build_dir}}/{{binary_name}}-linux-amd64  {{main_package}}
    GOOS=linux  GOARCH=arm64 go build {{ldflags}} -o {{build_dir}}/{{binary_name}}-linux-arm64  {{main_package}}
    # GOOS=darwin GOARCH=amd64 go build {{ldflags}} -o {{build_dir}}/{{binary_name}}-darwin-amd64 {{main_package}}
    GOOS=darwin GOARCH=arm64 go build {{ldflags}} -o {{build_dir}}/{{binary_name}}-darwin-arm64 {{main_package}}
    @echo "Binaries built in {{build_dir}}/"

# Install binary to GOPATH/bin
[group('build')]
install:
    @echo "Installing {{binary_name}}..."
    go install {{ldflags}} {{main_package}}

# Create a release build (override version with: just version=v1.0.0 release)
[group('build')]
release: check build-all
    @echo "Release {{version}} ready in {{build_dir}}/"

# Run tests
[group('test')]
test:
    go test ./...

# Run tests with coverage report
[group('test')]
test-coverage:
    @echo "Running tests with coverage..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Run tests in watch mode (requires entr: brew install entr)
[group('test')]
test-watch:
    @echo "Running tests in watch mode..."
    find . -name "*.go" | entr -c go test -v ./...

# Run the application directly with go run
[group('run')]
run:
    #!/usr/bin/env zsh
    echo "\n\n👇👇👇Clipboard 👇👇👇"
    pbpaste
    echo "\n👆👆👆Clipboard 👆👆👆\n"
    echo "👇👇👇Output 👇👇👇"
    go run {{main_package}}
    echo "\n👆👆👆Output 👆👆👆\n"

# Run the application in debug mode
[group('run')]
debug: build
    @echo "Running {{binary_name}} in debug mode..."
    ./{{binary_name}} --verbose

# Test application with sample inputs
[group('run')]
run-samples: build
    @echo "Testing GitHub URL transformation:"
    @echo "https://github.com/CompanyCam/Company-Cam-API/pull/15217" | ./{{binary_name}}
    @echo ""
    @echo "Testing JIRA key transformation:"
    @echo "PLAT-12345" | ./{{binary_name}}
    @echo ""
    @echo "Testing JIRA URL transformation:"
    @echo "https://companycam.atlassian.net/browse/PLAT-192" | ./{{binary_name}}
    @echo ""
    @echo "Testing Notion URL transformation:"
    @echo "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0" | ./{{binary_name}}
    @echo ""
    @echo "Testing generic URL transformation:"
    @echo "http://ww3.domain.tld/path/to/document?query=value#anchor" | ./{{binary_name}}

# Format code
[group('quality')]
fmt:
    @echo "Formatting code..."
    go fmt ./...

# Vet code
[group('quality')]
vet:
    @echo "Vetting code..."
    go vet ./...

# Run linter (requires golangci-lint)
[group('quality')]
lint:
    @echo "Running linter..."
    golangci-lint run

# Run all quality checks (fmt, vet, lint, test)
[group('quality')]
check: fmt vet lint test

# Clean build artifacts
[group('maintenance')]
clean:
    @echo "Cleaning build artifacts..."
    rm -f {{binary_name}}
    rm -rf {{build_dir}}
    rm -f coverage.out coverage.html

# Install dependencies
[group('maintenance')]
deps:
    @echo "Installing dependencies..."
    go mod download
    go mod tidy

# Setup development environment
[group('maintenance')]
dev-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Setting up development environment..."
    if ! command -v golangci-lint >/dev/null 2>&1; then
        echo "Installing golangci-lint..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    if ! command -v entr >/dev/null 2>&1; then
        echo "entr not found. Install with: brew install entr (for test-watch)"
    fi
    echo "Development setup complete!"

# Show available recipes
[private]
help:
    @just --list
