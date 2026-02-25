package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestCodeCommitLongParser_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name: "AWS CodeCommit long format",
			input: `Developer Tools
CodeCommit
Repositories
upserve-env
Pull requests
411
Developer Tools
CodeCommit
Repositories
upserve-env
Pull requests
411
411: SEC-12335: Pass SENDGRID_API_KEY Securley`,
			expected: true,
		},
		{
			name: "Short URL only",
			input: "https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1",
			expected: false,
		},
		{
			name: "Not enough lines",
			input: `CodeCommit
upserve-env`,
			expected: false,
		},
	}

	cfg := &types.Config{}
	parser := NewCodeCommitLongParser(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCodeCommitLongParser_Parse(t *testing.T) {
	cfg := &types.Config{}
	parser := NewCodeCommitLongParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectedType   types.ContentType
		expectedConf   int
		expectedRegion string
		expectedRepo   string
		expectedNumber string
		expectedTitle  string
	}{
		{
			name: "CodeCommit long format with region",
			input: `Developer Tools
CodeCommit
Repositories
upserve-env
Pull requests
411
Developer Tools
CodeCommit
Repositories
upserve-env
Pull requests
411
411: SEC-12335: Pass SENDGRID_API_KEY Securley`,
			expectedType:   types.ContentTypeCodeCommitLong,
			expectedConf:   90,
			expectedRegion: "us-east-1",
			expectedRepo:   "upserve-env",
			expectedNumber: "411",
			expectedTitle:  "SEC-12335: Pass SENDGRID_API_KEY Securley",
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

			if region := ctx.Metadata["region"]; region != tt.expectedRegion {
				t.Errorf("Metadata[region] = %v, want %v", region, tt.expectedRegion)
			}
			if repo := ctx.Metadata["repo"]; repo != tt.expectedRepo {
				t.Errorf("Metadata[repo] = %v, want %v", repo, tt.expectedRepo)
			}
			if number := ctx.Metadata["number"]; number != tt.expectedNumber {
				t.Errorf("Metadata[number] = %v, want %v", number, tt.expectedNumber)
			}
			if title := ctx.Metadata["title"]; title != tt.expectedTitle {
				t.Errorf("Metadata[title] = %v, want %v", title, tt.expectedTitle)
			}
		})
	}
}

func TestCodeCommitParser_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "AWS CodeCommit PR URL",
			input:    "https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1",
			expected: true,
		},
		{
			name:     "AWS CodeCommit PR URL different region",
			input:    "https://us-west-2.console.aws.amazon.com/codesuite/codecommit/repositories/my-repo/pull-requests/123/details?region=us-west-2",
			expected: true,
		},
		{
			name:     "GitHub URL",
			input:    "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expected: false,
		},
		{
			name:     "Non-URL text",
			input:    "PLAT-12345",
			expected: false,
		},
		{
			name:     "Plain text",
			input:    "hello world",
			expected: false,
		},
	}

	cfg := &types.Config{}
	parser := NewCodeCommitParser(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCodeCommitParser_Parse(t *testing.T) {
	cfg := &types.Config{}
	parser := NewCodeCommitParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectedType   types.ContentType
		expectedConf   int
		expectedRegion string
		expectedRepo   string
		expectedNumber string
	}{
		{
			name:           "CodeCommit PR URL us-east-1",
			input:          "https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1",
			expectedType:   types.ContentTypeCodeCommitURL,
			expectedConf:   90,
			expectedRegion: "us-east-1",
			expectedRepo:   "upserve-env",
			expectedNumber: "411",
		},
		{
			name:           "CodeCommit PR URL us-west-2",
			input:          "https://us-west-2.console.aws.amazon.com/codesuite/codecommit/repositories/my-repo/pull-requests/123/details?region=us-west-2",
			expectedType:   types.ContentTypeCodeCommitURL,
			expectedConf:   90,
			expectedRegion: "us-west-2",
			expectedRepo:   "my-repo",
			expectedNumber: "123",
		},
		{
			name:           "CodeCommit PR URL eu-west-1",
			input:          "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/test-repo/pull-requests/42/details?region=eu-west-1",
			expectedType:   types.ContentTypeCodeCommitURL,
			expectedConf:   90,
			expectedRegion: "eu-west-1",
			expectedRepo:   "test-repo",
			expectedNumber: "42",
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

			if region := ctx.Metadata["region"]; region != tt.expectedRegion {
				t.Errorf("Metadata[region] = %v, want %v", region, tt.expectedRegion)
			}
			if repo := ctx.Metadata["repo"]; repo != tt.expectedRepo {
				t.Errorf("Metadata[repo] = %v, want %v", repo, tt.expectedRepo)
			}
			if number := ctx.Metadata["number"]; number != tt.expectedNumber {
				t.Errorf("Metadata[number] = %v, want %v", number, tt.expectedNumber)
			}
		})
	}
}
