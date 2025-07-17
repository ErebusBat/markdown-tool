package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestURLParser_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"GitHub URL", "https://github.com/CompanyCam/Company-Cam-API/pull/15217", true},
		{"JIRA URL", "https://companycam.atlassian.net/browse/PLAT-192", true},
		{"Notion URL", "https://www.notion.so/companycam/VS-Code-Setup-654a6b070ae74ac3ad400c6d571507c0", true},
		{"Generic HTTP URL", "http://example.com/path", true},
		{"Generic HTTPS URL", "https://example.com/path", true},
		{"Non-URL text", "PLAT-12345", false},
		{"Plain text", "hello world", false},
		{"Empty string", "", false},
	}

	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Domain: "https://companycam.atlassian.net",
		},
	}
	parser := NewURLParser(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestURLParser_Parse_GitHub(t *testing.T) {
	cfg := &types.Config{}
	parser := NewURLParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectedType   types.ContentType
		expectedConf   int
		expectedOrg    string
		expectedRepo   string
		expectedNumber string
	}{
		{
			name:           "GitHub Pull Request",
			input:          "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedType:   types.ContentTypeGitHubURL,
			expectedConf:   90,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedNumber: "15217",
		},
		{
			name:           "GitHub Issue",
			input:          "https://github.com/CompanyCam/Company-Cam-API/issues/15217",
			expectedType:   types.ContentTypeGitHubURL,
			expectedConf:   90,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedNumber: "15217",
		},
		{
			name:           "GitHub Repository",
			input:          "https://github.com/pedropark99/zig-book",
			expectedType:   types.ContentTypeGitHubURL,
			expectedConf:   90,
			expectedOrg:    "pedropark99",
			expectedRepo:   "zig-book",
			expectedNumber: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := parser.Parse(tt.input)
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

			if org := ctx.Metadata["org"]; org != tt.expectedOrg {
				t.Errorf("Metadata[org] = %v, want %v", org, tt.expectedOrg)
			}
			if repo := ctx.Metadata["repo"]; repo != tt.expectedRepo {
				t.Errorf("Metadata[repo] = %v, want %v", repo, tt.expectedRepo)
			}
			if number := ctx.Metadata["number"]; number != tt.expectedNumber {
				// Handle case where number is nil (not set for repository URLs)
				if tt.expectedNumber == "" && number == nil {
					// This is expected for repository URLs without issue numbers
				} else {
					t.Errorf("Metadata[number] = %v, want %v", number, tt.expectedNumber)
				}
			}
		})
	}
}

func TestURLParser_Parse_JIRA(t *testing.T) {
	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Domain: "https://companycam.atlassian.net",
		},
	}
	parser := NewURLParser(cfg)

	tests := []struct {
		name         string
		input        string
		expectedType types.ContentType
		expectedKey  string
		hasComment   bool
	}{
		{
			name:         "JIRA Issue",
			input:        "https://companycam.atlassian.net/browse/PLAT-192",
			expectedType: types.ContentTypeJIRAURL,
			expectedKey:  "PLAT-192",
			hasComment:   false,
		},
		{
			name:         "JIRA Comment",
			input:        "https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266",
			expectedType: types.ContentTypeJIRAComment,
			expectedKey:  "PLAT-192",
			hasComment:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if ctx == nil {
				t.Fatal("Parse() returned nil context")
			}

			if ctx.DetectedType != tt.expectedType {
				t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, tt.expectedType)
			}

			if key := ctx.Metadata["issue_key"]; key != tt.expectedKey {
				t.Errorf("Metadata[issue_key] = %v, want %v", key, tt.expectedKey)
			}

			if tt.hasComment {
				if _, exists := ctx.Metadata["comment_id"]; !exists {
					t.Error("Expected comment_id in metadata for comment URL")
				}
			}
		})
	}
}

func TestURLParser_Parse_Notion(t *testing.T) {
	cfg := &types.Config{}
	parser := NewURLParser(cfg)

	input := "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0"
	ctx, err := parser.Parse(input)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if ctx == nil {
		t.Fatal("Parse() returned nil context")
	}

	if ctx.DetectedType != types.ContentTypeNotionURL {
		t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, types.ContentTypeNotionURL)
	}

	expectedTitle := "VS Code Setup for Standard rb RubyLSP"
	if title := ctx.Metadata["title"]; title != expectedTitle {
		t.Errorf("Metadata[title] = %v, want %v", title, expectedTitle)
	}
}

func TestURLParser_Parse_Generic(t *testing.T) {
	cfg := &types.Config{}
	parser := NewURLParser(cfg)

	input := "http://ww3.domain.tld/path/to/document?query=value#anchor"
	ctx, err := parser.Parse(input)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if ctx == nil {
		t.Fatal("Parse() returned nil context")
	}

	if ctx.DetectedType != types.ContentTypeURL {
		t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, types.ContentTypeURL)
	}

	expectedDomain := "ww3.domain.tld"
	if domain := ctx.Metadata["domain"]; domain != expectedDomain {
		t.Errorf("Metadata[domain] = %v, want %v", domain, expectedDomain)
	}
}
