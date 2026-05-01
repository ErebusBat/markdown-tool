package writer

import (
	"fmt"
	"net/url"
	"regexp"
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
	case types.ContentTypeJenkinsURL:
		return 90
	case types.ContentTypeYouTubeURL:
		return 95
	case types.ContentTypeCodeCommitURL:
		return 90
	case types.ContentTypeCodeCommitLong:
		return 95
	case types.ContentTypeNotionURL:
		return 85
	case types.ContentTypeMiniMaxURL:
		return 90
	case types.ContentTypeGeminiURL:
		return 90
	case types.ContentTypeCircleCI:
		return 90
	case types.ContentTypeChatGPT:
		return 90
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
	case types.ContentTypeJenkinsURL:
		return w.writeJenkinsURL(ctx)
	case types.ContentTypeYouTubeURL:
		return w.writeYouTubeURL(ctx)
	case types.ContentTypeCodeCommitURL:
		return w.writeCodeCommitURL(ctx)
	case types.ContentTypeCodeCommitLong:
		return w.writeCodeCommitLongURL(ctx)
	case types.ContentTypeNotionURL:
		return w.writeNotionURL(ctx)
	case types.ContentTypeMiniMaxURL:
		return w.writeMiniMaxURL(ctx)
	case types.ContentTypeGeminiURL:
		return w.writeGeminiURL(ctx)
	case types.ContentTypeCircleCI:
		return w.writeCircleCIURL(ctx)
	case types.ContentTypeChatGPT:
		return w.writeChatGPTURL(ctx)
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

	if org == "" || repo == "" {
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

	// If there's an issue/PR/commit number, format as org/repo#number
	// For commits, truncate hash to 7 characters in link text
	// Otherwise, format as org/repo for simple repository URLs
	var linkText string
	if number != "" {
		issueType, _ := ctx.Metadata["type"].(string)
		if issueType == "commit" && len(number) > 7 {
			// Truncate commit hash to 7 characters for display
			linkText = fmt.Sprintf("%s#%s", orgRepo, number[:7])
		} else {
			linkText = fmt.Sprintf("%s#%s", orgRepo, number)
		}
	} else {
		linkText = orgRepo
	}

	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeGitHubLongURL(ctx *types.ParseContext) (string, error) {
	org, _ := ctx.Metadata["org"].(string)
	repo, _ := ctx.Metadata["repo"].(string)
	title, _ := ctx.Metadata["title"].(string)
	number, _ := ctx.Metadata["number"].(string)
	issueType, _ := ctx.Metadata["type"].(string)

	title = stripLeadingJiraKey(title)

	if issueType == "" {
		issueType = "issues"
	}

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

func (w *URLWriter) writeJenkinsURL(ctx *types.ParseContext) (string, error) {
	jobName, _ := ctx.Metadata["job_name"].(string)
	buildNumber, _ := ctx.Metadata["build_number"].(string)

	if jobName == "" {
		return w.writeGenericURL(ctx)
	}

	// Format: jenkins/{job_name}#{build_number} or jenkins/{job_name} (if no build number)
	var linkText string
	if buildNumber != "" {
		linkText = fmt.Sprintf("jenkins/%s#%s", jobName, buildNumber)
	} else {
		linkText = fmt.Sprintf("jenkins/%s", jobName)
	}
	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeYouTubeURL(ctx *types.ParseContext) (string, error) {
	title, _ := ctx.Metadata["title"].(string)

	if title == "" {
		// Fallback to generic URL if title fetch failed
		return w.writeGenericURL(ctx)
	}

	youtubeType, _ := ctx.Metadata["youtube_type"].(string)
	icon := "📺"
	if youtubeType == "playlist" {
		icon = "🎥🗃️"
	}

	linkText := fmt.Sprintf("%s %s", icon, title)
	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeNotionURL(ctx *types.ParseContext) (string, error) {
	title, _ := ctx.Metadata["title"].(string)
	if title == "" {
		return w.writeGenericURL(ctx)
	}

	return fmt.Sprintf("[%s](%s)", title, ctx.OriginalInput), nil
}

func (w *URLWriter) writeMiniMaxURL(ctx *types.ParseContext) (string, error) {
	linkText := "🤖 MiniMax.io"
	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeGeminiURL(ctx *types.ParseContext) (string, error) {
	cleanURL, _ := ctx.Metadata["clean_url"].(string)
	if cleanURL == "" {
		cleanURL = ctx.OriginalInput
	}
	return fmt.Sprintf("[🤖 Gemini Chat](%s)", cleanURL), nil
}

func (w *URLWriter) writeCircleCIURL(ctx *types.ParseContext) (string, error) {
	org, _ := ctx.Metadata["org"].(string)
	repo, _ := ctx.Metadata["repo"].(string)
	pipelineNumber, _ := ctx.Metadata["pipeline_number"].(string)

	if org == "" || repo == "" || pipelineNumber == "" {
		return w.writeGenericURL(ctx)
	}

	linkText := fmt.Sprintf("🏗️ CircleCI %s/%s#%s", org, repo, pipelineNumber)
	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeChatGPTURL(ctx *types.ParseContext) (string, error) {
	return fmt.Sprintf("[🤖 ChatGPT](%s)", ctx.OriginalInput), nil
}

func (w *URLWriter) writeCodeCommitURL(ctx *types.ParseContext) (string, error) {
	region, _ := ctx.Metadata["region"].(string)
	repo, _ := ctx.Metadata["repo"].(string)
	number, _ := ctx.Metadata["number"].(string)

	if region == "" || repo == "" || number == "" {
		return w.writeGenericURL(ctx)
	}

	// Format: [region/repo#number](URL)
	linkText := fmt.Sprintf("%s/%s#%s", region, repo, number)
	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

func (w *URLWriter) writeCodeCommitLongURL(ctx *types.ParseContext) (string, error) {
	region, _ := ctx.Metadata["region"].(string)
	repo, _ := ctx.Metadata["repo"].(string)
	number, _ := ctx.Metadata["number"].(string)
	title, _ := ctx.Metadata["title"].(string)

	if region == "" || repo == "" || number == "" || title == "" {
		return ctx.OriginalInput, nil
	}

	// Build the CodeCommit URL
	codecommitURL := fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codecommit/repositories/%s/pull-requests/%s/details?region=%s",
		region, repo, number, region)

	// Format: [region/repo#number: title](URL)
	linkText := fmt.Sprintf("%s/%s#%s: %s", region, repo, number, title)
	return fmt.Sprintf("[%s](%s)", linkText, codecommitURL), nil
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
		// Convert domain to underscore format for lookup (e.g., mail.google.com -> mail_google_com)
		domainKey := strings.ReplaceAll(domain, ".", "_")

		// Try case-insensitive lookup since Viper lowercases map keys
		for key, mapped := range w.config.URL.DomainMappings {
			if strings.EqualFold(key, domainKey) {
				linkText = mapped
				break
			}
		}
	}

	return fmt.Sprintf("[%s](%s)", linkText, ctx.OriginalInput), nil
}

var leadingJiraKeyRegex = regexp.MustCompile(`^\s*(\[[A-Z][A-Z0-9]+-\d+\]\s*|[A-Z][A-Z0-9]+-\d+:\s*)`)

func stripLeadingJiraKey(title string) string {
	cleaned := leadingJiraKeyRegex.ReplaceAllString(title, "")
	return strings.TrimSpace(cleaned)
}
