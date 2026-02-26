# Agent Guide for markdown-tool

## Purpose
This repository is a Go CLI that converts raw text (stdin or clipboard) into markdown links using a Parse → Vote → Write pipeline. Parsers detect content and populate a context, writers vote, and the highest-confidence writer emits markdown.

## Key Architecture
- **Parsing Phase:** `internal/parser` produces `*types.ParseContext` objects.
- **Voting Phase:** `internal/writer` votes by content type and confidence.
- **Output Phase:** best writer formats markdown and writes to stdout (stderr for logs).
- **Config:** YAML via Viper in `~/.config/markdown-tool/config.yaml`.

## Build / Lint / Test

### Build
- `make build`
- `go build -o markdown-tool`
- `make build-all` (multi-platform; writes to `bin/`)

### Run
- `make run`
- `make debug` (build + run `--verbose`)
- `make run-samples` (sample inputs)

### Tests
- Full suite: `go test ./...`
- Verbose: `go test -v ./...`
- Coverage: `make test-coverage`
- Watch (requires `entr`): `make test-watch`

### Run a single test
- By name: `go test -run TestEndToEndTransformation ./...`
- In a package: `go test -run TestGitHubLongParser_Parse ./internal/parser`
- File-specific by package: `go test -run TestURLWriter_WriteGitHubLongURL ./internal/writer`

### Lint / Format / Vet
- Format: `go fmt ./...` or `make fmt`
- Vet: `go vet ./...` or `make vet`
- Lint: `golangci-lint run` or `make lint`
- All checks: `make check`

## Code Style & Conventions

### Imports
- Use Go’s standard grouping: stdlib, blank line, third-party, blank line, local.
- No unused imports; `go fmt` should keep ordering canonical.

### Formatting
- Always run `go fmt` on modified packages.
- Keep lines readable; prefer clear naming over cleverness.

### Naming
- Packages are short, lowercase (e.g., `parser`, `writer`).
- Types are PascalCase; variables/functions are camelCase.
- Test names are explicit and descriptive; use table-driven tests.

### Types & Data Flow
- Prefer explicit `types.ParseContext` metadata keys (strings) with documented usage.
- Parsers should return `nil, nil` when input doesn’t match.
- Writers should return original input when they cannot safely format.

### Error Handling
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- Avoid panics in core parsing/writing logic; return errors instead.
- In `cmd/root.go`, errors are surfaced via `log.Fatalf` in the CLI entry point.

### Parser Guidelines
- Keep `CanHandle` conservative; avoid false positives.
- Prefer regex helpers for clear intent; avoid overly permissive patterns.
- For multi-line GitHub UI content, keep org/repo detection stable.
- If a change affects parsing, add or update tests in `internal/parser/*_test.go`.

### Writer Guidelines
- Writer output must be stable and deterministic.
- Do not mutate input context outside of safe transformations.
- Use config mappings for GitHub org/repo and URL domains.
- Add tests in `internal/writer/*_test.go` for new output logic.

### Testing Strategy
- Favor unit tests in `internal/parser` and `internal/writer`.
- Use table-driven tests for multiple inputs.
- Avoid bash piping as test validation; use Go tests instead.
- Integration coverage lives in `integration_test.go`.

### Configuration
- YAML config uses keys `github`, `jira`, `url`, `jenkins`.
- Only configured JIRA project keys are transformed.
- GitHub mappings are case-insensitive due to Viper lowercasing.

## Repository-Specific Rules
- Use `/opt/homebrew/bin/git` for git commands in this environment.
- Structured logging goes to stderr, output goes to stdout.
- Unmatched input should pass through verbatim.

## Cursor / Copilot Rules
- No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` found in this repo.

## Helpful Files
- `CLAUDE.md` for project conventions and testing guidance.
- `Makefile` for build/test workflows.
- `README.md` for usage and examples.
