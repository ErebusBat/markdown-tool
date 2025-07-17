# Markdown Tool

A Go-based command-line tool for transforming text inputs into well-formatted markdown links. The tool intelligently detects URLs (GitHub, JIRA, Notion, generic), JIRA issue keys, and GitHub UI content, then converts them into proper markdown link format.

## Overview

The tool processes input through a three-phase architecture: Parse → Vote → Write. Multiple parsers analyze the input, writers vote on their confidence to handle the content, and the highest-confidence writer generates the transformed output.

| Input Type | Example Input | Example Output | Notes |
|------------|---------------|----------------|-------|
| **GitHub PR URL** | `https://github.com/CompanyCam/Company-Cam-API/pull/15217` | `[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)` | Supports org/repo name mapping |
| **GitHub Issue URL** | `https://github.com/CompanyCam/Company-Cam-API/issues/15217` | `[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/issues/15217)` | Supports org/repo name mapping |
| **GitHub Long Format** | Multi-line GitHub UI text | `[CompanyCam/mobile#6549: Issue Title](https://github.com/CompanyCam/companycam-mobile/issues/6549)` | Parses copied GitHub UI content |
| **JIRA Issue URL** | `https://companycam.atlassian.net/browse/PLAT-192` | `[PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)` | Domain configurable |
| **JIRA Comment URL** | `https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266` | `[PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)` | Detects comment URLs |
| **JIRA Key** | `PLAT-12345` | `[PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)` | Only configured projects |
| **JIRA Key + Description** | `PLAT-192`<br/><br/>`webhook proxy logs` | `[PLAT-192: webhook proxy logs](https://companycam.atlassian.net/browse/PLAT-192)` | Multi-line format support |
| **Notion URL** | `https://www.notion.so/companycam/VS-Code-Setup-654a6b07...` | `[VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-654a6b07...)` | Extracts title from URL |
| **Generic URL** | `https://www.example.com/path/to/page` | `[example.com](https://www.example.com/path/to/page)` | Strips www/ww* prefixes |
| **URL with Domain Mapping** | `https://companycam.slack.com/archives/D08UZ6X17MJ/...` | `[slack](https://companycam.slack.com/archives/D08UZ6X17MJ/...)` | Custom domain display names |
| **Plain Text** | `hello world` | `hello world` | Passed through unchanged |

## GitHub URLs

The tool handles GitHub pull requests and issues with intelligent org/repo name mapping.

### GitHub Pull Request URLs

**Input:**
```
https://github.com/CompanyCam/Company-Cam-API/pull/15217
```

**Output:**
```
[CompanyCam/API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)
```

### GitHub Issue URLs

**Input:**
```
https://github.com/someorg/somerepo/issues/42
```

**Output:**
```
[someorg/somerepo#42](https://github.com/someorg/somerepo/issues/42)
```

### GitHub Long Format (UI Content)

The tool can parse multi-line content copied from GitHub's web interface:

**Input:**
```
CompanyCam
companycam-mobile

Type / to search
Code
Issues
78
Pull requests
12
Actions
Projects
Wiki
Security
7
Insights
A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog #6549
```

**Output:**
```
[CompanyCam/mobile#6549: A specific Logger.error call in the SSO login workflow doesn't seem to log data to Datadog](https://github.com/CompanyCam/companycam-mobile/issues/6549)
```

#### GitHub Configuration

Org/repo names can be remapped via configuration:

```yaml
github:
  mappings:
    "companycam/company-cam-api": "CompanyCam/API"
    "companycam/companycam-mobile": "CompanyCam/mobile"
```

This allows displaying shorter, more readable names in the markdown output while preserving the actual repository URLs.

## JIRA

The tool supports multiple JIRA input formats and requires configuration to specify valid projects and domain.

### JIRA Issue URLs

**Input:**
```
https://companycam.atlassian.net/browse/PLAT-192
```

**Output:**
```
[PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)
```

### JIRA Comment URLs

The tool detects when a JIRA URL includes a comment focus parameter:

**Input:**
```
https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266
```

**Output:**
```
[PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)
```

### JIRA Keys (Standalone)

**Input:**
```
PLAT-12345
```

**Output:**
```
[PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)
```

### JIRA Keys with Description

The tool supports a multi-line format where a JIRA key is followed by a description:

**Input:**
```
PLAT-192

blinc - webhook proxy logs
```

