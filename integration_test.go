package main

import (
	"testing"

	"github.com/erebusbat/markdown-tool/internal/parser"
	"github.com/erebusbat/markdown-tool/internal/writer"
	"github.com/erebusbat/markdown-tool/pkg/types"
)

// TestEndToEndTransformation tests the complete pipeline from input to output
func TestEndToEndTransformation(t *testing.T) {
	// Create test configuration
	cfg := &types.Config{
		GitHub: types.GitHubConfig{
			DefaultOrg:  "CompanyCam",
			DefaultRepo: "Company-Cam-API",
			Mappings: map[string]string{
				"companycam/company-cam-api": "CompanyCam/API",
				"companycam/companycam-mobile": "CompanyCam/mobile",
			},
		},
		JIRA: types.JIRAConfig{
			Domain:   "https://companycam.atlassian.net",
			Projects: []string{"PLAT", "SPEED"},
		},
	}

	tests := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "GitHub Pull Request URL",
			input:          "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedOutput: "[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)",
		},
		{
			name:           "GitHub Issue URL",
			input:          "https://github.com/CompanyCam/Company-Cam-API/issues/15217",
			expectedOutput: "[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/issues/15217)",
		},
		{
			name:           "JIRA Issue URL",
			input:          "https://companycam.atlassian.net/browse/PLAT-192",
			expectedOutput: "[PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)",
		},
		{
			name:           "JIRA Comment URL",
			input:          "https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266",
			expectedOutput: "[PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)",
		},
		{
			name:           "Notion URL",
			input:          "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0",
			expectedOutput: "[VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0)",
		},
		{
			name:           "Generic URL",
			input:          "http://ww3.domain.tld/path/to/document?query=value#anchor",
			expectedOutput: "[domain.tld](http://ww3.domain.tld/path/to/document?query=value#anchor)",
		},
		{
			name:           "JIRA Key - PLAT",
			input:          "PLAT-12345",
			expectedOutput: "[PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)",
		},
		{
			name:           "JIRA Key - SPEED",
			input:          "SPEED-456",
			expectedOutput: "[SPEED-456](https://companycam.atlassian.net/browse/SPEED-456)",
		},
		{
			name:           "Unconfigured JIRA Key",
			input:          "INVALID-123",
			expectedOutput: "INVALID-123",
		},
		{
			name:           "Plain text",
			input:          "hello world",
			expectedOutput: "hello world",
		},
		{
			name:           "Mixed content",
			input:          "Check out PLAT-999",
			expectedOutput: "Check out PLAT-999", // Won't match because it's not standalone
		},
		{
			name: "GitHub Long format",
			input: `CompanyCam
companycam-mobile

Type / to search
Code
Issues
78
Pull requests
12
Actions
Projects
Wiki
Security
7
Insights
A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog #6549`,
			expectedOutput: "[CompanyCam/mobile#6549: A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog](https://github.com/CompanyCam/companycam-mobile/issues/6549)",
		},
		{
			name: "JIRA Key with description - PLAT",
			input: `PLAT-192

blinc - webhook proxy logs`,
			expectedOutput: "[PLAT-192: blinc - webhook proxy logs](https://companycam.atlassian.net/browse/PLAT-192)",
		},
		{
			name: "JIRA Key with description - SPEED",
			input: `SPEED-456

Optimize database query performance`,
			expectedOutput: "[SPEED-456: Optimize database query performance](https://companycam.atlassian.net/browse/SPEED-456)",
		},
		{
			name: "JIRA Key with multi-line description",
			input: `PLAT-789

Fix authentication issue with SSO
Additional details about the bug`,
			expectedOutput: "[PLAT-789: Fix authentication issue with SSO Additional details about the bug](https://companycam.atlassian.net/browse/PLAT-789)",
		},
		{
			name: "Unconfigured JIRA Key with description",
			input: `INVALID-123

This should not be transformed`,
			expectedOutput: `INVALID-123

This should not be transformed`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := processInput(t, cfg, tt.input)
			if output != tt.expectedOutput {
				t.Errorf("processInput(%q) = %q, want %q", tt.input, output, tt.expectedOutput)
			}
		})
	}
}

// processInput simulates the main application processing pipeline
func processInput(t *testing.T, cfg *types.Config, input string) string {
	// Parse input
	parsers := parser.GetParsers(cfg)
	contexts := make([]*types.ParseContext, 0)

	for _, p := range parsers {
		if ctx, err := p.Parse(input); err == nil && ctx != nil {
			contexts = append(contexts, ctx)
		}
	}

	// Vote on best writer
	writers := writer.GetWriters(cfg)
	bestWriter, bestScore := writer.Vote(writers, contexts)

	if bestWriter == nil || bestScore == 0 {
		// No writer wants to handle this, output verbatim
		return input
	}

	// Generate output
	if len(contexts) == 0 {
		return input
	}

	output, err := bestWriter.Write(contexts[0])
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	return output
}

// TestConfigurationIntegration tests that configuration affects processing
func TestConfigurationIntegration(t *testing.T) {
	// Test with different JIRA configurations
	cfg1 := &types.Config{
		JIRA: types.JIRAConfig{
			Domain:   "https://test1.atlassian.net",
			Projects: []string{"TST"},
		},
	}

	cfg2 := &types.Config{
		JIRA: types.JIRAConfig{
			Domain:   "https://test2.atlassian.net",
			Projects: []string{"DEV"},
		},
	}

	// Test that TST key works with cfg1 but not cfg2
	input := "TST-123"

	output1 := processInput(t, cfg1, input)
	expected1 := "[TST-123](https://test1.atlassian.net/browse/TST-123)"
	if output1 != expected1 {
		t.Errorf("Config1 output = %q, want %q", output1, expected1)
	}

	output2 := processInput(t, cfg2, input)
	expected2 := "TST-123" // Should be unchanged because TST not in cfg2 projects
	if output2 != expected2 {
		t.Errorf("Config2 output = %q, want %q", output2, expected2)
	}
}

// TestGitHubMappingIntegration tests GitHub org/repo mappings
func TestGitHubMappingIntegration(t *testing.T) {
	cfgWithMapping := &types.Config{
		GitHub: types.GitHubConfig{
			Mappings: map[string]string{
				"company/long-repo-name": "Company/Short",
			},
		},
	}

	cfgWithoutMapping := &types.Config{
		GitHub: types.GitHubConfig{
			Mappings: map[string]string{},
		},
	}

	input := "https://github.com/Company/Long-Repo-Name/pull/123"

	// With mapping (case-insensitive)
	output1 := processInput(t, cfgWithMapping, input)
	expected1 := "[Company/Short#123](https://github.com/Company/Long-Repo-Name/pull/123)"
	if output1 != expected1 {
		t.Errorf("With mapping: got %q, want %q", output1, expected1)
	}

	// Without mapping
	output2 := processInput(t, cfgWithoutMapping, input)
	expected2 := "[Company/Long-Repo-Name#123](https://github.com/Company/Long-Repo-Name/pull/123)"
	if output2 != expected2 {
		t.Errorf("Without mapping: got %q, want %q", output2, expected2)
	}
}
