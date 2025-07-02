package writer

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestPhoneWriter_Vote(t *testing.T) {
	cfg := &types.Config{}
	writer := NewPhoneWriter(cfg)

	tests := []struct {
		name         string
		contentType  types.ContentType
		confidence   int
		expectedVote int
	}{
		{"7-digit phone high confidence", types.ContentTypePhone7Digit, 95, 95},
		{"7-digit phone low confidence", types.ContentTypePhone7Digit, 60, 60},
		{"10-digit phone high confidence", types.ContentTypePhone10Digit, 95, 95},
		{"10-digit phone low confidence", types.ContentTypePhone10Digit, 70, 70},
		{"11-digit phone high confidence", types.ContentTypePhone11Digit, 95, 95},
		{"11-digit phone low confidence", types.ContentTypePhone11Digit, 65, 65},
		{"GitHub URL", types.ContentTypeGitHubURL, 90, 0},
		{"JIRA Key", types.ContentTypeJIRAKey, 95, 0},
		{"Generic URL", types.ContentTypeURL, 50, 0},
		{"Unknown", types.ContentTypeUnknown, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &types.ParseContext{
				DetectedType: tt.contentType,
				Confidence:   tt.confidence,
			}
			vote := writer.Vote(ctx)
			if vote != tt.expectedVote {
				t.Errorf("Vote() = %v, want %v", vote, tt.expectedVote)
			}
		})
	}
}

