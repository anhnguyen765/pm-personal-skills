# AGENTS.md — Setup & Onboarding Instructions

This file is for **AI agents** (Claude, GPT, Gemini, Copilot, etc.).
If a user says anything like *"set this up for me"*, *"onboard me"*, or *"get this working"*, follow the steps in this file exactly.

---

## What this project is

A suite of **Product Management tools** for ZaloPay's Payment Core team, exposed as:
- **CLI binaries** (Jira, Confluence, web fetcher) for use in terminal and Claude Code
- **An MCP server** (`pm-tools-mcp`) so Claude Desktop, Claude Web, and other AI clients get the same tools

The MCP server exposes five tools:

| Tool | What it does |
|---|---|
| `jira_get` | Fetch a Jira issue by key with full details and comments |
| `jira_search` | Run a JQL query and get a Markdown table of results |
| `confluence_search` | Full-text search across Confluence pages |
| `confluence_fetch` | Fetch a Confluence page by numeric ID |
| `web_fetch` | Fetch any URL as cleaned text, raw HTML, or JSON |

---

## Quick setup (for AI agents)

**Run one command:**
```bash
./setup.sh
```

That is all. The script handles everything below automatically. If the user needs non-interactive mode (credentials are already environment variables), run:
```bash
./setup.sh --non-interactive
```

---

## What `setup.sh` does (step by step)

If you need to perform setup manually — or the script fails and you need to fix it — follow these steps:

### 1. Build the MCP binary

Requirements: Go 1.19+ (`go version` to check).

```bash
cd pm-skills/tooling/mcp
go build -o pm-tools-mcp server.go
cd -
```

The binary will be at `pm-skills/tooling/mcp/pm-tools-mcp`.

### 2. Create the `.env` file

Copy `.env.example` to `.env` and fill in credentials **and domain configuration**. Ask the user for these values — do NOT guess or invent them.

```bash
cp .env.example .env
```

**Credentials** (required):

| Variable | Description |
|---|---|
| `JIRA_URL` | Base URL of the Jira instance, e.g. `https://jira.company.com` |
| `JIRA_EMAIL` | Email used to log into Jira |
| `JIRA_API_TOKEN` | Jira API token (from Jira → Profile → Security → API tokens) |
| `CONFLUENCE_URL` | Base URL of Confluence, e.g. `https://wiki.company.com` |
| `CONFLUENCE_EMAIL` | Email used to log into Confluence |
| `CONFLUENCE_API_TOKEN` | Confluence API token (same profile page as Jira) |

**Domain configuration** (required — tailors skills to the user's context):

| Variable | Description | Example |
|---|---|---|
| `DOMAIN_NAME` | Team/domain name | `Payment Core`, `Cross-border`, `FX Management` |
| `JIRA_PROJECT_KEY` | Primary Jira project key | `PCDPC`, `XB`, `FXM` |
| `CONFLUENCE_SPACES` | Comma-separated Confluence space keys | `PAYMENTS`, `PAYMENTS,ENGINE` |

**Optional:**

| Variable | Description |
|---|---|
| `WEB_USER_AGENT` | Custom User-Agent string for the web fetcher |
| `WEB_TIMEOUT` | Default timeout in seconds for web requests |

### 3. Configure Claude Code (`.claude/settings.json`)

Update the `pm-tools` entry with the **absolute path** to the binary just built.
Use `pwd` to get the current working directory if needed.

The `mcpServers.pm-tools.command` field must be an absolute path. Example using Python:

```bash
python3 - .claude/settings.json "$(pwd)/pm-skills/tooling/mcp/pm-tools-mcp" <<'EOF'
import json, sys
with open(sys.argv[1]) as f: cfg = json.load(f)
cfg.setdefault("mcpServers", {})
cfg["mcpServers"]["pm-tools"] = {
    "command": sys.argv[2],
    "env": {
        "JIRA_URL": "${JIRA_URL}", "JIRA_API_TOKEN": "${JIRA_API_TOKEN}",
        "CONFLUENCE_URL": "${CONFLUENCE_URL}", "CONFLUENCE_EMAIL": "${CONFLUENCE_EMAIL}",
        "CONFLUENCE_API_TOKEN": "${CONFLUENCE_API_TOKEN}"
    }
}
with open(sys.argv[1], "w") as f: json.dump(cfg, f, indent=2); f.write("\n")
print("Updated settings.json")
EOF
```

### 4. Configure Claude Desktop

Platform-specific config file:
- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux:** `~/.config/Claude/claude_desktop_config.json`
- **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

Merge this block into `mcpServers` in that file (create the file if it doesn't exist):

```json
{
  "mcpServers": {
    "pm-tools": {
      "command": "/absolute/path/to/pm-tools-mcp",
      "env": {
        "JIRA_URL": "https://jira.company.com",
        "JIRA_API_TOKEN": "...",
        "CONFLUENCE_URL": "https://wiki.company.com",
        "CONFLUENCE_EMAIL": "user@company.com",
        "CONFLUENCE_API_TOKEN": "..."
      }
    }
  }
}
```

Replace `"command"` with the actual absolute path to the binary.
Fill in the `env` values from the user's credentials (or from `.env`).

### 5. Verify

Send these two JSON-RPC messages to the binary over stdin and confirm it responds:

```bash
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}\n{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}\n' \
  | ./pm-skills/tooling/mcp/pm-tools-mcp
```

Expected: two JSON lines — the first confirms `protocolVersion: "2024-11-05"`, the second lists five tools.

---

## Troubleshooting guide for agents

| Symptom | Fix |
|---|---|
| `go: command not found` | Ask user to install Go from https://go.dev/dl/ |
| `required env var JIRA_URL is not set` | `.env` file is missing or not in the right directory. Check that `.env` exists at the project root |
| MCP server not appearing in Claude Desktop | Restart Claude Desktop after updating the config |
| MCP server not appearing in Claude Code | Start a new Claude Code session in the project directory |
| `HTTP 401` from Jira/Confluence tool calls | API token is wrong or expired — ask user to regenerate it |
| `HTTP 403` from Confluence | User does not have read access to that space |

---

## What NOT to do

- **Do not commit `.env`** — it contains credentials. It is in `.gitignore`.
- **Do not modify `server.go` during setup** — only build it.
- **Do not hard-code credentials** into any tracked file.
- **Do not run `setup.sh` with `sudo`** — it does not require elevated permissions.

---

## After setup

Tell the user:
1. **Claude Code:** restart the session. The `pm-tools` MCP server will be active.
2. **Claude Desktop:** quit and reopen the app. The five tools will appear.
3. **Verify credentials** by running a test search: ask the AI to call `jira_search` with `project = <their-project-key>`.
