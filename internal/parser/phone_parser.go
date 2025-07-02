package parser

import (
	"regexp"
	"strings"

	"github.com/erebusbat/markdown-tool/pkg/types"
)

type PhoneParser struct {
	config *types.Config
}

func NewPhoneParser(cfg *types.Config) *PhoneParser {
	return &PhoneParser{config: cfg}
}

func (p *PhoneParser) CanHandle(input string) bool {
	return p.detectPhoneNumber(input) != nil
}

func (p *PhoneParser) Parse(input string) (*types.ParseContext, error) {
	match := p.detectPhoneNumber(input)
	if match == nil {
		return nil, nil
	}

	// Determine confidence based on whether input is exactly the phone number
	trimmedInput := strings.TrimSpace(input)
	confidence := 60 // Default for embedded numbers
	if trimmedInput == match.RawNumber {
		confidence = 95 // High confidence for exact matches
	}

	ctx := &types.ParseContext{
		OriginalInput: input,
		DetectedType:  match.Type,
		Confidence:    confidence,
		Metadata: map[string]interface{}{
			"raw_number":        match.RawNumber,
			"formatted_display": match.FormattedDisplay,
			"tel_url":          match.TelURL,
			"is_exact_match":   trimmedInput == match.RawNumber,
		},
	}

	return ctx, nil
}

type phoneMatch struct {
	Type             types.ContentType
	RawNumber        string // The matched phone number from input
	FormattedDisplay string // How to display it: [123-4567]
	TelURL          string // tel: URL: tel:1234567
}

func (p *PhoneParser) detectPhoneNumber(input string) *phoneMatch {
	trimmed := strings.TrimSpace(input)
	
	// Try 7-digit patterns first
	if match := p.match7Digit(trimmed); match != nil {
		return match
	}
	
	// Try 10-digit patterns
	if match := p.match10Digit(trimmed); match != nil {
		return match
	}
	
	// Try 11-digit patterns (US and international)
	if match := p.match11Digit(trimmed); match != nil {
		return match
	}
	
	return nil
}

// match7Digit handles 7-digit phone numbers
func (p *PhoneParser) match7Digit(input string) *phoneMatch {
	patterns := []string{
		`^(\d{7})$`,           // 1234567
		`^(\d{3})-(\d{4})$`,   // 123-4567
		`^(\d{3})\.(\d{4})$`,  // 123.4567
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(input)
		if matches != nil {
			var digits string
			if len(matches) == 2 {
				// Single capture group (no separators)
				digits = matches[1]
			} else if len(matches) == 3 {
				// Two capture groups (with separators)
				digits = matches[1] + matches[2]
			}
			
			if len(digits) == 7 {
				return &phoneMatch{
					Type:             types.ContentTypePhone7Digit,
					RawNumber:        input,
					FormattedDisplay: digits[:3] + "-" + digits[3:],
					TelURL:          digits,
				}
			}
		}
	}
	
	return nil
}

// match10Digit handles 10-digit phone numbers
func (p *PhoneParser) match10Digit(input string) *phoneMatch {
	patterns := []string{
		`^(\d{10})$`,                      // 8901234567
		`^(\d{3})-(\d{3})-(\d{4})$`,       // 890-123-4567
		`^(\d{3})\.(\d{3})\.(\d{4})$`,     // 890.123.4567
		`^\((\d{3})\) (\d{3})-(\d{4})$`,   // (890) 123-4567
		`^\((\d{3})\)(\d{3})-(\d{4})$`,    // (890)123-4567
		`^\((\d{3})\)(\d{7})$`,            // (890)1234567
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(input)
		if matches != nil {
			var digits string
			// Extract all digit groups
			for i := 1; i < len(matches); i++ {
				digits += matches[i]
			}
			
			if len(digits) == 10 {
				return &phoneMatch{
					Type:             types.ContentTypePhone10Digit,
					RawNumber:        input,
					FormattedDisplay: digits[:3] + "-" + digits[3:6] + "-" + digits[6:],
					TelURL:          digits,
				}
			}
		}
	}
	
	return nil
}

// match11Digit handles 11-digit phone numbers (US and international)
func (p *PhoneParser) match11Digit(input string) *phoneMatch {
	// US patterns (country code 1)
	usPatterns := []string{
		`^(1)(\d{10})$`,                       // 18901234567
		`^(1)-(\d{3})-(\d{3})-(\d{4})$`,       // 1-890-123-4567
		`^(1)\.(\d{3})\.(\d{3})\.(\d{4})$`,    // 1.890.123.4567
		`^(1) \((\d{3})\) (\d{3})-(\d{4})$`,   // 1 (890) 123-4567
		`^(1)\((\d{3})\)(\d{3})-(\d{4})$`,     // 1(890)123-4567
		`^(1)\((\d{3})\)(\d{7})$`,             // 1(890)1234567
	}
	
	for _, pattern := range usPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(input)
		if matches != nil {
			var digits string
			// Extract all digit groups
			for i := 1; i < len(matches); i++ {
				digits += matches[i]
			}
			
			if len(digits) == 11 && digits[0] == '1' {
				return &phoneMatch{
					Type:             types.ContentTypePhone11Digit,
					RawNumber:        input,
					FormattedDisplay: "1-" + digits[1:4] + "-" + digits[4:7] + "-" + digits[7:],
					TelURL:          "+1" + digits[1:],
				}
			}
		}
	}
	
	// International patterns (must start with +)
	intlPatterns := []string{
		`^\+(\d)(\d{10})$`,                       // +78901234567
		`^\+(\d)-(\d{3})-(\d{3})-(\d{4})$`,       // +7-890-123-4567
		`^\+(\d)\.(\d{3})\.(\d{3})\.(\d{4})$`,    // +7.890.123.4567
		`^\+(\d) \((\d{3})\) (\d{3})-(\d{4})$`,   // +7 (890) 123-4567
		`^\+(\d)\((\d{3})\)(\d{3})-(\d{4})$`,     // +7(890)123-4567
		`^\+(\d)\((\d{3})\)(\d{7})$`,             // +7(890)1234567
	}
	
	for _, pattern := range intlPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(input)
		if matches != nil {
			var digits string
			// Extract all digit groups
			for i := 1; i < len(matches); i++ {
				digits += matches[i]
			}
			
			if len(digits) == 11 {
				countryCode := digits[0:1]
				phoneNumber := digits[1:]
				return &phoneMatch{
					Type:             types.ContentTypePhone11Digit,
					RawNumber:        input,
					FormattedDisplay: "+" + countryCode + "-" + phoneNumber[:3] + "-" + phoneNumber[3:6] + "-" + phoneNumber[6:],
					TelURL:          "+" + digits,
				}
			}
		}
	}
	
	return nil
}