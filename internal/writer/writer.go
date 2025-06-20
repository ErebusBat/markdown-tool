package writer

import (
	"github.com/erebusbat/markdown-tool/pkg/types"
)

// GetWriters returns all available writers
func GetWriters(cfg *types.Config) []types.Writer {
	return []types.Writer{
		NewURLWriter(cfg),
		NewJIRAKeyWithDescriptionWriter(cfg),
		NewJIRAWriter(cfg),
		NewPassthroughWriter(),
	}
}

// Vote determines which writer should handle the parsed contexts
func Vote(writers []types.Writer, contexts []*types.ParseContext) (types.Writer, int) {
	if len(contexts) == 0 {
		return nil, 0
	}

	var bestWriter types.Writer
	var bestScore int

	for _, ctx := range contexts {
		for _, writer := range writers {
			score := writer.Vote(ctx)
			if score > bestScore {
				bestScore = score
				bestWriter = writer
			}
		}
	}

	return bestWriter, bestScore
}
