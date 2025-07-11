package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestGitHubLongParser_CanHandle(t *testing.T) {
	cfg := &types.Config{}
	parser := NewGitHubLongParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name: "Valid GitHub issue text chunk",
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
			expected: true,
		},
		{
			name: "Valid GitHub PR text chunk",
			input: `MyOrg
my-awesome-repo
Code
Issues
Pull requests
Fix the login bug that affects mobile users #123`,
			expected: true,
		},
		{
			name: "Simple valid case",
			input: `TestOrg
test-repo
Some issue title here #42`,
			expected: true,
		},
		// New test cases for simple issue patterns
		{
			name:     "Simple issue title only",
			input:    "adds blinc ddagent file #15407",
			expected: true,
		},
		{
			name:     "Issue with username prefix",
			input:    "courtneylw adds blinc ddagent file #15407",
			expected: true,
		},
		{
			name:     "Issue with username and underscore",
			input:    "plat_188 adds blinc ddagent file #15407",
			expected: true,
		},
		{
			name: "Multi-line with username prefix",
			input: `CompanyCam
Company-Cam-API
courtneylw adds blinc ddagent file #15407`,
			expected: true,
		},
		{
			name:     "Too few lines",
			input:    "CompanyCam\ncompanycam-mobile",
			expected: false,
		},
		{
			name: "Missing issue title with number",
			input: `CompanyCam
companycam-mobile
Code
Issues
Just some text without number`,
			expected: false,
		},
		{
			name: "Missing org/repo pattern - multi-line without GitHub UI",
			input: `Some random text
More random text
Another line
But we have an issue title #123`,
			expected: false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGitHubLongParser_Parse(t *testing.T) {
	cfg := &types.Config{
		GitHub: types.GitHubConfig{
			DefaultOrg:  "CompanyCam",
			DefaultRepo: "Company-Cam-API",
		},
	}
	parser := NewGitHubLongParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectSuccess  bool
		expectedOrg    string
		expectedRepo   string
		expectedTitle  string
		expectedNumber string
	}{
		{
			name: "GitHub issue from example",
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
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "companycam-mobile",
			expectedTitle:  "A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog",
			expectedNumber: "6549",
		},
		{
			name: "Multi-line GitHub issue",
			input: `MyOrg
my-repo
Fix the authentication bug #123`,
			expectSuccess:  true,
			expectedOrg:    "MyOrg",
			expectedRepo:   "my-repo",
			expectedTitle:  "Fix the authentication bug",
			expectedNumber: "123",
		},
		{
			name: "GitHub issue with extra content",
			input: `TestOrg
test-repository
Code
Issues
Pull requests
Actions
Implement new feature for user management #999
Labels
Assignees`,
			expectSuccess:  true,
			expectedOrg:    "TestOrg",
			expectedRepo:   "test-repository",
			expectedTitle:  "Implement new feature for user management",
			expectedNumber: "999",
		},
		{
			name: "Invalid input - no issue number",
			input: `CompanyCam
companycam-mobile
Just some description without number`,
			expectSuccess: false,
		},
		{
			name: "Single line issue (should use defaults)",
			input: `Only one line with issue #123`,
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedTitle:  "Only one line with issue",
			expectedNumber: "123",
		},
		{
			name:          "Empty input",
			input:         "",
			expectSuccess: false,
		},
		// New test cases for simple issue patterns
		{
			name:           "Simple issue title with default org/repo",
			input:          "adds blinc ddagent file #15407",
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedTitle:  "adds blinc ddagent file",
			expectedNumber: "15407",
		},
		{
			name:           "Issue with username prefix",
			input:          "courtneylw adds blinc ddagent file #15407",
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedTitle:  "adds blinc ddagent file",
			expectedNumber: "15407",
		},
		{
			name:           "Issue with username and underscore",
			input:          "plat_188 adds blinc ddagent file #15407",
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedTitle:  "adds blinc ddagent file",
			expectedNumber: "15407",
		},
		{
			name: "Multi-line GitHub UI with username",
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
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedTitle:  "courtneylw adds blinc ddagent file",
			expectedNumber: "15407",
		},
		{
			name: "Multi-line GitHub UI without username in title",
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
			expectSuccess:  true,
			expectedOrg:    "CompanyCam",
			expectedRepo:   "Company-Cam-API",
			expectedTitle:  "plat_188 adds blinc ddagent file",
			expectedNumber: "15407",
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

				if ctx.DetectedType != types.ContentTypeGitHubLong {
					t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, types.ContentTypeGitHubLong)
				}

				// Simple issue patterns have confidence 95, multi-line patterns have 90
				expectedConfidence := 90
				if tt.name == "Simple issue title with default org/repo" ||
				   tt.name == "Issue with username prefix" ||
				   tt.name == "Issue with username and underscore" ||
				   tt.name == "Single line issue (should use defaults)" {
					expectedConfidence = 95
				}
				
				if ctx.Confidence != expectedConfidence {
					t.Errorf("Confidence = %v, want %v", ctx.Confidence, expectedConfidence)
				}

				if org := ctx.Metadata["org"]; org != tt.expectedOrg {
					t.Errorf("Metadata[org] = %v, want %v", org, tt.expectedOrg)
				}

				if repo := ctx.Metadata["repo"]; repo != tt.expectedRepo {
					t.Errorf("Metadata[repo] = %v, want %v", repo, tt.expectedRepo)
				}

				if title := ctx.Metadata["title"]; title != tt.expectedTitle {
					t.Errorf("Metadata[title] = %v, want %v", title, tt.expectedTitle)
				}

				if number := ctx.Metadata["number"]; number != tt.expectedNumber {
					t.Errorf("Metadata[number] = %v, want %v", number, tt.expectedNumber)
				}

				if issueType := ctx.Metadata["type"]; issueType != "issues" {
					t.Errorf("Metadata[type] = %v, want issues", issueType)
				}
			} else {
				if ctx != nil {
					t.Errorf("Expected nil context for invalid input, got %+v", ctx)
				}
			}
		})
	}
}

