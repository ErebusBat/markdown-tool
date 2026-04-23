package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestURLWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name         string
		contentType  types.ContentType
		expectedVote int
	}{
		{"GitHub URL", types.ContentTypeGitHubURL, 90},
		{"GitHub Long", types.ContentTypeGitHubLong, 95},
		{"JIRA URL", types.ContentTypeJIRAURL, 90},
		{"JIRA Comment", types.ContentTypeJIRAComment, 95},
		{"Jenkins URL", types.ContentTypeJenkinsURL, 90},
		{"YouTube URL", types.ContentTypeYouTubeURL, 95},
		{"CodeCommit URL", types.ContentTypeCodeCommitURL, 90},
		{"CodeCommit Long", types.ContentTypeCodeCommitLong, 95},
		{"Notion URL", types.ContentTypeNotionURL, 85},
		{"Generic URL", types.ContentTypeURL, 50},
		{"Gemini URL", types.ContentTypeGeminiURL, 90},
		{"JIRA Key", types.ContentTypeJIRAKey, 0},
		{"Unknown", types.ContentTypeUnknown, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				DetectedType: tt.contentType,
			}
			vote := writer.Vote(ctx)
			if vote != tt.expectedVote {
				t.Errorf("Vote() = %v, want %v", vote, tt.expectedVote)
			}
		})
	}
}

