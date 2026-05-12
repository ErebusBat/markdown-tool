# Markdown Tool — Language-Agnostic Specification

This document fully specifies the behavior of `markdown-tool` such that it can be
reimplemented in any language without reference to the original Go source.

---

## 1. Overview

`markdown-tool` is a CLI utility that takes raw text (URLs, JIRA keys, phone
numbers, multi-line clipboard dumps, custom URI schemes) and transforms them
into well-formatted markdown links suitable for knowledge-management tools
(Vimwiki, Obsidian, Notion, etc.).

**Design philosophy:**
- **Parse, Vote, Write pipeline** — parsers detect content; writers vote on
  which can best render it; the highest-confidence writer wins.
- **Deterministic output** — same input always produces same output.
- **Fail-safe** — unrecognized input passes through verbatim.
- **Config-driven** — behaviour for known services (GitHub, JIRA, Jenkins) is
  controlled by a YAML config file.

---

## 2. Architecture

```
                   +---------+
                   |  Input  |  (stdin or system clipboard)
                   +----+----+
                        |
                        v
                +-------+-------+
                |  Preprocess   |  (strip `tel:` prefix for phone URIs)
                +-------+-------+
                        |
                        v
              +---------+---------+
              |   ALL PARSERS     |  (run in registration order)
              |  produce Context[]|
              +---------+---------+
                        |
                        v
              +---------+---------+
              |   ALL WRITERS     |  (each votes on every context)
              |  highest wins     |
              +---------+---------+
                        |
                        v
               +--------+--------+
               |  Best Writer    |  (generates markdown)
               +--------+--------+
                        |
                        v
                  +----+-----+
                  |  stdout   |
                  +----------+
```

**Logging:** All diagnostic/log output goes to **stderr**. Only the transformed
text goes to **stdout**.

### 2.1 Parser Registration Order

Parsers are tried in this fixed order (first match wins for overlapping
detections):

| #  | Parser                        |
|----|-------------------------------|
| 1  | URLParser                     |
| 2  | GitHubLongParser              |
| 3  | CodeCommitLongParser          |
| 4  | CodeCommitParser              |
| 5  | JIRAKeyWithDescriptionParser  |
| 6  | JIRAKeyParser                 |
| 7  | PhoneParser                   |
| 8  | RaycastParser                 |
| 9  | OpenCodeSessionParser         |
| 10 | CodexParser                   |

### 2.2 Writer Voting

Each writer receives every parse context and returns an integer confidence
score (0 = "not interested"). The writer with the highest score wins. On ties,
the writer registered earlier wins.

Writers are registered in this order:

| #  | Writer                       | Fallback? |
|----|------------------------------|-----------|
| 1  | URLWriter                    | No        |
| 2  | JIRAKeyWithDescriptionWriter | No        |
| 3  | JIRAWriter                   | No        |
| 4  | PhoneWriter                  | No        |
| 5  | RaycastWriter                | No        |
| 6  | OpenCodeSessionWriter        | No        |
| 7  | CodexWriter                  | No        |
| 8  | PassthroughWriter            | Yes (always votes 1) |

If no writer scores > 0, output the input verbatim.

---

## 3. Content Types

### 3.1 Enumeration

```
ContentTypeUnknown               = 0
ContentTypeURL                   = 1   (generic/unrecognised URL)
ContentTypeGitHubURL             = 2
ContentTypeGitHubLong            = 3   (multi-line GitHub UI clipboard)
ContentTypeJIRAURL               = 4
ContentTypeJIRAComment           = 5   (JIRA URL with ?focusedCommentId=)
ContentTypeNotionURL             = 6
ContentTypeJenkinsURL            = 7
ContentTypeYouTubeURL            = 8
ContentTypeCodeCommitURL         = 9   (AWS CodeCommit console URL)
ContentTypeCodeCommitLong        = 10  (multi-line CodeCommit UI clipboard)
ContentTypeJIRAKey               = 11  (standalone PROJ-123)
ContentTypeJIRAKeyWithDescription = 12 (multi-line JIRA key + description)
ContentTypePhone7Digit           = 13
ContentTypePhone10Digit          = 14
ContentTypePhone11Digit          = 15
ContentTypeRaycastURI            = 16
ContentTypeOpenCodeSession       = 17
ContentTypeMiniMaxURL            = 18
ContentTypeGeminiURL             = 19
ContentTypeCodexThread           = 20
ContentTypeCircleCI              = 21
ContentTypeChatGPT               = 22
```

---

## 4. Core Interfaces

### 4.1 ParseContext

The data structure passed between parsers and writers:

```
ParseContext {
    OriginalInput  string              // unchanged raw input
    DetectedType   ContentType         // enum from §3
    Confidence     integer             // 0-100, parser's confidence
    Metadata       Map<string, any>    // key-value pairs
}
```

### 4.2 Parser Interface

```
Parser {
    // Returns (nil, nil) when input does not match.
    Parse(input string) -> (ParseContext | null, error | null)

    // Quick pre-check; must be cheap. Never has false negatives
    // relative to Parse().
    CanHandle(input string) -> boolean
}
```

### 4.3 Writer Interface

```
Writer {
    // Generate markdown output.
    Write(ctx ParseContext) -> (string, error)

    // Return confidence score 0-100 for this context.
    // 0 means "I cannot handle this."
    Vote(ctx ParseContext) -> integer

    // Human-readable name for logging.
    GetName() -> string
}
```

---

## 5. Parsers

### 5.1 URLParser

**What it detects:** Any `http://` or `https://` URL.

**CanHandle:** Parsable as URL AND starts with `http://` or `https://`.

