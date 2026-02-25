package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type OpenCodeSessionParser struct {
	config *types.Config
}

func NewOpenCodeSessionParser(cfg *types.Config) *OpenCodeSessionParser {
	return &OpenCodeSessionParser{config: cfg}
}

func (p *OpenCodeSessionParser) CanHandle(input string) bool {
	return p.findSessionToken(input) != ""
}

func (p *OpenCodeSessionParser) Parse(input string) (*types.ParseContext, error) {
	token := p.findSessionToken(input)
	if token == "" {
		return nil, nil
	}

	confidence := 70
	trimmed := strings.TrimSpace(input)
	if trimmed == token {
		confidence = 90
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeOpenCodeSession,
		Confidence:    confidence,
		Metadata: map[string]interface{}{
			"session_token":  token,
			"is_exact_match": trimmed == token,
		},
	}

	return ctx, nil
}

func (p *OpenCodeSessionParser) findSessionToken(input string) string {
	trimmed := strings.TrimSpace(input)
	re := regexp.MustCompile(`(?i)\bses_[a-z0-9]+\b`)
	match := re.FindString(trimmed)
	return match
}
