package parser

import (
	"net/url"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type RaycastParser struct {
	config *types.Config
}

func NewRaycastParser(cfg *types.Config) *RaycastParser {
	return &RaycastParser{config: cfg}
}

func (p *RaycastParser) CanHandle(input string) bool {
	// Check if it starts with raycast://
	if !strings.HasPrefix(input, "raycast://") {
		return false
	}
	
	// Verify it's a valid URI
	_, err := url.Parse(input)
	return err == nil
}

func (p *RaycastParser) Parse(input string) (*types.ParseContext, error) {
	if !p.CanHandle(input) {
		return nil, nil
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeRaycastURI,
		Confidence:    85,
		Metadata:      make(map[string]interface{}),
	}

	// Check if it's an AI Chat URI
	// The path should contain the full path including leading /
	isAIChat := strings.Contains(input, "extensions/raycast/raycast-ai/ai-chat")
	ctx.Metadata["isAIChat"] = isAIChat

	// Check if it's a Note URI
	isNote := strings.Contains(input, "extensions/raycast/raycast-notes/raycast-notes")
	ctx.Metadata["isNote"] = isNote

	return ctx, nil
}