func TestURLWriter_WriteGitHubURL(t *testing.T) {
	tests := []struct {
		name           string
		config         *types.Config
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name: "GitHub PR without mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "Company-Cam-API",
				"number": "15217",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedOutput: "[CompanyCam/Company-Cam-API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)",
		},
		{
			name: "GitHub PR with mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{
						"companycam/company-cam-api": "CompanyCam/API",
					},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "Company-Cam-API",
				"number": "15217",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
			expectedOutput: "[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)",
		},
		{
			name: "GitHub Issue",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "someorg",
				"repo":   "somerepo",
				"number": "42",
			},
			originalInput:  "https://github.com/someorg/somerepo/issues/42",
			expectedOutput: "[someorg/somerepo#42](https://github.com/someorg/somerepo/issues/42)",
		},
		{
			name: "GitHub Repository without mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":  "pedropark99",
				"repo": "zig-book",
			},
			originalInput:  "https://github.com/pedropark99/zig-book",
			expectedOutput: "[pedropark99/zig-book](https://github.com/pedropark99/zig-book)",
		},
		{
			name: "GitHub Repository with mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{
						"companycam/company-cam-api": "CompanyCam/API",
					},
				},
			},
			metadata: map[string]interface{}{
				"org":  "CompanyCam",
				"repo": "Company-Cam-API",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API",
			expectedOutput: "[CompanyCam/API](https://github.com/CompanyCam/Company-Cam-API)",
		},
		{
			name: "GitHub Commit Long Hash",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "ErebusBat",
				"repo":   "markdown-tool",
				"type":   "commit",
				"number": "aa062a602a02d33f4a6e7880809ac3609fe1417b",
			},
			originalInput:  "https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b",
			expectedOutput: "[ErebusBat/markdown-tool#aa062a6](https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b)",
		},
		{
			name: "GitHub Commit Short Hash",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "Company-Cam-API",
				"type":   "commit",
				"number": "abc123",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API/commit/abc123",
			expectedOutput: "[CompanyCam/Company-Cam-API#abc123](https://github.com/CompanyCam/Company-Cam-API/commit/abc123)",
		},
		{
			name: "GitHub Commit with mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{
						"companycam/company-cam-api": "CompanyCam/API",
					},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "Company-Cam-API",
				"type":   "commit",
				"number": "def456789abcdef123456789abcdef12345",
			},
			originalInput:  "https://github.com/CompanyCam/Company-Cam-API/commit/def456789abcdef123456789abcdef12345",
			expectedOutput: "[CompanyCam/API#def4567](https://github.com/CompanyCam/Company-Cam-API/commit/def456789abcdef123456789abcdef12345)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewURLWriter(tt.config)
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeGitHubURL,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteJIRAURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		contentType    types.ContentType
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name:        "JIRA Issue URL",
			contentType: types.ContentTypeJIRAURL,
			metadata: map[string]interface{}{
				"issue_key": "PLAT-192",
			},
			originalInput:  "https://companycam.atlassian.net/browse/PLAT-192",
			expectedOutput: "[PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)",
		},
		{
			name:        "JIRA Comment URL",
			contentType: types.ContentTypeJIRAComment,
			metadata: map[string]interface{}{
				"issue_key":  "PLAT-192",
				"comment_id": "20266",
			},
			originalInput:  "https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266",
			expectedOutput: "[PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  tt.contentType,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteJenkinsURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name:          "Jenkins build URL",
			originalInput: "https://jenkins.internal.upserve.com/job/app.swipely/114/",
			metadata: map[string]interface{}{
				"job_name":     "app.swipely",
				"build_number": "114",
			},
			expectedOutput: "[jenkins/app.swipely#114](https://jenkins.internal.upserve.com/job/app.swipely/114/)",
		},
		{
			name:          "Jenkins build URL with console text",
			originalInput: "https://jenkins.internal.upserve.com/job/app.swipely/114/consoleText",
			metadata: map[string]interface{}{
				"job_name":     "app.swipely",
				"build_number": "114",
			},
			expectedOutput: "[jenkins/app.swipely#114](https://jenkins.internal.upserve.com/job/app.swipely/114/consoleText)",
		},
		{
			name:          "Jenkins build URL with artifact path",
			originalInput: "https://jenkins.internal.upserve.com/job/my-project/42/artifact/build.log",
			metadata: map[string]interface{}{
				"job_name":     "my-project",
				"build_number": "42",
			},
			expectedOutput: "[jenkins/my-project#42](https://jenkins.internal.upserve.com/job/my-project/42/artifact/build.log)",
		},
		{
			name:          "Jenkins URL with lastBuild",
			originalInput: "https://jenkins.internal.upserve.com/job/app.swipely/lastBuild/",
			metadata: map[string]interface{}{
				"job_name": "app.swipely",
			},
			expectedOutput: "[jenkins/app.swipely](https://jenkins.internal.upserve.com/job/app.swipely/lastBuild/)",
		},
		{
			name:          "Jenkins URL without build identifier",
			originalInput: "https://jenkins.internal.upserve.com/job/my-project/",
			metadata: map[string]interface{}{
				"job_name": "my-project",
			},
			expectedOutput: "[jenkins/my-project](https://jenkins.internal.upserve.com/job/my-project/)",
		},
		{
			name:          "Jenkins URL with lastSuccessfulBuild and console text",
			originalInput: "https://jenkins.internal.upserve.com/job/app.swipely/lastSuccessfulBuild/consoleText",
			metadata: map[string]interface{}{
				"job_name": "app.swipely",
			},
			expectedOutput: "[jenkins/app.swipely](https://jenkins.internal.upserve.com/job/app.swipely/lastSuccessfulBuild/consoleText)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeJenkinsURL,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteYouTubeURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	ctx := &types.ParseContext{
		OriginalInput: "https://www.youtube.com/watch?v=fkT41ooKBuY",
		DetectedType:  types.ContentTypeYouTubeURL,
		Metadata: map[string]interface{}{
			"video_id": "fkT41ooKBuY",
			"title":    "Stop overpaying for OpenAI: Multi-model routing guide",
		},
	}

	expectedOutput := "[📺 Stop overpaying for OpenAI: Multi-model routing guide](https://www.youtube.com/watch?v=fkT41ooKBuY)"

	output, err := writer.Write(ctx)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Write() = %v, want %v", output, expectedOutput)
	}
}

func TestURLWriter_WriteNotionURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	ctx := &types.ParseContext{
		OriginalInput: "https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0",
		DetectedType:  types.ContentTypeNotionURL,
		Metadata: map[string]interface{}{
			"title": "VS Code Setup for Standard rb RubyLSP",
		},
	}

	expectedOutput := "[VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0)"

	output, err := writer.Write(ctx)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if output != expectedOutput {
		t.Errorf("Write() = %v, want %v", output, expectedOutput)
	}
}

