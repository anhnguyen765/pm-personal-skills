# PM Agents Documentation

This document provides a consolidated overview of the project instructions, active MCP servers, and available specialized skills for the Jora PM suite.

Supported AI runtimes: **Gemini CLI** and **Claude Code**.

---

## Claude Code Setup

### Project Instructions
See `jora/CLAUDE.md` — loaded automatically by Claude Code in this workspace.

### Custom Slash Commands
Skills are available as slash commands in Claude Code via `.claude/commands/`:

| Command | Skill |
|---|---|
| `/jira` | Velocity analysis + grooming-ready ticket writing |
| `/prd` | Technical PRD generation for Fintech Core Payment |
| `/advisor` | Confluence-grounded architectural guidance |
| `/memo` | Socratic debate + executive decision memo |
| `/ascii` | SA ASCII diagrams for system flows |
| `/product-audit` | Product quality + compliance auditing |
| `/translator` | Multi-view stakeholder communication |

### MCP Servers
Configured in `.claude/settings.json` (mirrors `.gemini/settings.json`):
- **atlassian**: `npx -y @atlassian/atlassian-mcp-server` — Jira & Confluence tools
- **fetch**: `npx -y @modelcontextprotocol/server-fetch` — URL fetching

---

## 1. Project Instructions (GEMINI.md)

# Gemini Project Instructions: Jora

This project is a personalised suite of Product Management (PM) skills and tools. It provides automated workflows for Jira, Confluence, and document generation.

## Project Structure
- `.gemini/skills/`: Active skills used by the Gemini CLI.
- `pm-skills/skills/`: Source code for PM skills.
- `pm-skills/tooling/`: Custom automation scripts (Jira/Confluence).
- `pm-skills/bundles/`: Packaged `.skill` files.
- `pm-skills/output/`: Default directory for generated artifacts (Memos, PRDs).

## Workflows

### 1. Skill Development
- When creating or modifying a skill, work in `pm-skills/skills/`.
- Use `./pm-skills/scripts/bundle_skills.sh` to package changes.
- Test in the sandbox folder (`pm-skills/output/sandbox/`).

### 2. MCP Integration (Jira & Confluence)
- This project uses **Model Context Protocol (MCP)** servers to interact with Jira and Confluence.
- **Benefits**: Robust link fetching, real-time search, and standardized tool schemas.
- **Configuration**: Defined in `.gemini/settings.json`.
- **Credentials**: Loaded from the repo-root `.env` file via environment variables.
- **Tools**: Prefixed with `mcp_atlassian_` (e.g., `mcp_atlassian_jira_get_issue`).

### 3. Jira Automation
- Use the scripts in `pm-skills/tooling/jira/` for daily reports and complex sprint analysis.
- Basic issue lookups and searches should prefer MCP tools.

### 4. Solution Architecture (SA) Diagrams
- Use the `ascii` skill for all high-level system flow visualisations.
- Follow the standardised shapes (`[ Component ]`, `( Database )`, `< Broker >`).

## Engineering Standards
- **Tone**: Professional, direct, and concise.
- **Language**: Australian English for all outputs (summarise, prioritise).
- **Documentation**: All new features or tools must be documented in `pm-skills/README.md`.

---

## 2. MCP Servers Configuration

The following MCP servers are configured in `.gemini/settings.json`:

### Atlassian (Jira & Confluence)
- **Command**: `npx -y @atlassian/atlassian-mcp-server`
- **Description**: Provides robust integration with Jira and Confluence, including issue management, page retrieval, and search capabilities.
- **Tools Prefix**: `mcp_atlassian_`

### Fetch
- **Command**: `npx -y @modelcontextprotocol/server-fetch`
- **Description**: Allows for fetching content from web URLs.

---

## 3. Specialized PM Skills

### advisor
**Description**: Synthesises Confluence technical documentation to provide guidance, architectural patterns, and system overviews.

#### Technical Advisor Skill
This skill acts as a Technical Advisor by synthesising architectural and technical information from Confluence to provide guidance, best practices, and system overviews.

**Workflow:**
1. **Identify Domain**: Determine the technical area the user is inquiring about.
2. **Search Confluence**: Use MCP tools (`mcp_atlassian_confluence_search`).
3. **Fetch Details**: Retrieve full content (`mcp_atlassian_confluence_get_page_content`).
4. **Advise**: Synthesise information into clear recommendations.

---

### ascii
**Description**: Create professional ASCII diagrams for system architecture and component interactions.

#### SA ASCII Diagrams
Standardised way to draw high-signal, professional ASCII diagrams mimicking Solution Architecture (SA) interaction diagrams.

**Shapes:** `[ Component ]`, `( Database )`, `< Broker >`, `[[ External ]]`.

---

### jira
**Description**: Unified Jira skill for velocity analysis, stuck task resolution, and writing technical grooming-ready ticket descriptions.

#### Jira Master
Combines advanced velocity analysis with technical ticket writing expertise for Payment Core.

**Capabilities:**
1. **Velocity Management**: Identifying bottlenecks and unblocking tasks.
2. **Technical Ticket Writing**: Drafting grooming-ready tickets for Fund Movement Core and Payment Engine.

---

### memo
**Description**: Assists in strategic product decision-making by applying Socratic debating skills to challenge assumptions.

#### Strategic Memo Helper
Acts as a Senior Strategic PO and Socratic Debater to pressure-test strategic thinking before drafting formal memos.

**Phases:** Socratic Debate, Memo Synthesis, Sandbox Review, Finalization.

---

### prd
**Description**: Generate and refine system-oriented Product Requirements Documents (PRD) for Fintech Core Payment platforms.

#### PRD Writer
Generates technical PRDs focused on transaction lifecycles, ledger consistency, and failure handling in the ZaloPay context.

---

### product-audit
**Description**: Systematically evaluate product quality, operational readiness, and compliance using the internal template.

#### Product Audit Guide
Systematic auditing through the "System-as-User" lens, anchoring findings in technical documentation.

---

### translator
**Description**: Translates a single topic into technical, product, and business views for stakeholder alignment.

#### Multi-View Translator
Bridges communication gaps between technical teams, product stakeholders, and business leadership by presenting three distinct perspectives.
