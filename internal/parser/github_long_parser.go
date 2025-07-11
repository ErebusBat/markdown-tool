package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type GitHubLongParser struct {
	config *types.Config
}

func NewGitHubLongParser(cfg *types.Config) *GitHubLongParser {
	return &GitHubLongParser{config: cfg}
}

func (p *GitHubLongParser) CanHandle(input string) bool {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	// Case 1: Simple issue title pattern (single line or few lines)
	if isSimpleIssueTitle(input) {
		return true
	}
	
	// Case 2 & 3: Multi-line GitHub UI content
	if len(lines) < 3 {
		return false
	}

	// Look for patterns that suggest GitHub UI content
	var hasOrgCandidate, hasRepoCandidate, hasIssueTitle bool
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check for potential organization name (simple word, often starts with capital)
		if isValidGitHubName(line) && !hasOrgCandidate {
			hasOrgCandidate = true
		}
		
		// Check for potential repository name (can contain hyphens, underscores)
		if isValidRepoName(line) && !hasRepoCandidate {
			hasRepoCandidate = true
		}
		
		// Check for issue/PR title with number at the end
		if hasIssueTitleWithNumber(line) {
			hasIssueTitle = true
		}
	}

	// We need at least org-like, repo-like, and issue title patterns
	return hasOrgCandidate && hasRepoCandidate && hasIssueTitle
}

func (p *GitHubLongParser) Parse(input string) (*types.ParseContext, error) {
	if !p.CanHandle(input) {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	var org, repo, issueTitle, issueNumber string
	
	// Check if this is a simple issue title pattern
	if isSimpleIssueTitle(input) {
		return p.parseSimpleIssueTitle(input)
	}
	
	// Handle multi-line GitHub UI content (existing logic)
	// First pass: find the issue title with number (most distinctive)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if hasIssueTitleWithNumber(line) {
			issueTitle, issueNumber = extractIssueTitleAndNumber(line)
			break
		}
	}
	
	if issueTitle == "" || issueNumber == "" {
		return nil, nil
	}
	
	// Second pass: find org and repo names
	// In GitHub UI copies, the first line is typically the org, second line is the repo
	var foundOrg, foundRepo bool
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, issueTitle) {
			continue
		}
		
		// First valid GitHub name should be the org
		if !foundOrg && isValidGitHubName(line) {
			org = line
			foundOrg = true
			continue
		}
		
		// Second valid name (that's different from org) should be the repo
		if foundOrg && !foundRepo && isValidRepoName(line) && line != org {
			repo = line
			foundRepo = true
			break
		}
	}
	
	if org == "" || repo == "" {
		return nil, nil
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeGitHubLong,
		Confidence:    90,
		Metadata: map[string]interface{}{
			"org":         org,
			"repo":        repo,
			"title":       issueTitle,
			"number":      issueNumber,
			"type":        "issues", // Default to issues, could be enhanced to detect PRs
		},
	}

	return ctx, nil
}

// parseSimpleIssueTitle handles simple issue title patterns using default org/repo from config
func (p *GitHubLongParser) parseSimpleIssueTitle(input string) (*types.ParseContext, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	var issueTitle, issueNumber string
	
	// Find the line with the issue title
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check for username prefix pattern
		if hasGitHubUsernamePrefix(line) {
			issueTitle, issueNumber = extractUsernameAndIssue(line)
			if issueTitle != "" && issueNumber != "" {
				break
			}
		}
		
		// Check for simple issue title with number
		if hasIssueTitleWithNumber(line) {
			issueTitle, issueNumber = extractIssueTitleAndNumber(line)
			if issueTitle != "" && issueNumber != "" {
				break
			}
		}
	}
	
	if issueTitle == "" || issueNumber == "" {
		return nil, nil
	}
	
	// Use default org and repo from config
	org := p.config.GitHub.DefaultOrg
	repo := p.config.GitHub.DefaultRepo
	
	if org == "" || repo == "" {
		return nil, nil
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeGitHubLong,
		Confidence:    95, // Higher confidence for simple patterns
		Metadata: map[string]interface{}{
			"org":         org,
			"repo":        repo,
			"title":       issueTitle,
			"number":      issueNumber,
			"type":        "issues",
		},
	}

	return ctx, nil
}

// extractUsernameAndIssue extracts the issue title and number from a line with username prefix
func extractUsernameAndIssue(s string) (title, number string) {
	// Pattern: username followed by title and #number
	// Example: "courtneylw adds blinc ddagent file #15407"
	re := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s+(.+)\s+#(\d+)\s*$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(s))
	if len(matches) == 4 {
		// Check if the first part is actually a valid GitHub username
		if isGitHubUsername(matches[1]) {
			return strings.TrimSpace(matches[2]), matches[3]
		}
	}
	return "", ""
}

