package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestPassthroughWriter_Vote(t *testing.T) {
	writer := NewPassthroughWriter()

	tests := []struct {
		name        string
		contentType types.ContentType
	}{
		{"GitHub URL", types.ContentTypeGitHubURL},
		{"JIRA Key", types.ContentTypeJIRAKey},
		{"Generic URL", types.ContentTypeURL},
		{"Unknown", types.ContentTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				DetectedType: tt.contentType,
			}
			vote := writer.Vote(ctx)
			// PassthroughWriter always votes 1
			if vote != 1 {
				t.Errorf("Vote() = %v, want 1", vote)
			}
		})
	}
}

func TestPassthroughWriter_Write(t *testing.T) {
	writer := NewPassthroughWriter()

	tests := []struct {
		name          string
		originalInput string
	}{
		{"Plain text", "hello world"},
		{"Unmatched JIRA key", "INVALID-123"},
		{"Random string", "some random text"},
		{"Number", "12345"},
		{"Empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				OriginalInput: tt.originalInput,
			}

			output, err := writer.Write(ctx)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			if output != tt.originalInput {
				t.Errorf("Write() = %v, want %v", output, tt.originalInput)
			}
		})
	}
}

func TestPassthroughWriter_GetName(t *testing.T) {
	writer := NewPassthroughWriter()
	name := writer.GetName()
	
	expected := "PassthroughWriter"
	if name != expected {
		t.Errorf("GetName() = %v, want %v", name, expected)
	}
}