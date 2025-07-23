package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestRaycastParser_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Raycast AI Chat URI", "raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D", true},
		{"Raycast generic extension URI", "raycast://extensions/other/extension", true},
		{"Raycast simple URI", "raycast://settings", true},
		{"Raycast with path", "raycast://extensions/raycast/window-management/center", true},
		{"Non-Raycast URI", "https://example.com", false},
		{"HTTP URL", "http://raycast.com", false},
		{"Plain text", "raycast is cool", false},
		{"Empty string", "", false},
		{"Just raycast://", "raycast://", true},
		{"Invalid URI after raycast://", "raycast://[invalid", false},
	}

	cfg := &types.Config{}
	parser := NewRaycastParser(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRaycastParser_Parse(t *testing.T) {
	cfg := &types.Config{}
	parser := NewRaycastParser(cfg)

	tests := []struct {
		name           string
		input          string
		expectedType   types.ContentType
		expectedConf   int
		expectedAIChat bool
		shouldParse    bool
	}{
		{
			name:           "Raycast AI Chat URI",
			input:          "raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D",
			expectedType:   types.ContentTypeRaycastURI,
			expectedConf:   85,
			expectedAIChat: true,
			shouldParse:    true,
		},
		{
			name:           "Raycast AI Chat URI without query",
			input:          "raycast://extensions/raycast/raycast-ai/ai-chat",
			expectedType:   types.ContentTypeRaycastURI,
			expectedConf:   85,
			expectedAIChat: true,
			shouldParse:    true,
		},
		{
			name:           "Raycast generic extension URI",
			input:          "raycast://extensions/other/extension",
			expectedType:   types.ContentTypeRaycastURI,
			expectedConf:   85,
			expectedAIChat: false,
			shouldParse:    true,
		},
		{
			name:           "Raycast settings URI",
			input:          "raycast://settings",
			expectedType:   types.ContentTypeRaycastURI,
			expectedConf:   85,
			expectedAIChat: false,
			shouldParse:    true,
		},
		{
			name:           "Non-Raycast URI",
			input:          "https://example.com",
			shouldParse:    false,
		},
		{
			name:           "Plain text",
			input:          "not a URI",
			shouldParse:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.input, err)
			}

			if !tt.shouldParse {
				if ctx != nil {
					t.Errorf("Parse(%q) = %v, want nil", tt.input, ctx)
				}
				return
			}

			if ctx == nil {
				t.Fatalf("Parse(%q) = nil, want non-nil", tt.input)
			}

			if ctx.DetectedType != tt.expectedType {
				t.Errorf("Parse(%q).DetectedType = %v, want %v", tt.input, ctx.DetectedType, tt.expectedType)
			}

			if ctx.Confidence != tt.expectedConf {
				t.Errorf("Parse(%q).Confidence = %v, want %v", tt.input, ctx.Confidence, tt.expectedConf)
			}

			if ctx.OriginalInput != tt.input {
				t.Errorf("Parse(%q).OriginalInput = %v, want %v", tt.input, ctx.OriginalInput, tt.input)
			}

			isAIChat, ok := ctx.Metadata["isAIChat"].(bool)
			if !ok {
				t.Errorf("Parse(%q).Metadata[\"isAIChat\"] not found or not bool", tt.input)
			} else if isAIChat != tt.expectedAIChat {
				t.Errorf("Parse(%q).Metadata[\"isAIChat\"] = %v, want %v", tt.input, isAIChat, tt.expectedAIChat)
			}
		})
	}
}