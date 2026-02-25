package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type CodeCommitLongParser struct {
	config *types.Config
}

func NewCodeCommitLongParser(cfg *types.Config) *CodeCommitLongParser {
	return &CodeCommitLongParser{config: cfg}
}

func (p *CodeCommitLongParser) CanHandle(input string) bool {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	// Need at least a few lines for long format
	if len(lines) < 3 {
		return false
	}

	// Look for AWS CodeCommit UI indicators
	hasCodeCommitIndicators := false
	hasRepoIndicator := false
	hasPRIndicator := false
	hasTitleWithNumber := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for AWS console UI indicators
		if strings.Contains(line, "Developer Tools") ||
			strings.Contains(line, "CodeCommit") {
			hasCodeCommitIndicators = true
		}

		// Check for "Repositories" indicator
		if strings.Contains(line, "Repositories") {
			hasRepoIndicator = true
		}

		// Check for "Pull requests" indicator
		if strings.Contains(line, "Pull requests") ||
			strings.Contains(line, "pull-requests") {
			hasPRIndicator = true
		}

		// Check for PR title with number pattern (e.g., "411: SEC-12335: Pass SENDGRID_API_KEY Securley")
		if hasCodeCommitPRTitle(line) {
			hasTitleWithNumber = true
		}
	}

	return hasCodeCommitIndicators && hasRepoIndicator && hasPRIndicator && hasTitleWithNumber
}

func (p *CodeCommitLongParser) Parse(input string) (*types.ParseContext, error) {
	if !p.CanHandle(input) {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")

	var region, repo, prNumber, prTitle string

	// Parse through lines to extract metadata
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to find region (usually appears early in the text)
		if region == "" && isAWSRegion(line) {
			region = line
			continue
		}

		// Try to find repository name
		// Repository name appears after "Repositories" and before "Pull requests"
		if repo == "" && i > 0 {
			prevLine := strings.TrimSpace(lines[i-1])
			if prevLine == "Repositories" || strings.Contains(prevLine, "CodeCommit") {
				if !strings.Contains(line, "Developer Tools") &&
					!strings.Contains(line, "CodeCommit") &&
					!strings.Contains(line, "Pull requests") &&
					!strings.Contains(line, "Repositories") &&
					isValidCodeCommitRepoName(line) {
					repo = line
					continue
				}
			}
		}

		// Try to find PR number (standalone number on a line)
		if prNumber == "" && isPRNumber(line) {
			prNumber = line
			continue
		}

		// Try to find PR title with number (e.g., "411: SEC-12335: Pass SENDGRID_API_KEY Securley")
		if prTitle == "" && hasCodeCommitPRTitle(line) {
			_, title := extractCodeCommitPRTitle(line)
			if title != "" {
				prTitle = title
				continue
			}
		}
	}

	// If we couldn't find a region in the text, try to infer from repository context
	if region == "" {
		region = "us-east-1" // Default to us-east-1 if not found
	}

	if repo == "" || prNumber == "" || prTitle == "" {
		return nil, nil
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeCodeCommitLong,
		Confidence:    90,
		Metadata: map[string]interface{}{
			"region": region,
			"repo":   repo,
			"number": prNumber,
			"title":  prTitle,
		},
	}

	return ctx, nil
}

// isAWSRegion checks if a string is a valid AWS region
func isAWSRegion(s string) bool {
	// Common AWS regions
	regions := []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"eu-west-1", "eu-west-2", "eu-west-3", "eu-central-1",
		"ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2",
		"ap-south-1", "sa-east-1", "ca-central-1",
	}

	for _, region := range regions {
		if s == region {
			return true
		}
	}
	return false
}

// isValidCodeCommitRepoName checks if a string could be a repository name
func isValidCodeCommitRepoName(s string) bool {
	if len(s) == 0 || len(s) > 100 {
		return false
	}

	// Repo names typically contain alphanumeric, hyphens, underscores
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	return re.MatchString(s)
}

// isPRNumber checks if a line is just a PR number
func isPRNumber(s string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(s) && len(s) < 10 // Reasonable PR number length
}

// hasCodeCommitPRTitle checks if a line contains a CodeCommit PR title
func hasCodeCommitPRTitle(s string) bool {
	// Pattern: number followed by colon and title
	// Example: "411: SEC-12335: Pass SENDGRID_API_KEY Securley"
	re := regexp.MustCompile(`^\d+:\s*.+`)
	return re.MatchString(s)
}

// extractCodeCommitPRTitle extracts the PR number and title from a line
func extractCodeCommitPRTitle(s string) (number, title string) {
	// Match pattern: number: title
	re := regexp.MustCompile(`^(\d+):\s*(.+)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(s))
	if len(matches) == 3 {
		return matches[1], strings.TrimSpace(matches[2])
	}
	return "", ""
}
