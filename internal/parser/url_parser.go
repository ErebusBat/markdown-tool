package parser

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type URLParser struct {
	config *types.Config
}

func NewURLParser(cfg *types.Config) *URLParser {
	return &URLParser{config: cfg}
}

func (p *URLParser) CanHandle(input string) bool {
	_, err := url.Parse(input)
	return err == nil && (strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://"))
}

func (p *URLParser) Parse(input string) (*types.ParseContext, error) {
	if !p.CanHandle(input) {
		return nil, nil
	}

	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		Metadata:      make(map[string]interface{}),
	}

	// Detect specific URL types
	switch {
	case p.isGitHubURL(u):
		ctx.DetectedType = types.ContentTypeGitHubURL
		ctx.Confidence = 90
		p.parseGitHubURL(u, ctx)
	case p.isJIRAURL(u):
		if p.isJIRACommentURL(u) {
			ctx.DetectedType = types.ContentTypeJIRAComment
			ctx.Confidence = 95
		} else {
			ctx.DetectedType = types.ContentTypeJIRAURL
			ctx.Confidence = 90
		}
		p.parseJIRAURL(u, ctx)
	case p.isNotionURL(u):
		ctx.DetectedType = types.ContentTypeNotionURL
		ctx.Confidence = 85
		p.parseNotionURL(u, ctx)
	default:
		ctx.DetectedType = types.ContentTypeURL
		ctx.Confidence = 50
		ctx.Metadata["domain"] = u.Host
	}

	return ctx, nil
}

func (p *URLParser) isGitHubURL(u *url.URL) bool {
	return u.Host == "github.com"
}

func (p *URLParser) isJIRAURL(u *url.URL) bool {
	if p.config.JIRA.Domain == "" {
		return false
	}
	jiraURL, _ := url.Parse(p.config.JIRA.Domain)
	return u.Host == jiraURL.Host
}

func (p *URLParser) isJIRACommentURL(u *url.URL) bool {
	return p.isJIRAURL(u) && u.Query().Get("focusedCommentId") != ""
}

func (p *URLParser) isNotionURL(u *url.URL) bool {
	return strings.Contains(u.Host, "notion.so")
}

func (p *URLParser) parseGitHubURL(u *url.URL, ctx *types.ParseContext) {
	// Extract org/repo and issue/PR number from GitHub URLs
	// Path format: /org/repo/pull/123 or /org/repo/issues/123
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) >= 4 {
		org := parts[0]
		repo := parts[1]
		issueType := parts[2] // "pull" or "issues"
		number := parts[3]

		ctx.Metadata["org"] = org
		ctx.Metadata["repo"] = repo
		ctx.Metadata["type"] = issueType
		ctx.Metadata["number"] = number
	}
}

func (p *URLParser) parseJIRAURL(u *url.URL, ctx *types.ParseContext) {
	// Extract JIRA issue key from URL
	// Path format: /browse/PROJ-123
	re := regexp.MustCompile(`/browse/([A-Z]+-\d+)`)
	matches := re.FindStringSubmatch(u.Path)
	if len(matches) > 1 {
		ctx.Metadata["issue_key"] = matches[1]
	}

	// Check if it's a comment URL
	if commentId := u.Query().Get("focusedCommentId"); commentId != "" {
		ctx.Metadata["comment_id"] = commentId
	}
}

func (p *URLParser) parseNotionURL(u *url.URL, ctx *types.ParseContext) {
	// Extract page title from Notion URL slug
	// Format: /workspace/Page-Title-with-Dashes-uuid
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) >= 2 {
		slug := parts[len(parts)-1]
		// Extract title part before the UUID
		re := regexp.MustCompile(`^(.+)-[a-f0-9]{32}`)
		matches := re.FindStringSubmatch(slug)
		if len(matches) > 1 {
			title := strings.ReplaceAll(matches[1], "-", " ")
			ctx.Metadata["title"] = title
		}
	}
}
