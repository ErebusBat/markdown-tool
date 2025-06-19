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
	// Allow multiple spaces before # and trailing spaces
	re := regexp.MustCompile(`^(.+?)\s+#(\d+)\s*$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(s))
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), matches[2]
	}
	return "", ""
}