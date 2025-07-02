package writer

import (
	"fmt"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type PhoneWriter struct {
	config *types.Config
}

func NewPhoneWriter(cfg *types.Config) *PhoneWriter {
	return &PhoneWriter{config: cfg}
}

func (w *PhoneWriter) GetName() string {
	return "PhoneWriter"
}

func (w *PhoneWriter) Vote(ctx *types.ParseContext) int {
	switch ctx.DetectedType {
	case types.ContentTypePhone7Digit:
		return ctx.Confidence
	case types.ContentTypePhone10Digit:
		return ctx.Confidence
	case types.ContentTypePhone11Digit:
		return ctx.Confidence
	default:
		return 0
	}
}

func (w *PhoneWriter) Write(ctx *types.ParseContext) (string, error) {
	// Only handle phone number content types
	switch ctx.DetectedType {
	case types.ContentTypePhone7Digit, types.ContentTypePhone10Digit, types.ContentTypePhone11Digit:
		return w.writePhoneNumber(ctx)
	default:
		// Return original input for non-phone content
		return ctx.OriginalInput, nil
	}
}

func (w *PhoneWriter) writePhoneNumber(ctx *types.ParseContext) (string, error) {
	formattedDisplay, ok := ctx.Metadata["formatted_display"].(string)
	if !ok {
		return ctx.OriginalInput, fmt.Errorf("missing formatted_display in phone context")
	}

	telURL, ok := ctx.Metadata["tel_url"].(string)
	if !ok {
		return ctx.OriginalInput, fmt.Errorf("missing tel_url in phone context")
	}

	// Generate markdown link: [formatted](tel:url)
	return fmt.Sprintf("[%s](tel:%s)", formattedDisplay, telURL), nil
}