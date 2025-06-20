package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestJIRAKeyWithDescriptionWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewJIRAKeyWithDescriptionWriter(cfg)

	tests := []struct {
		name         string
		contentType  types.ContentType
		expectedVote int
	}{
		{"JIRA Key with Description", types.ContentTypeJIRAKeyWithDescription, 98},
		{"JIRA Key (simple)", types.ContentTypeJIRAKey, 0},
		{"GitHub URL", types.ContentTypeGitHubURL, 0},
		{"JIRA URL", types.ContentTypeJIRAURL, 0},
		{"Generic URL", types.ContentTypeURL, 0},
		{"Unknown", types.ContentTypeUnknown, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				DetectedType: tt.contentType,
			}
			vote := writer.Vote(ctx)
			if vote != tt.expectedVote {
				t.Errorf("Vote() = %v, want %v", vote, tt.expectedVote)
			}
		})
	}
}

func TestJIRAKeyWithDescriptionWriter_Write(t *testing.T) {
	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Domain: "https://companycam.atlassian.net",
		},
	}
	writer := NewJIRAKeyWithDescriptionWriter(cfg)

	tests := []struct {
		name           string
		contentType    types.ContentType
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name:        "Valid JIRA Key with description",
			contentType: types.ContentTypeJIRAKeyWithDescription,
			metadata: map[string]interface{}{
				"issue_key":   "PLAT-12345",
				"project":     "PLAT",
				"description": "blinc - webhook proxy logs",
			},
			originalInput:  "PLAT-12345\n\nblinc - webhook proxy logs",
			expectedOutput: "[PLAT-12345: blinc - webhook proxy logs](https://companycam.atlassian.net/browse/PLAT-12345)",
		},
		{
			name:        "Another JIRA Key with description",
			contentType: types.ContentTypeJIRAKeyWithDescription,
			metadata: map[string]interface{}{
				"issue_key":   "SPEED-456",
				"project":     "SPEED",
				"description": "Optimize database query performance",
			},
			originalInput:  "SPEED-456\n\nOptimize database query performance",
			expectedOutput: "[SPEED-456: Optimize database query performance](https://companycam.atlassian.net/browse/SPEED-456)",
		},
		{
			name:        "JIRA Key with multi-line description",
			contentType: types.ContentTypeJIRAKeyWithDescription,
			metadata: map[string]interface{}{
				"issue_key":   "PLAT-789",
				"project":     "PLAT",
				"description": "Fix authentication issue with SSO Additional details about the bug",
			},
			originalInput:  "PLAT-789\n\nFix authentication issue with SSO\nAdditional details about the bug",
			expectedOutput: "[PLAT-789: Fix authentication issue with SSO Additional details about the bug](https://companycam.atlassian.net/browse/PLAT-789)",
		},
		{
			name:           "Non-JIRA content type returns original",
			contentType:    types.ContentTypeURL,
			metadata:       map[string]interface{}{},
			originalInput:  "https://example.com",
			expectedOutput: "https://example.com",
		},
		{
			name:           "JIRA Key with Description without metadata returns original",
			contentType:    types.ContentTypeJIRAKeyWithDescription,
			metadata:       map[string]interface{}{},
			originalInput:  "PLAT-999\n\nsome description",
			expectedOutput: "PLAT-999\n\nsome description",
		},
		{
			name:        "JIRA Key with Description missing issue_key returns original",
			contentType: types.ContentTypeJIRAKeyWithDescription,
			metadata: map[string]interface{}{
				"project":     "PLAT",
				"description": "some description",
			},
			originalInput:  "PLAT-999\n\nsome description",
			expectedOutput: "PLAT-999\n\nsome description",
		},
		{
			name:        "JIRA Key with Description missing description returns original",
			contentType: types.ContentTypeJIRAKeyWithDescription,
			metadata: map[string]interface{}{
				"issue_key": "PLAT-999",
				"project":   "PLAT",
			},
			originalInput:  "PLAT-999\n\nsome description",
			expectedOutput: "PLAT-999\n\nsome description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  tt.contentType,
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

func TestJIRAKeyWithDescriptionWriter_GetName(t *testing.T) {
	cfg := &types.Config{}
	writer := NewJIRAKeyWithDescriptionWriter(cfg)
	
	expected := "JIRAKeyWithDescriptionWriter"
	if name := writer.GetName(); name != expected {
		t.Errorf("GetName() = %v, want %v", name, expected)
	}
}