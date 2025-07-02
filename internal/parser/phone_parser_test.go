package parser

import (
	"testing"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

func TestPhoneParser_CanHandle(t *testing.T) {
	cfg := &types.Config{}
	parser := NewPhoneParser(cfg)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 7-digit matches
		{"7-digit plain", "1234567", true},
		{"7-digit dash", "123-4567", true},
		{"7-digit dot", "123.4567", true},
		
		// 7-digit non-matches
		{"7-digit space", "123 4567", false},
		{"7-digit comma", "123,4567", false},
		{"7-digit leading zero", "01234567", false},
		
		// 10-digit matches
		{"10-digit plain", "8901234567", true},
		{"10-digit dash", "890-123-4567", true},
		{"10-digit dot", "890.123.4567", true},
		{"10-digit parens space", "(890) 123-4567", true},
		{"10-digit parens no space", "(890)123-4567", true},
		{"10-digit parens plain", "(890)1234567", true},
		
		// 10-digit non-matches
		{"10-digit extra digit", "89012345670", false},
		{"10-digit spaces", "890 123 4567", false},
		{"10-digit mixed parens space", "(890) 123 4567", false},
		{"10-digit mixed format", "(890) 1234 567", false},
		{"10-digit incomplete", "(890)123-456", false},
		{"10-digit too many", "(890)12345679", false},
		
		// 11-digit US matches
		{"11-digit US plain", "18901234567", true},
		{"11-digit US dash", "1-890-123-4567", true},
		{"11-digit US dot", "1.890.123.4567", true},
		{"11-digit US parens space", "1 (890) 123-4567", true},
		{"11-digit US parens no space", "1(890)123-4567", true},
		{"11-digit US parens plain", "1(890)1234567", true},
		
		// 11-digit international matches
		{"11-digit intl plain", "+78901234567", true},
		{"11-digit intl dash", "+7-890-123-4567", true},
		{"11-digit intl dot", "+7.890.123.4567", true},
		{"11-digit intl parens space", "+7 (890) 123-4567", true},
		{"11-digit intl parens no space", "+7(890)123-4567", true},
		{"11-digit intl parens plain", "+7(890)1234567", true},
		
		// Non-phone inputs
		{"text", "hello world", false},
		{"URL", "https://example.com", false},
		{"JIRA key", "PLAT-123", false},
		{"email", "test@example.com", false},
		{"empty", "", false},
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

func TestPhoneParser_Parse(t *testing.T) {
	cfg := &types.Config{}
	parser := NewPhoneParser(cfg)

	tests := []struct {
		name                 string
		input                string
		expectedType         types.ContentType
		expectedConfidence   int
		expectedDisplay      string
		expectedTelURL       string
		expectedIsExact      bool
	}{
		// 7-digit tests
		{
			name:               "7-digit plain exact",
			input:              "1234567",
			expectedType:       types.ContentTypePhone7Digit,
			expectedConfidence: 95,
			expectedDisplay:    "123-4567",
			expectedTelURL:     "1234567",
			expectedIsExact:    true,
		},
		{
			name:               "7-digit dash exact",
			input:              "123-4567",
			expectedType:       types.ContentTypePhone7Digit,
			expectedConfidence: 95,
			expectedDisplay:    "123-4567",
			expectedTelURL:     "1234567",
			expectedIsExact:    true,
		},
		{
			name:               "7-digit dot exact",
			input:              "123.4567",
			expectedType:       types.ContentTypePhone7Digit,
			expectedConfidence: 95,
			expectedDisplay:    "123-4567",
			expectedTelURL:     "1234567",
			expectedIsExact:    true,
		},
		
		// 10-digit tests
		{
			name:               "10-digit plain exact",
			input:              "8901234567",
			expectedType:       types.ContentTypePhone10Digit,
			expectedConfidence: 95,
			expectedDisplay:    "890-123-4567",
			expectedTelURL:     "8901234567",
			expectedIsExact:    true,
		},
		{
			name:               "10-digit dash exact",
			input:              "890-123-4567",
			expectedType:       types.ContentTypePhone10Digit,
			expectedConfidence: 95,
			expectedDisplay:    "890-123-4567",
			expectedTelURL:     "8901234567",
			expectedIsExact:    true,
		},
		{
			name:               "10-digit dot exact",
			input:              "890.123.4567",
			expectedType:       types.ContentTypePhone10Digit,
			expectedConfidence: 95,
			expectedDisplay:    "890-123-4567",
			expectedTelURL:     "8901234567",
			expectedIsExact:    true,
		},
		{
			name:               "10-digit parens space exact",
			input:              "(890) 123-4567",
			expectedType:       types.ContentTypePhone10Digit,
			expectedConfidence: 95,
			expectedDisplay:    "890-123-4567",
			expectedTelURL:     "8901234567",
			expectedIsExact:    true,
		},
		{
			name:               "10-digit parens no space exact",
			input:              "(890)123-4567",
			expectedType:       types.ContentTypePhone10Digit,
			expectedConfidence: 95,
			expectedDisplay:    "890-123-4567",
			expectedTelURL:     "8901234567",
			expectedIsExact:    true,
		},
		{
			name:               "10-digit parens plain exact",
			input:              "(890)1234567",
			expectedType:       types.ContentTypePhone10Digit,
			expectedConfidence: 95,
			expectedDisplay:    "890-123-4567",
			expectedTelURL:     "8901234567",
			expectedIsExact:    true,
		},
		
		// 11-digit US tests
		{
			name:               "11-digit US plain exact",
			input:              "18901234567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "1-890-123-4567",
			expectedTelURL:     "+18901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit US dash exact",
			input:              "1-890-123-4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "1-890-123-4567",
			expectedTelURL:     "+18901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit US dot exact",
			input:              "1.890.123.4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "1-890-123-4567",
			expectedTelURL:     "+18901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit US parens space exact",
			input:              "1 (890) 123-4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "1-890-123-4567",
			expectedTelURL:     "+18901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit US parens no space exact",
			input:              "1(890)123-4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "1-890-123-4567",
			expectedTelURL:     "+18901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit US parens plain exact",
			input:              "1(890)1234567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "1-890-123-4567",
			expectedTelURL:     "+18901234567",
			expectedIsExact:    true,
		},
		
		// 11-digit international tests
		{
			name:               "11-digit intl plain exact",
			input:              "+78901234567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "+7-890-123-4567",
			expectedTelURL:     "+78901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit intl dash exact",
			input:              "+7-890-123-4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "+7-890-123-4567",
			expectedTelURL:     "+78901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit intl dot exact",
			input:              "+7.890.123.4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "+7-890-123-4567",
			expectedTelURL:     "+78901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit intl parens space exact",
			input:              "+7 (890) 123-4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "+7-890-123-4567",
			expectedTelURL:     "+78901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit intl parens no space exact",
			input:              "+7(890)123-4567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "+7-890-123-4567",
			expectedTelURL:     "+78901234567",
			expectedIsExact:    true,
		},
		{
			name:               "11-digit intl parens plain exact",
			input:              "+7(890)1234567",
			expectedType:       types.ContentTypePhone11Digit,
			expectedConfidence: 95,
			expectedDisplay:    "+7-890-123-4567",
			expectedTelURL:     "+78901234567",
			expectedIsExact:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if ctx == nil {
				t.Fatalf("Parse(%q) returned nil context", tt.input)
			}

			if ctx.DetectedType != tt.expectedType {
				t.Errorf("DetectedType = %v, want %v", ctx.DetectedType, tt.expectedType)
			}

			if ctx.Confidence != tt.expectedConfidence {
				t.Errorf("Confidence = %v, want %v", ctx.Confidence, tt.expectedConfidence)
			}

			if display := ctx.Metadata["formatted_display"]; display != tt.expectedDisplay {
				t.Errorf("formatted_display = %v, want %v", display, tt.expectedDisplay)
			}

			if telURL := ctx.Metadata["tel_url"]; telURL != tt.expectedTelURL {
				t.Errorf("tel_url = %v, want %v", telURL, tt.expectedTelURL)
			}

			if isExact := ctx.Metadata["is_exact_match"]; isExact != tt.expectedIsExact {
				t.Errorf("is_exact_match = %v, want %v", isExact, tt.expectedIsExact)
			}
		})
	}
}

func TestPhoneParser_ConfidenceScoring(t *testing.T) {
	cfg := &types.Config{}
	parser := NewPhoneParser(cfg)

	tests := []struct {
		name               string
		input              string
		expectedConfidence int
		description        string
	}{
		{
			name:               "exact 7-digit",
			input:              "1234567",
			expectedConfidence: 95,
			description:        "exact match should have high confidence",
		},
		{
			name:               "exact 10-digit",
			input:              "8901234567",
			expectedConfidence: 95,
			description:        "exact match should have high confidence",
		},
		{
			name:               "exact 11-digit",
			input:              "18901234567",
			expectedConfidence: 95,
			description:        "exact match should have high confidence",
		},
		// Note: embedded phone numbers would need a more complex test setup
		// since our parser currently only matches exact phone numbers
		// This is intentional based on the regex patterns designed to match full input
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.input, err)
			}
			if ctx == nil {
				t.Fatalf("Parse(%q) returned nil context", tt.input)
			}

			if ctx.Confidence != tt.expectedConfidence {
				t.Errorf("Confidence = %v, want %v (%s)", ctx.Confidence, tt.expectedConfidence, tt.description)
			}
		})
	}
}

