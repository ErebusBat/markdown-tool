package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestRaycastWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewRaycastWriter(cfg)

	tests := []struct {
		name         string
		ctx          *types.ParseContext
		expectedVote int
	}{
		{
			name: "Raycast URI",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeRaycastURI,
			},
			expectedVote: 85,
		},
		{
			name: "Non-Raycast content",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeURL,
			},
			expectedVote: 0,
		},
		{
			name: "JIRA content",
			ctx: &types.ParseContext{
				DetectedType: types.ContentTypeJIRAKey,
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

func TestRaycastWriter_Write(t *testing.T) {
	cfg := &types.Config{}
	writer := NewRaycastWriter(cfg)

	tests := []struct {
		name           string
		ctx            *types.ParseContext
		expectedOutput string
	}{
		{
			name: "Raycast AI Chat URI",
			ctx: &types.ParseContext{
				OriginalInput: "raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D",
				DetectedType:  types.ContentTypeRaycastURI,
				Metadata: map[string]interface{}{
					"isAIChat": true,
					"isNote":   false,
				},
			},
			expectedOutput: "[Raycast AI](raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D)",
		},
		{
			name: "Raycast Note URI",
			ctx: &types.ParseContext{
				OriginalInput: "raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22:%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D",
				DetectedType:  types.ContentTypeRaycastURI,
				Metadata: map[string]interface{}{
					"isAIChat": false,
					"isNote":   true,
				},
			},
			expectedOutput: "[Raycast Note](raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22:%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D)",
		},
		{
			name: "Raycast generic URI",
			ctx: &types.ParseContext{
				OriginalInput: "raycast://extensions/other/extension",
				DetectedType:  types.ContentTypeRaycastURI,
				Metadata: map[string]interface{}{
					"isAIChat": false,
					"isNote":   false,
				},
			},
			expectedOutput: "[Raycast](raycast://extensions/other/extension)",
		},
		{
			name: "Raycast settings URI",
			ctx: &types.ParseContext{
				OriginalInput: "raycast://settings",
				DetectedType:  types.ContentTypeRaycastURI,
				Metadata: map[string]interface{}{
					"isAIChat": false,
					"isNote":   false,
				},
			},
			expectedOutput: "[Raycast](raycast://settings)",
		},
		{
			name: "Non-Raycast content returns original",
			ctx: &types.ParseContext{
				OriginalInput: "not a raycast URI",
				DetectedType:  types.ContentTypeUnknown,
				Metadata:      map[string]interface{}{},
			},
			expectedOutput: "not a raycast URI",
		},
		{
			name: "Missing isAIChat metadata defaults to generic",
			ctx: &types.ParseContext{
				OriginalInput: "raycast://extensions/something",
				DetectedType:  types.ContentTypeRaycastURI,
				Metadata:      map[string]interface{}{},
			},
			expectedOutput: "[Raycast](raycast://extensions/something)",
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

func TestRaycastWriter_GetName(t *testing.T) {
	cfg := &types.Config{}
	writer := NewRaycastWriter(cfg)

	expectedName := "RaycastWriter"
	name := writer.GetName()
	if name != expectedName {
		t.Errorf("GetName() = %q, want %q", name, expectedName)
	}
}