package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestOpenCodeSessionWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewOpenCodeSessionWriter(cfg)

	tests := []struct {
		name         string
		ctx          *types.ParseContext
		expectedVote int
	}{
		{
			name: "OpenCode session",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeOpenCodeSession,
				Confidence:   90,
			},
			expectedVote: 90,
		},
		{
			name: "Non-OpenCode content",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeURL,
			},
			expectedVote: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vote := writer.Vote(tt.ctx)
			if vote != tt.expectedVote {
				t.Errorf("Vote() = %d, want %d", vote, tt.expectedVote)
			}
		})
	}
}

func TestOpenCodeSessionWriter_Write(t *testing.T) {
	cfg := &types.Config{}
	writer := NewOpenCodeSessionWriter(cfg)

	tests := []struct {
		name           string
		ctx            *types.ParseContext
		expectedOutput string
	}{
		{
			name: "OpenCode session token",
			ctx: &types.ParseContext{
				OriginalInput: "ses_36a7950aeffesS4WjOsOMX8XTq",
				DetectedType:  types.ContentTypeOpenCodeSession,
				Metadata: map[string]interface{}{
					"session_token": "ses_36a7950aeffesS4WjOsOMX8XTq",
				},
			},
			expectedOutput: "[🤖OpenCode](opencode://session/ses_36a7950aeffesS4WjOsOMX8XTq)",
		},
		{
			name: "Missing session token",
			ctx: &types.ParseContext{
				OriginalInput: "ses_abc123",
				DetectedType:  types.ContentTypeOpenCodeSession,
				Metadata:      map[string]interface{}{},
			},
			expectedOutput: "ses_abc123",
		},
		{
			name: "Non-OpenCode content",
			ctx: &types.ParseContext{
				OriginalInput: "no session",
				DetectedType:  types.ContentTypeUnknown,
			},
			expectedOutput: "no session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := writer.Write(tt.ctx)
			if err != nil {
				t.Fatalf("Write() returned error: %v", err)
			}
			if output != tt.expectedOutput {
				t.Errorf("Write() = %q, want %q", output, tt.expectedOutput)
			}
		})
	}
}

func TestOpenCodeSessionWriter_GetName(t *testing.T) {
	cfg := &types.Config{}
	writer := NewOpenCodeSessionWriter(cfg)

	expectedName := "OpenCodeSessionWriter"
	name := writer.GetName()
	if name != expectedName {
		t.Errorf("GetName() = %q, want %q", name, expectedName)
	}
}
