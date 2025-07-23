package writer

import (
	"fmt"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type RaycastWriter struct {
	config *types.Config
}

func NewRaycastWriter(cfg *types.Config) *RaycastWriter {
	return &RaycastWriter{config: cfg}
}

func (w *RaycastWriter) GetName() string {
	return "RaycastWriter"
}

func (w *RaycastWriter) Vote(ctx *types.ParseContext) int {
	if ctx.DetectedType == types.ContentTypeRaycastURI {
		return 85
	}
	return 0
}

func (w *RaycastWriter) Write(ctx *types.ParseContext) (string, error) {
	if ctx.DetectedType != types.ContentTypeRaycastURI {
		return ctx.OriginalInput, nil
	}

	// Check if it's a Note URI first
	isNote, ok := ctx.Metadata["isNote"].(bool)
	if !ok {
		isNote = false
	}

	// Check if it's an AI Chat URI
	isAIChat, ok := ctx.Metadata["isAIChat"].(bool)
	if !ok {
		isAIChat = false
	}

	linkText := "Raycast"
	if isNote {
		linkText = "Raycast Note"
	} else if isAIChat {
		linkText = "Raycast AI"
	}

	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}