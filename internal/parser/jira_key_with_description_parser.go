package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type JIRAKeyWithDescriptionParser struct {
	config *types.Config
}

func NewJIRAKeyWithDescriptionParser(cfg *types.Config) *JIRAKeyWithDescriptionParser {
	return &JIRAKeyWithDescriptionParser{config: cfg}
}

func (p *JIRAKeyWithDescriptionParser) CanHandle(input string) bool {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	// Must have at least 3 lines: JIRA key, empty line, description
	if len(lines) < 3 {
		return false
	}

	// First line should be a JIRA key
	firstLine := strings.TrimSpace(lines[0])
	jiraKeyRegex := regexp.MustCompile(`^[A-Z]+-\d+$`)
	if !jiraKeyRegex.MatchString(firstLine) {
		return false
	}

	// Second line should be empty (or only whitespace)
	secondLine := strings.TrimSpace(lines[1])
	if secondLine != "" {
		return false
	}

	// Must have non-empty description (check if there's any non-whitespace content after line 2)
	hasDescription := false
	for i := 2; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "" {
			hasDescription = true
			break
		}
	}

	return hasDescription
}

func (p *JIRAKeyWithDescriptionParser) Parse(input string) (*types.ParseContext, error) {
	if !p.CanHandle(input) {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")
	
	// Extract JIRA key from first line
	jiraKey := strings.TrimSpace(lines[0])
	
	// Extract project key
	parts := strings.Split(jiraKey, "-")
	if len(parts) < 2 {
		return nil, nil
	}
	projectKey := parts[0]

	// Check if this project is configured
	configured := false
	for _, validProject := range p.config.JIRA.Projects {
		if validProject == projectKey {
			configured = true
			break
		}
	}

	if !configured {
		return nil, nil // Don't handle unconfigured projects
	}

	// Extract description from lines 3 onwards
	descriptionLines := make([]string, 0)
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			descriptionLines = append(descriptionLines, line)
		}
	}

	if len(descriptionLines) == 0 {
		return nil, nil // No description found
	}

	description := strings.Join(descriptionLines, " ")

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeJIRAKeyWithDescription,
		Confidence:    98, // Higher confidence than simple JIRA key
		Metadata: map[string]interface{}{
			"issue_key":   jiraKey,
			"project":     projectKey,
			"description": description,
		},
	}

	return ctx, nil
}