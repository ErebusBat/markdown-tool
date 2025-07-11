package writer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type URLWriter struct {
	config *types.Config
}

func NewURLWriter(cfg *types.Config) *URLWriter {
	return &URLWriter{config: cfg}
}

func (w *URLWriter) GetName() string {
	return "URLWriter"
}

func (w *URLWriter) Vote(ctx *types.ParseContext) int {
	switch ctx.DetectedType {
	case types.ContentTypeGitHubURL:
		return 90
	case types.ContentTypeGitHubLong:
		return 95
	case types.ContentTypeJIRAURL:
		return 90
	case types.ContentTypeJIRAComment:
		return 95
	case types.ContentTypeNotionURL:
		return 85
	case types.ContentTypeURL:
		return 50
	default:
		return 0
	}
}

func (w *URLWriter) Write(ctx *types.ParseContext) (string, error) {
	switch ctx.DetectedType {
	case types.ContentTypeGitHubURL:
		return w.writeGitHubURL(ctx)
	case types.ContentTypeGitHubLong:
		return w.writeGitHubLongURL(ctx)
	case types.ContentTypeJIRAURL:
		return w.writeJIRAURL(ctx)
	case types.ContentTypeJIRAComment:
		return w.writeJIRACommentURL(ctx)
	case types.ContentTypeNotionURL:
		return w.writeNotionURL(ctx)
	case types.ContentTypeURL:
		return w.writeGenericURL(ctx)
	default:
		return ctx.OriginalInput, nil
	}
}

func (w *URLWriter) writeGitHubURL(ctx *types.ParseContext) (string, error) {
	org, _ := ctx.Metadata["org"].(string)
	repo, _ := ctx.Metadata["repo"].(string)
	number, _ := ctx.Metadata["number"].(string)

	if org == "" || repo == "" || number == "" {
		return w.writeGenericURL(ctx)
	}

	// Apply organization/repository mappings if configured
	orgRepo := fmt.Sprintf("%s/%s", org, repo)
	// Try case-insensitive lookup since Viper lowercases map keys
	for key, mapped := range w.config.GitHub.Mappings {
		if strings.EqualFold(key, orgRepo) {
			orgRepo = mapped
			break
		}
	}

	linkText := fmt.Sprintf("%s#%s", orgRepo, number)
	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeGitHubLongURL(ctx *types.ParseContext) (string, error) {
	org, _ := ctx.Metadata["org"].(string)
	repo, _ := ctx.Metadata["repo"].(string)
	title, _ := ctx.Metadata["title"].(string)
	number, _ := ctx.Metadata["number"].(string)
	issueType, _ := ctx.Metadata["type"].(string)

	if org == "" || repo == "" || title == "" || number == "" {
		return ctx.OriginalInput, nil
	}

	// Apply organization/repository mappings if configured
	orgRepo := fmt.Sprintf("%s/%s", org, repo)
	// Try case-insensitive lookup since Viper lowercases map keys
	for key, mapped := range w.config.GitHub.Mappings {
		if strings.EqualFold(key, orgRepo) {
			orgRepo = mapped
			break
		}
	}

	// Build the GitHub URL
	githubURL := fmt.Sprintf("https://github.com/%s/%s/%s/%s", org, repo, issueType, number)
	
	// Create the link text with org/repo#number: title format
	linkText := fmt.Sprintf("%s#%s: %s", orgRepo, number, title)
	return fmt.Sprintf("[%s](%s)", linkText, githubURL), nil
}

func (w *URLWriter) writeJIRAURL(ctx *types.ParseContext) (string, error) {
	issueKey, _ := ctx.Metadata["issue_key"].(string)
	if issueKey == "" {
		return w.writeGenericURL(ctx)
	}

	return fmt.Sprintf("[%s](%s)", issueKey, ctx.OriginalInput), nil
}

func (w *URLWriter) writeJIRACommentURL(ctx *types.ParseContext) (string, error) {
	issueKey, _ := ctx.Metadata["issue_key"].(string)
	if issueKey == "" {
		return w.writeGenericURL(ctx)
	}

	return fmt.Sprintf("[%s comment](%s)", issueKey, ctx.OriginalInput), nil
}

func (w *URLWriter) writeNotionURL(ctx *types.ParseContext) (string, error) {
	title, _ := ctx.Metadata["title"].(string)
	if title == "" {
		return w.writeGenericURL(ctx)
	}

	return fmt.Sprintf("[%s](%s)", title, ctx.OriginalInput), nil
}

func (w *URLWriter) writeGenericURL(ctx *types.ParseContext) (string, error) {
	u, err := url.Parse(ctx.OriginalInput)
	if err != nil {
		return ctx.OriginalInput, nil
	}

	// Extract domain, removing common prefixes
	domain := u.Host
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.TrimPrefix(domain, "ww3.")

	// Check for domain mappings first
	linkText := domain
	if w.config.URL.DomainMappings != nil {
		// Try case-insensitive lookup since Viper lowercases map keys
		for key, mapped := range w.config.URL.DomainMappings {
			if strings.EqualFold(key, domain) {
				linkText = mapped
				break
			}
		}
	}

	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}
