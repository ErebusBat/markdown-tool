package parser

import (
	"github.com/erebusbat/markdown-tool/pkg/types"
)

// ParseContext holds data collected during parsing
type ParseContext struct {
	OriginalInput string
	DetectedType  types.ContentType
	Confidence    int
	Metadata      map[string]interface{}
}

// GetParsers returns all available parsers
func GetParsers(cfg *types.Config) []types.Parser {
	return []types.Parser{
		NewURLParser(cfg),
		NewGitHubLongParser(cfg),
		NewJIRAKeyParser(cfg),
	}
}
