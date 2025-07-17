package cmd

import (
	"testing"
)

func TestPreprocessTelURIs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "7-digit tel URI",
			input:    "tel:1234567",
			expected: "1234567",
		},
		{
			name:     "10-digit tel URI",
			input:    "tel:8901234567",
			expected: "8901234567",
		},
		{
			name:     "11-digit tel URI",
			input:    "tel:18901234567",
			expected: "18901234567",
		},
		{
			name:     "International tel URI",
			input:    "tel:+18901234567",
			expected: "+18901234567",
		},
		{
			name:     "Tel URI with dashes",
			input:    "tel:890-123-4567",
			expected: "890-123-4567",
		},
		{
			name:     "Tel URI with parentheses",
			input:    "tel:(890)123-4567",
			expected: "(890)123-4567",
		},
		{
			name:     "Tel URI with spaces and parentheses",
			input:    "tel:1 (890) 123-4567",
			expected: "1 (890) 123-4567",
		},
		{
			name:     "Tel URI with whitespace",
			input:    "  tel:1234567  ",
			expected: "1234567",
		},
		{
			name:     "Non-tel URI (URL)",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "Plain phone number",
			input:    "8901234567",
			expected: "8901234567",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Tel URI with no phone number",
			input:    "tel:",
			expected: "",
		},
		{
			name:     "Tel URI case insensitive",
			input:    "TEL:1234567",
			expected: "TEL:1234567", // Should not match (case sensitive)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preprocessTelURIs(tt.input)
			if result != tt.expected {
				t.Errorf("preprocessTelURIs(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}