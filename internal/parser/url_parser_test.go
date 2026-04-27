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
		name              string
		input             string
		expectedType      types.ContentType
		expectedConf      int
		expectedOrg       string
		expectedRepo      string
		expectedNumber    string
		expectedIssueType string
	}{
		{
			name:              "GitHub Pull Request",
			input:             "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedType:      types.ContentTypeGitHubURL,
			expectedConf:      90,
			expectedOrg:       "CompanyCam",
			expectedRepo:      "Company-Cam-API",
			expectedNumber:    "15217",
			expectedIssueType: "pull",
		},
		{
			name:              "GitHub Issue",
			input:             "https://github.com/CompanyCam/Company-Cam-API/issues/15217",
			expectedType:      types.ContentTypeGitHubURL,
			expectedConf:      90,
			expectedOrg:       "CompanyCam",
			expectedRepo:      "Company-Cam-API",
			expectedNumber:    "15217",
			expectedIssueType: "issues",
		},
		{
			name:              "GitHub Repository",
			input:             "https://github.com/pedropark99/zig-book",
			expectedType:      types.ContentTypeGitHubURL,
			expectedConf:      90,
			expectedOrg:       "pedropark99",
			expectedRepo:      "zig-book",
			expectedNumber:    "",
			expectedIssueType: "",
		},
		{
			name:              "GitHub Commit Long Hash",
			input:             "https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b",
			expectedType:      types.ContentTypeGitHubURL,
			expectedConf:      90,
			expectedOrg:       "ErebusBat",
			expectedRepo:      "markdown-tool",
			expectedNumber:    "aa062a602a02d33f4a6e7880809ac3609fe1417b",
			expectedIssueType: "commit",
		},
		{
			name:              "GitHub Commit Short Hash",
			input:             "https://github.com/CompanyCam/Company-Cam-API/commit/abc123",
			expectedType:      types.ContentTypeGitHubURL,
			expectedConf:      90,
			expectedOrg:       "CompanyCam",
			expectedRepo:      "Company-Cam-API",
			expectedNumber:    "abc123",
			expectedIssueType: "commit",
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
			if issueType := ctx.Metadata["type"]; issueType != tt.expectedIssueType {
				// Handle case where type is nil (not set for repository URLs)
				if tt.expectedIssueType == "" && issueType == nil {
					// This is expected for repository URLs without type
				} else {
					t.Errorf("Metadata[type] = %v, want %v", issueType, tt.expectedIssueType)
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

func TestURLParser_Parse_Jenkins(t *testing.T) {
	cfg := &types.Config{
		Jenkins: types.JenkinsConfig{
			Domain: "https://jenkins.internal.upserve.com",
		},
	}
	parser := NewURLParser(cfg)

	tests := []struct {
		name                string
		input               string
		expectedType        types.ContentType
		expectedJobName     string
		expectedBuildNumber string
	}{
		{
			name:                "Jenkins build URL",
			input:               "https://jenkins.internal.upserve.com/job/app.swipely/114/",
			expectedType:        types.ContentTypeJenkinsURL,
			expectedJobName:     "app.swipely",
			expectedBuildNumber: "114",
		},
		{
			name:                "Jenkins build URL with console text",
			input:               "https://jenkins.internal.upserve.com/job/app.swipely/114/consoleText",
			expectedType:        types.ContentTypeJenkinsURL,
			expectedJobName:     "app.swipely",
			expectedBuildNumber: "114",
		},
		{
			name:                "Jenkins build URL with additional path",
			input:               "https://jenkins.internal.upserve.com/job/my-project/42/artifact/build.log",
			expectedType:        types.ContentTypeJenkinsURL,
			expectedJobName:     "my-project",
			expectedBuildNumber: "42",
		},
		{
			name:                "Jenkins URL with lastBuild",
			input:               "https://jenkins.internal.upserve.com/job/app.swipely/lastBuild/",
			expectedType:        types.ContentTypeJenkinsURL,
			expectedJobName:     "app.swipely",
			expectedBuildNumber: "",
		},
		{
			name:                "Jenkins URL without build identifier",
			input:               "https://jenkins.internal.upserve.com/job/my-project/",
			expectedType:        types.ContentTypeJenkinsURL,
			expectedJobName:     "my-project",
			expectedBuildNumber: "",
		},
		{
			name:                "Jenkins URL with lastSuccessfulBuild",
			input:               "https://jenkins.internal.upserve.com/job/app.swipely/lastSuccessfulBuild/consoleText",
			expectedType:        types.ContentTypeJenkinsURL,
			expectedJobName:     "app.swipely",
			expectedBuildNumber: "",
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

			if jobName := ctx.Metadata["job_name"]; jobName != tt.expectedJobName {
				t.Errorf("Metadata[job_name] = %v, want %v", jobName, tt.expectedJobName)
			}

			if buildNumber := ctx.Metadata["build_number"]; buildNumber != tt.expectedBuildNumber {
				// Handle case where build_number is nil (not set) vs empty string
				if tt.expectedBuildNumber == "" && buildNumber == nil {
					// This is expected for URLs without build numbers
				} else {
					t.Errorf("Metadata[build_number] = %v, want %v", buildNumber, tt.expectedBuildNumber)
				}
			}
		})
	}
}

func TestURLParser_Parse_CodeCommit(t *testing.T) {
	cfg := &types.Config{}
	parser := NewURLParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectedType   types.ContentType
		expectedRegion string
		expectedRepo   string
		expectedNumber string
	}{
		{
			name:           "CodeCommit PR URL us-east-1",
			input:          "https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1",
			expectedType:   types.ContentTypeCodeCommitURL,
			expectedRegion: "us-east-1",
			expectedRepo:   "upserve-env",
			expectedNumber: "411",
		},
		{
			name:           "CodeCommit PR URL us-west-2",
			input:          "https://us-west-2.console.aws.amazon.com/codesuite/codecommit/repositories/my-repo/pull-requests/123/details?region=us-west-2",
			expectedType:   types.ContentTypeCodeCommitURL,
			expectedRegion: "us-west-2",
			expectedRepo:   "my-repo",
			expectedNumber: "123",
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

func TestURLParser_Parse_CircleCI(t *testing.T) {
	cfg := &types.Config{}
	p := NewURLParser(cfg)

	tests := []struct {
		name              string
		input             string
		expectedType      types.ContentType
		expectedConf      int
		expectedVCS       string
		expectedOrg       string
		expectedRepo      string
		expectedPipeline  string
		expectedWorkflow  string
	}{
		{
			name:             "CircleCI pipeline URL",
			input:            "https://app.circleci.com/pipelines/github/upserve/swipely/96/workflows/17abd9c6-1190-49e9-a05f-4bf992a9d611",
			expectedType:     types.ContentTypeCircleCI,
			expectedConf:     90,
			expectedVCS:      "github",
			expectedOrg:      "upserve",
			expectedRepo:     "swipely",
			expectedPipeline: "96",
			expectedWorkflow: "17abd9c6-1190-49e9-a05f-4bf992a9d611",
		},
		{
			name:             "CircleCI with different org and repo",
			input:            "https://app.circleci.com/pipelines/github/CompanyCam/Company-Cam-API/15217/workflows/abc123de-4567-89ab-cdef-0123456789ab",
			expectedType:     types.ContentTypeCircleCI,
			expectedConf:     90,
			expectedVCS:      "github",
			expectedOrg:      "CompanyCam",
			expectedRepo:     "Company-Cam-API",
			expectedPipeline: "15217",
			expectedWorkflow: "abc123de-4567-89ab-cdef-0123456789ab",
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
			if vcs := ctx.Metadata["vcs"]; vcs != tt.expectedVCS {
				t.Errorf("Metadata[vcs] = %v, want %v", vcs, tt.expectedVCS)
			}
			if org := ctx.Metadata["org"]; org != tt.expectedOrg {
				t.Errorf("Metadata[org] = %v, want %v", org, tt.expectedOrg)
			}
			if repo := ctx.Metadata["repo"]; repo != tt.expectedRepo {
				t.Errorf("Metadata[repo] = %v, want %v", repo, tt.expectedRepo)
			}
			if pipeline := ctx.Metadata["pipeline_number"]; pipeline != tt.expectedPipeline {
				t.Errorf("Metadata[pipeline_number] = %v, want %v", pipeline, tt.expectedPipeline)
			}
			if workflow := ctx.Metadata["workflow_id"]; workflow != tt.expectedWorkflow {
				t.Errorf("Metadata[workflow_id] = %v, want %v", workflow, tt.expectedWorkflow)
			}
		})
	}
}

func TestURLParser_Parse_Gemini(t *testing.T) {
	cfg := &types.Config{}
	p := NewURLParser(cfg)

	tests := []struct {
		name             string
		input            string
		expectedType     types.ContentType
		expectedConf     int
		expectedChatID   string
		expectedCleanURL string
	}{
		{
			name:             "Gemini chat URL",
			input:            "https://gemini.google.com/app/ac9ebc9d76c30fc1",
			expectedType:     types.ContentTypeGeminiURL,
			expectedConf:     90,
			expectedChatID:   "ac9ebc9d76c30fc1",
			expectedCleanURL: "https://gemini.google.com/app/ac9ebc9d76c30fc1",
		},
		{
			name:             "Gemini chat URL with different ID",
			input:            "https://gemini.google.com/app/abcdef123456",
			expectedType:     types.ContentTypeGeminiURL,
			expectedConf:     90,
			expectedChatID:   "abcdef123456",
			expectedCleanURL: "https://gemini.google.com/app/abcdef123456",
		},
		{
			name:             "Gemini chat URL with trailing arrow",
			input:            "https://gemini.google.com/app/ac9ebc9d76c30fc1 →",
			expectedType:     types.ContentTypeGeminiURL,
			expectedConf:     90,
			expectedChatID:   "ac9ebc9d76c30fc1",
			expectedCleanURL: "https://gemini.google.com/app/ac9ebc9d76c30fc1",
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
				if cleanURL := ctx.Metadata["clean_url"]; cleanURL != tt.expectedCleanURL {
					t.Errorf("Metadata[clean_url] = %v, want %v", cleanURL, tt.expectedCleanURL)
				}
			}
		})
	}
}