**Output:**
```
[PLAT-192: blinc - webhook proxy logs](https://companycam.atlassian.net/browse/PLAT-192)
```

#### Multi-line Descriptions

Descriptions spanning multiple lines are concatenated with spaces:

**Input:**
```
PLAT-789

Fix authentication issue with SSO
Additional details about the bug
```

**Output:**
```
[PLAT-789: Fix authentication issue with SSO Additional details about the bug](https://companycam.atlassian.net/browse/PLAT-789)
```

#### JIRA Configuration

JIRA processing requires configuration specifying the domain and valid project keys:

```yaml
jira:
  domain: "https://companycam.atlassian.net"
  projects: ["PLAT", "SPEED"]
```

Only JIRA keys matching the configured projects will be transformed. Unconfigured projects are output verbatim:

**Input:** `INVALID-123` → **Output:** `INVALID-123`

## Notion URLs

The tool extracts page titles from Notion URLs by parsing the URL slug:

**Input:**
```
https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0
```

**Output:**
```
[VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0)
```

## Generic URLs

For URLs that don't match specific patterns, the tool creates a clean link using the domain name:

**Input:**
```
https://www.example.com/path/to/page
```

**Output:**
```
[example.com](https://www.example.com/path/to/page)
```

The tool intelligently strips common prefixes like `www.`, `ww2.`, `ww3.`, etc.:

**Input:**
```
http://ww3.domain.tld/path/to/document?query=value#anchor
```

**Output:**
```
[domain.tld](http://ww3.domain.tld/path/to/document?query=value#anchor)
```

### URL Domain Mappings

For frequently used domains, you can configure custom display names instead of using the domain name:

**Input:**
```
https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069
```

**Output with domain mapping:**
```
[slack](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)
```

**Output without domain mapping:**
```
[companycam.slack.com](https://companycam.slack.com/archives/D08UZ6X17MJ/p1752272874485069)
```

#### More Examples

| Input URL | Mapped Output | Unmapped Output |
|-----------|---------------|-----------------|
| `https://youtube.com/watch?v=abc123` | `[YouTube](https://youtube.com/watch?v=abc123)` | `[youtube.com](https://youtube.com/watch?v=abc123)` |
| `https://CompanyCam.Slack.com/archives/test` | `[slack](https://CompanyCam.Slack.com/archives/test)` | `[companycam.slack.com](https://CompanyCam.Slack.com/archives/test)` |

Domain mapping is case-insensitive, so `CompanyCam.Slack.com` matches the mapping for `companycam.slack.com`.

## Configuration

The tool uses YAML configuration stored in `~/.config/markdown-tool/config.yaml`:

```yaml
github:
  default_org: "CompanyCam"
  default_repo: "Company-Cam-API"
  mappings:
    "companycam/company-cam-api": "CompanyCam/API"
    "companycam/companycam-mobile": "CompanyCam/mobile"

jira:
  domain: "https://companycam.atlassian.net"
  projects: ["PLAT", "SPEED"]

url:
  domain_mappings:
    companycam_slack_com: "slack"
    youtube_com: "YouTube"
```

## Architecture

The tool follows a three-phase processing architecture:

1. **Parsing Phase**: Multiple parsers analyze the input and populate shared context objects
2. **Voting Phase**: Output writers vote on their confidence to handle the parsed content (0-100)
3. **Output Phase**: The highest-confidence writer generates the final transformed output

### Parsers

- **URLParser**: Handles all URL types (GitHub, JIRA, Notion, generic)
- **GitHubLongParser**: Processes multi-line GitHub UI content
- **JIRAKeyParser**: Handles standalone JIRA keys
- **JIRAKeyWithDescriptionParser**: Handles JIRA keys with descriptions

### Writers

- **URLWriter**: Transforms URL-based content (confidence: 50-95)
- **JIRAWriter**: Transforms standalone JIRA keys (confidence: 95)
- **JIRAKeyWithDescriptionWriter**: Transforms JIRA keys with descriptions (confidence: 98)
- **PassthroughWriter**: Outputs input unchanged (confidence: 1)

## Usage

```bash
# From stdin
echo "PLAT-12345" | markdown-tool

# From clipboard (fallback when no stdin)
markdown-tool

# Output
[PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)
```

The tool writes transformed output to stdout and debug information to stderr.
