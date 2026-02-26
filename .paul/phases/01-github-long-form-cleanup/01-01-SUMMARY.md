---
phase: 01-github-long-form-cleanup
plan: 01
subsystem: parsing
tags: [github, jira, markdown]

requires: []
provides:
  - GitHub long titles stripped of leading JIRA keys
  - Tests covering JIRA-key exclusion in GitHub long output
affects: [github-long]

tech-stack:
  added: []
  patterns: ["strip leading tokens before metadata persistence"]

key-files:
  created: []
  modified:
    - internal/parser/github_long_parser.go
    - internal/parser/github_long_parser_test.go
    - internal/writer/url_writer.go
    - internal/writer/url_writer_test.go

key-decisions: []
patterns-established:
  - "Normalize GitHub long titles to drop leading JIRA keys"

# Metrics
duration: 8min
started: 2026-02-26T10:06:56Z
completed: 2026-02-26T10:14:53Z
---

# Phase 1 Plan 01: GitHub Long Form Cleanup Summary

**GitHub long-form parsing now strips leading JIRA keys so markdown titles stay focused on the GitHub issue/PR title.**

## Performance

| Metric | Value |
|--------|-------|
| Duration | 8min |
| Started | 2026-02-26T10:06:56Z |
| Completed | 2026-02-26T10:14:53Z |
| Tasks | 2 completed |
| Files modified | 4 |

## Acceptance Criteria Results

| Criterion | Status | Notes |
|-----------|--------|-------|
| AC-1: Strip JIRA key from GitHub long titles | Pass | Parser and writer sanitize titles; tests added |
| AC-2: Preserve non-JIRA titles | Pass | No changes to non-JIRA title paths |
| AC-3: Parser does not treat JIRA keys as org/repo or issue number | Pass | Title stripped while number remains GitHub issue/PR |

## Accomplishments

- Added parser sanitization and tests to remove leading JIRA keys from GitHub long titles
- Guarded writer output with the same sanitization plus regression coverage

## Task Commits

Each task committed atomically:

| Task | Commit | Type | Description |
|------|--------|------|-------------|
| Task 1: Strip JIRA key from GitHub long titles in parser | n/a | test | Parser sanitization and tests |
| Task 2: Ensure GitHub long writer output reflects stripped titles | n/a | test | Writer sanitization and tests |

Plan metadata: n/a

## Files Created/Modified

| File | Change | Purpose |
|------|--------|---------|
| `internal/parser/github_long_parser.go` | Modified | Strip leading JIRA key tokens from titles |
| `internal/parser/github_long_parser_test.go` | Modified | Add JIRA-key title regression test |
| `internal/writer/url_writer.go` | Modified | Sanitize GitHub long titles before output |
| `internal/writer/url_writer_test.go` | Modified | Validate JIRA-key omission in output |

## Decisions Made

None - followed plan as specified.

## Deviations from Plan

### Summary

| Type | Count | Impact |
|------|-------|--------|
| Auto-fixed | 0 | None |
| Scope additions | 0 | None |
| Deferred | 0 | None |

**Total impact:** None — plan executed as written.

### Deferred Items

None - plan executed exactly as written.

## Issues Encountered

None

## Next Phase Readiness

**Ready:**
- GitHub long-form output now excludes leading JIRA keys

**Concerns:**
- None

**Blockers:**
- None

---
*Phase: 01-github-long-form-cleanup, Plan: 01*
*Completed: 2026-02-26*
