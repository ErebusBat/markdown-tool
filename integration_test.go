package main

import (
	"regexp"
	"strings"
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
		URL: types.URLConfig{
			DomainMappings: map[string]string{
				"companycam_slack_com": "slack",
				"youtube_com":          "YouTube",
			},
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
			name:           "GitHub Repository URL",
			input:          "https://github.com/pedropark99/zig-book",
			expectedOutput: "[pedropark99/zig-book](https://github.com/pedropark99/zig-book)",
		},
		{
			name:           "GitHub Repository URL with mapping",
			input:          "https://github.com/CompanyCam/Company-Cam-API",
			expectedOutput: "[CompanyCam/API](https://github.com/CompanyCam/Company-Cam-API)",
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
		// URL domain mapping tests
		{
			name:           "Slack URL with domain mapping",
			input:          "https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069",
			expectedOutput: "[slack](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)",
		},
		{
			name:           "YouTube URL with domain mapping",
			input:          "https://youtube.com/watch?v=abc123",
			expectedOutput: "[YouTube](https://youtube.com/watch?v=abc123)",
		},
		{
			name:           "Unmapped domain fallback behavior",
			input:          "https://example.com/path/to/page",
			expectedOutput: "[example.com](https://example.com/path/to/page)",
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
		// New GitHub issue transformation tests
		{
			name:           "Simple GitHub issue with default org/repo",
			input:          "adds blinc ddagent file #15407",
			expectedOutput: "[CompanyCam/API#15407: adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)",
		},
		{
			name:           "GitHub issue with username prefix",
			input:          "courtneylw adds blinc ddagent file #15407",
			expectedOutput: "[CompanyCam/API#15407: adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)",
		},
		{
			name:           "GitHub issue with username and underscore",
			input:          "plat_188 adds blinc ddagent file #15407",
			expectedOutput: "[CompanyCam/API#15407: adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)",
		},
		{
			name: "GitHub UI with username in title",
			input: `CompanyCam
Company-Cam-API

Type / to search
Code
Issues
209
Pull requests
67
Discussions
Actions
Projects
3
Wiki
Security
6
Insights
Settings
courtneylw adds blinc ddagent file #15407`,
			expectedOutput: "[CompanyCam/API#15407: courtneylw adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)",
		},
		{
			name: "GitHub UI without username in title", 
			input: `CompanyCam
Company-Cam-API

Type / to search
Code
Issues
209
Pull requests
67
Discussions
Actions
Projects
3
Wiki
Security
6
Insights
Settings
plat_188 adds blinc ddagent file #15407`,
			expectedOutput: "[CompanyCam/API#15407: plat_188 adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)",
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
		
		// Phone number tests - 7-digit
		{
			name:           "7-digit phone plain",
			input:          "1234567",
			expectedOutput: "📞 [123-4567](tel:1234567)",
		},
		{
			name:           "7-digit phone with dash",
			input:          "123-4567",
			expectedOutput: "📞 [123-4567](tel:1234567)",
		},
		{
			name:           "7-digit phone with dot",
			input:          "123.4567",
			expectedOutput: "📞 [123-4567](tel:1234567)",
		},
		
		// Phone number tests - 10-digit
		{
			name:           "10-digit phone plain",
			input:          "8901234567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "10-digit phone with dash",
			input:          "890-123-4567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "10-digit phone with dot",
			input:          "890.123.4567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "10-digit phone with parentheses and space",
			input:          "(890) 123-4567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "10-digit phone with parentheses no space",
			input:          "(890)123-4567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "10-digit phone with parentheses plain",
			input:          "(890)1234567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		
		// Phone number tests - 11-digit US
		{
			name:           "11-digit US phone plain",
			input:          "18901234567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "11-digit US phone with dash",
			input:          "1-890-123-4567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "11-digit US phone with dot",
			input:          "1.890.123.4567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "11-digit US phone with parentheses and space",
			input:          "1 (890) 123-4567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "11-digit US phone with parentheses no space",
			input:          "1(890)123-4567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "11-digit US phone with parentheses plain",
			input:          "1(890)1234567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		
		// Phone number tests - 11-digit international
		{
			name:           "11-digit international phone plain",
			input:          "+78901234567",
			expectedOutput: "📞 [+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:           "11-digit international phone with dash",
			input:          "+7-890-123-4567",
			expectedOutput: "📞 [+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:           "11-digit international phone with dot",
			input:          "+7.890.123.4567",
			expectedOutput: "📞 [+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:           "11-digit international phone with parentheses and space",
			input:          "+7 (890) 123-4567",
			expectedOutput: "📞 [+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:           "11-digit international phone with parentheses no space",
			input:          "+7(890)123-4567",
			expectedOutput: "📞 [+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:           "11-digit international phone with parentheses plain",
			input:          "+7(890)1234567",
			expectedOutput: "📞 [+7-890-123-4567](tel:+78901234567)",
		},
		
		// Phone number non-matches (should pass through unchanged)
		{
			name:           "7-digit with space - no match",
			input:          "123 4567",
			expectedOutput: "123 4567",
		},
		{
			name:           "7-digit with comma - no match",
			input:          "123,4567",
			expectedOutput: "123,4567",
		},
		{
			name:           "7-digit with leading zero - no match",
			input:          "01234567",
			expectedOutput: "01234567",
		},
		{
			name:           "10-digit with extra digit - no match",
			input:          "89012345670",
			expectedOutput: "89012345670",
		},
		{
			name:           "10-digit with spaces - no match",
			input:          "890 123 4567",
			expectedOutput: "890 123 4567",
		},
		{
			name:           "10-digit mixed separators - no match",
			input:          "(890) 123 4567",
			expectedOutput: "(890) 123 4567",
		},
		{
			name:           "10-digit wrong grouping - no match",
			input:          "(890) 1234 567",
			expectedOutput: "(890) 1234 567",
		},
		{
			name:           "10-digit incomplete - no match",
			input:          "(890)123-456",
			expectedOutput: "(890)123-456",
		},
		{
			name:           "10-digit too many digits - no match",
			input:          "(890)12345679",
			expectedOutput: "(890)12345679",
		},
		
		// Tel URI tests - should be preprocessed and handled by existing phone parsers
		{
			name:           "Tel URI 7-digit",
			input:          "tel:1234567",
			expectedOutput: "📞 [123-4567](tel:1234567)",
		},
		{
			name:           "Tel URI 10-digit",
			input:          "tel:8901234567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "Tel URI 11-digit US",
			input:          "tel:18901234567",
			expectedOutput: "📞 [1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "Tel URI international",
			input:          "tel:+18901234567",
			expectedOutput: "📞 [+1-890-123-4567](tel:+18901234567)",
		},
		{
			name:           "Tel URI with dashes",
			input:          "tel:890-123-4567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "Tel URI with parentheses",
			input:          "tel:(890)123-4567",
			expectedOutput: "📞 [890-123-4567](tel:8901234567)",
		},
		{
			name:           "Tel URI empty",
			input:          "tel:",
			expectedOutput: "",
		},
		{
			name:           "Raycast AI Chat URI",
			input:          "raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D",
			expectedOutput: "[Raycast AI](raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D)",
		},
		{
			name:           "Raycast Note URI",
			input:          "raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22:%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D",
			expectedOutput: "[Raycast Note](raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22:%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D)",
		},
		{
			name:           "Raycast generic extension URI",
			input:          "raycast://extensions/other/extension",
			expectedOutput: "[Raycast](raycast://extensions/other/extension)",
		},
		{
			name:           "Raycast settings URI",
			input:          "raycast://settings",
			expectedOutput: "[Raycast](raycast://settings)",
		},
		{
			name:           "Raycast window management URI",
			input:          "raycast://extensions/raycast/window-management/center",
			expectedOutput: "[Raycast](raycast://extensions/raycast/window-management/center)",
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

// preprocessTelURIs converts tel: URIs to phone numbers that can be processed by existing parsers
func preprocessTelURIs(input string) string {
	// Pattern to match tel: URIs
	telPattern := regexp.MustCompile(`^tel:(.*)$`)
	
	matches := telPattern.FindStringSubmatch(strings.TrimSpace(input))
	if matches != nil {
		// Extract the phone number part after "tel:"
		phoneNumber := matches[1]
		
		// Return the phone number without the tel: prefix
		// This allows existing phone parsers to handle all supported formats
		return phoneNumber
	}
	
	return input
}

// processInput simulates the main application processing pipeline
func processInput(t *testing.T, cfg *types.Config, input string) string {
	// Preprocess tel: URIs (same as in cmd/root.go)
	input = preprocessTelURIs(input)
	
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

// TestURLDomainMappingIntegration tests URL domain mappings
func TestURLDomainMappingIntegration(t *testing.T) {
	cfgWithMapping := &types.Config{
		URL: types.URLConfig{
			DomainMappings: map[string]string{
				"companycam_slack_com": "slack",
				"youtube_com":          "YouTube",
			},
		},
	}

	cfgWithoutMapping := &types.Config{
		URL: types.URLConfig{
			DomainMappings: map[string]string{},
		},
	}

	tests := []struct {
		name           string
		config         *types.Config
		input          string
		expectedOutput string
	}{
		{
			name:           "Slack URL with mapping",
			config:         cfgWithMapping,
			input:          "https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069",
			expectedOutput: "[slack](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)",
		},
		{
			name:           "YouTube URL with mapping",
			config:         cfgWithMapping,
			input:          "https://youtube.com/watch?v=abc123",
			expectedOutput: "[YouTube](https://youtube.com/watch?v=abc123)",
		},
		{
			name:           "Slack URL without mapping",
			config:         cfgWithoutMapping,
			input:          "https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069",
			expectedOutput: "[companycam.slack.com](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)",
		},
		{
			name:           "YouTube URL without mapping",
			config:         cfgWithoutMapping,
			input:          "https://youtube.com/watch?v=abc123",
			expectedOutput: "[youtube.com](https://youtube.com/watch?v=abc123)",
		},
		{
			name:           "Case-insensitive domain matching",
			config:         cfgWithMapping,
			input:          "https://CompanyCam.Slack.com/archives/test",
			expectedOutput: "[slack](https://CompanyCam.Slack.com/archives/test)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := processInput(t, tt.config, tt.input)
			if output != tt.expectedOutput {
				t.Errorf("processInput(%q) = %q, want %q", tt.input, output, tt.expectedOutput)
			}
		})
	}
}
