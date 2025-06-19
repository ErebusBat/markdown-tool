# Implementation Plan: markdown-tool

## Phase 1: Project Setup and Foundation

### 1.1 Go Module and Dependencies
- [ ] Initialize Go module: `go mod init github.com/erebusbat/markdown-tool`
- [ ] Add core dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/spf13/viper` - Configuration management
  - Logging library (e.g., `github.com/sirupsen/logrus` or `slog`)
  - Clipboard library (e.g., `github.com/atotto/clipboard`)

### 1.2 Project Structure
```
markdown-tool/
├── cmd/
│   └── root.go                 # Cobra root command
├── internal/
│   ├── config/
│   │   ├── config.go          # Configuration loading/validation
│   │   └── defaults.go        # Default configuration values
│   ├── parser/
│   │   ├── parser.go          # Parser interface
│   │   ├── url_parser.go      # URL detection and parsing
│   │   ├── jira_parser.go     # JIRA key detection
│   │   └── context.go         # Shared parsing context
│   ├── writer/
│   │   ├── writer.go          # Writer interface and voting
│   │   ├── url_writer.go      # URL transformation
│   │   ├── jira_writer.go     # JIRA transformation
│   │   └── passthrough_writer.go # Verbatim output
│   └── input/
│       ├── input.go           # Input handling (stdin/clipboard)
│       └── clipboard.go       # Clipboard integration
├── pkg/
│   └── types/
│       └── types.go           # Shared types and interfaces
├── testdata/                  # Test fixtures
├── main.go                    # Application entry point
├── go.mod
├── go.sum
└── README.md
```

### 1.3 Core Interfaces
- [ ] Define `Parser` interface for content detection
- [ ] Define `Writer` interface for output generation
- [ ] Define `Context` struct for sharing parsed data
- [ ] Define configuration structs

## Phase 2: Configuration System

### 2.1 Configuration Management
- [ ] Implement Viper-based config loading from `~/.config/markdown-tool/config.yaml`
- [ ] Create default configuration with CompanyCam values
- [ ] Implement config directory creation
- [ ] Add config validation
- [ ] Support for GitHub org/repo remapping

### 2.2 Default Configuration
```yaml
github:
  default_org: "CompanyCam"
  default_repo: "Company-Cam-API"
  mappings:
    "CompanyCam/Company-Cam-API": "CompanyCam/API"

jira:
  domain: "https://companycam.atlassian.net"
  projects:
    - "PLAT"
    - "SPEED"
```

## Phase 3: Input Handling

### 3.1 Input Sources
- [ ] Implement stdin reader
- [ ] Implement clipboard fallback when stdin is empty
- [ ] Add input validation and sanitization
- [ ] Handle empty input gracefully

## Phase 4: Parsing Engine

### 4.1 URL Parsers
- [ ] Generic URL parser (extract domain from any URL)
- [ ] GitHub URL parser (extract org/repo/issue from GitHub URLs)
- [ ] JIRA URL parser (extract issue keys from JIRA URLs)
- [ ] JIRA comment URL parser (detect focusedCommentId parameter)
- [ ] Notion URL parser (extract page titles from URL slugs)

### 4.2 Text Parsers
- [ ] JIRA key parser (detect standalone JIRA keys like "PLAT-12345")
- [ ] Validate JIRA keys against configured projects

### 4.3 Context Management
- [ ] Implement shared context for storing parsed data
- [ ] Add confidence scoring for each parser
- [ ] Store original input and parsed variants

## Phase 5: Output Writers

### 5.1 Writer Implementation
- [ ] URL writer (transform URLs to markdown links)
- [ ] JIRA writer (transform JIRA keys to markdown links)
- [ ] Passthrough writer (output verbatim for unmatched content)

### 5.2 Voting System
- [ ] Implement confidence-based voting
- [ ] Select highest-confidence writer
- [ ] Handle tie-breaking scenarios
- [ ] Fallback to passthrough for zero confidence

## Phase 6: CLI Interface

### 6.1 Cobra Integration
- [ ] Implement root command with proper help
- [ ] Add version command
- [ ] Add config validation command
- [ ] Support for debug/verbose logging flags

### 6.2 Logging
- [ ] Implement structured logging to stderr
- [ ] Add colorful, human-readable output format
- [ ] Support log levels (debug, info, warn, error)
- [ ] Log parsing decisions and confidence scores

## Phase 7: Testing

### 7.1 Unit Tests
- [ ] Parser tests for each content type
- [ ] Writer tests for each transformation
- [ ] Configuration loading tests
- [ ] Input handling tests
- [ ] Voting system tests

### 7.2 Integration Tests
- [ ] End-to-end processing tests
- [ ] Configuration file creation tests
- [ ] Error handling tests
- [ ] Edge case tests (empty input, malformed URLs, etc.)

### 7.3 Test Data
- [ ] Create comprehensive test fixtures
- [ ] Include real-world examples from PRD
- [ ] Add edge cases and error scenarios

## Phase 8: Build and Deployment

### 8.1 Build System
- [ ] Add Makefile or build scripts
- [ ] Configure cross-compilation for Linux/macOS
- [ ] Add version information to binary
- [ ] Optimize binary size

### 8.2 Documentation
- [ ] Update README with usage examples
- [ ] Add configuration documentation
- [ ] Include troubleshooting guide
- [ ] Document extension points for future development

## Implementation Order

1. **Foundation First**: Project setup, dependencies, and basic structure
2. **Configuration**: Get config loading working with defaults
3. **Input/Output**: Basic stdin/stdout flow
4. **Core Engine**: Three-phase architecture (parse → vote → output)
5. **Content Processors**: URL and JIRA transformations
6. **CLI Polish**: Help, logging, error handling
7. **Testing**: Comprehensive test coverage
8. **Build/Deploy**: Final packaging and documentation

## Success Criteria

Each phase should be considered complete when:
- All acceptance criteria from PRD are met
- Unit tests pass with good coverage
- Integration tests validate end-to-end functionality
- Code follows Go best practices and is well-documented
- CLI provides helpful error messages and usage information
