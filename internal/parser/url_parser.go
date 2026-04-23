package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	case p.isJenkinsURL(u):
		ctx.DetectedType = types.ContentTypeJenkinsURL
		ctx.Confidence = 90
		p.parseJenkinsURL(u, ctx)
	case p.isYouTubeURL(u):
		ctx.DetectedType = types.ContentTypeYouTubeURL
		ctx.Confidence = 90
		p.parseYouTubeURL(u, ctx)
	case p.isCodeCommitURL(u):
		ctx.DetectedType = types.ContentTypeCodeCommitURL
		ctx.Confidence = 90
		p.parseCodeCommitURL(u, ctx)
	case p.isNotionURL(u):
		ctx.DetectedType = types.ContentTypeNotionURL
		ctx.Confidence = 85
		p.parseNotionURL(u, ctx)
	case p.isMiniMaxURL(u):
		ctx.DetectedType = types.ContentTypeMiniMaxURL
		ctx.Confidence = 90
		p.parseMiniMaxURL(u, ctx)
	case p.isGeminiURL(u):
		ctx.DetectedType = types.ContentTypeGeminiURL
		ctx.Confidence = 90
		p.parseGeminiURL(u, ctx)
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

func (p *URLParser) isJenkinsURL(u *url.URL) bool {
	if p.config.Jenkins.Domain == "" {
		return false
	}
	jenkinsURL, _ := url.Parse(p.config.Jenkins.Domain)
	return u.Host == jenkinsURL.Host
}

func (p *URLParser) isYouTubeURL(u *url.URL) bool {
	return u.Host == "www.youtube.com" || u.Host == "youtube.com" || u.Host == "youtu.be" || u.Host == "m.youtube.com"
}

func (p *URLParser) isCodeCommitURL(u *url.URL) bool {
	return strings.Contains(u.Host, "console.aws.amazon.com") &&
		strings.Contains(u.Path, "/codesuite/codecommit/repositories/") &&
		strings.Contains(u.Path, "/pull-requests/")
}

func (p *URLParser) isNotionURL(u *url.URL) bool {
	return strings.Contains(u.Host, "notion.so")
}

func (p *URLParser) isMiniMaxURL(u *url.URL) bool {
	return u.Host == "agent.minimax.io"
}

func (p *URLParser) isGeminiURL(u *url.URL) bool {
	return u.Host == "gemini.google.com" && strings.HasPrefix(u.Path, "/app/")
}

func (p *URLParser) parseGeminiURL(u *url.URL, ctx *types.ParseContext) {
	re := regexp.MustCompile(`/app/([a-f0-9]+)`)
	matches := re.FindStringSubmatch(u.Path)
	if len(matches) > 1 {
		ctx.Metadata["chat_id"] = matches[1]
		cleanURL := &url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
			Path:   "/app/" + matches[1],
		}
		ctx.Metadata["clean_url"] = cleanURL.String()
	}
}

func (p *URLParser) parseMiniMaxURL(u *url.URL, ctx *types.ParseContext) {
	chatID := u.Query().Get("id")
	if chatID != "" {
		ctx.Metadata["chat_id"] = chatID
	}
}

func (p *URLParser) parseGitHubURL(u *url.URL, ctx *types.ParseContext) {
	// Extract org/repo and optionally issue/PR number from GitHub URLs
	// Path formats:
	// - /org/repo (simple repository URL)
	// - /org/repo/pull/123 or /org/repo/issues/123 (issue/PR URLs)
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) >= 2 {
		org := parts[0]
		repo := parts[1]

		ctx.Metadata["org"] = org
		ctx.Metadata["repo"] = repo

		// If there are 4+ parts, extract issue/PR information
		if len(parts) >= 4 {
			issueType := parts[2] // "pull" or "issues"
			number := parts[3]

			ctx.Metadata["type"] = issueType
			ctx.Metadata["number"] = number
		}
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

func (p *URLParser) parseJenkinsURL(u *url.URL, ctx *types.ParseContext) {
	// Extract job name and build number from Jenkins URLs
	// Path formats:
	// - /job/{job-name}/{build-number}/[optional-path]
	// - /job/{job-name}/lastBuild/
	// - /job/{job-name}/
	// Example: /job/app.swipely/114/consoleText

	// Try to match job name with numeric build number
	re := regexp.MustCompile(`^/job/([^/]+)/(\d+)`)
	matches := re.FindStringSubmatch(u.Path)
	if len(matches) > 2 {
		ctx.Metadata["job_name"] = matches[1]
		ctx.Metadata["build_number"] = matches[2]
		return
	}

	// If no numeric build number, just extract job name
	// This handles cases like /job/{job-name}/lastBuild/ or /job/{job-name}/
	re = regexp.MustCompile(`^/job/([^/]+)`)
	matches = re.FindStringSubmatch(u.Path)
	if len(matches) > 1 {
		ctx.Metadata["job_name"] = matches[1]
		// Don't set build_number - it will be empty/nil
	}
}

func (p *URLParser) parseYouTubeURL(u *url.URL, ctx *types.ParseContext) {
	// Extract video ID from YouTube URLs
	// Formats:
	// - https://www.youtube.com/watch?v=VIDEO_ID
	// - https://youtu.be/VIDEO_ID
	// - https://www.youtube.com/embed/VIDEO_ID

	var videoID string

	if u.Host == "youtu.be" {
		// Short format: https://youtu.be/VIDEO_ID
		videoID = strings.TrimPrefix(u.Path, "/")
	} else {
		// Standard format: https://www.youtube.com/watch?v=VIDEO_ID
		videoID = u.Query().Get("v")
	}

	if videoID == "" {
		return
	}

	ctx.Metadata["video_id"] = videoID

	// Fetch video title using YouTube oEmbed API (no API key required)
	title := p.fetchYouTubeTitle(videoID)
	if title != "" {
		ctx.Metadata["title"] = title
	}
}

func (p *URLParser) fetchYouTubeTitle(videoID string) string {
	// Use YouTube oEmbed API to get video title
	oembedURL := fmt.Sprintf("https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=%s&format=json", videoID)

	resp, err := http.Get(oembedURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var result struct {
		Title string `json:"title"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}

	return result.Title
}

func (p *URLParser) parseCodeCommitURL(u *url.URL, ctx *types.ParseContext) {
	// Extract region from subdomain (e.g., us-east-1.console.aws.amazon.com)
	parts := strings.Split(u.Host, ".")
	if len(parts) > 0 {
		ctx.Metadata["region"] = parts[0]
	}

	// Extract repository name and PR number from path
	// Path format: /codesuite/codecommit/repositories/{repo}/pull-requests/{number}/details
	re := regexp.MustCompile(`/repositories/([^/]+)/pull-requests/(\d+)`)
	matches := re.FindStringSubmatch(u.Path)
	if len(matches) > 2 {
		ctx.Metadata["repo"] = matches[1]
		ctx.Metadata["number"] = matches[2]
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
