# Implementation Plan: markdown-tool

## Status: Phase 1-5 Complete ✅

**Last Updated:** June 18, 2025

### Completed Features:
- ✅ **Core Application**: Fully functional markdown transformation tool
- ✅ **Configuration System**: Auto-creates `~/.config/markdown-tool/config.yaml` with CompanyCam defaults
- ✅ **URL Processing**: GitHub, JIRA, JIRA comments, Notion, and generic URLs
- ✅ **JIRA Key Detection**: Standalone keys (e.g., `PLAT-12345`) for configured projects
- ✅ **GitHub Mappings**: Org/repo name remapping (case-insensitive)
- ✅ **Input/Output**: stdin with clipboard fallback, stdout output
- ✅ **Three-Phase Architecture**: Parse → Vote → Output with confidence scoring
- ✅ **CLI Interface**: Cobra-based with help and basic flags

### Test Results:
```bash
# GitHub URL transformation
echo "https://github.com/CompanyCam/Company-Cam-API/pull/15217" | ./markdown-tool
# Output: [CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)

# JIRA key transformation
echo "PLAT-12345" | ./markdown-tool
# Output: [PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)

# Notion URL transformation
echo "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0#1c0d42d77c0b80268626fa64eb6ebdbe" | ./markdown-tool
# Output: [VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0#1c0d42d77c0b80268626fa64eb6ebdbe)
```

## Phase 1: Project Setup and Foundation

### 1.1 Go Module and Dependencies
- [x] Initialize Go module: `go mod init github.com/erebusbat/markdown-tool`
- [x] Add core dependencies:
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/spf13/viper` - Configuration management
  - `github.com/atotto/clipboard` - Clipboard library
  - *Note: Using standard library logging for now*

### 1.2 Project Structure ✅
```
markdown-tool/
├── cmd/
│   └── root.go                    # Cobra root command ✅
├── internal/
│   ├── config/
│   │   └── config.go             # Configuration loading/validation ✅
│   ├── parser/
│   │   ├── parser.go             # Parser interface ✅
│   │   ├── url_parser.go         # URL detection and parsing ✅
│   │   └── jira_parser.go        # JIRA key detection ✅
│   └── writer/
│       ├── writer.go             # Writer interface and voting ✅
│       ├── url_writer.go         # URL transformation ✅
│       ├── jira_writer.go        # JIRA transformation ✅
│       └── passthrough_writer.go # Verbatim output ✅
├── pkg/
│   └── types/
│       └── types.go              # Shared types and interfaces ✅
├── main.go                       # Application entry point ✅
├── go.mod                        # Go module file ✅
├── go.sum                        # Dependencies ✅
├── CLAUDE.md                     # Claude Code guidance ✅
├── IMPLEMENTATION_PLAN.md        # This file ✅
└── PRD.md                        # Product requirements ✅
```

**Note:** Input handling was integrated directly into `cmd/root.go` instead of separate `internal/input/` package for simplicity.

### 1.3 Core Interfaces
- [x] Define `Parser` interface for content detection
- [x] Define `Writer` interface for output generation
- [x] Define `ParseContext` struct for sharing parsed data
- [x] Define configuration structs (`Config`, `GitHubConfig`, `JIRAConfig`)

## Phase 2: Configuration System

### 2.1 Configuration Management
- [x] Implement Viper-based config loading from `~/.config/markdown-tool/config.yaml`
- [x] Create default configuration with CompanyCam values
- [x] Implement config directory creation
- [x] Add basic config validation
- [x] Support for GitHub org/repo remapping (case-insensitive)

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
- [x] Implement stdin reader
- [x] Implement clipboard fallback when stdin is empty
- [x] Add input validation and sanitization (whitespace trimming)
- [x] Handle empty input gracefully

## Phase 4: Parsing Engine

### 4.1 URL Parsers
- [x] Generic URL parser (extract domain from any URL)
- [x] GitHub URL parser (extract org/repo/issue from GitHub URLs)
- [x] JIRA URL parser (extract issue keys from JIRA URLs)
- [x] JIRA comment URL parser (detect focusedCommentId parameter)
- [x] Notion URL parser (extract page titles from URL slugs)

### 4.2 Text Parsers
- [x] JIRA key parser (detect standalone JIRA keys like "PLAT-12345")
- [x] Validate JIRA keys against configured projects

### 4.3 Context Management
- [x] Implement shared context for storing parsed data
- [x] Add confidence scoring for each parser
- [x] Store original input and parsed variants

## Phase 5: Output Writers

### 5.1 Writer Implementation
- [x] URL writer (transform URLs to markdown links)
- [x] JIRA writer (transform JIRA keys to markdown links)
- [x] Passthrough writer (output verbatim for unmatched content)

### 5.2 Voting System
- [x] Implement confidence-based voting
- [x] Select highest-confidence writer
- [x] Handle tie-breaking scenarios
- [x] Fallback to passthrough for zero confidence

## Phase 6: CLI Interface

### 6.1 Cobra Integration
- [x] Implement root command with proper help
- [ ] Add version command
- [ ] Add config validation command
- [x] Support for debug/verbose logging flags (basic implementation)

### 6.2 Logging
- [ ] Implement structured logging to stderr
- [ ] Add colorful, human-readable output format
- [ ] Support log levels (debug, info, warn, error)
- [ ] Log parsing decisions and confidence scores

*Note: Basic error logging implemented; structured logging pending*

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

## Next Steps (Remaining Work)

### Immediate Priorities:
1. **Phase 6.1**: Add version command and config validation command
2. **Phase 6.2**: Implement structured logging with color output
3. **Phase 7**: Comprehensive testing suite
4. **Phase 8**: Build system and documentation

### Recommended Development Order:
1. Add version command and improve CLI interface
2. Create unit tests for all parsers and writers
3. Add integration tests with test fixtures
4. Implement structured logging with debug output
5. Create build scripts and cross-compilation
6. Write comprehensive README and documentation

### Current State:
The application is **fully functional** and meets all core requirements from the PRD. It can be used immediately for markdown transformation tasks. The remaining work focuses on polish, testing, and operational concerns.
