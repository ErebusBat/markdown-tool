# Product Requirements Document: markdown-tool

## 1. Overview

### Product Vision
A lightweight, command-line Go application that processes small text inputs and transforms them into well-formatted markdown suitable for knowledge management tools like Vimwiki and Obsidian.

### Problem Statement
Users frequently need to convert plain text, code snippets, URLs, and other content into properly formatted markdown for their note-taking and knowledge management workflows. Current solutions are often web-based, feature-heavy, or don't integrate well with command-line workflows.

### Target Users
- Developers and technical writers
- Knowledge workers using markdown-based note-taking systems
- Users of Vimwiki, Obsidian, and similar tools
- Command-line enthusiasts
- Power users

## 2. Functional Requirements

### Core Features
1. **Store Specific Options in a YAML configuration file**
    - File should be stored in `~/.config/markdown-tool/`
    - Config should use spf13/viper and spf13/cobra
    - Should contain:
        - GitHub
            - Default Organization and Project names
            - A map to re-write descriptions of Organization and projects
                - For example we may want and entry so that we can later change the description for an pull request in `CompanyCam/Company-Cam-API` to something shorter like `CompanyCam/API`.
        - JIRA
            - Domain, i.e.: `https://companycam.atlassian.net`
            - List of valid projects: i.e.: `PLAT`

2. **Detect and Transform URLs**
    - If no other match succeeds then it should transform the URL into simple markdown:
       - Input: `http://ww3.domain.tld/path/to/document?query=value#anchor`
       - Output: `[domain.tld](http://ww3.domain.tld/path/to/document?query=value#anchor)`
    - Detect GitHub links and use the org/repo#issue number as the description:
        - For Pull Requests and Issues:
            - Inputs:
                - `https://github.com/CompanyCam/Company-Cam-API/pull/15217`
                - `https://github.com/CompanyCam/Company-Cam-API/issues/15217`
            - Output: `[CompanyCam/Company-Cam-API#15217](https://github.com/CompanyCam/Company-Cam-API/pull/15217)`
        - For Commits (truncate hash to 7 characters in link text):
            - Input: `https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b`
            - Output: `[ErebusBat/markdown-tool#aa062a6](https://github.com/ErebusBat/markdown-tool/commit/aa062a602a02d33f4a6e7880809ac3609fe1417b)`
    - Detect Jira Issue Links and use the issue number as the description:
        - Input: `https://companycam.atlassian.net/browse/PLAT-192`
        - Output: `[PLAT-192](https://companycam.atlassian.net/browse/PLAT-192)`
    - Detect Jira Comment Links and use the issue number and indicate it is a comment:
        - Input: `https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266`
        - Output: `[PLAT-192 comment](https://companycam.atlassian.net/browse/PLAT-192?focusedCommentId=20266)`
    - Detect Notion Links and extract the page title as the comment:
        - Input: `https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0#1c0d42d77c0b80268626fa64eb6ebdbe`
        - Output: `[VS Code Setup for Standard rb RubyLSP](https://www.notion.so/companycam/VS-Code-Setup-for-Standard-rb-RubyLSP-654a6b070ae74ac3ad400c6d571507c0#1c0d42d77c0b80268626fa64eb6ebdbe)`
    - Detect Raycast URIs and transform them into markdown links:
        - For AI Chat URIs (containing `extensions/raycast/raycast-ai/ai-chat`):
            - Input: `raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D`
            - Output: `[Raycast AI](raycast://extensions/raycast/raycast-ai/ai-chat?context=%7B%22id%22:%228926C709-D08B-4FFC-9FD8-7A0E5561156D%22%7D)`
        - For Note URIs (containing `extensions/raycast/raycast-notes/raycast-notes`):
            - Input: `raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22:%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D`
            - Output: `[Raycast Note](raycast://extensions/raycast/raycast-notes/raycast-notes?context=%7B%22id%22:%22C8411E30-ADD9-4BBA-BFA5-2B14AE3DB533%22%7D)`
        - For other Raycast URIs:
            - Input: `raycast://extensions/other/extension`
            - Output: `[Raycast](raycast://extensions/other/extension)`