**Sub-detection (in priority order):**
1. GitHub: host == `github.com`
2. JIRA: host matches configured JIRA domain
3. Jenkins: host matches configured Jenkins domain
4. YouTube: host in `[www.youtube.com, youtube.com, youtu.be, m.youtube.com]`
5. CodeCommit: host contains `console.aws.amazon.com` AND path contains
   `/codesuite/codecommit/repositories/` AND path contains `/pull-requests/`
6. Notion: host contains `notion.so`
7. MiniMax: host == `agent.minimax.io`
8. Gemini: host == `gemini.google.com` AND path starts with `/app/`
9. CircleCI: host == `app.circleci.com` AND path starts with `/pipelines/`
   AND path contains `/workflows/`
10. ChatGPT: host == `chatgpt.com` AND path length >= 3 AND path starts with `/c/`
11. Fallback: mark as generic `ContentTypeURL`

**Confidence scores:**
| Type             | Confidence |
|------------------|-----------|
| GitHub           | 90        |
| JIRA             | 90        |
| JIRA Comment     | 95        |
| Jenkins          | 90        |
| YouTube          | 90        |
| CodeCommit       | 90        |
| Notion           | 85        |
| MiniMax          | 90        |
| Gemini           | 90        |
| CircleCI         | 90        |
| ChatGPT          | 90        |
| Generic URL      | 50        |

#### 5.1.1 GitHub URL Metadata

Extract from path `/org/repo[/type/number]`:
- `org` (string)
- `repo` (string)
- `type` (string) — "pull", "issues", or "commit" if >= 4 path segments
- `number` (string) — issue/PR number or commit hash

#### 5.1.2 JIRA URL Metadata

Regex from path: `/browse/([A-Z]+-\d+)`
- `issue_key` (string)
- `comment_id` (string) — from query param `focusedCommentId`

#### 5.1.3 Jenkins URL Metadata

Regex: `^/job/([^/]+)/(\d+)` for build number, or `^/job/([^/]+)` for job
name only.
- `job_name` (string)
- `build_number` (string) — absent if `/lastBuild/`, `/lastSuccessfulBuild/`,
  or bare `/job/name/`

#### 5.1.4 YouTube URL Metadata

**Video:** Extract `v` query param, or strip leading `/` from `youtu.be`.
**Playlist:** Path == `/playlist`, extract `list` query param.

Metadata:
- `youtube_type` — "video" or "playlist"
- `video_id` or `playlist_id` (string)
- `title` (string) — fetched via YouTube oEmbed API:
  `https://www.youtube.com/oembed?url=<encoded>&format=json` → `title` field

If title fetch fails or returns empty, `title` is absent (empty string).

#### 5.1.5 CodeCommit URL Metadata

Extract region from subdomain (first dot-separated segment of host).
Extract from path: `/repositories/([^/]+)/pull-requests/(\d+)`
- `region` (string)
- `repo` (string)
- `number` (string)

#### 5.1.6 Notion URL Metadata

Extract from last path segment: `^(.+)-[a-f0-9]{32}$`
Convert dashes to spaces for title.
- `title` (string)

#### 5.1.7 MiniMax URL Metadata

Extract `id` query param.
- `chat_id` (string)

#### 5.1.8 Gemini URL Metadata

Extract from path: `/app/([a-f0-9]+)`
- `chat_id` (string)
- `clean_url` (string) — reconstructed clean URL (scheme + host + `/app/` + id)

Note: The parser matches up to the first non-alphanumeric character after the
hex ID, so trailing characters like ` →` in `ac9ebc9d76c30fc1 →` are ignored.

#### 5.1.9 ChatGPT URL Metadata

Extract from path: `^/c/([a-f0-9\-]+)`
- `chat_id` (string)

#### 5.1.10 CircleCI URL Metadata

Extract from path:
`^/pipelines/([^/]+)/([^/]+)/([^/]+)/(\d+)/workflows/` for vcs/org/repo/pipeline
Extract workflow ID: `/workflows/([a-f0-9\-]+)`
- `vcs` (string)
- `org` (string)
- `repo` (string)
- `pipeline_number` (string)
- `workflow_id` (string)

#### 5.1.11 Generic URL Metadata

- `domain` (string) — the hostname

---

### 5.2 GitHubLongParser

**What it detects:** Multi-line clipboard text copied from a GitHub issue or
PR page, AND simple issue-title lines (with or without username prefix).

**CanHandle:** Three cases:
1. **Simple issue title** — single line matching issue-title-with-number pattern
   OR 2-3 lines where exactly one looks like an issue title and no GitHub UI
   indicators are present (see §5.2.1).
2. **Multi-line GitHub UI dump** — ≥3 lines containing at least: one
   org-name-like line, one repo-name-like line, and one issue-title line.

#### 5.2.1 Simple Issue Title (Cases 1)

A line matches simple-issue-title when it matches one of:
- **Username prefix:** `^[a-zA-Z0-9_-]+\s+.+\s+#\d+\s*$` where the first word
  is a valid GitHub username (see §5.2.2).
- **Title + number:** `^.+\s*#\d+\s*$`

If input has ≤3 lines, no GitHub UI indicators (lines containing "Type / to
search", "Pull requests", "Discussions", "Actions", "Projects", "Wiki",
"Security", "Insights", "Settings"), and ≤1 org-repo-like line, classify as
simple.

#### 5.2.2 GitHub Username Validation

A valid GitHub username:
- Length 1-39
- Matches `^[a-zA-Z0-9]([a-zA-Z0-9_-]*[a-zA-Z0-9])?$` (cannot start/end with
  hyphen, can contain hyphen and underscore)
