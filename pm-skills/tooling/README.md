# Tooling

Custom Go binaries providing direct access to Jira, Confluence, and web content — both as standalone CLI tools and as an MCP server for Claude Desktop / Claude web.

## MCP Server (`mcp/`)

`pm-tools-mcp` is a unified MCP server that exposes all tools below to **Claude Desktop**, **Claude Code**, or any MCP-compatible client. It is the recommended way to share this tooling with teammates.

**Tools exposed:**

| Tool | Description |
|---|---|
| `jira_get` | Get a Jira issue by key (full details + comments) |
| `jira_search` | JQL search — returns a Markdown table |
| `confluence_search` | Full-text search across Confluence pages |
| `confluence_fetch` | Fetch a Confluence page by ID (cleaned text) |
| `web_fetch` | Fetch any URL (text / html / json output) |

**Build:**
```bash
cd mcp && ./build.sh
# or: go build -o pm-tools-mcp server.go
```

**Claude Desktop setup** (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):
```json
{
  "mcpServers": {
    "pm-tools": {
      "command": "/absolute/path/to/pm-tools-mcp",
      "env": {
        "JIRA_URL": "https://jira.example.com",
        "JIRA_API_TOKEN": "your-token",
        "CONFLUENCE_URL": "https://confluence.example.com",
        "CONFLUENCE_EMAIL": "you@example.com",
        "CONFLUENCE_API_TOKEN": "your-token"
      }
    }
  }
}
```

A ready-to-edit template is at `mcp/claude_desktop_config.json`.

**Credentials:** The MCP server loads credentials from environment variables passed in the config above, or from the nearest `.env` file in the working directory (same lookup as the CLI tools).

---

## Installation (CLI tools)

All binaries are pre-compiled for Apple Silicon (ARM64). To rebuild:

```bash
# Jira
cd jira && go build -o jira jira.go

# Confluence
cd confluence && go build -o confluence confluence.go

# Web
cd web && go build -o web web.go

# MCP server
cd mcp && go build -o pm-tools-mcp server.go
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
