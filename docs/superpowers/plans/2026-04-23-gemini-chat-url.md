# Gemini Chat URL Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add support for converting `gemini.google.com/app/<id>` URLs into `[🤖 Gemini Chat](<url>)` markdown links.

**Architecture:** Extend the existing URLParser + URLWriter mega-pattern. Add a new ContentType, parser detection method, and writer formatting method following the same pattern as MiniMax/YouTube/Notion.

**Tech Stack:** Go, table-driven tests, existing parser/writer pipeline.

---

## Chunk 1: Types and Parser

### Task 1: Add ContentType constant

**Files:**
- Modify: `pkg/types/types.go:33`

- [ ] **Step 1: Add `ContentTypeGeminiURL` constant after `ContentTypeMiniMaxURL`**

In `pkg/types/types.go`, add after line 33 (`ContentTypeMiniMaxURL`):

```go
	ContentTypeGeminiURL
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`
Expected: compiles successfully (unused constant is fine in Go)

- [ ] **Step 3: Commit**

```bash
git add pkg/types/types.go
git commit -m "feat: add ContentTypeGeminiURL constant"
```

### Task 2: Add parser detection and tests

**Files:**
- Modify: `internal/parser/url_parser.go`
- Modify: `internal/parser/url_parser_test.go`

- [ ] **Step 1: Write the failing parser test**

Add to `internal/parser/url_parser_test.go`:

```go
func TestURLParser_Parse_Gemini(t *testing.T) {
	cfg := &types.Config{}
	p := NewURLParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectedType   types.ContentType
		expectedConf   int
		expectedChatID string
	}{
		{
			name:           "Gemini chat URL",
			input:          "https://gemini.google.com/app/ac9ebc9d76c30fc1",
			expectedType:   types.ContentTypeGeminiURL,
			expectedConf:   90,
			expectedChatID: "ac9ebc9d76c30fc1",
		},
		{
			name:           "Gemini chat URL with different ID",
			input:          "https://gemini.google.com/app/abcdef123456",
			expectedType:   types.ContentTypeGeminiURL,
			expectedConf:   90,
			expectedChatID: "abcdef123456",
		},
		{
			name:         "Gemini root URL is not a chat",
			input:        "https://gemini.google.com/app",
			expectedType: types.ContentTypeURL,
			expectedConf: 50,
		},
		{
			name:         "Gemini non-app URL is generic",
			input:        "https://gemini.google.com/about",
			expectedType: types.ContentTypeURL,
			expectedConf: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if ctx == nil {
				t.Fatal("Parse() returned nil context")
			}
			if ctx.DetectedType != tt.expectedType {
				t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, tt.expectedType)
			}
			if ctx.Confidence != tt.expectedConf {
				t.Errorf("Confidence = %v, want %v", ctx.Confidence, tt.expectedConf)
			}
			if tt.expectedChatID != "" {
				if chatID := ctx.Metadata["chat_id"]; chatID != tt.expectedChatID {
					t.Errorf("Metadata[chat_id] = %v, want %v", chatID, tt.expectedChatID)
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -run TestURLParser_Parse_Gemini ./internal/parser/ -v`
Expected: FAIL — Gemini URLs detected as `ContentTypeURL` (generic), not `ContentTypeGeminiURL`

- [ ] **Step 3: Add `isGeminiURL` and `parseGeminiURL` methods to URLParser**

Add to `internal/parser/url_parser.go`, after the `isMiniMaxURL` method (around line 127):

```go
func (p *URLParser) isGeminiURL(u *url.URL) bool {
	return u.Host == "gemini.google.com" && strings.HasPrefix(u.Path, "/app/")
}

func (p *URLParser) parseGeminiURL(u *url.URL, ctx *types.ParseContext) {
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) >= 2 {
		ctx.Metadata["chat_id"] = parts[1]
	}
}
```

Add a case in the `Parse` method's switch (after the `isMiniMaxURL` case, before `default`):

