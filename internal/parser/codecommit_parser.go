package parser

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type CodeCommitParser struct {
	config *types.Config
}

func NewCodeCommitParser(cfg *types.Config) *CodeCommitParser {
	return &CodeCommitParser{config: cfg}
}

func (p *CodeCommitParser) CanHandle(input string) bool {
	u, err := url.Parse(input)
	if err != nil {
		return false
	}

	// Check if it's an AWS CodeCommit console URL
	// Format: https://{region}.console.aws.amazon.com/codesuite/codecommit/repositories/{repo}/pull-requests/{number}/details?region={region}
	return strings.Contains(u.Host, "console.aws.amazon.com") &&
		strings.Contains(u.Path, "/codesuite/codecommit/repositories/") &&
		strings.Contains(u.Path, "/pull-requests/")
}

func (p *CodeCommitParser) Parse(input string) (*types.ParseContext, error) {
	if !p.CanHandle(input) {
		return nil, nil
	}

	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  types.ContentTypeCodeCommitURL,
		Confidence:    90,
		Metadata:      make(map[string]interface{}),
	}

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

	return ctx, nil
}
