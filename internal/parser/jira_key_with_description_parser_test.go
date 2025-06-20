package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestJIRAKeyWithDescriptionParser_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name: "Valid JIRA key with description",
			input: `PLAT-12345

blinc - webhook proxy logs`,
			expected: true,
		},
		{
			name: "Valid JIRA key with multi-line description",
			input: `SPEED-456

Fix authentication issue
Additional details here`,
			expected: true,
		},
		{
			name: "Valid JIRA key with description and extra whitespace",
			input: `  PLAT-789  

  Description text  `,
			expected: true,
		},
		{
			name: "Invalid - no empty line separator",
			input: `PLAT-12345
blinc - webhook proxy logs`,
			expected: false,
		},
		{
			name: "Invalid - no description",
			input: `PLAT-12345

`,
			expected: false,
		},
		{
			name: "Invalid - simple JIRA key only",
			input: "PLAT-12345",
			expected: false,
		},
		{
			name: "Invalid - not a JIRA key",
			input: `not-a-jira-key

some description`,
			expected: false,
		},
		{
			name: "Invalid - lowercase JIRA key",
			input: `plat-12345

description`,
			expected: false,
		},
		{
			name: "Empty string",
			input: "",
			expected: false,
		},
	}

	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Projects: []string{"PLAT", "SPEED"},
		},
	}
	parser := NewJIRAKeyWithDescriptionParser(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestJIRAKeyWithDescriptionParser_Parse(t *testing.T) {
	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Projects: []string{"PLAT", "SPEED"},
		},
	}
	parser := NewJIRAKeyWithDescriptionParser(cfg)

	tests := []struct {
		name                string
		input               string
		expectSuccess       bool
		expectedKey         string
		expectedProject     string
		expectedDescription string
	}{
		{
			name: "Valid PLAT key with description",
			input: `PLAT-12345

blinc - webhook proxy logs`,
			expectSuccess:       true,
			expectedKey:         "PLAT-12345",
			expectedProject:     "PLAT",
			expectedDescription: "blinc - webhook proxy logs",
		},
		{
			name: "Valid SPEED key with description",
			input: `SPEED-456

Fix authentication issue`,
			expectSuccess:       true,
			expectedKey:         "SPEED-456",
			expectedProject:     "SPEED",
			expectedDescription: "Fix authentication issue",
		},
		{
			name: "Valid key with multi-line description",
			input: `PLAT-789

Fix authentication issue with SSO
Additional details about the bug`,
			expectSuccess:       true,
			expectedKey:         "PLAT-789",
			expectedProject:     "PLAT",
			expectedDescription: "Fix authentication issue with SSO Additional details about the bug",
		},
		{
			name: "Valid key with whitespace cleanup",
			input: `  PLAT-999  

  Description with spaces  
  Second line  `,
			expectSuccess:       true,
			expectedKey:         "PLAT-999",
			expectedProject:     "PLAT",
			expectedDescription: "Description with spaces Second line",
		},
		{
			name: "Unconfigured project",
			input: `INVALID-123

This should not be transformed`,
			expectSuccess: false,
		},
		{
			name: "Invalid format - no empty line",
			input: `PLAT-12345
description`,
			expectSuccess: false,
		},
		{
			name: "Invalid format - no description",
			input: `PLAT-12345

`,
			expectSuccess: false,
		},
		{
			name: "Non-JIRA text",
			input: `hello world

some description`,
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

				if ctx.DetectedType != types.ContentTypeJIRAKeyWithDescription {
					t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, types.ContentTypeJIRAKeyWithDescription)
				}

				if ctx.Confidence != 98 {
					t.Errorf("Confidence = %v, want 98", ctx.Confidence)
				}

				if key := ctx.Metadata["issue_key"]; key != tt.expectedKey {
					t.Errorf("Metadata[issue_key] = %v, want %v", key, tt.expectedKey)
				}

				if project := ctx.Metadata["project"]; project != tt.expectedProject {
					t.Errorf("Metadata[project] = %v, want %v", project, tt.expectedProject)
				}

				if description := ctx.Metadata["description"]; description != tt.expectedDescription {
					t.Errorf("Metadata[description] = %v, want %v", description, tt.expectedDescription)
				}
			} else {
				if ctx != nil {
					t.Errorf("Expected nil context for invalid input, got %+v", ctx)
				}
			}
		})
	}
}