// isValidGitHubName checks if a string could be a GitHub organization/user name
func isValidGitHubName(s string) bool {
	if len(s) == 0 || len(s) > 39 {
		return false
	}
	
	// GitHub names can contain alphanumeric and hyphens, but not start/end with hyphen
	re := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`)
	return re.MatchString(s)
}

// isValidRepoName checks if a string could be a GitHub repository name
func isValidRepoName(s string) bool {
	if len(s) == 0 || len(s) > 100 {
		return false
	}
	
	// Repo names can contain alphanumeric, hyphens, underscores, dots
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	return re.MatchString(s)
}

// hasIssueTitleWithNumber checks if a line contains an issue title ending with #number
func hasIssueTitleWithNumber(s string) bool {
	// Look for text ending with #number pattern (allowing multiple spaces)
	re := regexp.MustCompile(`^.+\s+#\d+\s*$`)
	return re.MatchString(strings.TrimSpace(s))
}

// extractIssueTitleAndNumber extracts the title and number from an issue title line
func extractIssueTitleAndNumber(s string) (title, number string) {
	// Match everything before the last space followed by #number
	re := regexp.MustCompile(`^(.+)\s+#(\d+)\s*$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(s))
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), matches[2]
	}
	return "", ""
}

// isSimpleIssueTitle checks if input is a simple issue title pattern (just title + #number)
func isSimpleIssueTitle(input string) bool {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	// For simple patterns, we allow single line or up to 3 lines
	// But we need to ensure it's not a multi-line GitHub UI pattern
	if len(lines) > 3 {
		return false
	}
	
	// If it's a single line, check for issue patterns
	if len(lines) == 1 {
		line := strings.TrimSpace(lines[0])
		return hasGitHubUsernamePrefix(line) || hasIssueTitleWithNumber(line)
	}
	
	// For multi-line, it should not have typical GitHub UI indicators
	// Also, if we have org/repo-like lines, it should be handled as multi-line
	hasGitHubUIIndicators := false
	var issueLineCount int
	var orgRepoLikeLineCount int
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check for GitHub UI indicators
		if strings.Contains(line, "Type / to search") || 
		   strings.Contains(line, "Pull requests") ||
		   strings.Contains(line, "Discussions") ||
		   strings.Contains(line, "Actions") ||
		   strings.Contains(line, "Projects") ||
		   strings.Contains(line, "Wiki") ||
		   strings.Contains(line, "Security") ||
		   strings.Contains(line, "Insights") ||
		   strings.Contains(line, "Settings") {
			hasGitHubUIIndicators = true
		}
		
		// Count lines with issue patterns
		if hasGitHubUsernamePrefix(line) || hasIssueTitleWithNumber(line) {
			issueLineCount++
		} else if isValidGitHubName(line) || isValidRepoName(line) {
			// Count potential org/repo lines
			orgRepoLikeLineCount++
		}
	}
	
	// If we have GitHub UI indicators, this should be handled by multi-line logic
	if hasGitHubUIIndicators {
		return false
	}
	
	// If we have org/repo-like lines, treat as multi-line GitHub UI content
	if orgRepoLikeLineCount >= 2 {
		return false
	}
	
	// We should have exactly one issue line for simple patterns
	return issueLineCount == 1
}

// hasGitHubUsernamePrefix checks if a line starts with a GitHub username followed by issue title
func hasGitHubUsernamePrefix(s string) bool {
	// Pattern: username followed by text and ending with #number
	// Examples: "courtneylw adds blinc ddagent file #15407", "plat_188 adds blinc ddagent file #15407"
	// Must not match simple issue titles like "adds blinc ddagent file #15407"
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+\s+.+\s+#\d+\s*$`)
	line := strings.TrimSpace(s)
	if !re.MatchString(line) {
		return false
	}
	
	// Make sure the first word is actually a valid GitHub username
	parts := strings.Fields(line)
	if len(parts) < 3 { // username + at least one word + #number
		return false
	}
	
	return isGitHubUsername(parts[0])
}

// isGitHubUsername checks if a string looks like a GitHub username
func isGitHubUsername(s string) bool {
	if len(s) == 0 || len(s) > 39 {
		return false
	}
	
	// GitHub usernames can contain alphanumeric, hyphens, and underscores
	// but cannot start or end with hyphens
	re := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9_-]*[a-zA-Z0-9])?$`)
	if !re.MatchString(s) {
		return false
	}
	
	// Additional check: common English words are unlikely to be usernames
	// This helps distinguish "adds" from actual usernames like "user123"
	commonWords := []string{"adds", "fixes", "updates", "removes", "creates", "deletes", "implements", "enhances", "refactors", "only", "some", "the", "and", "with", "for", "from", "this", "that"}
	for _, word := range commonWords {
		if strings.EqualFold(s, word) {
			return false
		}
	}
	
	return true
}