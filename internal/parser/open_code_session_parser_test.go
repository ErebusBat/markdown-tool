package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestOpenCodeSessionParser_CanHandle(t *testing.T) {
	cfg := &types.Config{}
	parser := NewOpenCodeSessionParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Exact session token lowercase",
			input:    "ses_36a7950aeffesS4WjOsOMX8XTq",
			expected: true,
		},
		{
			name:     "Exact session token uppercase",
			input:    "SES_ABC123XYZ",
			expected: true,
		},
		{
			name:     "Session token in text",
			input:    "Use ses_abc123 now",
			expected: true,
		},
		{
			name:     "Missing token",
			input:    "no session here",
			expected: false,
		},
		{
			name:     "Token followed by punctuation",
			input:    "ses_abc-123",
			expected: true,
		},
		{
			name:     "Prefix without body",
			input:    "ses_",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.CanHandle(tt.input)
			if result != tt.expected {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestOpenCodeSessionParser_Parse(t *testing.T) {
	cfg := &types.Config{}
	parser := NewOpenCodeSessionParser(cfg)

	tests := []struct {
		name          string
		input         string
		expectedToken string
		expectedType  types.ContentType
		expectedConf  int
		shouldParse   bool
		isExactMatch  bool
	}{
		{
			name:          "Exact token preserves casing",
			input:         "ses_36a7950aeffesS4WjOsOMX8XTq",
			expectedToken: "ses_36a7950aeffesS4WjOsOMX8XTq",
			expectedType:  types.ContentTypeOpenCodeSession,
			expectedConf:  90,
			shouldParse:   true,
			isExactMatch:  true,
		},
		{
			name:          "Mixed casing preserves token",
			input:         "Use SES_ABC123xyz now",
			expectedToken: "SES_ABC123xyz",
			expectedType:  types.ContentTypeOpenCodeSession,
			expectedConf:  70,
			shouldParse:   true,
			isExactMatch:  false,
		},
		{
			name:        "No token",
			input:       "nothing to see",
			shouldParse: false,
		},
		{
			name:          "Token surrounded by punctuation",
			input:         "(ses_ABC123)",
			expectedToken: "ses_ABC123",
			expectedType:  types.ContentTypeOpenCodeSession,
			expectedConf:  70,
			shouldParse:   true,
			isExactMatch:  false,
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

			sessionToken, ok := ctx.Metadata["session_token"].(string)
			if !ok {
				t.Fatalf("Parse(%q) missing session_token metadata", tt.input)
			}

			if sessionToken != tt.expectedToken {
				t.Errorf("Parse(%q).Metadata[session_token] = %v, want %v", tt.input, sessionToken, tt.expectedToken)
			}

			isExactMatch, ok := ctx.Metadata["is_exact_match"].(bool)
			if !ok {
				t.Fatalf("Parse(%q) missing is_exact_match metadata", tt.input)
			}

			if isExactMatch != tt.isExactMatch {
				t.Errorf("Parse(%q).Metadata[is_exact_match] = %v, want %v", tt.input, isExactMatch, tt.isExactMatch)
			}
		})
	}
}
