# Claude Project Instructions: PM Skills Suite

This project is a **customisable suite of Product Management (PM) skills and tools**. It provides automated workflows for Jira, Confluence, and document generation tailored to your domain (e.g., Payment Core, Cross-border, FX Management, etc.).

## Getting Started

**Quick setup (recommended):**
```bash
./setup.sh
```

The script will:
- Build the `pm-tools-mcp` binary (provides 5 MCP tools for Claude Desktop, Code, and web)
- Create `.env` from your credentials
- Wire up `.claude/settings.json` and Claude Desktop config
- Ask you which domain/project you are working in (so skills are contextualised correctly)

**Manual setup:**
1. Edit `.env` with your Jira/Confluence credentials and domain scope
2. Update `.claude/settings.json` to point `pm-tools.command` to your local `pm-skills/tooling/mcp/pm-tools-mcp` binary (use absolute path)
3. For Claude Desktop: merge `pm-skills/tooling/mcp/claude_desktop_config.json` into `~/Library/Application Support/Claude/claude_desktop_config.json`
4. Restart Claude Code / Desktop

See `AGENTS.md` for AI-assisted onboarding.

## Project Structure
- `.claude/commands/`: Custom slash commands for each PM skill (invoke with `/skill-name`).
- `pm-skills/skills/`: Source SKILL.md definitions and assets for each PM skill.
- `pm-skills/tooling/`: Custom automation scripts (Jira reports, Confluence fetcher).
- `pm-skills/bundles/`: Packaged `.skill` files for distribution.
- `pm-skills/output/`: Default directory for generated artifacts (Memos, PRDs, Audit reports).
- `pm-skills/output/sandbox/`: Staging area before final publication.

## Available Skills (Slash Commands)
Invoke these with `/` prefix in Claude Code:

| Command | Description |
|---|---|
| `/jira` | Velocity analysis, stuck task resolution, and grooming-ready ticket writing (configured for your domain during setup) |
| `/prd` | Generate and refine technical PRDs for your platform |
| `/data-analysis` | Synthesise platform metrics, KPIs, and identify next actionable items (configured for your domain) |
| `/advisor` | Synthesise Confluence technical documentation for architectural guidance |
| `/memo` | Socratic debate + executive decision memo drafting |
| `/grilling` | Interview relentlessly about a plan or design — stress-test assumptions before building (integrates with memo/prd) |
| `/grill-me` | Quick trigger for a grilling session to sharpen your thinking before finalising a memo or PRD |
| `/ascii` | Professional ASCII diagrams for system architecture and component flows |
| `/product-audit` | Systematic product quality and compliance auditing using the internal Excel template |
| `/translator` | Translate a single topic into technical, product, and business stakeholder views |
| `/verdict` | Executive product brief — synthesise execution and strategy signals for PO and leadership |
| `/new-skill` | Reference for writing and editing PM skills well — design, structure, and principles for predictable automation |

## Tooling
Go binaries in `pm-skills/tooling/` provide direct Jira, Confluence, and web access via `.env` credentials — no MCP required.

### Jira & Confluence
| Command | Description |
|---|---|
| `./pm-skills/tooling/jira/jira get <key>` | Get full issue details |
| `./pm-skills/tooling/jira/jira search "<jql>" [--limit N]` | JQL search |
| `./pm-skills/tooling/confluence/confluence search "<q>" [--space KEY] [--limit N]` | Full-text search |
| `./pm-skills/tooling/confluence/confluence fetch <page_id>` | Get page content |

Credentials come from `.env` (see `.env.example`): `JIRA_URL`, `JIRA_API_TOKEN`, `CONFLUENCE_URL`, `CONFLUENCE_EMAIL`, `CONFLUENCE_API_TOKEN`.

### Web Fetcher
| Command | Description |
|---|---|
| `./pm-skills/tooling/web/web fetch <url>` | Fetch and clean web content (default: plain text) |
| `./pm-skills/tooling/web/web fetch <url> --format html` | Fetch raw HTML |
| `./pm-skills/tooling/web/web fetch <url> --format json` | Fetch and parse JSON |
| `./pm-skills/tooling/web/web fetch <url> --timeout SEC` | Custom timeout (default: 30s) |

Optional `.env` variables: `WEB_USER_AGENT`, `WEB_TIMEOUT`, `WEB_CACHE_DIR`.

See `pm-skills/tooling/WEB_FETCHER.md` for detailed usage and examples.

## Workflows

### Skill Development
1. Work in `pm-skills/skills/<skill-name>/SKILL.md`.
2. Mirror changes to `.claude/commands/<skill-name>.md`.
3. Run `./pm-skills/scripts/bundle_skills.sh` to package.
4. Test artifacts in `pm-skills/output/sandbox/`.

### Jira Automation
- Use the `jira_get` and `jira_search` MCP tools for issue lookups and JQL searches.
- Domain context is loaded from `.env` during setup (e.g., your project key, board ID, team prefix).
- JQL examples: `project = YOUR_PROJECT AND status = 'In Progress'`, `assignee = currentUser() AND updated >= -7d`.

### Confluence Research
- Use the `confluence_search` and `confluence_fetch` MCP tools to find and retrieve pages.
- Spaces and keywords are configured in `.env` based on your domain (set during setup).

### SA Diagrams
- Use the `/ascii` skill for all system flow visualisations.
- Shapes: `[ Component ]`, `( Database )`, `< Broker >`, `[[ External ]]`.

## Engineering Standards
- **Tone**: Professional, direct, and concise.
- **Language**: Australian English (summarise, prioritise, centralised).
- **Documentation**: New tools or features must be documented in `pm-skills/README.md`.
- **Sandbox first**: Always stage artifacts in `pm-skills/output/sandbox/` before final publication.

## Common Commands
```bash
# Daily Jira report
./pm-skills/tooling/jira/daily_report

# Sprint velocity analysis
./pm-skills/skills/jira/scripts/velocity_analyzer

# Fetch Confluence page
./pm-skills/tooling/confluence/fetch_full <page_id>

# Fetch and clean web content
./pm-skills/tooling/web/web fetch "https://example.com"
./pm-skills/tooling/web/web fetch "https://api.example.com/data" --format json

# Product audit Excel (example: FX Management)
go run pm-skills/skills/product-audit/scripts/excel_helper.go \
  -mode write -json-file pm-skills/skills/product-audit/scripts/fx_management_audit.json
```
