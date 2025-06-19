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

Since this is a new project, standard Go commands will be used:

```bash
# Initialize Go module
go mod init github.com/erebusbat/markdown-tool

# Build the application
go build -o markdown-tool

# Run tests
go test ./...

# Run specific test
go test -run TestSpecificFunction ./pkg/package

# Install dependencies
go mod tidy

# Format code
go fmt ./...

# Lint (if golangci-lint is available)
golangci-lint run
```

**IMPORTANT Git Configuration:**
When running git commands via Claude Code tools, always use the full path: `/opt/homebrew/bin/git` instead of just `git`. This is required for this specific development environment.

## Configuration Structure

The tool expects YAML configuration with:
- GitHub settings (org, repo, description mappings)
- JIRA settings (domain, valid project keys)
- Default values: CompanyCam org, Company-Cam-API repo, companycam.atlassian.net domain, PLAT/SPEED projects

## Testing Strategy

Tests should cover:
- URL transformation scenarios for each supported platform
- JIRA key detection and transformation
- Configuration loading and validation
- Input/output handling (stdin, clipboard, stdout)
- Error handling for edge cases
- Three-phase processing architecture

## Implementation Notes

- Unmatched input should be output verbatim
- Only JIRA keys matching configured projects should be transformed
- GitHub org/repo names can be remapped via configuration
- Structured logging uses stderr for debug output
- Application is designed as a short-lived process
