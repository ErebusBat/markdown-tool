# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based command-line tool for transforming text inputs into well-formatted markdown. The tool is designed to detect and transform URLs (GitHub, JIRA, Notion, generic), JIRA issue keys, and other content types into appropriate markdown links.

## Architecture

The application follows a three-phase processing architecture:

1. **Parsing Phase**: All parsers analyze input and populate a shared context object
2. **Voting Phase**: Output writers vote on their confidence to handle the parsed content
3. **Output Phase**: The highest-confidence writer generates the final transformed output

## Key Components

- **Configuration**: Uses spf13/viper for YAML config management stored in `~/.config/markdown-tool/`
- **CLI Interface**: Built with spf13/cobra for command-line interaction
- **Processors**: Modular design with separate processors for different content types (URLs, JIRA keys, etc.)
- **Input Sources**: Supports both stdin and clipboard fallback
- **Output**: Writes to stdout with structured logging to stderr

## Development Commands

Standard Go commands for development:

```bash
# Initialize Go module
go mod init github.com/erebusbat/markdown-tool

# Build the application
go build -o markdown-tool

# Run tests (preferred method for feature validation)
go test ./...

# Run specific test
go test -run TestSpecificFunction ./pkg/package

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Install dependencies
go mod tidy

# Format code
go fmt ./...

# Lint (if golangci-lint is available)
golangci-lint run
```

## Development Workflow

When adding new features or fixing bugs:

1. **Write tests first**: Create comprehensive test cases in appropriate `*_test.go` files
2. **Implement the feature**: Write the minimum code to make tests pass
3. **Validate with tests**: Use `go test ./...` instead of manual bash piping
4. **Add integration tests**: Ensure end-to-end functionality works
5. **Run full test suite**: Verify no regressions were introduced

**Avoid using bash piping for feature validation** - instead, add proper test cases that can be run repeatedly and automatically.

**IMPORTANT Git Configuration:**
When running git commands via Claude Code tools, always use the full path: `/opt/homebrew/bin/git` instead of just `git`. This is required for this specific development environment.

## Configuration Structure

The tool expects YAML configuration with:
- GitHub settings (org, repo, description mappings)
- JIRA settings (domain, valid project keys)
- Default values: CompanyCam org, Company-Cam-API repo, companycam.atlassian.net domain, PLAT/SPEED projects

## Testing Strategy

**IMPORTANT: Use Go test files, not bash piping for feature validation**

All feature testing and validation should be done through proper Go test files (`*_test.go`) rather than bash commands with piping. This ensures:
- Reliable, repeatable tests
- Proper test isolation and setup
- Better error reporting and debugging
- Integration with Go's testing framework
- Automated test execution in CI/CD

### Test Structure

Tests should follow this hierarchy:
1. **Unit Tests**: Test individual components in isolation
   - Parser tests: `internal/parser/*_test.go`
   - Writer tests: `internal/writer/*_test.go`
   - Config tests: `internal/config/*_test.go`

2. **Integration Tests**: Test complete input-to-output pipeline
   - End-to-end tests: `integration_test.go`
   - Multi-component interaction tests

3. **Test Coverage Areas**:
   - URL transformation scenarios for each supported platform
   - JIRA key detection and transformation
   - Configuration loading and validation
   - Input/output handling (stdin, clipboard, stdout)
   - Error handling for edge cases
   - Three-phase processing architecture
   - Regression prevention for bug fixes

### Testing Best Practices

- **Add tests for new features**: Every new feature should include comprehensive test coverage
- **Test both positive and negative cases**: Include valid inputs and edge cases
- **Use table-driven tests**: For multiple similar test scenarios
- **Mock external dependencies**: Use test doubles for configuration, file system, etc.
- **Test preprocessing logic**: Include tests for any input preprocessing (e.g., tel: URI handling)
- **Avoid bash piping in tests**: Use Go's testing framework instead of shell commands

## Implementation Notes

- Unmatched input should be output verbatim
- Only JIRA keys matching configured projects should be transformed
- GitHub org/repo names can be remapped via configuration
- Structured logging uses stderr for debug output
- Application is designed as a short-lived process
