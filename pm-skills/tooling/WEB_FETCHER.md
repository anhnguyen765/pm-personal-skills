# Web Fetcher Tool

A lightweight CLI for fetching and parsing web content from arbitrary URLs.

## Quick Start

```bash
# Fetch and clean HTML to plain text
./web/web fetch "https://docs.example.com/guide"

# Fetch raw HTML
./web/web fetch "https://example.com" --format html

# Fetch and parse JSON
./web/web fetch "https://api.example.com/v1/data" --format json

# Custom timeout
./web/web fetch "https://slow-api.example.com" --timeout 60
```

## Features

- **Multiple Output Formats**: text (default), html, raw JSON
- **Smart Cleaning**: Removes scripts, styles, HTML tags, and collapses whitespace
- **Metadata Extraction**: Pulls title, meta description, and links
- **Customizable**: User-Agent, timeout, and optional caching
- **Zero Dependencies**: Pure Go, no external tools needed

## Configuration

Optional `.env` variables:

```bash
# Override default User-Agent
WEB_USER_AGENT="Custom User Agent 1.0"

# Request timeout (default: 30 seconds)
WEB_TIMEOUT=60

# Cache directory (future feature)
WEB_CACHE_DIR="/tmp/web-cache"
```

## Use Cases

### 1. Document Research
```bash
./web/web fetch "https://developer.example.com/api-reference" > /tmp/api_ref.md
cat /tmp/api_ref.md | claude "Summarize the key endpoints"
```

### 2. External Data Integration
```bash
./web/web fetch "https://public-api.example.com/status" --format json | \
  jq '.services[] | select(.status=="critical")'
```

### 3. Multi-source Content Gathering
```bash
echo "# Research Summary" > /tmp/research.md
./web/web fetch "https://docs.example.com/architecture" >> /tmp/research.md
./web/web fetch "https://blog.example.com/latest-post" >> /tmp/research.md
# Feed to Claude for synthesis
```

### 4. Monitoring & Alerts
```bash
# Check service status page and extract key info
./web/web fetch "https://status.example.com" | grep -i "degraded"
```

## Output Format

The tool outputs markdown-formatted text ready for:
- Piping to Claude Code (`claude` command)
- Saving to files for later processing
- Integration with other command-line tools
- Feeding into document generation workflows

**Example output:**
```
# Page Title

**URL:** https://example.com | **Content-Type:** text/html

**Description:** Brief description from meta tag

------------------------------------------------------------

[Cleaned page content...]
```

## Error Handling

- **HTTP errors**: Outputs the error status code and response
- **Timeout**: Respects the `--timeout` flag (max 30s default)
- **Invalid format**: Returns error for JSON parsing if format is json but content isn't valid JSON

## Performance

- Single-threaded, suitable for sequential fetching
- For bulk operations, consider parallel execution:
  ```bash
  cat urls.txt | xargs -I {} ./web/web fetch {} > /tmp/results.md
  ```

## Limitations

- Does not execute JavaScript (static HTML parsing only)
- No cookie/session management
- No proxy support (future feature)
- Respects `robots.txt` conventions only by honour (configure timeout/UA)

## Building from Source

```bash
cd pm-skills/tooling/web
go build -o web web.go
```

Requires Go 1.16+.
