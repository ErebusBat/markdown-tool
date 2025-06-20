package writer

import (
	"fmt"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type JIRAKeyWithDescriptionWriter struct {
	config *types.Config
}

func NewJIRAKeyWithDescriptionWriter(cfg *types.Config) *JIRAKeyWithDescriptionWriter {
	return &JIRAKeyWithDescriptionWriter{config: cfg}
}

func (w *JIRAKeyWithDescriptionWriter) GetName() string {
	return "JIRAKeyWithDescriptionWriter"
}

func (w *JIRAKeyWithDescriptionWriter) Vote(ctx *types.ParseContext) int {
	if ctx.DetectedType == types.ContentTypeJIRAKeyWithDescription {
		return 98 // Higher confidence than simple JIRA key writer
	}
	return 0
}

func (w *JIRAKeyWithDescriptionWriter) Write(ctx *types.ParseContext) (string, error) {
	if ctx.DetectedType != types.ContentTypeJIRAKeyWithDescription {
		return ctx.OriginalInput, nil
	}

	issueKey, issueKeyOk := ctx.Metadata["issue_key"].(string)
	description, descriptionOk := ctx.Metadata["description"].(string)

	if !issueKeyOk || !descriptionOk || issueKey == "" || description == "" {
		return ctx.OriginalInput, nil
	}

	// Build JIRA URL
	jiraURL := fmt.Sprintf("%s/browse/%s", w.config.JIRA.Domain, issueKey)
	return fmt.Sprintf("[%s: %s](%s)", issueKey, description, jiraURL), nil
}