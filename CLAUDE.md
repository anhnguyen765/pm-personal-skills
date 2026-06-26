# Claude Project Instructions: Jora

This project is a personalised suite of Product Management (PM) skills and tools. It provides automated workflows for Jira, Confluence, and document generation for the **ZaloPay Payment Core** platform.

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
| `/jira` | Velocity analysis, stuck task resolution, and grooming-ready ticket writing for Payment Core |
| `/prd` | Generate and refine technical PRDs for Fintech Core Payment platforms |
| `/data-analysis` | Synthesise platform metrics, cross-border KPIs, and identify next actionable items for Payment Core and Cross-border PO |
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
- Prefer MCP tools for single issue lookups and sprint searches.
- Use `pm-skills/tooling/jira/daily_report` for bulk daily reporting.
- Sprint prefix is always **"PCDPC"** (e.g., "PCDPC - Sprint 26.06.A").

### Confluence Research
- Search via MCP: `mcp__atlassian__confluence_search`
- Priority spaces: "Zalopay Technology Management", "Payment Engine".
- High-signal keywords: `FMC/PE`, `AEv2`, `ZAS`, `Idempotency`, `ACv2`.

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
