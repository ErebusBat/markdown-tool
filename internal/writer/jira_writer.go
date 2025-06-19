package writer

import (
	"fmt"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type JIRAWriter struct {
	config *types.Config
}

func NewJIRAWriter(cfg *types.Config) *JIRAWriter {
	return &JIRAWriter{config: cfg}
}

func (w *JIRAWriter) GetName() string {
	return "JIRAWriter"
}

func (w *JIRAWriter) Vote(ctx *types.ParseContext) int {
	if ctx.DetectedType == types.ContentTypeJIRAKey {
		return 95
	}
	return 0
}

func (w *JIRAWriter) Write(ctx *types.ParseContext) (string, error) {
	if ctx.DetectedType != types.ContentTypeJIRAKey {
		return ctx.OriginalInput, nil
	}

	issueKey, _ := ctx.Metadata["issue_key"].(string)
	if issueKey == "" {
		return ctx.OriginalInput, nil
	}

	// Build JIRA URL
	jiraURL := fmt.Sprintf("%s/browse/%s", w.config.JIRA.Domain, issueKey)
	return fmt.Sprintf("[%s](%s)", issueKey, jiraURL), nil
}
