# Tooling

Custom Go binaries providing direct access to Jira, Confluence, and web content without MCP.

## Installation

All binaries are pre-compiled for Apple Silicon (ARM64). To rebuild:

```bash
# Jira
cd jira && go build -o jira jira.go

# Confluence
cd confluence && go build -o confluence confluence.go

# Web
cd web && go build -o web web.go
```

## Jira

Query and fetch Jira issues directly.

**Setup:**
```bash
export JIRA_URL="https://jira.example.com"
export JIRA_API_TOKEN="your_token_here"
```

**Usage:**
```bash
./jira/jira get PCDPC-123
./jira/jira search "sprint = 'PCDPC - Sprint 26.06.A' AND status != Done" --limit 20
```

## Confluence

Search and fetch Confluence pages.

**Setup:**
```bash
export CONFLUENCE_URL="https://confluence.example.com"
export CONFLUENCE_EMAIL="your_email@example.com"
export CONFLUENCE_API_TOKEN="your_token_here"
```

**Usage:**
```bash
./confluence/confluence search "Payment Engine" --space "Payment Engine" --limit 10
./confluence/confluence fetch "123456789"
```

## Web

Fetch and clean content from arbitrary web URLs.

**Features:**
- Fetches HTML and JSON from any URL
- Strips HTML and cleans content
- Extracts title, meta tags, and links
- Supports custom User-Agent and timeout
- Multiple output formats: text (default), html, json

**Optional Setup (via .env):**
```bash
# Custom User-Agent string
WEB_USER_AGENT="Custom Bot 1.0"

# Request timeout in seconds
WEB_TIMEOUT=60

# Cache directory for fetched content
WEB_CACHE_DIR="./cache"
```

**Usage:**

Fetch and clean HTML content (default):
```bash
./web/web fetch "https://example.com"
./web/web fetch "https://docs.example.com/api" --timeout 30
```

Fetch raw HTML:
```bash
./web/web fetch "https://example.com" --format html
```

Fetch and parse JSON:
```bash
./web/web fetch "https://api.example.com/data" --format json
```

**Examples:**

```bash
# Fetch documentation and extract plain text
./web/web fetch "https://docs.zalopay.com/payment-api"

# Fetch API schema as JSON
./web/web fetch "https://api.zalopay.com/v1/schema" --format json

# Fetch with custom timeout
./web/web fetch "https://slow-service.example.com" --timeout 60
```

## Output

All tools output markdown-formatted text suitable for piping into Claude prompts or saving to files.

Example workflow:
```bash
# Fetch documentation
./web/web fetch "https://docs.example.com/spec" > /tmp/doc.md

# Fetch Confluence architecture guide
./confluence/confluence fetch "123456" > /tmp/architecture.md

# Fetch Jira sprint status
./jira/jira search "sprint = 'PCDPC - Sprint 26.06.A'" > /tmp/sprint.md

# Feed into Claude
cat /tmp/*.md | claude
```