func TestPhoneParser_InvalidFormats(t *testing.T) {
	cfg := &types.Config{}
	parser := NewPhoneParser(cfg)

	invalidInputs := []string{
		// From the examples - these should NOT match
		"123 4567",        // 7-digit with space
		"123,4567",        // 7-digit with comma
		"01234567",        // 7-digit with leading zero
		"89012345670",     // 10-digit with extra digit
		"890 123 4567",    // 10-digit with spaces
		"(890) 123 4567",  // 10-digit mixed separators
		"(890) 1234 567",  // 10-digit wrong grouping
		"(890)123-456",    // 10-digit incomplete
		"(890)12345679",   // 10-digit too many digits
		
		// Other invalid formats
		"",                // empty
		"abc",             // text
		"123-45-6789",     // SSN format
		"12345",           // too short
		"123456789012",    // too long
		"+",               // just plus
		"++1234567890",    // double plus
	}

	for _, input := range invalidInputs {
		t.Run("invalid_"+input, func(t *testing.T) {
			if parser.CanHandle(input) {
				t.Errorf("CanHandle(%q) = true, want false (should not match invalid format)", input)
			}

			ctx, err := parser.Parse(input)
			if err != nil {
				t.Errorf("Parse(%q) should not error on invalid input, got: %v", input, err)
			}
			if ctx != nil {
				t.Errorf("Parse(%q) = %+v, want nil (should not match invalid format)", input, ctx)
			}
		})
	}
}