- Is NOT a case-insensitive match for any common English verb:
  `adds`, `fixes`, `updates`, `removes`, `creates`, `deletes`, `implements`,
  `enhances`, `refactors`, `only`, `some`, `the`, `and`, `with`, `for`,
  `from`, `this`, `that`

#### 5.2.3 Org/Repo Name Validation

**Org name:** Length 1-39, matches `^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`
**Repo name:** Length 1-100, matches `^[a-zA-Z0-9._-]+$`

#### 5.2.4 Issue Title & Number Extraction

**Title with number:** Regex `^(.+?)\s*#(\d+)\s*$`
- Group 1 → title (trimmed)
- Group 2 → number

**Username prefix:** Regex `^([a-zA-Z0-9_-]+)\s+(.+)\s+#(\d+)\s*$`
- Group 2 → title (trimmed)
- Group 3 → number
(The username in group 1 is discarded.)

**Standalone number line:** Regex `^\s*#\s*(\d+)\s*$`
The preceding non-empty line becomes the title.

**JIRA key stripping:** Before using a title, strip any leading JIRA key
pattern matching `^\s*(\[[A-Z][A-Z0-9]+-\d+\]\s*|[A-Z][A-Z0-9]+-\d+:\s*)`

#### 5.2.5 Issue Type Detection

Scan all lines (case-insensitive):
- If any line contains "review requested" or "requested your review" → `pull`
- If "pull request" (singular) appears without "pull requests" (plural) → `pull`
- If "Agents" AND "Pull requests" AND "Repository navigation" all appear → `pull`
- Otherwise → `issues`

#### 5.2.6 Simple Issue Title Output

When classified as simple:
- Uses `default_org` and `default_repo` from config
- Returns `ContentTypeGitHubLong` with confidence **95**
- Metadata: `org`, `repo`, `title`, `number`, `type` (always "issues")

#### 5.2.7 Multi-line GitHub UI Output

When classified as multi-line:
- Scans lines to find: org name, repo name, issue title, issue number
- First valid GitHub name becomes `org`, next distinct valid repo name becomes
  `repo`
- Returns `ContentTypeGitHubLong` with confidence **90**

---

### 5.3 CodeCommitLongParser

**What it detects:** Multi-line clipboard text copied from an AWS CodeCommit
PR detail page.

**CanHandle:** ≥3 lines containing ALL of:
- A line with "Developer Tools" or "CodeCommit"
- A line with "Repositories"
- A line with "Pull requests" or "pull-requests"
- A line matching CodeCommit PR title pattern (see below)

**PR title pattern:** `^\d+:\s*.+` (e.g., `411: SEC-12335: Pass SENDGRID_API_KEY`)

**Extraction:**
- `region` — the first line matching a known AWS region identifier
  (`us-east-1`, `us-west-2`, `eu-west-1`, etc. — see §5.3.1)
- `repo` — the line immediately after "Repositories" or "CodeCommit", if it
  matches `^[a-zA-Z0-9._-]+$` and is not a known UI keyword
- `number` — standalone numeric line (regex `^\d+$`, length < 10)
- `title` — from PR title line, regex `^(\d+):\s*(.+)`; the text after the
  colon is the title

**Confidence:** 90 (content type: `ContentTypeCodeCommitLong`)

**Fallback region:** If no region found, default to `us-east-1`.

#### 5.3.1 Known AWS Regions

```
us-east-1, us-east-2, us-west-1, us-west-2,
eu-west-1, eu-west-2, eu-west-3, eu-central-1,
ap-northeast-1, ap-northeast-2, ap-southeast-1, ap-southeast-2,
ap-south-1, sa-east-1, ca-central-1
```

---

### 5.4 CodeCommitParser

**What it detects:** AWS CodeCommit console PR URLs (short URL format).

**CanHandle:** URL where host contains `console.aws.amazon.com` AND path
contains `/codesuite/codecommit/repositories/` AND path contains
`/pull-requests/`.

**Extraction:** Same as §5.1.5 (region from subdomain, repo/number from path).

**Confidence:** 90 (content type: `ContentTypeCodeCommitURL`)

---

### 5.5 JIRAKeyWithDescriptionParser

**What it detects:** Multi-line input with a JIRA key on the first line,
blank second line, and description text following.

**CanHandle:** ≥3 lines where:
1. First line matches `^[A-Z]+-\d+$`
2. Second line is empty or whitespace-only
3. At least one non-empty line exists after line 2

**Project gate:** The project prefix (text before the first `-`) MUST be in
the configured `jira.projects` list. If not configured, returns `(nil, nil)`.

