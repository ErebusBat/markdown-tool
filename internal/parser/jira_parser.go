package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type JIRAKeyParser struct {
	config *types.Config
}

func NewJIRAKeyParser(cfg *types.Config) *JIRAKeyParser {
	return &JIRAKeyParser{config: cfg}
}

func (p *JIRAKeyParser) CanHandle(input string) bool {
	// Check if input looks like a JIRA key (PROJECT-123)
	re := regexp.MustCompile(`^[A-Z]+-\d+$`)
	return re.MatchString(strings.TrimSpace(input))
}

func (p *JIRAKeyParser) Parse(input string) (*types.ParseContext, error) {
	trimmed := strings.TrimSpace(input)
	if !p.CanHandle(trimmed) {
		return nil, nil
	}

	// Extract project key
	parts := strings.Split(trimmed, "-")
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

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeJIRAKey,
		Confidence:    95,
		Metadata: map[string]interface{}{
			"issue_key": trimmed,
			"project":   projectKey,
		},
	}

	return ctx, nil
}
