package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestCodexParser_CanHandle(t *testing.T) {
	cfg := &types.Config{}
	p := NewCodexParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid codex thread URI",
			input:    "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
			expected: true,
		},
		{
			name:     "Valid codex thread URI with leading/trailing whitespace",
			input:    "  codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf  ",
			expected: true,
		},
		{
			name:     "Missing thread ID",
			input:    "codex://threads/",
			expected: false,
		},
		{
			name:     "Non-UUID thread ID",
			input:    "codex://threads/notauuid",
			expected: false,
		},
		{
			name:     "Wrong scheme",
			input:    "opencode://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
			expected: false,
		},
		{
			name:     "Plain text",
			input:    "no match here",
			expected: false,
		},
		{
			name:     "URL with extra path",
			input:    "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf/extra",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCodexParser_Parse(t *testing.T) {
	cfg := &types.Config{}
	p := NewCodexParser(cfg)

	tests := []struct {
		name             string
		input            string
		shouldParse      bool
		expectedThreadID string
		expectedURL      string
		expectedConf     int
	}{
		{
			name:             "Valid codex thread URI",
			input:            "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
			shouldParse:      true,
			expectedThreadID: "019dcc44-e7b8-7c23-816c-34c194bdb3cf",
			expectedURL:      "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
			expectedConf:     90,
		},
		{
			name:             "URI with surrounding whitespace",
			input:            "  codex://threads/aabbccdd-1122-3344-5566-778899aabbcc  ",
			shouldParse:      true,
			expectedThreadID: "aabbccdd-1122-3344-5566-778899aabbcc",
			expectedURL:      "codex://threads/aabbccdd-1122-3344-5566-778899aabbcc",
			expectedConf:     90,
		},
		{
			name:        "No match",
			input:       "hello world",
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := p.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if !tt.shouldParse {
				if ctx != nil {
					t.Errorf("Parse(%q) = %v, want nil", tt.input, ctx)
				}
				return
			}

			if ctx == nil {
				t.Fatalf("Parse(%q) = nil, want non-nil", tt.input)
				return
			}

			if ctx.DetectedType != types.ContentTypeCodexThread {
				t.Errorf("DetectedType = %v, want ContentTypeCodexThread", ctx.DetectedType)
			}

			if ctx.Confidence != tt.expectedConf {
				t.Errorf("Confidence = %d, want %d", ctx.Confidence, tt.expectedConf)
			}

			threadID, ok := ctx.Metadata["thread_id"].(string)
			if !ok {
				t.Fatalf("missing thread_id metadata")
			}
			if threadID != tt.expectedThreadID {
				t.Errorf("thread_id = %q, want %q", threadID, tt.expectedThreadID)
			}

			rawURL, ok := ctx.Metadata["url"].(string)
			if !ok {
				t.Fatalf("missing url metadata")
			}
			if rawURL != tt.expectedURL {
				t.Errorf("url = %q, want %q", rawURL, tt.expectedURL)
			}
		})
	}
}