func TestGitHubLongParser_HelperFunctions(t *testing.T) {
	tests := []struct {
		name     string
		function string
		input    string
		expected bool
	}{
		// isValidGitHubName tests
		{"Valid org name", "isValidGitHubName", "CompanyCam", true},
		{"Valid user name", "isValidGitHubName", "octocat", true},
		{"Valid name with hyphen", "isValidGitHubName", "my-org", true},
		{"Invalid - starts with hyphen", "isValidGitHubName", "-invalid", false},
		{"Invalid - ends with hyphen", "isValidGitHubName", "invalid-", false},
		{"Invalid - too long", "isValidGitHubName", "thisnameistoolongforgihtuborganizationnames", false},
		{"Invalid - empty", "isValidGitHubName", "", false},
		{"Invalid - special chars", "isValidGitHubName", "invalid@name", false},

		// isValidRepoName tests
		{"Valid repo name", "isValidRepoName", "companycam-mobile", true},
		{"Valid repo with dots", "isValidRepoName", "my.awesome.repo", true},
		{"Valid repo with underscores", "isValidRepoName", "my_repo_name", true},
		{"Valid mixed", "isValidRepoName", "test-repo_v2.0", true},
		{"Invalid - empty", "isValidRepoName", "", false},
		{"Invalid - special chars", "isValidRepoName", "repo@name", false},

		// hasIssueTitleWithNumber tests
		{"Valid issue title", "hasIssueTitleWithNumber", "Fix the bug #123", true},
		{"Valid complex title", "hasIssueTitleWithNumber", "A specific Logger.error call doesn't log #6549", true},
		{"Invalid - no number", "hasIssueTitleWithNumber", "Just some text", false},
		{"Invalid - number not at end", "hasIssueTitleWithNumber", "Issue #123 with more text", false},
		{"Invalid - no hash", "hasIssueTitleWithNumber", "Issue 123", false},

		// hasGitHubUsernamePrefix tests
		{"Valid username prefix", "hasGitHubUsernamePrefix", "courtneylw adds blinc ddagent file #15407", true},
		{"Valid username with underscore", "hasGitHubUsernamePrefix", "plat_188 adds blinc ddagent file #15407", true},
		{"Valid username with hyphen", "hasGitHubUsernamePrefix", "user-name fixes issue #123", true},
		{"Invalid - no issue number", "hasGitHubUsernamePrefix", "courtneylw adds blinc ddagent file", false},
		{"Invalid - no username", "hasGitHubUsernamePrefix", "adds blinc ddagent file #15407", false},
		{"Invalid - number not at end", "hasGitHubUsernamePrefix", "user adds #123 more text", false},

		// isGitHubUsername tests  
		{"Valid username", "isGitHubUsername", "courtneylw", true},
		{"Valid username with underscore", "isGitHubUsername", "plat_188", true},
		{"Valid username with hyphen", "isGitHubUsername", "user-name", true},
		{"Valid single char", "isGitHubUsername", "a", true},
		{"Invalid - starts with hyphen", "isGitHubUsername", "-invalid", false},
		{"Invalid - ends with hyphen", "isGitHubUsername", "invalid-", false},
		{"Invalid - too long", "isGitHubUsername", "thisusernameistoolongforgihtubusernameextra", false},
		{"Invalid - empty", "isGitHubUsername", "", false},
		{"Invalid - special chars", "isGitHubUsername", "user@name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch tt.function {
			case "isValidGitHubName":
				result = isValidGitHubName(tt.input)
			case "isValidRepoName":
				result = isValidRepoName(tt.input)
			case "hasIssueTitleWithNumber":
				result = hasIssueTitleWithNumber(tt.input)
			case "hasGitHubUsernamePrefix":
				result = hasGitHubUsernamePrefix(tt.input)
			case "isGitHubUsername":
				result = isGitHubUsername(tt.input)
			default:
				t.Fatalf("Unknown function: %s", tt.function)
			}

			if result != tt.expected {
				t.Errorf("%s(%q) = %v, want %v", tt.function, tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractIssueTitleAndNumber(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedTitle   string
		expectedNumber  string
	}{
		{
			name:           "Simple issue",
			input:          "Fix the bug #123",
			expectedTitle:  "Fix the bug",
			expectedNumber: "123",
		},
		{
			name:           "Complex issue title",
			input:          "A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog #6549",
			expectedTitle:  "A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog",
			expectedNumber: "6549",
		},
		{
			name:           "Issue with extra spaces",
			input:          "  Fix authentication   #999  ",
			expectedTitle:  "Fix authentication",
			expectedNumber: "999",
		},
		{
			name:           "Invalid format",
			input:          "No number here",
			expectedTitle:  "",
			expectedNumber: "",
		},
		{
			name:           "Number not at end",
			input:          "Issue #123 with more text",
			expectedTitle:  "",
			expectedNumber: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, number := extractIssueTitleAndNumber(tt.input)
			
			if title != tt.expectedTitle {
				t.Errorf("extractIssueTitleAndNumber() title = %q, want %q", title, tt.expectedTitle)
			}
			
			if number != tt.expectedNumber {
				t.Errorf("extractIssueTitleAndNumber() number = %q, want %q", number, tt.expectedNumber)
			}
		})
	}
}

func TestExtractUsernameAndIssue(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedTitle   string
		expectedNumber  string
	}{
		{
			name:           "Username with simple issue",
			input:          "courtneylw adds blinc ddagent file #15407",
			expectedTitle:  "adds blinc ddagent file",
			expectedNumber: "15407",
		},
		{
			name:           "Username with underscore",
			input:          "plat_188 adds blinc ddagent file #15407",
			expectedTitle:  "adds blinc ddagent file",
			expectedNumber: "15407",
		},
		{
			name:           "Username with hyphen",
			input:          "user-name fixes authentication bug #999",
			expectedTitle:  "fixes authentication bug",
			expectedNumber: "999",
		},
		{
			name:           "Username with complex issue title",
			input:          "developer123 implement new feature for user management system #1234",
			expectedTitle:  "implement new feature for user management system",
			expectedNumber: "1234",
		},
		{
			name:           "Issue with extra spaces",
			input:          "  user_123   fix   bug   #456  ",
			expectedTitle:  "fix   bug",
			expectedNumber: "456",
		},
		{
			name:           "Invalid format - no username (common word)",
			input:          "adds blinc ddagent file #15407",
			expectedTitle:  "",
			expectedNumber: "",
		},
		{
			name:           "Invalid format - no number",
			input:          "courtneylw adds blinc ddagent file",
			expectedTitle:  "",
			expectedNumber: "",
		},
		{
			name:           "Invalid format - number not at end",
			input:          "user adds #123 more text",
			expectedTitle:  "",
			expectedNumber: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, number := extractUsernameAndIssue(tt.input)
			
			if title != tt.expectedTitle {
				t.Errorf("extractUsernameAndIssue() title = %q, want %q", title, tt.expectedTitle)
			}
			
			if number != tt.expectedNumber {
				t.Errorf("extractUsernameAndIssue() number = %q, want %q", number, tt.expectedNumber)
			}
		})
	}
}