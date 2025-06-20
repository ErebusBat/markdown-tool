# JIRA Key with Description Examples

This document shows examples of JIRA key inputs with descriptions and their expected markdown transformations.

## Input Format

The tool accepts JIRA keys followed by descriptions in this format:
- JIRA key (e.g., PLAT-192) on the first line
- Empty line
- Description text on subsequent lines

## Examples

### Example 1: Basic JIRA key with description

**Input:**
```
PLAT-192

blinc - webhook proxy logs
```

**Expected Output:**
```
[PLAT-192: blinc - webhook proxy logs](https://companycam.atlassian.net/browse/PLAT-192)
```

### Example 2: SPEED project with description

**Input:**
```
SPEED-456

Optimize database query performance
```

**Expected Output:**
```
[SPEED-456: Optimize database query performance](https://companycam.atlassian.net/browse/SPEED-456)
```

### Example 3: Multi-line description

**Input:**
```
PLAT-789

Fix authentication issue with SSO
Additional details about the bug
```

**Expected Output:**
```
[PLAT-789: Fix authentication issue with SSO Additional details about the bug](https://companycam.atlassian.net/browse/PLAT-789)
```

## Notes

- Only configured JIRA projects (PLAT, SPEED) will be transformed
- Unconfigured projects will be output verbatim
- The description is concatenated with spaces if it spans multiple lines
- Simple JIRA keys without descriptions continue to work as before