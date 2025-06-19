package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestVote(t *testing.T) {
	cfg := &types.Config{
		JIRA: types.JIRAConfig{
			Domain: "https://companycam.atlassian.net",
		},
	}

	writers := GetWriters(cfg)

	tests := []struct {
		name               string
		contexts           []*types.ParseContext
		expectedWriterName string
		expectedScore      int
	}{
		{
			name: "GitHub URL should select URLWriter",
			contexts: []*types.ParseContext{
				{
					OriginalInput: "https://github.com/CompanyCam/Company-Cam-API/pull/15217",
					DetectedType:  types.ContentTypeGitHubURL,
					Confidence:    90,
					Metadata: map[string]interface{}{
						"org":    "CompanyCam",
						"repo":   "Company-Cam-API",
						"number": "15217",
					},
				},
			},
			expectedWriterName: "URLWriter",
			expectedScore:      90,
		},
		{
			name: "JIRA Key should select JIRAWriter",
			contexts: []*types.ParseContext{
				{
					OriginalInput: "PLAT-12345",
					DetectedType:  types.ContentTypeJIRAKey,
					Confidence:    95,
					Metadata: map[string]interface{}{
						"issue_key": "PLAT-12345",
						"project":   "PLAT",
					},
				},
			},
			expectedWriterName: "JIRAWriter",
			expectedScore:      95,
		},
		{
			name: "JIRA Comment URL should select URLWriter with high score",
			contexts: []*types.ParseContext{
				{
					OriginalInput: "https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266",
					DetectedType:  types.ContentTypeJIRAComment,
					Confidence:    95,
					Metadata: map[string]interface{}{
						"issue_key":  "PLAT-192",
						"comment_id": "20266",
					},
				},
			},
			expectedWriterName: "URLWriter",
			expectedScore:      95,
		},
		{
			name: "Unknown content should select PassthroughWriter",
			contexts: []*types.ParseContext{
				{
					OriginalInput: "some random text",
					DetectedType:  types.ContentTypeUnknown,
					Confidence:    0,
					Metadata:      map[string]interface{}{},
				},
			},
			expectedWriterName: "PassthroughWriter",
			expectedScore:      1,
		},
		{
			name:               "Empty contexts should return nil",
			contexts:           []*types.ParseContext{},
			expectedWriterName: "",
			expectedScore:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bestWriter, bestScore := Vote(writers, tt.contexts)

			if tt.expectedScore == 0 {
				if bestWriter != nil {
					t.Errorf("Expected nil writer, got %v", bestWriter.GetName())
				}
				if bestScore != 0 {
					t.Errorf("Expected score 0, got %v", bestScore)
				}
				return
			}

			if bestWriter == nil {
				t.Fatal("Expected non-nil writer")
			}

			if bestWriter.GetName() != tt.expectedWriterName {
				t.Errorf("Expected writer %v, got %v", tt.expectedWriterName, bestWriter.GetName())
			}

			if bestScore != tt.expectedScore {
				t.Errorf("Expected score %v, got %v", tt.expectedScore, bestScore)
			}
		})
	}
}

func TestGetWriters(t *testing.T) {
	cfg := &types.Config{}
	writers := GetWriters(cfg)

	if len(writers) != 3 {
		t.Errorf("Expected 3 writers, got %v", len(writers))
	}

	expectedNames := []string{"URLWriter", "JIRAWriter", "PassthroughWriter"}
	for i, writer := range writers {
		if writer.GetName() != expectedNames[i] {
			t.Errorf("Expected writer %v at index %v, got %v", expectedNames[i], i, writer.GetName())
		}
	}
}