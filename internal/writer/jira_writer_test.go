package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestJIRAWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewJIRAWriter(cfg)

	tests := []struct {
		name         string
		contentType  types.ContentType
		expectedVote int
	}{
		{"JIRA Key", types.ContentTypeJIRAKey, 95},
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

func TestJIRAWriter_Write(t *testing.T) {
	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Domain: "https://companycam.atlassian.net",
		},
	}
	writer := NewJIRAWriter(cfg)

	tests := []struct {
		name           string
		contentType    types.ContentType
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name:        "Valid JIRA Key",
			contentType: types.ContentTypeJIRAKey,
			metadata: map[string]interface{}{
				"issue_key": "PLAT-12345",
				"project":   "PLAT",
			},
			originalInput:  "PLAT-12345",
			expectedOutput: "[PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)",
		},
		{
			name:        "Another JIRA Key",
			contentType: types.ContentTypeJIRAKey,
			metadata: map[string]interface{}{
				"issue_key": "SPEED-456",
				"project":   "SPEED",
			},
			originalInput:  "SPEED-456",
			expectedOutput: "[SPEED-456](https://companycam.atlassian.net/browse/SPEED-456)",
		},
		{
			name:           "Non-JIRA content type returns original",
			contentType:    types.ContentTypeURL,
			metadata:       map[string]interface{}{},
			originalInput:  "https://example.com",
			expectedOutput: "https://example.com",
		},
		{
			name:        "JIRA Key without metadata returns original",
			contentType: types.ContentTypeJIRAKey,
			metadata:    map[string]interface{}{},
			originalInput:  "PLAT-999",
			expectedOutput: "PLAT-999",
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