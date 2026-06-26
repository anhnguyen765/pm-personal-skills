# PM Skills Pack

A comprehensive suite of AI-powered Product Management tools and automations for teams using Jira, Confluence, and document-driven workflows.

## Quick Links
- **[Setup Guide](./SETUP.md)** — Installation, configuration, and troubleshooting
- **[Tooling Reference](./tooling/README.md)** — CLI utilities for Jira, Confluence, and web access
- **[Skill Documentation](./skills/)** — Individual skill guides and usage examples

## Available Skills

Invoke skills with `/` in Claude Code:

| Skill | Command | Purpose |
|-------|---------|---------|
| **Jira** | `/jira` | Velocity analysis, sprint reports, issue grooming |
| **PRD** | `/prd` | Generate system-oriented Product Requirements Documents |
| **Product Audit** | `/product-audit` | Quality and compliance auditing with Excel export |
| **Advisor** | `/advisor` | Synthesise Confluence documentation for architectural guidance |
| **Memo** | `/memo` | Socratic debate + executive decision memos |
| **Grilling** | `/grilling` | Stress-test assumptions and plans before building |
| **ASCII** | `/ascii` | Professional ASCII diagrams for system flows |
| **Translator** | `/translator` | Translate features to technical, product, and business views |
| **Data Analysis** | `/data-analysis` | Synthesise KPIs and identify actionable insights |
| **Rapid** | `/rapid` | Weekly activity reports and team summaries |

## Quick Start

```bash
# 1. See SETUP.md for installation and configuration
# 2. Create your .env file with Jira/Confluence credentials
# 3. Invoke a skill:
/jira          # Jira velocity analysis
/prd           # Generate a PRD
/ascii         # Create ASCII diagrams
```

## Core Capabilities

- **Jira Automation**: Velocity analysis, sprint reports, daily standups, blockers, grooming
- **PRD Generation**: System-oriented PRDs with technical completeness and edge-case handling
- **Product Auditing**: Quality, compliance, and operational readiness evaluation
- **Architecture Guidance**: Synthesised Confluence documentation and patterns
- **Decision Making**: Socratic debate and executive memo drafting
- **Stakeholder Communication**: Technical, product, and business translations
- **Reporting**: Weekly activity summaries, KPI analysis, trend identification

## Project Structure

```
pm-skills/
├── SETUP.md                    # Setup guide (start here!)
├── README.md                   # This file — overview
├── skills/                     # Individual PM skills
│   ├── jira/
│   ├── prd/
│   ├── product-audit/
│   ├── advisor/
│   ├── memo/
│   └── ... (see SETUP.md for full list)
├── tooling/                    # CLI utilities for Jira, Confluence, web
│   ├── jira/                   # Jira binaries
│   ├── confluence/             # Confluence binaries
│   └── web/                    # Web fetcher
├── output/
│   ├── sandbox/                # Staging area for artifacts
│   └── reports/                # Final published reports
└── bundles/                    # Pre-packaged skill artifacts
```

## Setup & Configuration

See the **[Setup Guide](./SETUP.md)** for:
- Installation instructions
- Jira and Confluence credential configuration
- Skill usage examples
- Troubleshooting
- Common workflows

## Sandbox Environment

All artifacts are staged in `output/sandbox/` before publication. This allows for review and refinement before sharing externally.

```bash
# Example: Generate a daily report in sandbox
./pm-skills/tooling/jira/daily_report

# Review in output/sandbox/
# Then move to output/reports/ when ready
```

## License & Attribution

This PM Skills Pack is provided for Product Management teams. Ensure compliance with your organisation's policies when sharing or distributing artifacts.

---

**For detailed setup instructions, see [SETUP.md](./SETUP.md).**