3. **Detect JIRA Issue Keys**
    - Should detect a JIRA key and transform
        - Input: `PLAT-12345`
        - Output: `[PLAT-12345](https://companycam.atlassian.net/browse/PLAT-12345)`
    - Should only detect keys that are setup in the configuration file

4. **Detect and Transform Phone Numbers**
    - Should detect phone numbers in various formats and transform them to markdown tel: links
    - **7-digit numbers**: Match only specific separators, output in standard format
        - Matches: `1234567`, `123-4567`, `123.4567`
        - Does not match: `123 4567`, `123,4567`, `01234567` (leading zero)
        - Output: `📞 [123-4567](tel:1234567)`
    - **10-digit numbers**: Match various common formats, output in standard format
        - Matches: `8901234567`, `890-123-4567`, `890.123.4567`, `(890) 123-4567`, `(890)123-4567`, `(890)1234567`
        - Does not match: `89012345670` (extra digit), `890 123 4567` (spaces), `(890) 123 4567` (mixed separators)
        - Output: `📞 [890-123-4567](tel:8901234567)`
    - **11-digit numbers with country code**: Support US and international formats
        - US numbers (country code 1): `18901234567`, `1-890-123-4567`, `1.890.123.4567`, `1 (890) 123-4567`, `1(890)123-4567`, `1(890)1234567`
        - International numbers (must have + prefix): `+78901234567`, `+7-890-123-4567`, `+7.890.123.4567`, `+7 (890) 123-4567`, `+7(890)123-4567`, `+7(890)1234567`
        - US output: `📞 [1-890-123-4567](tel:+18901234567)`
        - International output: `📞 [+7-890-123-4567](tel:+78901234567)`
    - **Confidence scoring**: Phone number detection should use confidence levels that affect voting
        - Exact phone number match (input is only the phone number): High confidence (95)
        - Phone number embedded in other text: Lower confidence (60-70)
        - This ensures phone numbers don't override other content types inappropriately

5. If nothing matches then the text should be output verbatim

### Input Methods
- Standard input (stdin) for piping
- Clipboard; only if nothing was specified on stdin

### Output Options
- Standard output (stdout) for piping
- Standard error (stderr) for logging and debugging

## 3. Technical Requirements

### Language and Runtime
- **Language**: Go (latest stable version)
- **Target Platforms**: Linux, macOS
- **Distribution**: Single binary executable

### Dependencies
- Minimal external dependencies
- Must use the following dependencies:
    - https://github.com/spf13/Viper
        - For configuration
    - https://github.com/spf13/Cobra
        - For application structure / Command line interface
- If not specified then Standard library preferred
- Each item should use structured logging
    - Default log output should be configured to write to stderr in a pretty / colorful human readable way

### Performance
- This will be a quick executable and then it will die; we don't need to worry about resource leakage as the OS will handle that for us

### Architecture
- CLI interface using cobra
- Modular design with separate processors for different content types
- Plugin-like architecture for extensibility
    - Requiring a re-compile for different "plugins" is fine
- Tests should be written to support each given input and output case
- There should be three phases:
    - **Parsing phase**
        - All parsers will run and add information to a global context object
    - **Voting phase**
        - Each output writer will "vote" on if it thinks it should output
    - **Output phase**
        - After all output writers have voted then the writer with the highest vote count will get the context and write the transformed text

## 4. User Stories

### As a Developer
- I want to process URLs into proper markdown links for documentation
- I want to convert GitHub and JIRA issue keys into full markdown links for documentation
- I want to convert phone numbers into clickable tel: links for documentation and notes

## 5. Acceptance Criteria

### Configuration Management
- [ ] Application creates configuration directory at `~/.config/markdown-tool/` if it doesn't exist
- [ ] Configuration file is properly loaded using Viper with GitHub and JIRA settings
- [ ] Invalid configuration displays helpful error messages
- [ ] If no configuration file is present then a default one should be created;
    - [ ] GitHub Org: `CompanyCam`
    - [ ] GitHub: `Company-Cam-API`
    - [ ] JIRA Domain: `companycam.atlassian.net`
    - [ ] JIRA PROJECTS: [`PLAT`, `SPEED`]
- [ ] Missing configuration falls back to reasonable defaults

