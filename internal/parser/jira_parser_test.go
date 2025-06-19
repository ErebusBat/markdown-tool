package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestJIRAKeyParser_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid JIRA key", "PLAT-12345", true},
		{"Valid JIRA key with spaces", "  PLAT-12345  ", true},
		{"Another valid JIRA key", "SPEED-456", true},
		{"Invalid - no number", "PLAT-", false},
		{"Invalid - no dash", "PLAT12345", false},
		{"Invalid - lowercase", "plat-12345", false},
		{"Invalid - URL", "https://example.com", false},
		{"Invalid - contains extra text", "PLAT-12345 and more", false},
		{"Empty string", "", false},
	}

	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Projects: []string{"PLAT", "SPEED"},
		},
	}
	parser := NewJIRAKeyParser(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestJIRAKeyParser_Parse(t *testing.T) {
	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Projects: []string{"PLAT", "SPEED"},
		},
	}
	parser := NewJIRAKeyParser(cfg)

	tests := []struct {
		name            string
		input           string
		expectSuccess   bool
		expectedKey     string
		expectedProject string
	}{
		{
			name:            "Valid PLAT key",
			input:           "PLAT-12345",
			expectSuccess:   true,
			expectedKey:     "PLAT-12345",
			expectedProject: "PLAT",
		},
		{
			name:            "Valid SPEED key",
			input:           "SPEED-456",
			expectSuccess:   true,
			expectedKey:     "SPEED-456",
			expectedProject: "SPEED",
		},
		{
			name:            "Valid key with whitespace",
			input:           "  PLAT-789  ",
			expectSuccess:   true,
			expectedKey:     "PLAT-789",
			expectedProject: "PLAT",
		},
		{
			name:          "Unconfigured project",
			input:         "INVALID-123",
			expectSuccess: false,
		},
		{
			name:          "Non-JIRA text",
			input:         "hello world",
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := parser.Parse(tt.input)

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if tt.expectSuccess {
				if ctx == nil {
					t.Fatal("Expected successful parse but got nil context")
				}

				if ctx.DetectedType != types.ContentTypeJIRAKey {
					t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, types.ContentTypeJIRAKey)
				}

				if ctx.Confidence != 95 {
					t.Errorf("Confidence = %v, want 95", ctx.Confidence)
				}

				if key := ctx.Metadata["issue_key"]; key != tt.expectedKey {
					t.Errorf("Metadata[issue_key] = %v, want %v", key, tt.expectedKey)
				}

				if project := ctx.Metadata["project"]; project != tt.expectedProject {
					t.Errorf("Metadata[project] = %v, want %v", project, tt.expectedProject)
				}
			} else {
				if ctx != nil {
					t.Errorf("Expected nil context for invalid input, got %+v", ctx)
				}
			}
		})
	}
}
