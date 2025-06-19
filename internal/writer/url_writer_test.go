package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestURLWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name         string
		contentType  types.ContentType
		expectedVote int
	}{
		{"GitHub URL", types.ContentTypeGitHubURL, 90},
		{"JIRA URL", types.ContentTypeJIRAURL, 90},
		{"JIRA Comment", types.ContentTypeJIRAComment, 95},
		{"Notion URL", types.ContentTypeNotionURL, 85},
		{"Generic URL", types.ContentTypeURL, 50},
		{"JIRA Key", types.ContentTypeJIRAKey, 0},
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

func TestURLWriter_WriteGitHubURL(t *testing.T) {
	tests := []struct {
		name           string
		config         *types.Config
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name: "GitHub PR without mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "Company-Cam-API",
				"number": "15217",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedOutput: "[CompanyCam/Company-Cam-API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)",
		},
		{
			name: "GitHub PR with mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{
						"companycam/company-cam-api": "CompanyCam/API",
					},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "Company-Cam-API",
				"number": "15217",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedOutput: "[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)",
		},
		{
			name: "GitHub Issue",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "someorg",
				"repo":   "somerepo",
				"number": "42",
			},
			originalInput:  "https://github.com/someorg/somerepo/issues/42",
			expectedOutput: "[someorg/somerepo#42](https://github.com/someorg/somerepo/issues/42)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewURLWriter(tt.config)
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeGitHubURL,
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

func TestURLWriter_WriteJIRAURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		contentType    types.ContentType
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name:        "JIRA Issue URL",
			contentType: types.ContentTypeJIRAURL,
			metadata: map[string]interface{}{
				"issue_key": "PLAT-192",
			},
			originalInput:  "https://companycam.atlassian.net/browse/PLAT-192",
			expectedOutput: "[PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)",
		},
		{
			name:        "JIRA Comment URL",
			contentType: types.ContentTypeJIRAComment,
			metadata: map[string]interface{}{
				"issue_key":  "PLAT-192",
				"comment_id": "20266",
			},
			originalInput:  "https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266",
			expectedOutput: "[PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)",
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

func TestURLWriter_WriteNotionURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	ctx := &types.ParseContext{
		OriginalInput: "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0",
		DetectedType:  types.ContentTypeNotionURL,
		Metadata: map[string]interface{}{
			"title": "VS Code Setup for Standard rb RubyLSP",
		},
	}

	expectedOutput := "[VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0)"

	output, err := writer.Write(ctx)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Write() = %v, want %v", output, expectedOutput)
	}
}

func TestURLWriter_WriteGenericURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name:          "Generic URL with www",
			originalInput: "https://www.example.com/path/to/page",
			metadata: map[string]interface{}{
				"domain": "www.example.com",
			},
			expectedOutput: "[example.com](https://www.example.com/path/to/page)",
		},
		{
			name:          "Generic URL with ww3",
			originalInput: "http://ww3.domain.tld/path/to/document?query=value#anchor",
			metadata: map[string]interface{}{
				"domain": "ww3.domain.tld",
			},
			expectedOutput: "[domain.tld](http://ww3.domain.tld/path/to/document?query=value#anchor)",
		},
		{
			name:          "Simple domain",
			originalInput: "https://example.org/page",
			metadata: map[string]interface{}{
				"domain": "example.org",
			},
			expectedOutput: "[example.org](https://example.org/page)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeURL,
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
