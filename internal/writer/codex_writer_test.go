package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestCodexWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	w := NewCodexWriter(cfg)

	tests := []struct {
		name         string
		ctx          *types.ParseContext
		expectedVote int
	}{
		{
			name: "Codex thread",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeCodexThread,
				Confidence:   90,
			},
			expectedVote: 90,
		},
		{
			name: "Non-Codex content",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeURL,
				Confidence:   80,
			},
			expectedVote: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vote := w.Vote(tt.ctx)
			if vote != tt.expectedVote {
				t.Errorf("Vote() = %d, want %d", vote, tt.expectedVote)
			}
		})
	}
}

func TestCodexWriter_Write(t *testing.T) {
	cfg := &types.Config{}
	w := NewCodexWriter(cfg)

	tests := []struct {
		name           string
		ctx            *types.ParseContext
		expectedOutput string
	}{
		{
			name: "Codex thread URI",
			ctx: &types.ParseContext{
				OriginalInput: "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
				DetectedType:  types.ContentTypeCodexThread,
				Metadata: map[string]interface{}{
					"thread_id": "019dcc44-e7b8-7c23-816c-34c194bdb3cf",
					"url":       "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
				},
			},
			expectedOutput: "[🤖 Codex](codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf)",
		},
		{
			name: "Missing url metadata",
			ctx: &types.ParseContext{
				OriginalInput: "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
				DetectedType:  types.ContentTypeCodexThread,
				Metadata:      map[string]interface{}{},
			},
			expectedOutput: "codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf",
		},
		{
			name: "Non-Codex content",
			ctx: &types.ParseContext{
				OriginalInput: "something else",
				DetectedType:  types.ContentTypeUnknown,
			},
			expectedOutput: "something else",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := w.Write(tt.ctx)
			if err != nil {
				t.Fatalf("Write() returned error: %v", err)
			}
			if output != tt.expectedOutput {
				t.Errorf("Write() = %q, want %q", output, tt.expectedOutput)
			}
		})
	}
}

func TestCodexWriter_GetName(t *testing.T) {
	cfg := &types.Config{}
	w := NewCodexWriter(cfg)

	expectedName := "CodexWriter"
	name := w.GetName()
	if name != expectedName {
		t.Errorf("GetName() = %q, want %q", name, expectedName)
	}
}
