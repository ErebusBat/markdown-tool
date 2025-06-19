package writer

import (
	"github.com/erebusbat/markdown-tool/pkg/types"
)

type PassthroughWriter struct{}

func NewPassthroughWriter() *PassthroughWriter {
	return &PassthroughWriter{}
}

func (w *PassthroughWriter) GetName() string {
	return "PassthroughWriter"
}

func (w *PassthroughWriter) Vote(ctx *types.ParseContext) int {
	// Always willing to handle input, but with lowest priority
	return 1
}

func (w *PassthroughWriter) Write(ctx *types.ParseContext) (string, error) {
	// Output the original input unchanged
	return ctx.OriginalInput, nil
}