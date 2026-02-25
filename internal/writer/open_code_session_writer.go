package writer

import (
	"fmt"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type OpenCodeSessionWriter struct {
	config *types.Config
}

func NewOpenCodeSessionWriter(cfg *types.Config) *OpenCodeSessionWriter {
	return &OpenCodeSessionWriter{config: cfg}
}

func (w *OpenCodeSessionWriter) GetName() string {
	return "OpenCodeSessionWriter"
}

func (w *OpenCodeSessionWriter) Vote(ctx *types.ParseContext) int {
	if ctx.DetectedType == types.ContentTypeOpenCodeSession {
		return ctx.Confidence
	}
	return 0
}

func (w *OpenCodeSessionWriter) Write(ctx *types.ParseContext) (string, error) {
	if ctx.DetectedType != types.ContentTypeOpenCodeSession {
		return ctx.OriginalInput, nil
	}

	sessionToken, ok := ctx.Metadata["session_token"].(string)
	if !ok || sessionToken == "" {
		return ctx.OriginalInput, nil
	}

	return fmt.Sprintf("[🤖OpenCode](opencode -s %s)", sessionToken), nil
}