### URL Processing
- [ ] Generic URLs are converted to `[domain.tld](full-url)` format
- [ ] GitHub PR/issue URLs are converted to `[org/repo#number](url)` format
- [ ] GitHub commit URLs are converted to `[org/repo#hash7chars](url)` format with truncated hash
- [ ] JIRA issue URLs are converted to `[ISSUE-KEY](url)` format
- [ ] JIRA comment URLs are converted to `[ISSUE-KEY comment](url)` format
- [ ] Notion URLs extract page titles from URL slugs
- [ ] Raycast AI Chat URIs (containing `extensions/raycast/raycast-ai/ai-chat`) are converted to `[Raycast AI](url)` format
- [ ] Raycast Note URIs (containing `extensions/raycast/raycast-notes/raycast-notes`) are converted to `[Raycast Note](url)` format
- [ ] Other Raycast URIs are converted to `[Raycast](url)` format
- [ ] Organization/project name mappings from config are applied correctly

### JIRA Key Detection
- [ ] Standalone JIRA keys (e.g., `PLAT-12345`) are converted to full markdown links
- [ ] Only configured project keys are recognized
- [ ] Invalid or unconfigured project keys are left unchanged

### Phone Number Processing
- [ ] 7-digit phone numbers are converted to `📞 [123-4567](tel:1234567)` format
- [ ] 10-digit phone numbers are converted to `📞 [890-123-4567](tel:8901234567)` format
- [ ] 11-digit US numbers are converted to `📞 [1-890-123-4567](tel:+18901234567)` format
- [ ] International numbers with + prefix are converted to `📞 [+7-890-123-4567](tel:+78901234567)` format
- [ ] Only specific separators are recognized (none, dash, dot, parentheses)
- [ ] Space separators and mixed formats are rejected
- [ ] Exact phone number matches have high confidence (95)
- [ ] Embedded phone numbers have lower confidence (60-70)
- [ ] Invalid formats are left unchanged

### Input/Output Handling
- [ ] Processes text from stdin when available
- [ ] Falls back to clipboard content when stdin is empty
- [ ] Outputs transformed markdown to stdout
- [ ] Logs debug information to stderr
- [ ] Handles empty input gracefully

### Parsing Architecture
- [ ] Parsing phase: All parsers analyze input and populate context
- [ ] Voting phase: Output writers vote on confidence levels
- [ ] Output phase: Highest-confidence writer generates final output
- [ ] Unmatched text is output verbatim

### Error Handling
- [ ] Invalid URLs are left unchanged
- [ ] Network failures (if any) don't crash the application
- [ ] Malformed input produces helpful error messages
- [ ] Edge cases (empty input, special characters) are handled gracefully

### Command Line Interface
- [ ] Application runs without flags and processes input/output
- [ ] Help documentation is clear and comprehensive
- [ ] Version information is available
- [ ] Structured logging outputs to stderr in human-readable format

### Testing
- [ ] Unit tests cover all URL transformation scenarios
- [ ] Unit tests cover all phone number transformation scenarios
- [ ] Integration tests verify end-to-end functionality
- [ ] Edge case tests ensure robust error handling
- [ ] Configuration loading tests validate Viper integration
- [ ] Phone number confidence scoring tests validate voting behavior

## 7. Non-Functional Requirements

### Usability
- Intuitive command-line interface
- Clear help documentation
- Ability to be called from other tools and work
    - i.e. if invoked without `--help` or `help` it should just try to parse/output

### Reliability
- Handle edge cases gracefully
- Provide meaningful error messages
- Output unmatched text as is

### Maintainability
- Well-documented code
- Comprehensive test coverage
- Modular architecture for easy extension
- Ability to easily add more test cases for the given use cases without having to modify this PRD document

## 8. Future Considerations

### Potential Enhancements
- Ability to detect larger blocks of copied text and extract GitHub / JIRA information from that
- Connect to websites to get the actual titles of items
    - Would require an auth system
    - For example:
        - Youtube
        - Notion
        - Confluence
- Plugin system for custom processors
- Advanced formatting options

---

*This PRD serves as the foundation for developing the markdown-tool application and should be reviewed and updated as requirements evolve.*