func TestURLWriter_WriteGitHubLongURL(t *testing.T) {
	tests := []struct {
		name           string
		config         *types.Config
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name: "GitHub Long without mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "companycam-mobile",
				"title":  "A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog",
				"number": "6549",
				"type":   "issues",
			},
			originalInput:  "GitHub UI text chunk",
			expectedOutput: "[CompanyCam/companycam-mobile#6549: A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog](https://github.com/CompanyCam/companycam-mobile/issues/6549)",
		},
		{
			name: "GitHub Long with mapping",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{
						"companycam/companycam-mobile": "CompanyCam/API",
					},
				},
			},
			metadata: map[string]interface{}{
				"org":    "CompanyCam",
				"repo":   "companycam-mobile",
				"title":  "Fix authentication bug",
				"number": "123",
				"type":   "issues",
			},
			originalInput:  "GitHub UI text chunk",
			expectedOutput: "[CompanyCam/API#123: Fix authentication bug](https://github.com/CompanyCam/companycam-mobile/issues/123)",
		},
		{
			name: "GitHub Long strips JIRA key from title",
			config: &types.Config{
				GitHub: types.GitHubConfig{
					Mappings: map[string]string{},
				},
			},
			metadata: map[string]interface{}{
				"org":    "upserve",
				"repo":   "tokenizer",
				"title":  "[HQ-13237] Update rack to 2.2.22",
				"number": "250",
				"type":   "pull",
			},
			originalInput:  "GitHub UI text chunk",
			expectedOutput: "[upserve/tokenizer#250: Update rack to 2.2.22](https://github.com/upserve/tokenizer/pull/250)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewURLWriter(tt.config)
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeGitHubLong,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteGenericURL(t *testing.T) {
	tests := []struct {
		name           string
		config         *types.Config
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name: "Generic URL with www (no mapping)",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{},
				},
			},
			originalInput: "https://www.example.com/path/to/page",
			metadata: map[string]interface{}{
				"domain": "www.example.com",
			},
			expectedOutput: "[example.com](https://www.example.com/path/to/page)",
		},
		{
			name: "Generic URL with ww3 (no mapping)",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{},
				},
			},
			originalInput: "http://ww3.domain.tld/path/to/document?query=value#anchor",
			metadata: map[string]interface{}{
				"domain": "ww3.domain.tld",
			},
			expectedOutput: "[domain.tld](http://ww3.domain.tld/path/to/document?query=value#anchor)",
		},
		{
			name: "Simple domain (no mapping)",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{},
				},
			},
			originalInput: "https://example.org/page",
			metadata: map[string]interface{}{
				"domain": "example.org",
			},
			expectedOutput: "[example.org](https://example.org/page)",
		},
		{
			name: "Slack URL with domain mapping",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{
						"companycam_slack_com": "slack",
					},
				},
			},
			originalInput: "https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069",
			metadata: map[string]interface{}{
				"domain": "companycam.slack.com",
			},
			expectedOutput: "[slack](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)",
		},
		{
			name: "YouTube URL with domain mapping",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{
						"youtube_com": "YouTube",
					},
				},
			},
			originalInput: "https://youtube.com/watch?v=abc123",
			metadata: map[string]interface{}{
				"domain": "youtube.com",
			},
			expectedOutput: "[YouTube](https://youtube.com/watch?v=abc123)",
		},
		{
			name: "Case-insensitive domain mapping",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{
						"companycam_slack_com": "slack",
					},
				},
			},
			originalInput: "https://CompanyCam.Slack.com/archives/test",
			metadata: map[string]interface{}{
				"domain": "CompanyCam.Slack.com",
			},
			expectedOutput: "[slack](https://CompanyCam.Slack.com/archives/test)",
		},
		{
			name: "Domain with mapping but different domain",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: map[string]string{
						"companycam_slack_com": "slack",
					},
				},
			},
			originalInput: "https://example.com/path",
			metadata: map[string]interface{}{
				"domain": "example.com",
			},
			expectedOutput: "[example.com](https://example.com/path)",
		},
		{
			name: "Domain mapping with nil DomainMappings",
			config: &types.Config{
				URL: types.URLConfig{
					DomainMappings: nil,
				},
			},
			originalInput: "https://example.com/path",
			metadata: map[string]interface{}{
				"domain": "example.com",
			},
			expectedOutput: "[example.com](https://example.com/path)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := NewURLWriter(tt.config)
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeURL,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteCodeCommitURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name:          "CodeCommit PR URL us-east-1",
			originalInput: "https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1",
			metadata: map[string]interface{}{
				"region": "us-east-1",
				"repo":   "upserve-env",
				"number": "411",
			},
			expectedOutput: "[us-east-1/upserve-env#411](https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1)",
		},
		{
			name:          "CodeCommit PR URL us-west-2",
			originalInput: "https://us-west-2.console.aws.amazon.com/codesuite/codecommit/repositories/my-repo/pull-requests/123/details?region=us-west-2",
			metadata: map[string]interface{}{
				"region": "us-west-2",
				"repo":   "my-repo",
				"number": "123",
			},
			expectedOutput: "[us-west-2/my-repo#123](https://us-west-2.console.aws.amazon.com/codesuite/codecommit/repositories/my-repo/pull-requests/123/details?region=us-west-2)",
		},
		{
			name:          "CodeCommit PR URL eu-west-1",
			originalInput: "https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/test-repo/pull-requests/42/details?region=eu-west-1",
			metadata: map[string]interface{}{
				"region": "eu-west-1",
				"repo":   "test-repo",
				"number": "42",
			},
			expectedOutput: "[eu-west-1/test-repo#42](https://eu-west-1.console.aws.amazon.com/codesuite/codecommit/repositories/test-repo/pull-requests/42/details?region=eu-west-1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeCodeCommitURL,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteCodeCommitLongURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name:          "CodeCommit long format",
			originalInput: "AWS Console text chunk",
			metadata: map[string]interface{}{
				"region": "us-east-1",
				"repo":   "upserve-env",
				"number": "411",
				"title":  "SEC-12335: Pass SENDGRID_API_KEY Securley",
			},
			expectedOutput: "[us-east-1/upserve-env#411: SEC-12335: Pass SENDGRID_API_KEY Securley](https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1)",
		},
		{
			name:          "CodeCommit long format different region",
			originalInput: "AWS Console text chunk",
			metadata: map[string]interface{}{
				"region": "us-west-2",
				"repo":   "my-repo",
				"number": "123",
				"title":  "Fix authentication bug",
			},
			expectedOutput: "[us-west-2/my-repo#123: Fix authentication bug](https://us-west-2.console.aws.amazon.com/codesuite/codecommit/repositories/my-repo/pull-requests/123/details?region=us-west-2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeCodeCommitLong,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestURLWriter_WriteGeminiURL(t *testing.T) {
	cfg := &types.Config{}
	writer := NewURLWriter(cfg)

	tests := []struct {
		name           string
		originalInput  string
		metadata       map[string]interface{}
		expectedOutput string
	}{
		{
			name:          "Gemini chat URL",
			originalInput: "https://gemini.google.com/app/ac9ebc9d76c30fc1",
			metadata: map[string]interface{}{
				"chat_id": "ac9ebc9d76c30fc1",
			},
			expectedOutput: "[🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)",
		},
		{
			name:          "Gemini chat URL with different ID",
			originalInput: "https://gemini.google.com/app/abcdef123456",
			metadata: map[string]interface{}{
				"chat_id": "abcdef123456",
			},
			expectedOutput: "[🤖 Gemini Chat](https://gemini.google.com/app/abcdef123456)",
		},
		{
			name:          "Gemini chat URL with trailing arrow uses clean_url",
			originalInput: "https://gemini.google.com/app/ac9ebc9d76c30fc1 →",
			metadata: map[string]interface{}{
				"chat_id":   "ac9ebc9d76c30fc1",
				"clean_url": "https://gemini.google.com/app/ac9ebc9d76c30fc1",
			},
			expectedOutput: "[🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
				DetectedType:  types.ContentTypeGeminiURL,
				Metadata:      tt.metadata,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}
			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}