func TestPhoneWriter_Write(t *testing.T) {
	cfg := &types.Config{}
	writer := NewPhoneWriter(cfg)

	tests := []struct {
		name           string
		contentType    types.ContentType
		metadata       map[string]interface{}
		originalInput  string
		expectedOutput string
	}{
		{
			name:        "7-digit phone",
			contentType: types.ContentTypePhone7Digit,
			metadata: map[string]interface{}{
				"formatted_display": "123-4567",
				"tel_url":          "1234567",
			},
			originalInput:  "1234567",
			expectedOutput: "[123-4567](tel:1234567)",
		},
		{
			name:        "7-digit phone with dashes",
			contentType: types.ContentTypePhone7Digit,
			metadata: map[string]interface{}{
				"formatted_display": "123-4567",
				"tel_url":          "1234567",
			},
			originalInput:  "123-4567",
			expectedOutput: "[123-4567](tel:1234567)",
		},
		{
			name:        "7-digit phone with dots",
			contentType: types.ContentTypePhone7Digit,
			metadata: map[string]interface{}{
				"formatted_display": "123-4567",
				"tel_url":          "1234567",
			},
			originalInput:  "123.4567",
			expectedOutput: "[123-4567](tel:1234567)",
		},
		{
			name:        "10-digit phone",
			contentType: types.ContentTypePhone10Digit,
			metadata: map[string]interface{}{
				"formatted_display": "890-123-4567",
				"tel_url":          "8901234567",
			},
			originalInput:  "8901234567",
			expectedOutput: "[890-123-4567](tel:8901234567)",
		},
		{
			name:        "10-digit phone with dashes",
			contentType: types.ContentTypePhone10Digit,
			metadata: map[string]interface{}{
				"formatted_display": "890-123-4567",
				"tel_url":          "8901234567",
			},
			originalInput:  "890-123-4567",
			expectedOutput: "[890-123-4567](tel:8901234567)",
		},
		{
			name:        "10-digit phone with dots",
			contentType: types.ContentTypePhone10Digit,
			metadata: map[string]interface{}{
				"formatted_display": "890-123-4567",
				"tel_url":          "8901234567",
			},
			originalInput:  "890.123.4567",
			expectedOutput: "[890-123-4567](tel:8901234567)",
		},
		{
			name:        "10-digit phone with parentheses and space",
			contentType: types.ContentTypePhone10Digit,
			metadata: map[string]interface{}{
				"formatted_display": "890-123-4567",
				"tel_url":          "8901234567",
			},
			originalInput:  "(890) 123-4567",
			expectedOutput: "[890-123-4567](tel:8901234567)",
		},
		{
			name:        "10-digit phone with parentheses no space",
			contentType: types.ContentTypePhone10Digit,
			metadata: map[string]interface{}{
				"formatted_display": "890-123-4567",
				"tel_url":          "8901234567",
			},
			originalInput:  "(890)123-4567",
			expectedOutput: "[890-123-4567](tel:8901234567)",
		},
		{
			name:        "10-digit phone with parentheses plain",
			contentType: types.ContentTypePhone10Digit,
			metadata: map[string]interface{}{
				"formatted_display": "890-123-4567",
				"tel_url":          "8901234567",
			},
			originalInput:  "(890)1234567",
			expectedOutput: "[890-123-4567](tel:8901234567)",
		},
		{
			name:        "11-digit US phone",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "1-890-123-4567",
				"tel_url":          "+18901234567",
			},
			originalInput:  "18901234567",
			expectedOutput: "[1-890-123-4567](tel:+18901234567)",
		},
		{
			name:        "11-digit US phone with dashes",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "1-890-123-4567",
				"tel_url":          "+18901234567",
			},
			originalInput:  "1-890-123-4567",
			expectedOutput: "[1-890-123-4567](tel:+18901234567)",
		},
		{
			name:        "11-digit US phone with dots",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "1-890-123-4567",
				"tel_url":          "+18901234567",
			},
			originalInput:  "1.890.123.4567",
			expectedOutput: "[1-890-123-4567](tel:+18901234567)",
		},
		{
			name:        "11-digit US phone with parentheses and space",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "1-890-123-4567",
				"tel_url":          "+18901234567",
			},
			originalInput:  "1 (890) 123-4567",
			expectedOutput: "[1-890-123-4567](tel:+18901234567)",
		},
		{
			name:        "11-digit US phone with parentheses no space",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "1-890-123-4567",
				"tel_url":          "+18901234567",
			},
			originalInput:  "1(890)123-4567",
			expectedOutput: "[1-890-123-4567](tel:+18901234567)",
		},
		{
			name:        "11-digit US phone with parentheses plain",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "1-890-123-4567",
				"tel_url":          "+18901234567",
			},
			originalInput:  "1(890)1234567",
			expectedOutput: "[1-890-123-4567](tel:+18901234567)",
		},
		{
			name:        "11-digit international phone",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "+7-890-123-4567",
				"tel_url":          "+78901234567",
			},
			originalInput:  "+78901234567",
			expectedOutput: "[+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:        "11-digit international phone with dashes",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "+7-890-123-4567",
				"tel_url":          "+78901234567",
			},
			originalInput:  "+7-890-123-4567",
			expectedOutput: "[+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:        "11-digit international phone with dots",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "+7-890-123-4567",
				"tel_url":          "+78901234567",
			},
			originalInput:  "+7.890.123.4567",
			expectedOutput: "[+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:        "11-digit international phone with parentheses and space",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "+7-890-123-4567",
				"tel_url":          "+78901234567",
			},
			originalInput:  "+7 (890) 123-4567",
			expectedOutput: "[+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:        "11-digit international phone with parentheses no space",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "+7-890-123-4567",
				"tel_url":          "+78901234567",
			},
			originalInput:  "+7(890)123-4567",
			expectedOutput: "[+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:        "11-digit international phone with parentheses plain",
			contentType: types.ContentTypePhone11Digit,
			metadata: map[string]interface{}{
				"formatted_display": "+7-890-123-4567",
				"tel_url":          "+78901234567",
			},
			originalInput:  "+7(890)1234567",
			expectedOutput: "[+7-890-123-4567](tel:+78901234567)",
		},
		{
			name:           "Non-phone content type returns original",
			contentType:    types.ContentTypeURL,
			metadata:       map[string]interface{}{},
			originalInput:  "https://example.com",
			expectedOutput: "https://example.com",
		},
		{
			name:           "Phone type without metadata returns original",
			contentType:    types.ContentTypePhone7Digit,
			metadata:       map[string]interface{}{},
			originalInput:  "1234567",
			expectedOutput: "1234567",
		},
		{
			name:        "Phone type missing formatted_display returns original",
			contentType: types.ContentTypePhone7Digit,
			metadata: map[string]interface{}{
				"tel_url": "1234567",
			},
			originalInput:  "1234567",
			expectedOutput: "1234567",
		},
		{
			name:        "Phone type missing tel_url returns original",
			contentType: types.ContentTypePhone7Digit,
			metadata: map[string]interface{}{
				"formatted_display": "123-4567",
			},
			originalInput:  "1234567",
			expectedOutput: "1234567",
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
			if err != nil && tt.contentType >= types.ContentTypePhone7Digit && tt.contentType <= types.ContentTypePhone11Digit {
				// Only expect errors for phone types with missing metadata
				if len(tt.metadata) >= 2 {
					t.Fatalf("Write() error = %v", err)
				}
			}

			if output != tt.expectedOutput {
				t.Errorf("Write() = %v, want %v", output, tt.expectedOutput)
			}
		})
	}
}

func TestPhoneWriter_GetName(t *testing.T) {
	cfg := &types.Config{}
	writer := NewPhoneWriter(cfg)

	expected := "PhoneWriter"
	if name := writer.GetName(); name != expected {
		t.Errorf("GetName() = %v, want %v", name, expected)
	}
}