```go
	case p.isGeminiURL(u):
		ctx.DetectedType = types.ContentTypeGeminiURL
		ctx.Confidence = 90
		p.parseGeminiURL(u, ctx)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -run TestURLParser_Parse_Gemini ./internal/parser/ -v`
Expected: PASS

- [ ] **Step 5: Run full parser tests**

Run: `go test ./internal/parser/ -v`
Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add internal/parser/url_parser.go internal/parser/url_parser_test.go
git commit -m "feat: add Gemini chat URL parser"
```

---

## Chunk 2: Writer

### Task 3: Add writer vote and formatting with tests

**Files:**
- Modify: `internal/writer/url_writer.go`
- Modify: `internal/writer/url_writer_test.go`

- [ ] **Step 1: Write the failing writer tests**

Add to `internal/writer/url_writer_test.go`:

First, add a Gemini entry to the `TestURLWriter_Vote` table (after the `"Generic URL"` entry):

```go
		{"Gemini URL", types.ContentTypeGeminiURL, 90},
```

Then add a new test function:

```go
func TestURLWriter_WriteGeminiURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name:          "Gemini chat URL",
			originalInput: "https://gemini.google.com/app/ac9ebc9d76c30fc1",
			metadata: map[string]interface{}{
				"chat_id": "ac9ebc9d76c30fc1",
			},
			expectedOutput: "[🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)",
		},
		{
			name:          "Gemini chat URL with different ID",
			originalInput: "https://gemini.google.com/app/abcdef123456",
			metadata: map[string]interface{}{
				"chat_id": "abcdef123456",
			},
			expectedOutput: "[🤖 Gemini Chat](https://gemini.google.com/app/abcdef123456)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeGeminiURL,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}
			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -run "TestURLWriter_Vote|TestURLWriter_WriteGeminiURL" ./internal/writer/ -v`
Expected: FAIL — vote returns 0 for ContentTypeGeminiURL, write falls through to default

- [ ] **Step 3: Add vote, write, and `writeGeminiURL` to URLWriter**

In `internal/writer/url_writer.go`:

Add to the `Vote` switch (after `case types.ContentTypeMiniMaxURL: return 90`):

```go
	case types.ContentTypeGeminiURL:
		return 90
```

Add to the `Write` switch (after `case types.ContentTypeMiniMaxURL:`):

```go
	case types.ContentTypeGeminiURL:
		return w.writeGeminiURL(ctx)
```

Add the method (after `writeMiniMaxURL`):

```go
func (w *URLWriter) writeGeminiURL(ctx *types.ParseContext) (string, error) {
	return fmt.Sprintf("[🤖 Gemini Chat](%s)", ctx.OriginalInput), nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test -run "TestURLWriter_Vote|TestURLWriter_WriteGeminiURL" ./internal/writer/ -v`
Expected: PASS

- [ ] **Step 5: Run full writer tests**

Run: `go test ./internal/writer/ -v`
Expected: all PASS

- [ ] **Step 6: Commit**

```bash
git add internal/writer/url_writer.go internal/writer/url_writer_test.go
git commit -m "feat: add Gemini chat URL writer"
```

---

## Chunk 3: Integration Test

### Task 4: Add end-to-end integration test

**Files:**
- Modify: `integration_test.go`

- [ ] **Step 1: Write the failing integration test**

In `integration_test.go`, find the test table in `TestEndToEndTransformation`. Add after the "Notion URL" test case (around line 108):

```go
		{
			name:           "Gemini Chat URL",
			input:          "https://gemini.google.com/app/ac9ebc9d76c30fc1",
			expectedOutput: "[🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)",
		},
```

- [ ] **Step 2: Run the integration test**

Run: `go test -run TestEndToEndTransformation -v`
Expected: PASS (parser and writer already implemented)

- [ ] **Step 3: Run full test suite**

Run: `go test ./... -v`
Expected: all PASS

- [ ] **Step 4: Run lint/format/vet**

Run: `go fmt ./... && go vet ./... && make lint`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add integration_test.go
git commit -m "test: add Gemini chat URL integration test"
```
