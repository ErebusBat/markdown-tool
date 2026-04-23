# Gemini Chat URL Support

## Summary

Add support for converting Gemini chat URLs into markdown links with a robot emoji prefix.

**Input:** `https://gemini.google.com/app/ac9ebc9d76c30fc1`
**Output:** `[🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)`

## Approach

Extend the existing URLParser + URLWriter mega-pattern (same as YouTube, MiniMax, Notion).

## Files to Modify

1. `pkg/types/types.go` — Add `ContentTypeGeminiURL` constant
2. `internal/parser/url_parser.go` — Add `isGeminiURL()` + `parseGeminiURL()` + switch case
3. `internal/writer/url_writer.go` — Add vote score (90) + `writeGeminiURL()` + switch case
4. `internal/parser/url_parser_test.go` — Add test cases
5. `internal/writer/url_writer_test.go` — Add test cases
6. `integration_test.go` — Add end-to-end test

## Parser Logic

- Detect host `gemini.google.com` with path starting with `/app/`
- Extract chat ID from path segment after `/app/`
- Set confidence to 90

## Writer Logic

- Always emit `[🤖 Gemini Chat](<url>)` with hardcoded display text
- No config needed
- Vote score: 90 (standard, same as GitHub/JIRA/MiniMax)
