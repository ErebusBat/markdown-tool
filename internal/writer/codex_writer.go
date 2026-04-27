package writer

import (
	"fmt"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type CodexWriter struct {
	config *types.Config
}

func NewCodexWriter(cfg *types.Config) *CodexWriter {
	return &CodexWriter{config: cfg}
}

func (w *CodexWriter) GetName() string {
	return "CodexWriter"
}

func (w *CodexWriter) Vote(ctx *types.ParseContext) int {
	if ctx.DetectedType == types.ContentTypeCodexThread {
		return ctx.Confidence
	}
	return 0
}

func (w *CodexWriter) Write(ctx *types.ParseContext) (string, error) {
	if ctx.DetectedType != types.ContentTypeCodexThread {
		return ctx.OriginalInput, nil
	}

	rawURL, ok := ctx.Metadata["url"].(string)
	if !ok || rawURL == "" {
		return ctx.OriginalInput, nil
	}

	return fmt.Sprintf("[🤖 Codex](%s)", rawURL), nil
}
