package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

var codexThreadRe = regexp.MustCompile(`^codex://threads/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)

type CodexParser struct {
	config *types.Config
}

func NewCodexParser(cfg *types.Config) *CodexParser {
	return &CodexParser{config: cfg}
}

func (p *CodexParser) CanHandle(input string) bool {
	return codexThreadRe.MatchString(strings.TrimSpace(input))
}

func (p *CodexParser) Parse(input string) (*types.ParseContext, error) {
	trimmed := strings.TrimSpace(input)
	matches := codexThreadRe.FindStringSubmatch(trimmed)
	if matches == nil {
		return nil, nil
	}

	threadID := matches[1]

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeCodexThread,
		Confidence:    90,
		Metadata: map[string]interface{}{
			"thread_id": threadID,
			"url":       trimmed,
		},
	}

	return ctx, nil
}