**Description extraction:**
- Scan from line index 2 onwards
- Skip any leading numeric-only or comma-formatted-number line
  (e.g., JIRA's unread count of `"1"` or `"1,234"`) that is followed by
  non-empty content
- Join remaining non-empty lines with spaces

**Metadata:**
- `issue_key` (string) — e.g., `"PLAT-192"`
- `project` (string) — e.g., `"PLAT"`
- `description` (string)

**Confidence:** 98 (content type: `ContentTypeJIRAKeyWithDescription`)

---

### 5.6 JIRAKeyParser

**What it detects:** A single JIRA key (`PROJECT-123`) on its own line.

**CanHandle:** Input matches `^[A-Z]+-\d+$` (trimmed).

**Project gate:** Project prefix MUST be in `jira.projects`. Returns
`(nil, nil)` for unconfigured projects.

**Metadata:**
- `issue_key` (string)
- `project` (string)

**Confidence:** 95 (content type: `ContentTypeJIRAKey`)

---

### 5.7 PhoneParser

**What it detects:** 7, 10, or 11-digit phone numbers in specific formats.

**Detection order:** 7-digit first, then 10-digit, then 11-digit. The first
match wins.

**Note:** Leading zeros on a 7-digit number are NOT matched (e.g.,
`01234567` passes through unchanged). Spaces as separators are NOT
supported — only dashes, dots, and parentheses.

#### 5.7.1 7-Digit Patterns

| Pattern                | Example      | Type              |
|------------------------|--------------|-------------------|
| `^(\d{7})$`            | `1234567`    | Phone7Digit       |
| `^(\d{3})-(\d{4})$`   | `123-4567`   | Phone7Digit       |
| `^(\d{3})\.(\d{4})$`  | `123.4567`   | Phone7Digit       |

Metadata:
- `raw_number` — exact matched text
- `formatted_display` — `XXX-XXXX` (digits only, with inserted dash)
- `tel_url` — bare digits (no separators)

#### 5.7.2 10-Digit Patterns

| Pattern                         | Example          | Type              |
|---------------------------------|------------------|-------------------|
| `^(\d{10})$`                    | `8901234567`     | Phone10Digit      |
| `^(\d{3})-(\d{3})-(\d{4})$`    | `890-123-4567`   | Phone10Digit      |
| `^(\d{3})\.(\d{3})\.(\d{4})$`  | `890.123.4567`   | Phone10Digit      |
| `^\((\d{3})\) (\d{3})-(\d{4})$`| `(890) 123-4567` | Phone10Digit      |
| `^\((\d{3})\)(\d{3})-(\d{4})$` | `(890)123-4567`  | Phone10Digit      |
| `^\((\d{3})\)(\d{7})$`         | `(890)1234567`   | Phone10Digit      |

Metadata:
- `formatted_display` — `XXX-XXX-XXXX`
- `tel_url` — bare 10 digits

#### 5.7.3 11-Digit US Patterns (country code 1)

| Pattern                               | Example            | Type              |
|---------------------------------------|--------------------|-------------------|
| `^(1)(\d{10})$`                       | `18901234567`      | Phone11Digit      |
| `^(1)-(\d{3})-(\d{3})-(\d{4})$`      | `1-890-123-4567`   | Phone11Digit      |
| `^(1)\.(\d{3})\.(\d{3})\.(\d{4})$`   | `1.890.123.4567`   | Phone11Digit      |
| `^(1) \((\d{3})\) (\d{3})-(\d{4})$`  | `1 (890) 123-4567` | Phone11Digit      |
| `^(1)\((\d{3})\)(\d{3})-(\d{4})$`    | `1(890)123-4567`   | Phone11Digit      |
| `^(1)\((\d{3})\)(\d{7})$`            | `1(890)1234567`    | Phone11Digit      |

Metadata:
- `formatted_display` — `1-XXX-XXX-XXXX`
- `tel_url` — `+1` + 10-digit number

#### 5.7.4 11-Digit International Patterns (leading +)

| Pattern                                 | Example              | Type              |
|-----------------------------------------|----------------------|-------------------|
| `^\+(\d)(\d{10})$`                      | `+78901234567`       | Phone11Digit      |
| `^\+(\d)-(\d{3})-(\d{3})-(\d{4})$`    | `+7-890-123-4567`    | Phone11Digit      |
| `^\+(\d)\.(\d{3})\.(\d{3})\.(\d{4})$` | `+7.890.123.4567`    | Phone11Digit      |
| `^\+(\d) \((\d{3})\) (\d{3})-(\d{4})$` | `+7 (890) 123-4567`  | Phone11Digit      |
| `^\+(\d)\((\d{3})\)(\d{3})-(\d{4})$`  | `+7(890)123-4567`    | Phone11Digit      |
| `^\+(\d)\((\d{3})\)(\d{7})$`          | `+7(890)1234567`     | Phone11Digit      |

Metadata:
- `formatted_display` — `+C-XXX-XXX-XXXX` (C = country code digit)
- `tel_url` — `+` + all 11 digits

#### 5.7.5 Phone Confidence

| Condition                  | Confidence |
|----------------------------|-----------|
| Exact match (input == raw) | 95        |
| Embedded in other text     | 60        |

Additional metadata:
- `raw_number` — matched text
- `is_exact_match` — boolean

---

### 5.8 RaycastParser

**What it detects:** `raycast://` URIs.

**CanHandle:** Input starts with `raycast://` and is a valid URI.

**Metadata:**
- `isAIChat` (boolean) — true if URI contains `extensions/raycast/raycast-ai/ai-chat`
- `isNote` (boolean) — true if URI contains `extensions/raycast/raycast-notes/raycast-notes`

**Confidence:** 85 (content type: `ContentTypeRaycastURI`)

---

### 5.9 OpenCodeSessionParser

**What it detects:** An OpenCode session token (`ses_` prefix).

**CanHandle:** Input contains a token matching `(?i)\bses_[a-z0-9]+\b`
(case-insensitive).

**Metadata:**
- `session_token` (string)
- `is_exact_match` (boolean)

**Confidence:** 90 if exact match, 70 if embedded in other text.

---

### 5.10 CodexParser

**What it detects:** `codex://threads/UUID` URIs.

**CanHandle:** Matches regex:
`^codex://threads/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`
(case-insensitive hex).

**Metadata:**
- `thread_id` (string) — the UUID
- `url` (string) — the full URI

**Confidence:** 90 (content type: `ContentTypeCodexThread`)

---

## 6. Writers

### 6.1 URLWriter

Handles ALL URL-based content types. Vote scores:

| Content Type          | Vote |
|-----------------------|------|
| GitHubLong            | 95   |
| GitHubURL             | 90   |
| JIRAComment           | 95   |
| JIRAURL               | 90   |
| JenkinsURL            | 90   |
| YouTubeURL            | 95   |
| CodeCommitLong        | 95   |
| CodeCommitURL         | 90   |
| NotionURL             | 85   |
| MiniMaxURL            | 90   |
| GeminiURL             | 90   |
| CircleCI              | 90   |
| ChatGPT               | 90   |
| Generic URL           | 50   |
| Everything else       | 0    |

#### 6.1.1 GitHub URL Output

**Format:** `[linkText](originalURL)`

Link text construction:
- If `number` is present and `type` is `"commit"` and len(number) > 7:
  `orgRepo#7charHash` (truncated to first 7 chars)
- If `number` is present: `orgRepo#number`
- Otherwise: `orgRepo`

`orgRepo`:
- Start with `org/repo` from metadata
- Perform case-insensitive lookup against `github.mappings`
- If found, replace with mapped value

#### 6.1.2 GitHub Long Output

**Format:** `[linkText](https://github.com/{org}/{repo}/{type}/{number})`

Link text: `{mappedOrgRepo}#{number}: {title}`

- Apply `github.mappings` (case-insensitive) to org/repo
- Strip leading JIRA key from title (same regex as §5.2.4)
- Default `type` to `"issues"` if empty

If any required metadata is missing (org, repo, title, number), return original
input unchanged.

#### 6.1.3 JIRA URL Output

- Standard: `[issueKey](originalURL)`
- Comment: `[issueKey comment](originalURL)`

#### 6.1.4 Jenkins URL Output

**Format:** `[linkText](originalURL)`

- With build number: `jenkins/{jobName}#{buildNumber}`
- Without: `jenkins/{jobName}`

#### 6.1.5 YouTube URL Output

**Format:** `[linkText](originalURL)`
- Video: `📺 {title}`
- Playlist: `🎥🗃️ {title}`
- If title is empty, fall back to generic URL writer

#### 6.1.6 CodeCommit URL Output

- Short URL: `[{region}/{repo}#{number}](originalURL)`
- Long format: `[{region}/{repo}#{number}: {title}](constructedURL)`

Constructed URL for long format:
`https://{region}.console.aws.amazon.com/codesuite/codecommit/repositories/{repo}/pull-requests/{number}/details?region={region}`

#### 6.1.7 Notion URL Output

`[{title}](originalURL)`
- Falls back to generic URL if title is empty

#### 6.1.8 MiniMax URL Output

`[🤖 MiniMax.io](originalURL)`

#### 6.1.9 Gemini URL Output

`[🤖 Gemini Chat](cleanURL)`
- Uses `clean_url` from metadata if present; otherwise falls back to
  `originalInput`

#### 6.1.10 ChatGPT URL Output

`[🤖 ChatGPT](originalURL)`

#### 6.1.11 CircleCI URL Output

`[🏗️ CircleCI {org}/{repo}#{pipelineNumber}](originalURL)`
- Falls back to generic URL if org, repo, or pipeline_number is missing

#### 6.1.12 Generic URL Output

`[{domain}](originalURL)`

Domain construction:
1. Extract host from URL
2. Strip `www.` prefix
3. Strip `ww3.` prefix
4. Look up domain in `url.domain_mappings`:
   - Convert domain dots to underscores for key lookup
   - Perform **case-insensitive** matching
   - If found, use mapped value as link text
5. Otherwise use the domain as link text

---

### 6.2 JIRAKeyWithDescriptionWriter

**Vote:** 98 for `ContentTypeJIRAKeyWithDescription`, 0 otherwise.

**Output:** `[{issueKey}: {description}]({jiraDomain}/browse/{issueKey})`

---

### 6.3 JIRAWriter

**Vote:** 95 for `ContentTypeJIRAKey`, 0 otherwise.

**Output:** `[{issueKey}]({jiraDomain}/browse/{issueKey})`

---

### 6.4 PhoneWriter

**Vote:** Returns `ctx.Confidence` for Phone7/10/11Digit types; 0 otherwise.

**Output:** `📞 [{formatted_display}](tel:{tel_url})`

---

### 6.5 RaycastWriter

**Vote:** 85 for `ContentTypeRaycastURI`, 0 otherwise.

**Output:** `[{linkText}](originalURI)`

Link text:
- `"Raycast Note"` if `isNote` is true
- `"Raycast AI"` if `isAIChat` is true
- `"Raycast"` otherwise

---

### 6.6 OpenCodeSessionWriter

**Vote:** Returns `ctx.Confidence` for `ContentTypeOpenCodeSession`; 0 otherwise.

**Output:** `[🤖 OpenCode](opencode://session/{sessionToken})`

---

### 6.7 CodexWriter

**Vote:** Returns `ctx.Confidence` for `ContentTypeCodexThread`; 0 otherwise.

**Output:** `[🤖 Codex]({url})`

---

### 6.8 PassthroughWriter

**Vote:** Always returns 1 (lowest priority).

**Write:** Returns `ctx.OriginalInput` unchanged.

---

## 7. Configuration

### 7.1 Location and Loading

Default config path: `~/.config/markdown-tool/config.yaml`

On first run, if no config file exists, one is created with defaults.
Override with `--config` flag.

### 7.2 YAML Schema

```yaml
github:
  default_org: "CompanyCam"            # fallback org for simple issue titles
  default_repo: "Company-Cam-API"      # fallback repo for simple issue titles
  mappings:                            # case-insensitive org/repo -> display name
    "CompanyCam/Company-Cam-API": "CompanyCam/API"

jira:
  domain: "https://companycam.atlassian.net"
  projects:                            # only these project keys are transformed
    - "PLAT"
    - "SPEED"

jenkins:
  domain: "https://jenkins.internal.upserve.com"

url:
  domain_mappings:                     # case-insensitive domain -> display name
    companycam_slack_com: "slack"      # dots in domain become underscores
    youtube_com: "YouTube"
```

### 7.3 Key Behaviors

- **GitHub mappings:** Lookup is case-insensitive (because YAML keys are
  lowercased by Viper). If `org/repo` matches a mapping key (case-insensitive),
  the mapped value replaces it for display.
- **JIRA domain:** Required for JIRA URL detection; host comparison only.
  If domain is empty, no JIRA URLs are detected.
- **Jenkins domain:** Required for Jenkins URL detection. Empty domain = no
  Jenkins URLs detected.
- **JIRA projects:** Only configured project keys produce output. Unconfigured
  keys pass through unchanged. The parser enforces this gate, not the writer.
- **URL domain mappings:** Keys use underscores instead of dots. Lookup is
  case-insensitive. Falls back to hostname if no mapping found.
- **Viper key delimiter:** The config loader uses `::` as the key delimiter
  instead of `.` to prevent domain names with dots from being interpreted as
  nested YAML structures.

### 7.4 Default Config (written on first run)

```yaml
github:
  default_org: "CompanyCam"
  default_repo: "Company-Cam-API"
  mappings:
    "CompanyCam/Company-Cam-API": "CompanyCam/API"

jira:
  domain: "https://companycam.atlassian.net"
  projects:
    - "PLAT"
    - "SPEED"

jenkins:
  domain: "https://jenkins.internal.upserve.com"

url:
  domain_mappings:
    companycam_slack_com: "slack"
    youtube_com: "YouTube"
```

---

## 8. CLI Flow

### 8.1 Input Sources

1. Check if data is piped via stdin (`os.ModeCharDevice` not set)
2. If piped, read all from stdin
3. Otherwise, read from system clipboard

### 8.2 Preprocessing

Before parsing, strip any `tel:` URI prefix:

```
Regex: ^tel:(.*)$
Action: Replace with capture group 1 (the phone number)
```

This is case-sensitive — `TEL:1234567` is NOT stripped. If the result is
empty (e.g., input was `tel:`), the pipeline stops and outputs nothing.

### 8.3 Processing

```
input = trim(input)
if input == "": exit 0

input = preprocessTelURIs(input)
if input == "": exit 0

contexts = []
for parser in parsers:
    ctx = parser.Parse(input)
    if ctx != nil:
        contexts.append(ctx)

bestWriter, bestScore = vote(writers, contexts)

if bestWriter == nil or bestScore == 0 or len(contexts) == 0:
    print(input)       // verbatim
    exit 0

output = bestWriter.Write(contexts[0])
print(output)
```

### 8.4 CLI Flags

| Flag          | Short | Description                                      |
|---------------|-------|--------------------------------------------------|
| `--verbose`   | `-v`  | Enable verbose logging to stderr                 |
| `--config`    |       | Override config file path                        |

### 8.5 Error Handling

- Config load failure → exit with error message to stderr
- Input read failure → exit with error
- Write error → exit with error
- Parse errors from individual parsers are silenced (the parser just returns
  nil)

---

## 9. Test Vectors

Each entry shows: input → expected output. All tests use the default config
from §7.4 unless otherwise noted.

### 9.1 GitHub

```
// PR URL with mapping
https://github.com/CompanyCam/Company-Cam-API/pull/15217
→ [CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)

// Issue URL with mapping
https://github.com/CompanyCam/Company-Cam-API/issues/15217
→ [CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/issues/15217)

// Commit long hash (truncated)
https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b
→ [ErebusBat/markdown-tool#aa062a6](https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b)

// Commit short hash
https://github.com/CompanyCam/Company-Cam-API/commit/abc123
→ [CompanyCam/API#abc123](https://github.com/CompanyCam/Company-Cam-API/commit/abc123)

// Repository URL (no mapping)
https://github.com/pedropark99/zig-book
→ [pedropark99/zig-book](https://github.com/pedropark99/zig-book)

// Repository with mapping
https://github.com/CompanyCam/Company-Cam-API
→ [CompanyCam/API](https://github.com/CompanyCam/Company-Cam-API)

// Case-insensitive mapping
https://github.com/Company/Long-Repo-Name/pull/123
(with mapping "company/long-repo-name": "Company/Short")
→ [Company/Short#123](https://github.com/Company/Long-Repo-Name/pull/123)
```

### 9.2 GitHub Long Format (multi-line UI)

```
// Standard issue
CompanyCam
companycam-mobile

Type / to search
Code
Issues
78
Pull requests
12
...
A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog #6549
→ [CompanyCam/companycam-mobile#6549: A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog](https://github.com/CompanyCam/companycam-mobile/issues/6549)

// PR with split number and mapping
dotswipely
weekly_digest_email_pipeline
Repository navigation
...
Agents
...
Update list of merchants to validate against
 #286
→ [dotswipely/weekly_digest_email_pipeline#286: Update list of merchants to validate against](https://github.com/dotswipely/weekly_digest_email_pipeline/pull/286)
```

### 9.3 Simple Issue Titles

```
// Plain
adds blinc ddagent file #15407
→ [CompanyCam/API#15407: adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)

// With username prefix
courtneylw adds blinc ddagent file #15407
→ [CompanyCam/API#15407: adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)

// Username with underscore
plat_188 adds blinc ddagent file #15407
→ [CompanyCam/API#15407: adds blinc ddagent file](https://github.com/CompanyCam/Company-Cam-API/issues/15407)
```

### 9.4 JIRA

```
// URL
https://companycam.atlassian.net/browse/PLAT-192
→ [PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)

// Comment URL
https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266
→ [PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)

// Standalone key
PLAT-12345
→ [PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)

// Unconfigured key (passes through)
INVALID-123
→ INVALID-123

// Key with description
PLAT-192

blinc - webhook proxy logs
→ [PLAT-192: blinc - webhook proxy logs](https://companycam.atlassian.net/browse/PLAT-192)

// Multi-line description
PLAT-789

Fix authentication issue with SSO
Additional details about the bug
→ [PLAT-789: Fix authentication issue with SSO Additional details about the bug](https://companycam.atlassian.net/browse/PLAT-789)

// Unconfigured key with description (passes through)
INVALID-123

This should not be transformed
→ INVALID-123\n\nThis should not be transformed
```

### 9.5 Jenkins

```
// Build URL
https://jenkins.internal.upserve.com/job/app.swipely/114/
→ [jenkins/app.swipely#114](https://jenkins.internal.upserve.com/job/app.swipely/114/)

// Build with console text
https://jenkins.internal.upserve.com/job/app.swipely/114/consoleText
→ [jenkins/app.swipely#114](https://jenkins.internal.upserve.com/job/app.swipely/114/consoleText)

// lastBuild (no number)
https://jenkins.internal.upserve.com/job/app.swipely/lastBuild/
→ [jenkins/app.swipely](https://jenkins.internal.upserve.com/job/app.swipely/lastBuild/)

// Bare job URL
https://jenkins.internal.upserve.com/job/my-project/
→ [jenkins/my-project](https://jenkins.internal.upserve.com/job/my-project/)

// lastSuccessfulBuild
https://jenkins.internal.upserve.com/job/app.swipely/lastSuccessfulBuild/consoleText
→ [jenkins/app.swipely](https://jenkins.internal.upserve.com/job/app.swipely/lastSuccessfulBuild/consoleText)
```

### 9.6 YouTube

```
// Video (with oEmbed title fetch)
https://www.youtube.com/watch?v=fkT41ooKBuY
→ [📺 Stop overpaying for OpenAI: Multi-model routing guide](https://www.youtube.com/watch?v=fkT41ooKBuY)

// Playlist
https://www.youtube.com/playlist?list=PLCC34OHNcOtpcgR9LEYSdi9r7XIbpkpK1
→ [🎥🗃️ Deep Learning With PyTorch](https://www.youtube.com/playlist?list=PLCC34OHNcOtpcgR9LEYSdi9r7XIbpkpK1)
```

### 9.7 Notion

```
https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0
→ [VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0)
```

### 9.8 Generic URLs and Domain Mappings

```
// Plain generic URL
http://ww3.domain.tld/path/to/document?query=value#anchor
→ [domain.tld](http://ww3.domain.tld/path/to/document?query=value#anchor)

// Domain mapping (dots → underscores)
https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069
(with mapping: companycam_slack_com → "slack")
→ [slack](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)

// Case-insensitive domain mapping
https://CompanyCam.Slack.com/archives/test
(with mapping: companycam_slack_com → "slack")
→ [slack](https://CompanyCam.Slack.com/archives/test)

// Unmapped domain
https://example.com/path/to/page
→ [example.com](https://example.com/path/to/page)
```

### 9.9 Special URLs

```
// Gemini Chat
https://gemini.google.com/app/ac9ebc9d76c30fc1
→ [🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)

// Gemini with trailing junk (uses clean_url)
https://gemini.google.com/app/ac9ebc9d76c30fc1 →
→ [🤖 Gemini Chat](https://gemini.google.com/app/ac9ebc9d76c30fc1)

// Gemini root (stays generic)
https://gemini.google.com/app
→ [gemini.google.com](https://gemini.google.com/app)

// ChatGPT
https://chatgpt.com/c/69efd1c6-a230-83e8-8778-5dc7754dcdd3
→ [🤖 ChatGPT](https://chatgpt.com/c/69efd1c6-a230-83e8-8778-5dc7754dcdd3)

// CircleCI
https://app.circleci.com/pipelines/github/upserve/swipely/96/workflows/17abd9c6-1190-49e9-a05f-4bf992a9d611
→ [🏗️ CircleCI upserve/swipely#96](https://app.circleci.com/pipelines/github/upserve/swipely/96/workflows/17abd9c6-1190-49e9-a05f-4bf992a9d611)

// MiniMax
https://agent.minimax.io/?id=abc123
→ [🤖 MiniMax.io](https://agent.minimax.io/?id=abc123)

// CodeCommit short URL
https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1
→ [us-east-1/upserve-env#411](https://us-east-1.console.aws.amazon.com/codesuite/codecommit/repositories/upserve-env/pull-requests/411/details?region=us-east-1)
```

### 9.10 Phone Numbers

```
// 7-digit
1234567     → 📞 [123-4567](tel:1234567)
123-4567    → 📞 [123-4567](tel:1234567)
123.4567    → 📞 [123-4567](tel:1234567)

// 10-digit
8901234567       → 📞 [890-123-4567](tel:8901234567)
890-123-4567     → 📞 [890-123-4567](tel:8901234567)
890.123.4567     → 📞 [890-123-4567](tel:8901234567)
(890) 123-4567   → 📞 [890-123-4567](tel:8901234567)
(890)123-4567    → 📞 [890-123-4567](tel:8901234567)
(890)1234567     → 📞 [890-123-4567](tel:8901234567)

// 11-digit US
18901234567         → 📞 [1-890-123-4567](tel:+18901234567)
1-890-123-4567      → 📞 [1-890-123-4567](tel:+18901234567)
1.890.123.4567      → 📞 [1-890-123-4567](tel:+18901234567)
1 (890) 123-4567    → 📞 [1-890-123-4567](tel:+18901234567)
1(890)123-4567      → 📞 [1-890-123-4567](tel:+18901234567)
1(890)1234567       → 📞 [1-890-123-4567](tel:+18901234567)

// 11-digit international
+78901234567        → 📞 [+7-890-123-4567](tel:+78901234567)
+7-890-123-4567     → 📞 [+7-890-123-4567](tel:+78901234567)
+7.890.123.4567     → 📞 [+7-890-123-4567](tel:+78901234567)
+7 (890) 123-4567   → 📞 [+7-890-123-4567](tel:+78901234567)
+7(890)123-4567     → 📞 [+7-890-123-4567](tel:+78901234567)
+7(890)1234567      → 📞 [+7-890-123-4567](tel:+78901234567)

// Non-matches (pass through)
123 4567        → 123 4567
123,4567        → 123,4567
01234567        → 01234567
89012345670     → 89012345670
890 123 4567    → 890 123 4567
(890) 123 4567  → (890) 123 4567
(890) 1234 567  → (890) 1234 567
(890)123-456    → (890)123-456
(890)12345679   → (890)12345679
```

### 9.11 Tel URI Preprocessing

```
tel:1234567          → 📞 [123-4567](tel:1234567)
tel:8901234567       → 📞 [890-123-4567](tel:8901234567)
tel:18901234567      → 📞 [1-890-123-4567](tel:+18901234567)
tel:+18901234567      → 📞 [+1-890-123-4567](tel:+18901234567)
tel:890-123-4567     → 📞 [890-123-4567](tel:8901234567)
tel:(890)123-4567    → 📞 [890-123-4567](tel:8901234567)
tel:                 → (empty output)
TEL:1234567          → TEL:1234567  (case-sensitive, no match)
```

### 9.12 Raycast

```
// AI Chat
raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22%3A%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D
→ [Raycast AI](raycast://extensions/raycast/raycast-ai/ai-chat?context=...)

// Note
raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22%3A%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D
→ [Raycast Note](raycast://extensions/raycast/raycast-notes/raycast-notes?context=...)

// Generic extension
raycast://extensions/other/extension
→ [Raycast](raycast://extensions/other/extension)

// Settings
raycast://settings
→ [Raycast](raycast://settings)
```

### 9.13 Codex

```
codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf
→ [🤖 Codex](codex://threads/019dcc44-e7b8-7c23-816c-34c194bdb3cf)
```

### 9.14 OpenCode Session

```
ses_2017f15ceffeK5CZjD3EX3fHnW
→ [🤖 OpenCode](opencode://session/ses_2017f15ceffeK5CZjD3EX3fHnW)
```

### 9.15 Non-Matches (Verbatim)

```
hello world     → hello world
Check out PLAT-999  → Check out PLAT-999  (embedded text, not standalone)
```

---

## 10. Extensibility

### 10.1 Adding a New Parser

1. Define a new `ContentType` constant
2. Implement the `Parser` interface (§4.2):
   - `CanHandle(input) -> boolean` — quick, cheap check
   - `Parse(input) -> (ParseContext | null, error | null)` — return `(nil, nil)`
     when input doesn't match
3. Register in the parser list (§2.1) at the appropriate position

### 10.2 Adding a New Writer

1. Implement the `Writer` interface (§4.3):
   - `Vote(ctx) -> integer` — return 0 for unrelated content types
   - `Write(ctx) -> string` — return `ctx.OriginalInput` as safe fallback
   - `GetName() -> string` — human-readable for logging
2. Register in the writer list (§2.2)

### 10.3 Metadata Keys Reference

| Key                 | Used By                          | Type     |
|---------------------|----------------------------------|----------|
| `org`               | URLWriter (GitHub, CircleCI)     | string   |
| `repo`              | URLWriter (GitHub, CodeCommit, CircleCI) | string |
| `number`            | URLWriter (GitHub, CodeCommit)   | string   |
| `type`              | URLWriter (GitHub issue type)    | string   |
| `title`             | URLWriter (GitHub Long, YouTube, Notion, CodeCommit Long) | string |
| `issue_key`         | URLWriter, JIRAWriter, JIRAKWDW  | string   |
| `project`           | (parser gate only)               | string   |
| `description`       | JIRAKeyWithDescriptionWriter     | string   |
| `comment_id`        | URLWriter (JIRA comment)         | string   |
| `job_name`          | URLWriter (Jenkins)              | string   |
| `build_number`      | URLWriter (Jenkins)              | string   |
| `youtube_type`      | URLWriter (YouTube)              | string   |
| `video_id`          | URLWriter (YouTube)              | string   |
| `playlist_id`       | URLWriter (YouTube)              | string   |
| `domain`            | URLWriter (generic)              | string   |
| `region`            | URLWriter (CodeCommit)           | string   |
| `vcs`               | URLWriter (CircleCI)             | string   |
| `pipeline_number`   | URLWriter (CircleCI)             | string   |
| `workflow_id`       | URLWriter (CircleCI)             | string   |
| `chat_id`           | URLWriter (MiniMax, Gemini, ChatGPT) | string |
| `clean_url`         | URLWriter (Gemini)               | string   |
| `raw_number`        | PhoneWriter                      | string   |
| `formatted_display` | PhoneWriter                      | string   |
| `tel_url`           | PhoneWriter                      | string   |
| `is_exact_match`    | OpenCodeSessionWriter, PhoneWriter | bool   |
| `isAIChat`          | RaycastWriter                    | bool     |
| `isNote`            | RaycastWriter                    | bool     |
| `session_token`     | OpenCodeSessionWriter            | string   |
| `thread_id`         | CodexWriter                      | string   |
| `url`               | CodexWriter                      | string   |

---
> Generated from the Go source at `github.com/ErebusBat/markdown-tool` — version as of 2026-05-12.
